import { buttonVariants } from "@/components/ui/button";
import { useAuth } from "@/contexts/AuthContext";
import { FileQuestionMark, Home, Loader2 } from "lucide-react";
import { useEffect, useState, type ReactNode } from "react";
import { Link, Navigate, useLocation, useParams } from "react-router-dom";
import { Editor, loadFullPost } from "./editor";
import type { EnrichedGetResponse } from "./mapper";
import { Viewer } from "./viewer";

export function PostPage(props: { editMode: boolean }) {
  const { pathname } = useLocation();
  return (
    <PostLoader
      key={pathname}
      {...props}
    />
  );
}

function PostLoader({ editMode }: { editMode: boolean }) {
  const { slug } = useParams();
  const { user } = useAuth();
  const [result, setResult] = useState<{ doc: EnrichedGetResponse } | { error: string }>();

  useEffect(() => {
    if (!slug) return;
    loadFullPost(slug)
      .then((res) => setResult(res.ok ? { doc: res } : { error: res.message }))
      .catch(() => setResult({ error: "Couldn't reach the server." }));
  }, [slug]);

  if ((!slug && !editMode) || (editMode && !user))
    return (
      <Navigate
        to="/"
        replace
      />
    );

  if (result && "error" in result) {
    let msg = result.error;

    if (result.error === "post.NotFound") {
      msg = "The post you are looking for does not exist.";
      return (
        <div className="flex w-full flex-1 flex-col items-center justify-center gap-4 px-4 text-center">
          <p className="text-muted-foreground text-sm font-medium">
            <FileQuestionMark className="text-muted-foreground size-12" />
          </p>
          <h1 className="text-3xl font-bold tracking-tight">Post not found</h1>
          <p className="text-muted-foreground max-w-sm">
            This post doesn't exist, or it may have been moved or deleted.
          </p>
          <Link
            to="/"
            className={buttonVariants({ variant: "outline" })}
          >
            <Home />
            Back to home
          </Link>
        </div>
      );
    }
    return (
      <Centered>
        <p className="text-red-500">{msg}</p>
      </Centered>
    );
  }

  if (slug && !result)
    return (
      <Centered>
        <Loader2 className="size-10 animate-spin" />
      </Centered>
    );

  const doc = result?.doc;
  return editMode ? <Editor remoteResponse={doc} /> : <Viewer doc={doc!} />;
}

const Centered = ({ children }: { children: ReactNode }) => (
  <div className="flex w-full flex-1 items-center justify-center">{children}</div>
);
