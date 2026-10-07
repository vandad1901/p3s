import { Button } from "@/components/ui/button";
import { useSortable } from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { BlockType } from "@gen/api/postpb/v1/post";
import { cn } from "cn";
import { GripVertical, ImagePlus, Loader2, RotateCw, Trash, Type } from "lucide-react";
import { type ActionDispatch, type KeyboardEvent } from "react";
import { imageSizeItems, textSizeItems, type Action, type Block, type MediaBlock } from "./model";

type SortableBlockProps = {
  index: number;
  block: Block;
  dispatch: ActionDispatch<[a: Action]>;
  onTextKeyDown: (e: KeyboardEvent<HTMLTextAreaElement>, i: number) => void;
  addText: (i: number) => void;
  insertAt: React.RefObject<number>;
  fileInput: React.RefObject<HTMLInputElement | null>;
  startUpload: (b: MediaBlock) => void;
  setFocusedId: (id: string | null) => void;
};

export function SortableBlock({
  index,
  block,
  dispatch,
  onTextKeyDown,
  addText,
  insertAt,
  fileInput,
  startUpload,
  setFocusedId,
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
      data-id={block.id}
      className={"group relative flex flex-col items-start gap-1 rounded-md"}
    >
      <div
        className={cn("relative flex w-full rounded-md p-1", isDragging && "opacity-20")}
        style={{
          transform: CSS.Translate.toString(transform),
          transition,
        }}
        onClick={() => setFocusedId(block.id)}
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
            className={cn(
              "placeholder:text-muted-foreground/60 field-sizing-content w-full flex-1 resize-none overflow-hidden outline-none",
              textSizeItems.find((i) => i.value === block.metadata.size)?.className,
            )}
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
        onTextAdd={() => addText(index + 1)}
        onImageAdd={() => {
          insertAt.current = index + 1;
          fileInput.current?.click();
        }}
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
          "w-full rounded-md object-contain",
          imageSizeItems.find((i) => i.value === block.metadata.size)?.className,
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
  onTextAdd: () => void;
  onImageAdd: () => void;
};

export function AddBlockButtons({ onTextAdd, onImageAdd }: AddBlockButtonsProps) {
  return (
    <div className="flex w-full flex-row justify-center gap-4 opacity-0 focus-within:opacity-50 hover:opacity-50">
      <Button
        variant="ghost"
        size="sm"
        onClick={onTextAdd}
      >
        <Type /> Text
      </Button>
      <Button
        variant="ghost"
        size="sm"
        onClick={onImageAdd}
      >
        <ImagePlus /> Image
      </Button>
    </div>
  );
}
