import { api, type APIResponse, type WithAPIResponse } from "@/api/client";
import type { CreateRequest, CreateResponse, PostService } from "@gen/api/postpb/v1/post";

async function CreatePost(request: CreateRequest): Promise<APIResponse<CreateResponse>> {
  return await api("/api/v1/post", {
    method: "POST",
    authenticated: true,
    body: JSON.stringify(request),
  });
}

export const postService = {
  Create: CreatePost,
} satisfies Pick<WithAPIResponse<PostService>, "Create">;
