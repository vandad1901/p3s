import { postService } from "@/api/api.service";
import {
  AlertDialog,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import { Button, buttonVariants } from "@/components/ui/button";
import { useAuth } from "@/contexts/AuthContext";
import { BlockType, type GetResponse } from "@gen/api/postpb/v1/post";
import { cn } from "cn";
import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { getSrcSet, type EnrichedPostBlock } from "./mapper";
import { imageSizeItems, textSizeItems, type ImageMetadata, type TextMetadata } from "./model";

export function Viewer({ doc }: { doc: GetResponse }) {
  const { user } = useAuth();
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();

  return (
    <div className="mx-auto flex w-full max-w-2xl flex-col items-center pb-12">
      <div className="bg-background sticky top-0 z-5 flex w-full flex-col gap-2 border-b px-2 pt-6 pb-3 md:px-0">
        <div className="flex items-end justify-between">
          <h1
            dir="auto"
            className="placeholder:text-muted-foreground/60 resize-none text-3xl font-bold outline-none"
            aria-label="Post title"
          >
            {doc.post?.title}
          </h1>
          {user && user.userId === doc.post?.createdBy && (
            <div className="flex items-center gap-2">
              <Link
                to={`/editor/${doc.post?.slug}`}
                className={cn(buttonVariants({ variant: "outline" }))}
              >
                Edit
              </Link>
              <AlertDialog>
                <AlertDialogTrigger render={<Button variant="destructive" />}>
                  Delete
                </AlertDialogTrigger>
                <AlertDialogContent>
                  <AlertDialogHeader>
                    <AlertDialogTitle>Are you absolutely sure?</AlertDialogTitle>
                    <AlertDialogDescription>
                      This action cannot be undone. This will permanently delete your post.
                    </AlertDialogDescription>
                  </AlertDialogHeader>
                  <AlertDialogFooter>
                    <AlertDialogCancel>Cancel</AlertDialogCancel>
                    <AlertDialogCancel
                      variant="destructive"
                      onClick={async () => {
                        try {
                          const res = await postService.Delete({
                            id: doc.post!.id!,
                            updatedAt: doc.post?.updatedAt,
                          });
                          if (!res.ok) {
                            switch (res.message) {
                              case "post.NotFound":
                                setError(
                                  "The post has been updated or deleted by someone else. Please refresh the page and try again.",
                                );

                                return;
                            }
                            setError("Failed to delete the post");
                            return;
                          }
                        } catch (e) {
                          setError("Couldn't reach the server.");
                          return;
                        }

                        return navigate("/");
                      }}
                    >
                      Delete
                    </AlertDialogCancel>
                  </AlertDialogFooter>
                </AlertDialogContent>
              </AlertDialog>
            </div>
          )}
        </div>
        {error && (
          <p
            className="text-destructive/80 text-sm"
            aria-live="polite"
          >
            {error}
          </p>
        )}
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
          "w-full rounded-md object-contain",
          imageSizeItems.find((i) => i.value === imageMeta.size)?.className,
        )}
      />
    );
  }
}
