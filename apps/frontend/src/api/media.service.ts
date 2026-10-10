import { api, type APIResponse, type WithAPIResponse } from "./client";
import type {
  MediaIngestedRequest,
  MediaIngestedResponse,
  GetMediaRequest,
  GetMediaResponse,
  MediaService,
} from "@gen/media/mediapb/v1/media";

async function MediaIngested(
  request: MediaIngestedRequest,
): Promise<APIResponse<MediaIngestedResponse>> {
  return await api("/media/v1/ingested", {
    method: "POST",
    authenticated: true,
    body: JSON.stringify(request),
  });
}

async function GetMedia(request: GetMediaRequest): Promise<APIResponse<GetMediaResponse>> {
  return await api("/media/v1/getMedia", {
    method: "POST",
    body: JSON.stringify(request),
  });
}

export const mediaService = {
  MediaIngested,
  GetMedia,
} satisfies WithAPIResponse<MediaService>;
