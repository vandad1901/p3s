import { api, type APIResponse, type WithAPIResponse } from "@/api/client";
import type {
  CreateRequest,
  CreateResponse,
  DeleteRequest,
  DeleteResponse,
  GetRequest,
  GetResponse,
  PostService,
  UpdateRequest,
  UpdateResponse,
} from "@gen/api/postpb/v1/post";

async function CreatePost(request: CreateRequest): Promise<APIResponse<CreateResponse>> {
  return await api("/api/v1/post", {
    method: "POST",
    authenticated: true,
    body: JSON.stringify(request),
  });
}

async function GetPost(request: GetRequest): Promise<APIResponse<GetResponse>> {
  return await api(`/api/v1/post/${request.slug}`, {
    method: "GET",
    authenticated: true,
  });
}

async function UpdatePost(request: UpdateRequest): Promise<APIResponse<UpdateResponse>> {
  return await api("/api/v1/post", {
    method: "PATCH",
    authenticated: true,
    body: JSON.stringify(request),
  });
}

async function DeletePost(request: DeleteRequest): Promise<APIResponse<DeleteResponse>> {
  return await api(
    `/api/v1/post/${request.id}?updatedAt=${encodeURIComponent(request.updatedAt!)}`,
    {
      method: "DELETE",
      authenticated: true,
    },
  );
}

export const postService = {
  Create: CreatePost,
  Get: GetPost,
  Update: UpdatePost,
  Delete: DeletePost,
} satisfies WithAPIResponse<PostService>;
