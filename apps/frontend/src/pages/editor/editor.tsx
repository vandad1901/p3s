import { postService } from "@/api/api.service";
import { uploadMedia } from "@/api/upload.service";
import { Button } from "@/components/ui/button";
import { ButtonGroup } from "@/components/ui/button-group";
import { Input } from "@base-ui/react/input";
import {
  closestCenter,
  DndContext,
  DragOverlay,
  KeyboardSensor,
  PointerSensor,
  useSensor,
  useSensors,
  type DragEndEvent,
  type DragStartEvent,
} from "@dnd-kit/core";
import { restrictToVerticalAxis } from "@dnd-kit/modifiers";
import {
  SortableContext,
  sortableKeyboardCoordinates,
  useSortable,
  verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { BlockType, PostStatus, type CreateRequest } from "@gen/api/postpb/v1/post";
import { cn } from "cn";
import { GripVertical, ImagePlus, Loader2, Redo, RotateCw, Trash, Type, Undo } from "lucide-react";
import {
  useEffect,
  useReducer,
  useRef,
  useState,
  type ActionDispatch,
  type ClipboardEvent,
  type KeyboardEvent,
} from "react";
import { useNavigate } from "react-router-dom";

const MAX_IMAGE_BYTES = 20 << 20;

// types

type TextBlock = {
  id: string;
  blockType: BlockType.BLOCK_TYPE_TEXT;
  text: string;
};

type MediaBlock = {
  id: string;
  blockType: BlockType.BLOCK_TYPE_MEDIA;
  file: File;
  src: string;
  status: "uploading" | "ready" | "error";
  error?: string;
};

type Block = TextBlock | MediaBlock;

type Doc = {
  title: string;
  slug: string;
  blocks: Block[];
};

// History Management

type State = {
  past: Doc[];
  present: Doc;
  future: Doc[];
  key: string | null;
};
type Action =
  | { type: "title"; title: string }
  | { type: "slug"; slug: string }
  | { type: "text"; id: string; text: string } // editing text block with id
  | { type: "insert"; index: number; blocks: Block[] }
  | { type: "remove"; id: string }
  | { type: "move"; from: number; to: number }
  | { type: "upload"; id: string; status: MediaBlock["status"]; error?: string } // documenting upload status changes
  | { type: "undo" }
  | { type: "redo" };

function commit(s: State, doc: Doc, key: string | null = null): State {
  if (key !== null && key === s.key) return { ...s, present: doc, future: [] };
  return { past: [...s.past, s.present].slice(-100), present: doc, future: [], key };
}

const newText = (): TextBlock => ({
  id: crypto.randomUUID(),
  blockType: BlockType.BLOCK_TYPE_TEXT,
  text: "",
});
const newMedia = (file: File, src: string): MediaBlock => ({
  id: crypto.randomUUID(),
  blockType: BlockType.BLOCK_TYPE_MEDIA,
  file,
  src,
  status: "uploading",
});

function reducer(s: State, a: Action): State {
  const d = s.present;
  const editBlockByID = (doc: Doc, id: string, fn: (b: Block) => Block): Doc => ({
    ...doc,
    blocks: doc.blocks.map((b) => (b.id === id ? fn(b) : b)),
  });

  switch (a.type) {
    case "title": {
      const slugify = (s: string) =>
        s
          .toLowerCase()
          .normalize("NFKD")
          .replace(/[^a-z0-9]+/g, "-")
          .replace(/^-+|-+$/g, "");

      const syncSlug = d.slug === slugify(d.title);
      return commit(
        s,
        { ...d, title: a.title, slug: syncSlug ? slugify(a.title) : d.slug },
        "title",
      );
    }
    case "slug": {
      return commit(s, { ...d, slug: a.slug }, "slug");
    }
    case "text": {
      return commit(
        s,
        editBlockByID(d, a.id, (b) => ({ ...b, text: a.text })),
        `text:${a.id}`,
      );
    }
    case "insert": {
      const i = Math.max(0, Math.min(a.index, d.blocks.length));
      return commit(s, {
        ...d,
        blocks: [...d.blocks.slice(0, i), ...a.blocks, ...d.blocks.slice(i)],
      });
    }
    case "remove": {
      const remainingBlocks = d.blocks.filter((b) => b.id !== a.id);
      return commit(s, {
        ...d,
        blocks: remainingBlocks.length !== 0 ? remainingBlocks : [newText()],
      });
    }
    case "move": {
      const to = Math.max(0, Math.min(a.to, d.blocks.length));
      if (a.from === to) return s;
      return commit(s, {
        ...d,
        blocks: d.blocks.toSpliced(a.from, 1).toSpliced(to, 0, d.blocks[a.from]),
      });
    }
    case "upload": {
      const patch = (doc: Doc) =>
        editBlockByID(doc, a.id, (b) => ({
          ...b,
          status: a.status,
          error: a.error,
        }));
      return { ...s, past: s.past.map(patch), present: patch(d), future: s.future.map(patch) };
    }
    case "undo": {
      const prevDoc = s.past.at(-1);
      return prevDoc !== undefined
        ? { past: s.past.slice(0, -1), present: prevDoc, future: [d, ...s.future], key: null }
        : s;
    }
    case "redo": {
      const nextDoc = s.future.at(0);
      return nextDoc !== undefined
        ? { past: [...s.past, d], present: nextDoc, future: s.future.slice(1), key: null }
        : s;
    }
  }
}

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

function toCreateRequest(doc: Doc, publish: boolean): CreateRequest {
  const blocks = doc.blocks.filter(
    (b) => b.blockType === BlockType.BLOCK_TYPE_MEDIA || b.text.trim(),
  );
  return {
    post: {
      id: 0,
      title: doc.title.trim(),
      slug: doc.slug,
      postStatus: publish ? PostStatus.POST_STATUS_PUBLISHED : PostStatus.POST_STATUS_DRAFT,
      createdAt: undefined,
      createdBy: 0,
      updatedAt: undefined,
      updatedBy: 0,
    },
    postBlocks: blocks.map((b, i) => ({
      id: 0,
      postId: 0,
      position: i,
      blockType: b.blockType,
      text: b.blockType === BlockType.BLOCK_TYPE_TEXT ? b.text.trim() : "",
      media: b.blockType === BlockType.BLOCK_TYPE_MEDIA ? b.id : "",
      metadata: "{}",
    })),
  };
}

function focusBlock(root: HTMLElement | null, id: string) {
  const elem = root?.querySelector<HTMLInputElement | HTMLTextAreaElement>(`[data-focus="${id}"]`);
  elem?.focus();
  elem?.setSelectionRange(elem.value.length, elem.value.length);
}

export function Editor() {
  const navigate = useNavigate();

  const [state, dispatch] = useReducer(reducer, undefined, () => ({
    past: [],
    future: [],
    key: null,
    present: {
      title: "",
      slug: "",
      blocks: [newText()],
    },
  }));
  const { present: doc, past, future } = state;
  const [error, setError] = useState<string | null>(null);
  const [busy, setBusy] = useState(false);
  const [activeId, setActiveId] = useState<string | null>(null);
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
    // const last = doc.blocks.at(-1);
    // if (last?.blockType === BlockType.BLOCK_TYPE_TEXT && last.text === "")
    //   return focusBlock(root.current, last.id);
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

    try {
      const res = await postService.Create(toCreateRequest(doc, publish));
      if (res.ok) navigate(`/posts/${doc.slug}`);
      else setError(`couldn't save the post (${res.message}).`);
    } catch {
      setError("couldn't reach the server. Try again later.");
    } finally {
      setBusy(false);
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
      className="mx-auto flex max-w-2xl flex-col items-center pb-12"
      onPaste={onPaste}
      onDragOver={(e) => e.dataTransfer.types.includes("Files") && e.preventDefault()}
      onDrop={(e) => {
        if (e.dataTransfer.files.length === 0) return;
        e.preventDefault();
        addFiles(Array.from(e.dataTransfer.files), doc.blocks.length);
      }}
    >
      <div className="bg-background sticky top-0 z-5 w-full border-b">
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
        <div className="flex w-full flex-col gap-1">
          <div className="flex w-full justify-between p-2">
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
              <Button
                variant={"outline"}
                onClick={() => submit(false)}
                disabled={busy}
              >
                Draft
              </Button>
              <Button
                onClick={() => submit(true)}
                disabled={busy}
              >
                Publish
              </Button>
            </div>
          </div>
          <p className="text-destructive/80 px-2 text-sm">{error}</p>
        </div>
      </div>
      <DndContext
        sensors={sensors}
        collisionDetection={(p) => {
          return closestCenter(p);
          // if (p.pointerCoordinates) {
          //   p.pointerCoordinates.x = p.droppableContainers[0].rect.current?.left ?? 0;
          //   const topY = p.droppableContainers.reduce(
          //     (y, c) => Math.min(y, c.rect.current?.top ?? 0),
          //     999999999,
          //   );
          //   const botY = p.droppableContainers.reduce(
          //     (y, c) => Math.max(y, c.rect.current?.bottom ?? 0),
          //     -999999999,
          //   );

          //   p.pointerCoordinates.y = Math.max(topY, Math.min(botY, p.pointerCoordinates.y));
          // }

          // return pointerWithin(p);
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
              addText={addText}
              insertAtIndex={0}
              insertAt={insertAt}
              fileInput={fileInput}
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
                  className="field-sizing-content w-full resize-none outline-none"
                  placeholder="Start typing... (Enter adds a block, Shift+Enter a line break"
                  value={activeBlock.text}
                ></textarea>
              ) : (
                <img
                  src={activeBlock.src}
                  className="max-h-96 max-w-full min-w-0 flex-1 rounded-md object-contain pr-8"
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

type SortableBlockProps = {
  index: number;
  block: Block;
  dispatch: ActionDispatch<[a: Action]>;
  onTextKeyDown: (e: KeyboardEvent<HTMLTextAreaElement>, i: number) => void;
  addText: (i: number) => void;
  insertAt: React.RefObject<number>;
  fileInput: React.RefObject<HTMLInputElement | null>;
  startUpload: (b: MediaBlock) => void;
};

function SortableBlock({
  index,
  block,
  dispatch,
  onTextKeyDown,
  addText,
  insertAt,
  fileInput,
  startUpload,
}: SortableBlockProps) {
  const {
    attributes,
    listeners,
    setNodeRef,
    setActivatorNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({
    id: block.id,
  });

  return (
    <li
      ref={setNodeRef}
      data-index={index}
      className={"group relative flex flex-col items-start gap-1 rounded-md"}
    >
      <div
        className={cn("relative flex w-full rounded-md p-1", isDragging && "opacity-20")}
        style={{
          transform: CSS.Translate.toString(transform),
          transition,
        }}
      >
        <div
          ref={setActivatorNodeRef}
          {...attributes}
          {...listeners}
          className={cn(
            "text-muted-foreground/50 cursor-grab self-start",
            isDragging && "outline-none",
          )}
        >
          <GripVertical className="size-5" />
        </div>
        {block.blockType === BlockType.BLOCK_TYPE_TEXT && (
          <textarea
            dir="auto"
            data-focus={block.id}
            className="placeholder:text-muted-foreground/60 field-sizing-content w-full flex-1 resize-none outline-none"
            placeholder="Start typing... (Enter adds a block, Shift+Enter a line break"
            value={block.text}
            onChange={(e) =>
              dispatch({
                type: "text",
                id: block.id,
                text: e.target.value,
              })
            }
            onKeyDown={(e) => onTextKeyDown(e, index)}
          ></textarea>
        )}
        {block.blockType === BlockType.BLOCK_TYPE_MEDIA && (
          <MediaView
            block={block}
            onRetry={() => startUpload(block)}
          />
        )}
        <Button
          className="hover:text-destructive text-destructive/80 self-start bg-transparent opacity-0 group-hover:opacity-50 hover:bg-transparent"
          variant={"destructive"}
          onClick={() => {
            dispatch({ type: "remove", id: block.id });
          }}
        >
          <Trash className="size-4" />
        </Button>
      </div>
      <AddBlockButtons
        addText={addText}
        insertAtIndex={index + 1}
        insertAt={insertAt}
        fileInput={fileInput}
      />
    </li>
  );
}

type MediaViewProps = {
  block: MediaBlock;
  onRetry: () => void;
};

function MediaView({ block, onRetry }: MediaViewProps) {
  return (
    <div className="relative w-full overflow-hidden">
      <img
        src={block.src}
        className={cn(
          "max-h-96 w-full rounded-md object-contain",
          block.status !== "ready" && "opacity-60",
        )}
      />
      {block.status === "uploading" && (
        <div className="absolute inset-0 grid place-items-center">
          <Loader2
            className="size-6 animate-spin"
            aria-label="Uploading"
          />
        </div>
      )}
      {block.status === "error" && (
        <div className="bg-background/90 text-destructive absolute inset-x-0 bottom-0 flex items-center justify-between gap-2 p-2 text-sm">
          <span className="min-w-0 truncate">Upload failed ({block.error})</span>
          <Button
            size="sm"
            variant="outline"
            onClick={onRetry}
          >
            <RotateCw /> Retry
          </Button>
        </div>
      )}
    </div>
  );
}

type AddBlockButtonsProps = {
  addText: (i: number) => void;
  insertAtIndex: number;
  insertAt: React.RefObject<number>;
  fileInput: React.RefObject<HTMLInputElement | null>;
};

function AddBlockButtons({ addText, insertAtIndex, insertAt, fileInput }: AddBlockButtonsProps) {
  return (
    <div className="flex w-full flex-row justify-center gap-4 opacity-0 focus-within:opacity-50 hover:opacity-50">
      <Button
        variant="ghost"
        size="sm"
        onClick={() => addText(insertAtIndex)}
      >
        <Type /> Text
      </Button>
      <Button
        variant="ghost"
        size="sm"
        onClick={() => {
          insertAt.current = insertAtIndex;
          fileInput.current?.click();
        }}
      >
        <ImagePlus /> Image
      </Button>
    </div>
  );
}
