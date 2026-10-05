import { api, type APIResponse, type WithAPIResponse } from "@/api/client";
import type {
  FinalizeUploadRequest,
  FinalizeUploadResponse,
  GenerateURLRequest,
  GenerateURLResponse,
  UploadService,
} from "@gen/upload/uploadpb/v1/upload";

async function GenerateURL(request: GenerateURLRequest): Promise<APIResponse<GenerateURLResponse>> {
  return await api("/upload/v1/init", {
    method: "POST",
    authenticated: true,
    body: JSON.stringify(request),
  });
}

async function FinalizeUpload(
  request: FinalizeUploadRequest,
): Promise<APIResponse<FinalizeUploadResponse>> {
  return await api("/upload/v1/finalize", {
    method: "POST",
    authenticated: true,
    body: JSON.stringify(request),
  });
}

export const uploadService = {
  GenerateURL,
  FinalizeUpload,
} satisfies WithAPIResponse<UploadService>;

export async function uploadMedia(
  mediaKey: string,
  file: File,
): Promise<APIResponse<FinalizeUploadResponse>> {
  const init = await uploadService.GenerateURL({ mediaKey });
  if (!init.ok) return init;

  const form = new FormData();
  for (const [name, value] of Object.entries(init.fields)) {
    form.append(name, value);
  }
  form.append("file", file); // must come after the policy fields

  const stored = await fetch(init.url, { method: "POST", body: form });
  if (!stored.ok) {
    return { ok: false, code: stored.status, message: "upload.storageRejected" };
  }

  return await uploadService.FinalizeUpload({ mediaKey });
}
