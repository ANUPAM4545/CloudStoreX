"use client";

import React, { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "next/navigation";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { authApi } from "./api";
import { mapAuthDTOToUserViewModel } from "./mappers";
import { loginSchema, LoginFormData } from "./schemas";
import { useAuthStore } from "./store";

export interface LoginFormProps {
  onSwitchToRegister: () => void;
}

export function LoginForm({ onSwitchToRegister }: LoginFormProps) {
  const [serverError, setServerError] = useState("");
  const router = useRouter();
  const setSession = useAuthStore((state) => state.setSession);

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<LoginFormData>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      email: "",
      password: "",
    },
  });

  const onSubmit = async (data: LoginFormData) => {
    setServerError("");
    try {
      const res = await authApi.login(data);
      if (res && res.token) {
        const userVM = mapAuthDTOToUserViewModel(res, data.email);
        setSession(res.token, userVM);
        toast.success("Welcome back to CloudStoreX!");
        router.push("/dashboard/overview");
      } else {
        throw new Error("Invalid token received from server.");
      }
    } catch (err: any) {
      const errorMsg = err?.message || "Invalid credentials. Please try again.";
      setServerError(errorMsg);
      toast.error("Authentication failed", { description: errorMsg });
    }
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4" noValidate>
      {serverError && (
        <div className="p-3 text-sm text-rose-700 dark:text-rose-400 bg-rose-500/15 border border-rose-500/30 rounded-md animate-fade-in">
          {serverError}
        </div>
      )}

      <div className="space-y-1.5">
        <label
          htmlFor="email"
          className="text-xs font-medium uppercase tracking-wider text-muted-foreground"
        >
          Email
        </label>
        <Input
          id="email"
          type="email"
          placeholder="user@example.com"
          autoComplete="email"
          {...register("email")}
          className={errors.email ? "border-destructive focus-visible:ring-destructive" : ""}
        />
        {errors.email && (
          <p className="text-xs text-destructive mt-1">{errors.email.message}</p>
        )}
      </div>

      <div className="space-y-1.5">
        <label
          htmlFor="password"
          className="text-xs font-medium uppercase tracking-wider text-muted-foreground"
        >
          Password
        </label>
        <Input
          id="password"
          type="password"
          placeholder="••••••••"
          autoComplete="current-password"
          {...register("password")}
          className={errors.password ? "border-destructive focus-visible:ring-destructive" : ""}
        />
        {errors.password && (
          <p className="text-xs text-destructive mt-1">{errors.password.message}</p>
        )}
      </div>

      <Button type="submit" className="w-full h-10 font-medium" disabled={isSubmitting}>
        {isSubmitting ? "Signing in..." : "Sign in"}
      </Button>

      <div className="text-sm text-center text-muted-foreground pt-2">
        Don&apos;t have an account?{" "}
        <button
          type="button"
          onClick={onSwitchToRegister}
          className="text-primary hover:underline font-semibold"
        >
          Create account
        </button>
      </div>
    </form>
  );
}
