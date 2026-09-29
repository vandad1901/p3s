const API_URL = import.meta.env.VITE_API_URL;

let accessToken: string | null = null;

export function setAccessToken(token: string | null) {
  accessToken = token;
}

type ErrorResponse = {
  code: number;
  message: string;
};

export type APIResponse<T> =
  | (T & { ok: true })
  | (ErrorResponse & { ok: false });

export type WithAPIResponse<T> = {
  [K in keyof T]: T[K] extends (...args: infer Args) => Promise<infer Response>
    ? (...args: Args) => Promise<APIResponse<Response>>
    : T[K];
};

export async function api<T>(
  path: string,
  options?: RequestInit & { authenticated?: boolean },
): Promise<APIResponse<T>> {
  const { authenticated = false, headers, ...fetchOptions } = options ?? {};

  const requestHeaders = new Headers(headers);
  requestHeaders.set("Content-Type", "application/json");

  if (authenticated && accessToken) {
    requestHeaders.set("Authorization", `Bearer ${accessToken}`);
  }

  const response = await fetch(`${API_URL}${path}`, {
    ...fetchOptions,
    headers: requestHeaders,
  });

  return {
    ...(await response.json()),
    ok: response.ok,
  };
}
