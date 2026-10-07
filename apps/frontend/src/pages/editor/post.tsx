import { Loader2 } from "lucide-react";
import { useState, useEffect } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { loadFullPost, Editor } from "./editor";
import type { EnrichedGetResponse } from "./mapper";
import { Viewer } from "./viewer";

export function PostPage({ editMode }: { editMode: boolean }) {
  const { slug } = useParams();
  const navigate = useNavigate();
  const [initialDoc, setInitialDoc] = useState<EnrichedGetResponse | undefined>(undefined);

  if (!slug && !editMode) {
    navigate("/");
    return null;
  }

  useEffect(() => {
    setInitialDoc(undefined);

    if (!slug) {
      return;
    }

    loadFullPost(slug)
      .then((res) => {
        setInitialDoc(res);
      })
      .catch((err) => {
        console.error("Failed to load post:", err);
      });
  }, [editMode, slug]);

  if (slug && !initialDoc) {
    return (
      <div className="flex h-screen w-full items-center justify-center">
        <Loader2 className="size-10 animate-spin" />
      </div>
    );
  }

  if (editMode) {
    return (
      <Editor
        key={slug}
        remoteResponse={initialDoc}
      />
    );
  } else {
    return (
      <Viewer
        key={slug}
        doc={initialDoc!}
      ></Viewer>
    );
  }
}
