import { postService } from "@/api/api.service";
import { uploadMedia } from "@/api/upload.service";
import { Button } from "@/components/ui/button";
import { ButtonGroup } from "@/components/ui/button-group";
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Input } from "@base-ui/react/input";
import {
  DndContext,
  DragOverlay,
  KeyboardSensor,
  PointerSensor,
  pointerWithin,
  useSensor,
  useSensors,
  type DragEndEvent,
  type DragStartEvent,
} from "@dnd-kit/core";
import { restrictToVerticalAxis } from "@dnd-kit/modifiers";
import {
  SortableContext,
  sortableKeyboardCoordinates,
  verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import { BlockType, PostStatus } from "@gen/api/postpb/v1/post";
import { cn } from "cn";
import { GripVertical, Redo, Undo } from "lucide-react";
import {
  useEffect,
  useReducer,
  useRef,
  useState,
  type ClipboardEvent,
  type KeyboardEvent,
} from "react";
import { useNavigate } from "react-router-dom";
import { AddBlockButtons, SortableBlock } from "./editor-block";
import {
  mapRemoteToDoc,
  toCreateRequest,
  toUpdateRequest,
  type EnrichedGetResponse,
} from "./mapper";
import {
  MAX_IMAGE_BYTES,
  imageSizeItems,
  newMedia,
  newText,
  reducer,
  textSizeItems,
  type Action,
  type Doc,
  type MediaBlock,
  type TextBlock,
  type TextSizeType,
} from "./model";

function validate(doc: Doc): string | null {
  if (!doc.title.trim()) return "Add a title.";

  if (!/^[a-z0-9-]+$/.test(doc.slug)) return "Only lowercase letters, numbers or dashes in slug.";

  const mediaBlocks = doc.blocks.filter((b) => b.blockType === BlockType.BLOCK_TYPE_MEDIA);

  if (mediaBlocks.some((m) => m.status === "uploading"))
    return "Please wait for medias to finish uploading.";

  if (mediaBlocks.some((m) => m.status === "error"))
    return "Retry or remove the image that failed to upload.";

  if (!doc.blocks.some((b) => b.blockType === BlockType.BLOCK_TYPE_MEDIA || b.text.trim()))
    return "Add some content.";

  return null;
}

function focusBlock(root: HTMLElement | null, id: string) {
  const elem = root?.querySelector<HTMLInputElement | HTMLTextAreaElement>(`[data-focus="${id}"]`);
  elem?.focus();
  elem?.setSelectionRange(elem.value.length, elem.value.length);
}

export function Editor({ remoteResponse }: { remoteResponse?: EnrichedGetResponse }) {
  const navigate = useNavigate();

  const [state, dispatch] = useReducer(reducer, undefined, () => ({
    past: [],
    future: [],
    key: null,
    present: remoteResponse
      ? mapRemoteToDoc(remoteResponse)
      : {
          title: "",
          slug: "",
          blocks: [newText()],
        },
  }));
  const { present: doc, past, future } = state;
  const [error, setError] = useState<string | null>(null);
  const [busy, setBusy] = useState(false);
  const [activeId, setActiveId] = useState<string | null>(null);
  const [focusedId, setFocusedId] = useState<string | null>(null);
  const focusedBlock = doc.blocks.find((b) => b.id === focusedId);
  const activeBlock = doc.blocks.find((b) => b.id === activeId);

  const root = useRef<HTMLDivElement>(null);
  const fileInput = useRef<HTMLInputElement>(null);
  const insertAt = useRef(0);
  const focusAfterRender = useRef<string | null>(null);
  const blobUrls = useRef<string[]>([]);

  const media = doc.blocks.filter((b) => b.blockType === BlockType.BLOCK_TYPE_MEDIA);
  const uploading = media.filter((b) => b.status === "uploading").length;
  const failed = media.filter((b) => b.status === "error").length;
  const words = doc.blocks
    .filter((b) => b.blockType === BlockType.BLOCK_TYPE_TEXT)
    .reduce((n, b) => n + (b.text.match(/\S+/g)?.length ?? 0), 0);

  useEffect(() => {
    const id = focusAfterRender.current;
    focusAfterRender.current = null;
    if (id) focusBlock(root.current, id);
  }, [doc]);

  const unsaved = uploading > 0 || past.length > 0;
  useEffect(() => {
    if (!unsaved) return;
    const warn = (e: BeforeUnloadEvent) => e.preventDefault();

    window.addEventListener("beforeunload", warn);
    return () => window.removeEventListener("beforeunload", warn);
  }, [unsaved]);

  useEffect(() => {
    const urls = blobUrls.current;
    return () => urls.forEach(URL.revokeObjectURL);
  }, []);

  useEffect(() => {
    function onDocumentKeyDown(e: globalThis.KeyboardEvent) {
      if (document.activeElement && !root.current?.contains(document.activeElement)) return;

      if (!(e.metaKey || e.ctrlKey)) return;

      const k = e.key.toLowerCase();
      if (k !== "z" && k !== "y") return;
      e.preventDefault();
      dispatch({ type: k === "y" || e.shiftKey ? "redo" : "undo" });
    }

    window.addEventListener("keydown", onDocumentKeyDown);
    return () => window.removeEventListener("keydown", onDocumentKeyDown);
  }, []);

  async function startUpload(b: MediaBlock) {
    if (b.file === undefined) return;
    dispatch({
      type: "upload",
      id: b.id,
      status: "uploading",
    });
    try {
      const res = await uploadMedia(b.id, b.file);
      dispatch(
        res.ok
          ? { type: "upload", id: b.id, status: "ready" }
          : { type: "upload", id: b.id, status: "error", error: res.message },
      );
    } catch {
      dispatch({
        type: "upload",
        id: b.id,
        status: "error",
        error: "Network error",
      });
    }
  }

  function addFiles(files: File[], index: number) {
    setError(null);
    const added: MediaBlock[] = [];
    files.forEach((file) => {
      if (!file.type.startsWith("image/")) {
        setError(`${file.name} is not an image.`);
        return;
      }

      if (file.size > MAX_IMAGE_BYTES) {
        setError(`${file.name} is over ${MAX_IMAGE_BYTES >> 20}MB.`);
        return;
      }

      const src = URL.createObjectURL(file);
      blobUrls.current.push(src);
      added.push(newMedia(file, src));
    });

    if (added.length === 0) return;

    dispatch({
      type: "insert",
      index,
      blocks: added,
    });

    added.forEach(startUpload);
  }

  function addText(i: number) {
    const b = newText();
    focusAfterRender.current = b.id;
    dispatch({
      type: "insert",
      index: i,
      blocks: [b],
    });
  }

  async function submit(publish: boolean) {
    if (busy) {
      return;
    }

    const problem = validate(doc);
    if (problem) return setError(problem);

    setBusy(true);
    setError(null);

    if (remoteResponse === undefined) {
      try {
        const res = await postService.Create(toCreateRequest(doc, publish));
        if (res.ok) {
          navigate(`/post/${doc.slug}`);
          return;
        }

        switch (res.message) {
          case "post.validation.SlugConflict":
            setError("Slug already exists. Please choose another one.");
            break;
          default:
            setError(`couldn't save the post (${res.message}).`);
        }
      } catch {
        setError("couldn't reach the server. Try again later.");
      } finally {
        setBusy(false);
      }
    } else {
      try {
        const res = await postService.Update(
          toUpdateRequest(remoteResponse.post!, remoteResponse.postBlocks, doc, publish),
        );
        if (res.ok) navigate(`/post/${doc.slug}`);
        else setError(`couldn't save the post (${res.message}).`);
      } catch {
        setError("couldn't reach the server. Try again later.");
      } finally {
        setBusy(false);
      }
    }
  }

  function onTextKeyDown(e: KeyboardEvent<HTMLTextAreaElement>, i: number) {
    if (e.nativeEvent.isComposing) return;
    const el = e.currentTarget;
    const caretAtStart = el.selectionStart === 0 && el.selectionEnd === 0;
    const caretAtEnd = el.selectionStart === el.value.length && el.selectionEnd === el.value.length;

    if (e.key === "Enter" && !e.shiftKey && caretAtEnd) {
      e.preventDefault();
      const b = newText();
      focusAfterRender.current = b.id;
      dispatch({
        type: "insert",
        index: i + 1,
        blocks: [b],
      });
      return;
    }

    if (e.key === "Backspace" && el.value === "" && doc.blocks.length > 1) {
      e.preventDefault();
      focusAfterRender.current = (doc.blocks[i - 1] ?? doc.blocks[i + 1]).id;
      dispatch({
        type: "remove",
        id: doc.blocks[i].id,
      });
    }

    if (e.key === "ArrowUp" && caretAtStart && i > 0) {
      e.preventDefault();
      focusBlock(root.current, doc.blocks[i - 1].id);
    }

    if (e.key === "ArrowDown" && caretAtEnd && i < doc.blocks.length - 1) {
      e.preventDefault();
      focusBlock(root.current, doc.blocks[i + 1].id);
    }
  }

  function onPaste(e: ClipboardEvent) {
    const images = Array.from(e.clipboardData.files).filter((f) => f.type.startsWith("image/"));
    if (images.length === 0) return;
    e.preventDefault();
    const at = (e.target as HTMLElement).closest<HTMLElement>("[data-index]")?.dataset.index;
    addFiles(images, at === undefined ? doc.blocks.length : Number(at) + 1);
  }

  const sensors = useSensors(
    useSensor(PointerSensor, { activationConstraint: { distance: 5 } }),
    useSensor(KeyboardSensor, { coordinateGetter: sortableKeyboardCoordinates }),
  );

  function onDragEnd({ active, over }: DragEndEvent) {
    setActiveId(null);
    if (!over || active.id === over.id) return;
    const from = doc.blocks.findIndex((b) => b.id === active.id);
    const to = doc.blocks.findIndex((b) => b.id === over.id);
    if (from === -1 || to === -1) return;
    dispatch({ type: "move", from, to });
  }

  return (
    <div
      ref={root}
      className="mx-auto flex w-full max-w-2xl flex-col items-center pb-12"
      onPaste={onPaste}
      onDragOver={(e) => e.dataTransfer.types.includes("Files") && e.preventDefault()}
      onDrop={(e) => {
        if (e.dataTransfer.files.length === 0) return;
        e.preventDefault();
        addFiles(Array.from(e.dataTransfer.files), doc.blocks.length);
      }}
    >
      <div className="bg-background sticky top-0 z-5 w-full border-b px-2 md:px-0">
        <Input
          dir="auto"
          className="placeholder:text-muted-foreground/60 mt-10 resize-none text-3xl font-bold outline-none"
          aria-label="Post title"
          placeholder="Title"
          maxLength={200}
          value={doc.title}
          onChange={(e) =>
            dispatch({
              type: "title",
              title: e.target.value,
            })
          }
          onKeyDown={(e) => {
            if (e.key === "Enter" && !e.nativeEvent.isComposing) {
              e.preventDefault();
              focusBlock(root.current, doc.blocks[0].id);
            }
          }}
        />
        <label
          dir="ltr"
          className="text-muted-foreground flex w-full"
        >
          /
          <Input
            className="flex-1 resize-none text-xs outline-none"
            aria-label="Post slug"
            placeholder="post-slug"
            maxLength={200}
            value={doc.slug}
            onChange={(e) =>
              dispatch({
                type: "slug",
                slug: e.target.value.toLowerCase().replace(/[^a-z0-9-]/g, ""),
              })
            }
          />
        </label>
        <div className="flex w-full flex-col gap-2 p-2">
          <div className="flex w-full flex-wrap justify-between gap-y-3">
            <div className="flex items-center gap-2">
              <ButtonGroup>
                <Button
                  variant={"outline"}
                  onClick={() => dispatch({ type: "undo" })}
                  disabled={past.length === 0}
                >
                  Undo <Undo />
                </Button>
                <Button
                  variant={"outline"}
                  onClick={() => dispatch({ type: "redo" })}
                  disabled={future.length === 0}
                >
                  Redo <Redo />
                </Button>
              </ButtonGroup>
              {focusedBlock?.blockType === BlockType.BLOCK_TYPE_TEXT && (
                <TextSizeSelect
                  focusedBlock={focusedBlock}
                  dispatch={dispatch}
                />
              )}
              {focusedBlock?.blockType === BlockType.BLOCK_TYPE_MEDIA && (
                <ImageSizeSelect
                  focusedBlock={focusedBlock}
                  dispatch={dispatch}
                />
              )}
            </div>
            <div className="flex items-center gap-1">
              <p
                aria-live="polite"
                className={cn(
                  "text-muted-foreground min-w-0 flex-1 truncate px-2 text-end text-xs",
                  failed > 0 && "text-destructive",
                )}
              >
                {uploading > 0
                  ? `Uploading ${uploading} image${uploading > 1 ? "s" : ""}…`
                  : failed > 0
                    ? `${failed} upload${failed > 1 ? "s" : ""} failed`
                    : `${words} word${words === 1 ? "" : "s"}`}
              </p>
              {remoteResponse?.post?.postStatus !== PostStatus.POST_STATUS_PUBLISHED && (
                <Button
                  variant={"outline"}
                  onClick={() => submit(false)}
                  disabled={busy}
                >
                  {remoteResponse === undefined ? "Save draft" : "Update draft"}
                </Button>
              )}
              <Button
                onClick={() => submit(true)}
                disabled={busy}
              >
                {remoteResponse === undefined
                  ? "Publish"
                  : remoteResponse.post?.postStatus === PostStatus.POST_STATUS_PUBLISHED
                    ? "Update"
                    : "Publish"}
              </Button>
            </div>
          </div>
          <p
            className="text-destructive/80 text-sm"
            aria-live="polite"
          >
            {error}
          </p>
        </div>
      </div>
      <DndContext
        sensors={sensors}
        collisionDetection={(p) => {
          if (p.pointerCoordinates) {
            p.pointerCoordinates.x = p.droppableContainers[0].rect.current?.left ?? 0;
            const topY = p.droppableContainers.reduce(
              (y, c) => Math.min(y, c.rect.current?.top ?? 0),
              999999999,
            );
            const botY = p.droppableContainers.reduce(
              (y, c) => Math.max(y, c.rect.current?.bottom ?? 0),
              -999999999,
            );

            p.pointerCoordinates.y = Math.max(topY, Math.min(botY, p.pointerCoordinates.y));
          }

          return pointerWithin(p);
        }}
        onDragStart={({ active }: DragStartEvent) => {
          setActiveId(String(active.id));
        }}
        onDragEnd={onDragEnd}
        onDragCancel={() => setActiveId(null)}
      >
        <SortableContext
          items={doc.blocks.map((b) => b.id)}
          strategy={verticalListSortingStrategy}
        >
          <ul className="w-full gap-1">
            <AddBlockButtons
              onTextAdd={() => addText(0)}
              onImageAdd={() => {
                insertAt.current = 0;
                fileInput.current?.click();
              }}
            />
            {doc.blocks.map((block, index) => (
              <SortableBlock
                index={index}
                block={block}
                dispatch={dispatch}
                onTextKeyDown={onTextKeyDown}
                addText={addText}
                insertAt={insertAt}
                fileInput={fileInput}
                key={block.id}
                startUpload={startUpload}
                setFocusedId={setFocusedId}
              />
            ))}
          </ul>
        </SortableContext>

        <DragOverlay modifiers={[restrictToVerticalAxis]}>
          {activeBlock ? (
            <div className="bg-background flex w-full rounded-md p-1 opacity-90 shadow-xl">
              <GripVertical className="text-muted-foreground/50 size-5 shrink-0" />
              {activeBlock.blockType === BlockType.BLOCK_TYPE_TEXT ? (
                <textarea
                  readOnly
                  dir="auto"
                  className={cn(
                    "field-sizing-content w-full resize-none overflow-hidden outline-none",
                    textSizeItems.find((i) => i.value === activeBlock.metadata.size)?.className,
                  )}
                  placeholder="Start typing... (Enter adds a block, Shift+Enter a line break"
                  value={activeBlock.text}
                ></textarea>
              ) : (
                <img
                  src={activeBlock.src}
                  srcSet={activeBlock.srcSet}
                  className={cn(
                    "max-w-full min-w-0 flex-1 rounded-md object-contain pr-8",
                    imageSizeItems.find((i) => i.value === activeBlock.metadata.size)?.className,
                  )}
                />
              )}
            </div>
          ) : null}
        </DragOverlay>
      </DndContext>
      <input
        ref={fileInput}
        type="file"
        accept="image/*"
        multiple
        hidden
        onChange={(e) => {
          addFiles(Array.from(e.target.files ?? []), insertAt.current);
          e.target.value = "";
        }}
      />
    </div>
  );
}

function TextSizeSelect({
  focusedBlock,
  dispatch,
}: {
  focusedBlock: TextBlock;
  dispatch: React.ActionDispatch<[a: Action]>;
}) {
  return (
    <Select
      items={focusedBlock.blockType === BlockType.BLOCK_TYPE_TEXT ? textSizeItems : imageSizeItems}
      value={focusedBlock.metadata.size}
      onValueChange={(value) => {
        if (value === null) return;
        dispatch({
          type: "updateTextMeta",
          id: focusedBlock.id,
          metadata: { size: value as TextSizeType },
        });
      }}
    >
      <SelectTrigger className="w-30 focus-within:opacity-100 hover:opacity-100">
        <SelectValue placeholder="Size" />
      </SelectTrigger>
      <SelectContent>
        <SelectGroup>
          {(focusedBlock.blockType === BlockType.BLOCK_TYPE_TEXT
            ? textSizeItems
            : imageSizeItems
          ).map((item) => (
            <SelectItem
              key={item.value}
              value={item.value}
            >
              {item.label}
            </SelectItem>
          ))}
        </SelectGroup>
      </SelectContent>
    </Select>
  );
}

function ImageSizeSelect({
  focusedBlock,
  dispatch,
}: {
  focusedBlock: MediaBlock;
  dispatch: React.ActionDispatch<[a: Action]>;
}) {
  return (
    <Select
      items={imageSizeItems}
      value={focusedBlock.metadata.size}
      onValueChange={(value) => {
        if (value === null) return;
        dispatch({
          type: "updateImageMeta",
          id: focusedBlock.id,
          metadata: { size: value },
        });
      }}
    >
      <SelectTrigger className="w-30 focus-within:opacity-100 hover:opacity-100">
        <SelectValue placeholder="Size" />
      </SelectTrigger>
      <SelectContent>
        <SelectGroup>
          {imageSizeItems.map((item) => (
            <SelectItem
              key={item.value}
              value={item.value}
            >
              {item.label}
            </SelectItem>
          ))}
        </SelectGroup>
      </SelectContent>
    </Select>
  );
}
