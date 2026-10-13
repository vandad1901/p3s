import { postService } from "@/api/post.service";
import {
  AlertDialog,
  AlertDialogAction,
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
import { useTheme } from "@/contexts/ThemeContext";
import { BlockType, type GetResponse } from "@gen/api/postpb/v1/post";
import { cn } from "cn";
import "katex/dist/katex.min.css";
import {
  isValidElement,
  lazy,
  Suspense,
  useState,
  type ComponentProps,
  type ReactNode,
} from "react";
import Markdown from "react-markdown";
import { Link, useNavigate } from "react-router-dom";
import rehypeKatex from "rehype-katex";
import rehypeSanitize, { defaultSchema } from "rehype-sanitize";
import remarkMath from "remark-math";
import { getSrcSet, type EnrichedPostBlock } from "./mapper";
import { imageSizeItems, textSizeItems, type ImageMetadata, type TextMetadata } from "./model";

const CodeBlock = lazy(() => import("./code-block"));

export function Viewer({ doc }: { doc: GetResponse }) {
  const { user } = useAuth();
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();

  const [deleteDialogueOpen, setDeleteDialogueOpen] = useState(false);

  return (
    <div className="mx-auto flex w-full max-w-2xl flex-col items-center pb-12">
      <div className="bg-background sticky top-0 z-5 flex w-full flex-col gap-2 border-b px-2 pt-6 pb-3 md:px-0">
        <div className="flex items-end justify-between">
          <h1
            dir="auto"
            className="placeholder:text-muted-foreground/60 w-full resize-none text-3xl font-bold outline-none"
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
              <AlertDialog
                open={deleteDialogueOpen}
                onOpenChange={setDeleteDialogueOpen}
              >
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
                    <AlertDialogAction
                      variant="destructive"
                      onClick={async () => {
                        setDeleteDialogueOpen(false);
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
                    </AlertDialogAction>
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
          <li key={block.id}>
            <PostBlockComponent
              index={index}
              block={block}
            />
          </li>
        ))}
      </ul>
    </div>
  );
}

const schema = {
  ...defaultSchema,
  attributes: {
    ...defaultSchema.attributes,
    code: [...(defaultSchema.attributes?.code || []), ["className", "math-inline", "math-display"]],
  },
};

function PostBlockComponent({ block }: { block: EnrichedPostBlock; index: number }) {
  const { theme } = useTheme();
  if (block.blockType === BlockType.BLOCK_TYPE_TEXT) {
    const textMeta: TextMetadata = JSON.parse(block.metadata || "{}");
    const textClassName = textSizeItems.find((i) => i.value === textMeta.size)!.className;

    const processedText = block.text.replace(/^(#{1,6})(?=\s)/gm, "\\$1");

    return (
      <div>
        <Markdown
          remarkPlugins={[remarkMath]}
          rehypePlugins={[[rehypeSanitize, schema], rehypeKatex]}
          disallowedElements={textMeta.size && ["h1", "h2", "h3", "h4", "h5", "h6"]}
          unwrapDisallowed
          urlTransform={(url) => {
            if (/^(https?:\/\/|mailto:)/i.test(url)) return url;
            if (url.startsWith("#")) return url;
            if (url.startsWith("/")) return url;
            return "";
          }}

          components={{
            p: ({ node, ...props }) => (
              <p
                dir="auto"
                className={cn(
                  "w-full flex-1 resize-none overflow-hidden text-start whitespace-pre-wrap outline-none",
                  textClassName,
                )}
                {...props}
              />
            ),
            a: ({ node, ...props }) => {
              const href = props.href || "";
              if (href === "") return "INVALID URL";
              if ((href.startsWith("/") && !href.startsWith("//")) || href.startsWith("#")) {
                return (
                  <Link
                    {...props}
                    to={href}
                    className="text-primary visited:text-muted-foreground inline cursor-pointer underline"
                  />
                );
              }
              return <ExternalURLWarning url={href}>{props.children}</ExternalURLWarning>;
            },
            code: ({ node, ...props }) => {
              return (
                <InlineCode
                  code={String(props.children!)}
                  textClass={textClassName}
                />
              );
            },
            pre: ({ node, ...props }) => {
              const child = props.children;
              if (!isValidElement<ComponentProps<"code">>(child)) return null;
              const language = /language-([\w+-]+)/.exec(child.props.className ?? "")?.[1];
              const code = String(child.props.children! ?? "").replace(/\n$/, "");

              return (
                <div>
                  <div
                    className="text-muted-foreground mt-3 -mb-2 w-fit rounded-t-[0.3rem] px-4 pt-1 text-sm"
                    style={{ backgroundColor: theme === "dark" ? "#282c34" : "#fafafa" }}
                  >
                    {language}
                  </div>
                  <Suspense
                    fallback={
                      <pre
                        className="my-2 overflow-x-auto rounded-[0.3em] rounded-tl-none p-4 ps-10"
                        style={{ backgroundColor: theme === "dark" ? "#282c34" : "#fafafa" }}
                      >
                        <code>{code}</code>
                      </pre>
                    }
                  >
                    <CodeBlock
                      code={code}
                      language={language ?? ""}
                    />
                  </Suspense>
                </div>
              );
            },
            ul: ({ node, ...props }) => (
              <ul
                dir="auto"
                className={cn(
                  "flex list-inside list-disc flex-col py-2 whitespace-pre-wrap",
                  textClassName,
                )}
                {...props}
              />
            ),
            ol: ({ node, ...props }) => (
              <ol
                dir="auto"
                className={cn(
                  "flex list-inside list-decimal flex-col py-2 whitespace-pre-wrap",
                  textClassName,
                )}
                {...props}
              />
            ),
            li: ({ node, ...props }) => (
              <li
                dir="auto"
                className={cn("resize-none text-start whitespace-pre-wrap", textClassName)}
                {...props}
              />
            ),
          }}
        >
          {processedText}
        </Markdown>
      </div>
    );
  }

  if (block.blockType === BlockType.BLOCK_TYPE_MEDIA) {
    const imageMeta: ImageMetadata = JSON.parse(block.metadata || "{}");
    const imageClassName = imageSizeItems.find((i) => i.value === imageMeta.size)?.className;

    return (
      <img
        srcSet={getSrcSet(block.mediaInfo)}
        className={cn("w-full rounded-md object-contain", imageClassName)}
      />
    );
  }
}

function ExternalURLWarning({ url, children }: { url: string; children: ReactNode }) {
  const [open, setOpen] = useState(false);
  const [openInNewTab, setOpenInNewTab] = useState(false);

  return (
    <AlertDialog
      open={open}
      onOpenChange={setOpen}
    >
      <AlertDialogTrigger
        render={
          <a
            href={url}
            className="text-primary visited:text-muted-foreground inline cursor-pointer underline"
            onClick={(e) => {
              e.preventDefault();
              setOpen(true);
              if (e.ctrlKey || e.metaKey) {
                setOpenInNewTab(true);
              } else {
                setOpenInNewTab(false);
              }
            }}
            onAuxClick={(e) => {
              if (e.button === 1) {
                e.preventDefault();
                setOpen(true);
                setOpenInNewTab(true);
              }
            }}
          />
        }
      >
        {children}
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>External Link</AlertDialogTitle>
          <AlertDialogDescription>
            You are about to open an external link. Do you want to proceed?
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel>Cancel</AlertDialogCancel>
          <AlertDialogAction
            variant="default"
            onClick={(e) => {
              const newTab = openInNewTab || e.ctrlKey || e.metaKey;
              setOpen(false);
              window.open(url, newTab ? "_blank" : "_self", "noopener noreferrer");
            }}
          >
            Open Link
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}

function InlineCode({ code, textClass }: { code: string; textClass: string }) {
  return (
    <code className={cn("bg-muted text-muted-foreground rounded-md px-1", textClass)}>{code}</code>
  );
}
