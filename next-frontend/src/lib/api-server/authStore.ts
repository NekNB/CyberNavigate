import { create } from "zustand";

interface AuthStore {
  isAuth: boolean;
  setIsAuth: (value: boolean) => void;
}

export const useAuthStore = create<AuthStore>((set) => ({
  isAuth: false,
  setIsAuth: (value) => {
    set({ isAuth: value });
  },
}));
