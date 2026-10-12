import { api, type APIResponse, type WithAPIResponse } from "@/api/client";
import type {
  AuthnService,
  LoginRequest,
  LoginResponse,
  LogoutRequest,
  LogoutResponse,
  RefreshJWTRequest,
  RefreshJWTResponse,
  RegisterRequest,
  RegisterResponse,
} from "@gen/auth/authnpb/v1/authn";

async function Register(request: RegisterRequest): Promise<APIResponse<RegisterResponse>> {
  return await api("/auth/v1/authn/register", {
    method: "POST",
    credentials: "include",
    body: JSON.stringify(request),
  });
}

async function Login(request: LoginRequest): Promise<APIResponse<LoginResponse>> {
  return await api("/auth/v1/authn/login", {
    method: "POST",
    credentials: "include",
    body: JSON.stringify(request),
  });
}

async function RefreshJWT(request: RefreshJWTRequest): Promise<APIResponse<RefreshJWTResponse>> {
  return await api("/auth/v1/authn/refresh", {
    method: "POST",
    credentials: "include",
    body: JSON.stringify(request),
  });
}

async function Logout(request: LogoutRequest): Promise<APIResponse<LogoutResponse>> {
  return await api("/auth/v1/authn/logout", {
    method: "POST",
    credentials: "include",
    body: JSON.stringify(request),
  });
}

export const authnService = {
  Register,
  Login,
  RefreshJWT,
  Logout,
} satisfies WithAPIResponse<AuthnService>;
