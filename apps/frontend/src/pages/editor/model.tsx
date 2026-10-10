import { BlockType } from "@gen/api/postpb/v1/post";

export const MAX_IMAGE_BYTES = 20 << 20;

// types

export const textSizeItems = [
  {
    label: "Heading 1",
    value: "h1",
    className: "text-3xl font-bold",
  },
  {
    label: "Heading 2",
    value: "h2",
    className: "text-2xl font-semibold",
  },
  {
    label: "Heading 3",
    value: "h3",
    className: "text-xl font-semibold",
  },
  {
    label: "Paragraph",
    value: "p",
    className: "text-base",
  },
] as const;
export type TextSizeType = (typeof textSizeItems)[number]["value"];

export type TextMetadata = {
  size: TextSizeType;
};

export type TextBlock = {
  remoteID: string;
  id: string;
  position: number;
  blockType: BlockType.BLOCK_TYPE_TEXT;
  text: string;
  metadata: TextMetadata;
};

export const imageSizeItems = [
  {
    label: "Small",
    value: "sm",
    className: "max-h-48",
  },
  {
    label: "Medium",
    value: "md",
    className: "max-h-96",
  },
  {
    label: "Large",
    value: "lg",
    className: "max-h-128",
  },
] as const;
export type ImageSizeType = (typeof imageSizeItems)[number]["value"];

export type ImageMetadata = {
  size: ImageSizeType;
};

export type MediaBlock = {
  remoteID: string;
  id: string;
  position: number;
  blockType: BlockType.BLOCK_TYPE_MEDIA;
  file?: File;
  src?: string;
  srcSet: string;
  status: "uploading" | "ready" | "error";
  error?: string;
  metadata: ImageMetadata;
};

export type Block = TextBlock | MediaBlock;

export type Doc = {
  title: string;
  slug: string;
  blocks: Block[];
};

// History Management

export type State = {
  past: Doc[];
  present: Doc;
  future: Doc[];
  key: string | null;
};

export type Action =
  | { type: "title"; title: string }
  | { type: "slug"; slug: string }
  | { type: "text"; id: string; text: string }
  | { type: "insert"; index: number; blocks: Block[] }
  | { type: "updateTextMeta"; metadata: TextMetadata; id: string }
  | { type: "updateImageMeta"; metadata: ImageMetadata; id: string }
  | { type: "remove"; id: string }
  | { type: "move"; from: number; to: number }
  | { type: "upload"; id: string; status: MediaBlock["status"]; error?: string } // documenting upload status changes
  | { type: "undo" }
  | { type: "redo" };

export function commit(s: State, doc: Doc, key: string | null = null): State {
  if (key !== null && key === s.key) return { ...s, present: doc, future: [] };
  return { past: [...s.past, s.present].slice(-100), present: doc, future: [], key };
}

export const newText = (): TextBlock => ({
  remoteID: "0",
  id: crypto.randomUUID(),
  position: 0,
  blockType: BlockType.BLOCK_TYPE_TEXT,
  text: "",
  metadata: { size: "p" },
});
export const newMedia = (file: File, src: string): MediaBlock => ({
  remoteID: "0",
  id: crypto.randomUUID(),
  position: 0,
  blockType: BlockType.BLOCK_TYPE_MEDIA,
  file,
  src,
  srcSet: "",
  status: "uploading",
  metadata: { size: "lg" },
});

export function reducer(s: State, a: Action): State {
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
    case "updateTextMeta": {
      return commit(
        s,
        editBlockByID(d, a.id, (b) =>
          b.blockType === BlockType.BLOCK_TYPE_TEXT ? { ...b, metadata: a.metadata } : b,
        ),
        `metadata:${a.id}`,
      );
    }
    case "updateImageMeta": {
      return commit(
        s,
        editBlockByID(d, a.id, (b) =>
          b.blockType === BlockType.BLOCK_TYPE_MEDIA ? { ...b, metadata: a.metadata } : b,
        ),
        `metadata:${a.id}`,
      );
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
