import { z } from "zod";

export const loginSchema = z.object({
  email: z.string().email("Please enter a valid work email address."),
  password: z.string().min(1, "Password is required."),
  rememberMe: z.boolean().optional(),
});

const passwordRules = z.string()
  .min(12, "Password must be at least 12 characters.")
  .regex(/[A-Z]/, "Must contain at least one uppercase letter.")
  .regex(/[a-z]/, "Must contain at least one lowercase letter.")
  .regex(/[0-9]/, "Must contain at least one number.")
  .regex(/[^A-Za-z0-9]/, "Must contain at least one special character.");

export const registerSchema = z.object({
  name: z.string().min(2, "Full name is required."),
  email: z.string().email("Please enter a valid work email address."),
  password: passwordRules,
});

export const resetPasswordSchema = z.object({
  email: z.string().email("Please enter a valid email address."),
});

export const updatePasswordSchema = z.object({
  password: passwordRules,
  confirmPassword: z.string().min(1, "Please confirm your password."),
}).refine((data) => data.password === data.confirmPassword, {
  message: "Passwords do not match.",
  path: ["confirmPassword"],
});

export type LoginFormData = z.infer<typeof loginSchema>;
export type RegisterFormData = z.infer<typeof registerSchema>;
export type ResetPasswordFormData = z.infer<typeof resetPasswordSchema>;
export type UpdatePasswordFormData = z.infer<typeof updatePasswordSchema>;
