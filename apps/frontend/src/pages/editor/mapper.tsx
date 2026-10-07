import {
  BlockType,
  PostStatus,
  type CreateRequest,
  type GetResponse,
  type Post,
  type PostBlock,
  type UpdateRequest,
} from "@gen/api/postpb/v1/post";
import type { MediaInfo } from "@gen/media/mediapb/v1/media";
import type { Block, Doc } from "./model";

const PositionGap = 1000;
const INT32_MAX = 2 ** 31 - 1;

type EnrichedPostBlock = PostBlock & { mediaInfo?: MediaInfo };

export type EnrichedGetResponse = Omit<GetResponse, "postBlocks"> & {
  postBlocks: EnrichedPostBlock[];
};

export function mapRemoteToDoc(res: EnrichedGetResponse): Doc {
  const { post: remotePost, postBlocks: remotePostBlocks } = res;

  return {
    title: remotePost!.title,
    slug: remotePost!.slug,
    blocks: remotePostBlocks.map((b) => {
      if (b.blockType === BlockType.BLOCK_TYPE_TEXT) {
        return {
          remoteID: b.id,
          id: b.media || crypto.randomUUID(),
          position: b.position,
          blockType: b.blockType,
          text: b.text,
          metadata: JSON.parse(b.metadata || "{}"),
        };
      } else if (b.blockType === BlockType.BLOCK_TYPE_MEDIA) {
        const srcSet = b.mediaInfo?.derivatives.map((d) => `${d.url} ${d.width}w`).join(", ") || "";
        debugger;

        return {
          remoteID: b.id,
          id: b.media || crypto.randomUUID(),
          position: b.position,
          blockType: b.blockType,
          src: undefined,
          srcSet,
          status: "ready",
          metadata: JSON.parse(b.metadata || "{}"),
        };
      } else {
        throw new Error(`Unknown block type: ${b.blockType}`);
      }
    }),
  };
}

function mapToRemotePostBlock(b: Block, postId: string): PostBlock {
  return {
    id: b.remoteID,
    postId: postId,
    position: b.position,
    blockType: b.blockType,
    text: b.blockType === BlockType.BLOCK_TYPE_TEXT ? b.text.trim() : "",
    media: b.blockType === BlockType.BLOCK_TYPE_MEDIA ? b.id : "",
    metadata: JSON.stringify(b.metadata),
  };
}

export function toCreateRequest(doc: Doc, publish: boolean): CreateRequest {
  const blocks = doc.blocks.filter(hasContent);

  const positionRichBlocks: Block[] = blocks.map((b, i) => ({ ...b, position: i * PositionGap }));

  return {
    post: {
      id: "0",
      title: doc.title.trim(),
      slug: doc.slug,
      postStatus: publish ? PostStatus.POST_STATUS_PUBLISHED : PostStatus.POST_STATUS_DRAFT,
      createdAt: undefined,
      createdBy: "0",
      updatedAt: undefined,
      updatedBy: "0",
    },
    postBlocks: positionRichBlocks.map((b) => mapToRemotePostBlock(b, "0")),
  };
}

function canonical(value: unknown): string {
  return JSON.stringify(value, (_, v) =>
    v && typeof v === "object" && !Array.isArray(v)
      ? Object.fromEntries(Object.entries(v).sort(([a], [b]) => (a < b ? -1 : a > b ? 1 : 0)))
      : v,
  );
}

const sameJson = (a: string, b: string) =>
  canonical(JSON.parse(a || "{}")) === canonical(JSON.parse(b || "{}"));

export function toUpdateRequest(
  remotePost: Post,
  remotePostBlocks: PostBlock[],
  doc: Doc,
  publish: boolean,
): UpdateRequest {
  const blocks = withPositions(doc.blocks.filter(hasContent));

  const deletedBlocks = remotePostBlocks
    .filter((remoteBlock) => !blocks.some((b) => b.remoteID === remoteBlock.id))
    .map((b) => b.id);
  const updatedBlocks = blocks.filter((b) => {
    const remoteBlock = remotePostBlocks.find((remoteBlock) => remoteBlock.id === b.remoteID);

    if (!remoteBlock) return false;
    if (b.blockType !== remoteBlock.blockType) return true;
    if (b.position !== remoteBlock.position) return true;

    if (b.blockType === BlockType.BLOCK_TYPE_TEXT) {
      return (
        b.text.trim() !== remoteBlock.text.trim() ||
        !sameJson(JSON.stringify(b.metadata), remoteBlock.metadata)
      );
    }
    if (b.blockType === BlockType.BLOCK_TYPE_MEDIA) {
      return JSON.stringify(b.metadata) !== remoteBlock.metadata;
    }
  });
  const insertedBlocks = blocks.filter(
    (b) => !remotePostBlocks.some((remoteBlock) => remoteBlock.id === b.remoteID),
  );

  return {
    post: {
      id: remotePost.id,
      title: doc.title.trim(),
      slug: doc.slug,

      postStatus: publish ? PostStatus.POST_STATUS_PUBLISHED : PostStatus.POST_STATUS_DRAFT,
      createdAt: remotePost.createdAt,
      createdBy: remotePost.createdBy,
      updatedAt: remotePost.updatedAt,
      updatedBy: remotePost.updatedBy,
    },
    postBlockRequest: {
      inserted: insertedBlocks.map((b) => mapToRemotePostBlock(b, remotePost.id)),
      updated: updatedBlocks.map((b) => mapToRemotePostBlock(b, remotePost.id)),
      deleted: deletedBlocks,
    },
  };
}

type Anchor = {
  index: number;
  position: number;
};

function longestIncreasing(items: Anchor[]): Anchor[] {
  const longest = (chains: Anchor[][]) =>
    chains.reduce((best, chain) => (chain.length > best.length ? chain : best), []);

  const chainsEndingAt: Anchor[][] = [];
  items.forEach((item, i) => {
    const earlierAndSmaller = chainsEndingAt.filter((_, j) => items[j].position < item.position);
    chainsEndingAt[i] = [...longest(earlierAndSmaller), item];
  });

  return longest(chainsEndingAt);
}

function spread(lo: number, hi: number, count: number): number[] | null {
  if (count === 0) return [];
  const gap = hi - lo;
  if (gap <= count) return null;
  return Array.from({ length: count }, (_, k) => lo + Math.floor((gap * (k + 1)) / (count + 1)));
}

const isFilled = (segment: number[] | null): segment is number[] => segment !== null;

function assignPositions(old: (number | undefined)[]): number[] {
  const known = old.flatMap((position, index) =>
    position === undefined ? [] : [{ index, position }],
  );
  const kept = longestIncreasing(known);

  const starts: Anchor[] = [{ index: -1, position: 0 }, ...kept];
  const segments = starts.map((lo, i) => {
    const hi = kept.at(i);
    const count = (hi?.index ?? old.length) - lo.index - 1;
    const upper = hi?.position ?? lo.position + PositionGap * (count + 1);
    const fresh = spread(lo.position, upper, count);
    return fresh && [...fresh, ...(hi ? [hi.position] : [])];
  });

  const positions = segments.every(isFilled) ? segments.flat() : null;
  return positions?.every((p) => p <= INT32_MAX)
    ? positions
    : old.map((_, i) => (i + 1) * PositionGap);
}

const isSaved = (b: Block) => b.remoteID !== "0";
const hasContent = (b: Block) => b.blockType === BlockType.BLOCK_TYPE_MEDIA || b.text.trim() !== "";

function withPositions(blocks: Block[]): Block[] {
  const positions = assignPositions(blocks.map((b) => (isSaved(b) ? b.position : undefined)));
  return blocks.map((b, i) => ({ ...b, position: positions[i] }));
}
