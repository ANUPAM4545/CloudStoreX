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
import { registerSchema, RegisterFormData } from "./schemas";
import { useAuthStore } from "./store";

export interface RegisterFormProps {
  onSwitchToLogin: () => void;
}

export function RegisterForm({ onSwitchToLogin }: RegisterFormProps) {
  const [serverError, setServerError] = useState("");
  const router = useRouter();
  const setSession = useAuthStore((state) => state.setSession);

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<RegisterFormData>({
    resolver: zodResolver(registerSchema),
    defaultValues: {
      email: "",
      password: "",
      confirmPassword: "",
    },
  });

  const onSubmit = async (data: RegisterFormData) => {
    setServerError("");
    try {
      await authApi.register({ email: data.email, password: data.password });
      const res = await authApi.login({ email: data.email, password: data.password });
      if (res && res.token) {
        const userVM = mapAuthDTOToUserViewModel(res, data.email);
        setSession(res.token, userVM);
        toast.success("Account created!", { description: "Welcome to CloudStoreX." });
        router.push("/dashboard/overview");
      } else {
        throw new Error("Registration successful, but failed to automatically log in.");
      }
    } catch (err: any) {
      const errorMsg = err?.message || "Failed to create account. Email may already be registered.";
      setServerError(errorMsg);
      toast.error("Registration failed", { description: errorMsg });
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
          htmlFor="reg-email"
          className="text-xs font-medium uppercase tracking-wider text-muted-foreground"
        >
          Email
        </label>
        <Input
          id="reg-email"
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
          htmlFor="reg-password"
          className="text-xs font-medium uppercase tracking-wider text-muted-foreground"
        >
          Password
        </label>
        <Input
          id="reg-password"
          type="password"
          placeholder="••••••••"
          autoComplete="new-password"
          {...register("password")}
          className={errors.password ? "border-destructive focus-visible:ring-destructive" : ""}
        />
        {errors.password && (
          <p className="text-xs text-destructive mt-1">{errors.password.message}</p>
        )}
      </div>

      <div className="space-y-1.5">
        <label
          htmlFor="reg-confirmPassword"
          className="text-xs font-medium uppercase tracking-wider text-muted-foreground"
        >
          Confirm Password
        </label>
        <Input
          id="reg-confirmPassword"
          type="password"
          placeholder="••••••••"
          autoComplete="new-password"
          {...register("confirmPassword")}
          className={
            errors.confirmPassword
              ? "border-destructive focus-visible:ring-destructive"
              : ""
          }
        />
        {errors.confirmPassword && (
          <p className="text-xs text-destructive mt-1">
            {errors.confirmPassword.message}
          </p>
        )}
      </div>

      <Button type="submit" className="w-full h-10 font-medium" disabled={isSubmitting}>
        {isSubmitting ? "Creating account..." : "Create Account"}
      </Button>

      <div className="text-sm text-center text-muted-foreground pt-2">
        Already have an account?{" "}
        <button
          type="button"
          onClick={onSwitchToLogin}
          className="text-primary hover:underline font-semibold"
        >
          Sign in
        </button>
      </div>
    </form>
  );
}
