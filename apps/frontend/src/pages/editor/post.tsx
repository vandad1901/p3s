import { Loader2 } from "lucide-react";
import { useState, useEffect } from "react";
import { useParams, Navigate } from "react-router-dom";
import { loadFullPost, Editor } from "./editor";
import type { EnrichedGetResponse } from "./mapper";
import { Viewer } from "./viewer";

export function PostPage({ editMode }: { editMode: boolean }) {
  const { slug } = useParams();
  const [loaded, setLoaded] = useState<
    { slug: string; editMode: boolean; doc: EnrichedGetResponse } | undefined
  >();

  useEffect(() => {
    if (!slug) return;
    let ignore = false;
    loadFullPost(slug)
      .then((doc) => {
        if (!ignore) setLoaded({ slug, editMode, doc });
      })
      .catch((err) => console.error("Failed to load post:", err));
    return () => {
      ignore = true;
    };
  }, [slug, editMode]);

  const doc = loaded?.slug === slug && loaded?.editMode === editMode ? loaded.doc : undefined;

  if (!slug && !editMode)
    return (
      <Navigate
        to="/"
        replace
      />
    );

  if (slug && !doc) {
    return (
      <div className="flex h-screen w-full items-center justify-center">
        <Loader2 className="size-10 animate-spin" />
      </div>
    );
  }

  return editMode ? (
    <Editor
      key={slug ?? "new"}
      remoteResponse={doc}
    />
  ) : (
    <Viewer
      key={slug}
      doc={doc!}
    />
  );
}
