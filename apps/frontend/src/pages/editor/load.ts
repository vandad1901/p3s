import { postService } from "@/api/api.service";
import type { APIResponse } from "@/api/client";
import { mediaService } from "@/api/media.service";
import { BlockType } from "@gen/api/postpb/v1/post";
import type { EnrichedGetResponse } from "./mapper";

export async function loadFullPost(slug: string): Promise<APIResponse<EnrichedGetResponse>> {
  const res = await postService.Get({ slug });
  if (!res.ok) {
    return { ok: false, message: res.message, code: res.code };
  }

  const mediaKeys = res.postBlocks
    .filter((b) => b.blockType === BlockType.BLOCK_TYPE_MEDIA)
    .map((b) => `${res.post?.createdBy}/${b.media}`);

  const mediaInfo = await mediaService.GetMedia({ mediaKeys: mediaKeys });
  if (!mediaInfo.ok) {
    return { ok: false, message: mediaInfo.message, code: mediaInfo.code };
  }

  const enrichedBlocks = res.postBlocks.map((b) => {
    if (b.blockType === BlockType.BLOCK_TYPE_MEDIA) {
      const info = mediaInfo.media.find((m) => m.mediaKey === `${res.post?.createdBy}/${b.media}`);
      return { ...b, mediaInfo: info };
    }
    return b;
  });

  return { ...res, postBlocks: enrichedBlocks };
}
