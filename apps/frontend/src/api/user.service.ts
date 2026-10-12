import type {
  GetUserRequest,
  GetUserResponse,
  GetUsersRequest,
  GetUsersResponse,
  UpdateUserRequest,
  UpdateUserResponse,
  UserService,
} from "@gen/auth/userpb/v1/user";
import { api, type APIResponse, type WithAPIResponse } from "./client";

async function GetUser(request: GetUserRequest): Promise<APIResponse<GetUserResponse>> {
  return await api(`/auth/v1/user/${request.id}`, {
    method: "GET",
    authenticated: true,
  });
}

async function GetUsers(request: GetUsersRequest): Promise<APIResponse<GetUsersResponse>> {
  return await api("/auth/v1/users", {
    method: "POST",
    authenticated: true,
    body: JSON.stringify(request),
  });
}

async function UpdateUser(request: UpdateUserRequest): Promise<APIResponse<UpdateUserResponse>> {
  return await api("/auth/v1/user/update", {
    method: "POST",
    authenticated: true,
    body: JSON.stringify(request),
  });
}

export const userService = {
  GetUser,
  GetUsers,
  UpdateUser,
} satisfies WithAPIResponse<UserService>;
