import { Loader2 } from "lucide-react";
import { useState, useEffect, type ReactNode } from "react";
import { useParams, Navigate, useLocation } from "react-router-dom";
import { loadFullPost, Editor } from "./editor";
import type { EnrichedGetResponse } from "./mapper";
import { Viewer } from "./viewer";
import { useAuth } from "@/contexts/AuthContext";

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

  if (result && "error" in result)
    return (
      <Centered>
        <p className="text-red-500">{result.error}</p>
      </Centered>
    );

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
  <div className="flex h-screen w-full items-center justify-center">{children}</div>
);
