import { apiClient } from "@/lib/api/client";
import { APIResponse } from "@/lib/api/types";
import { AuthDTO, LoginRequestDTO, RegisterRequestDTO } from "./types";

export const authApi = {
  async login(payload: LoginRequestDTO): Promise<AuthDTO> {
    const res = await apiClient.post<any, APIResponse<AuthDTO>>("/auth/login", payload);
    if (!res.data && (res as any).token) {
      return res as unknown as AuthDTO;
    }
    return res.data || { token: "" };
  },

  async register(payload: RegisterRequestDTO): Promise<AuthDTO> {
    const res = await apiClient.post<any, APIResponse<AuthDTO>>("/auth/register", payload);
    if (!res.data && (res as any).token) {
      return res as unknown as AuthDTO;
    }
    return res.data || { token: "" };
  },

  async getMe(): Promise<{ email: string }> {
    const res = await apiClient.get<any, APIResponse<{ email: string }>>("/auth/me");
    return res.data || { email: "user@example.com" };
  },
};
