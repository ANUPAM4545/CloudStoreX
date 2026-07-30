import { AuthDTO, UserViewModel } from "./types";

export function mapAuthDTOToUserViewModel(dto: AuthDTO, fallbackEmail = "user@example.com"): UserViewModel {
  const email = dto.user?.email || fallbackEmail;
  const id = dto.user?.id || "usr_" + Math.random().toString(36).substring(2, 9);
  const fullName = dto.user?.full_name || email.split("@")[0] || "User";

  return {
    id,
    email,
    fullName,
    avatarUrl: `https://api.dicebear.com/7.x/initials/svg?seed=${encodeURIComponent(fullName)}`,
  };
}
