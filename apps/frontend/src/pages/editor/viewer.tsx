import { BlockType, type GetResponse } from "@gen/api/postpb/v1/post";
import { cn } from "cn";
import { getSrcSet, type EnrichedPostBlock } from "./mapper";
import { imageSizeItems, textSizeItems, type ImageMetadata, type TextMetadata } from "./model";

export function Viewer({ doc }: { doc: GetResponse }) {
  return (
    <div className="mx-auto flex max-w-2xl flex-col items-center pb-12">
      <div className="bg-background sticky top-0 z-5 w-full border-b px-2 md:px-0">
        <h1
          dir="auto"
          className="placeholder:text-muted-foreground/60 mt-6 mb-3 resize-none text-3xl font-bold outline-none"
          aria-label="Post title"
        >
          {doc.post?.title}
        </h1>
      </div>
      <ul className="flex w-full flex-col gap-4 px-4 pt-6 md:px-0">
        {doc.postBlocks.map((block, index) => (
          <PostBlockComponent
            index={index}
            block={block}
            key={block.id}
          />
        ))}
      </ul>
    </div>
  );
}

function PostBlockComponent({ block }: { block: EnrichedPostBlock; index: number }) {
  if (block.blockType === BlockType.BLOCK_TYPE_TEXT) {
    const textMeta: TextMetadata = JSON.parse(block.metadata || "{}");

    return (
      <textarea
        dir="auto"
        readOnly
        className={cn(
          "placeholder:text-muted-foreground/60 field-sizing-content w-full flex-1 resize-none overflow-hidden outline-none",
          textSizeItems.find((i) => i.value === textMeta.size)?.className,
        )}
      >
        {block.text}
      </textarea>
    );
  }

  if (block.blockType === BlockType.BLOCK_TYPE_MEDIA) {
    const imageMeta: ImageMetadata = JSON.parse(block.metadata || "{}");

    return (
      <img
        srcSet={getSrcSet(block.mediaInfo)}
        className={cn(
          "rounded-md w-full object-contain",
          imageSizeItems.find((i) => i.value === imageMeta.size)?.className,
        )}
      />
    );
  }
}
