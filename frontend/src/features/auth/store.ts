import { create } from 'zustand';
import { UserViewModel } from './types';

interface AuthState {
  token: string | null;
  user: UserViewModel | null;
  isAuthenticated: boolean;
  setToken: (token: string) => void;
  setSession: (token: string, user: UserViewModel) => void;
  logout: () => void;
}

const getDefaultUser = (email = "user@example.com"): UserViewModel => ({
  id: "usr_1",
  email,
  fullName: email.split("@")[0] || "User",
  avatarUrl: `https://api.dicebear.com/7.x/initials/svg?seed=${encodeURIComponent(email)}`,
});

export const useAuthStore = create<AuthState>((set) => ({
  token: typeof window !== 'undefined' ? localStorage.getItem('token') : null,
  user: typeof window !== 'undefined' && localStorage.getItem('token') ? getDefaultUser() : null,
  isAuthenticated: typeof window !== 'undefined' ? !!localStorage.getItem('token') : false,
  setToken: (token: string) => {
    if (typeof window !== 'undefined') {
      localStorage.setItem('token', token);
    }
    set({ token, isAuthenticated: true, user: getDefaultUser() });
  },
  setSession: (token: string, user: UserViewModel) => {
    if (typeof window !== 'undefined') {
      localStorage.setItem('token', token);
    }
    set({ token, user, isAuthenticated: true });
  },
  logout: () => {
    if (typeof window !== 'undefined') {
      localStorage.removeItem('token');
    }
    set({ token: null, user: null, isAuthenticated: false });
  },
}));
