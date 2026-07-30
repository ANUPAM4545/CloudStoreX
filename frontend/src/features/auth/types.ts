export interface LoginRequestDTO {
  email: string;
  password?: string;
}

export interface RegisterRequestDTO {
  email: string;
  password?: string;
  full_name?: string;
}

export interface AuthDTO {
  token: string;
  user?: {
    id: string;
    email: string;
    full_name?: string;
  };
}

export interface UserViewModel {
  id: string;
  email: string;
  fullName: string;
  avatarUrl: string;
}
