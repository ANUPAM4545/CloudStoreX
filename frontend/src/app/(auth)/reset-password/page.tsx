"use client";

import React, { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { toast } from "sonner";
import { motion, AnimatePresence } from "framer-motion";
import { Key, ArrowRight, ArrowLeft, Mail, AlertTriangle } from "lucide-react";
import { FloatingInput } from "@/features/auth/components/FloatingInput";
import { Button } from "@/components/ui/button";

import { authApi } from "@/features/auth/api";
import { resetPasswordSchema, ResetPasswordFormData } from "@/features/auth/schemas";

export default function ResetPasswordPage() {
  const [success, setSuccess] = useState(false);
  const [serverError, setServerError] = useState("");

  const {
    register,
    handleSubmit,
    watch,
    formState: { errors, isSubmitting, dirtyFields },
  } = useForm<ResetPasswordFormData>({
    resolver: zodResolver(resetPasswordSchema),
    defaultValues: {
      email: "",
    },
    mode: "onChange",
  });

  const emailValue = watch("email");
  const isEmailValid = dirtyFields.email && !errors.email && emailValue?.length > 0;

  const onSubmit = async (data: ResetPasswordFormData) => {
    setServerError("");
    try {
      await authApi.resetPassword(data.email);
      setSuccess(true);
    } catch (err: any) {
      setServerError(err?.message || "Failed to send reset email. Please try again.");
    }
  };

  return (
    <div className="w-full">
      <div className="max-w-sm mx-auto space-y-6">
        <a href="/login" className="inline-flex items-center gap-2 text-sm font-medium text-muted-foreground hover:text-foreground transition-colors group">
          <ArrowLeft size={16} className="group-hover:-translate-x-1 transition-transform" />
          Back to login
        </a>

        <div className="space-y-2">
          <div className="w-12 h-12 bg-primary/10 rounded-xl flex items-center justify-center mb-6">
            <Key className="text-primary" size={24} />
          </div>
          <h1 className="text-2xl font-bold tracking-tight">Reset password</h1>
          <p className="text-sm text-muted-foreground">
            Enter your work email address and we'll send you a link to reset your password.
          </p>
        </div>

        <AnimatePresence mode="wait">
          {success ? (
            <motion.div
              key="success"
              initial={{ opacity: 0, scale: 0.95 }}
              animate={{ opacity: 1, scale: 1 }}
              className="bg-emerald-500/10 border border-emerald-500/20 rounded-xl p-6 text-center space-y-4"
            >
              <div className="w-12 h-12 bg-emerald-500/20 rounded-full flex items-center justify-center mx-auto">
                <Mail className="text-emerald-500" size={24} />
              </div>
              <div>
                <h3 className="font-semibold text-emerald-600 mb-1">Check your inbox</h3>
                <p className="text-sm text-emerald-600/80">
                  We've sent a password reset link to <span className="font-medium">{emailValue}</span>
                </p>
              </div>
              <Button 
                variant="outline" 
                className="w-full mt-4 bg-transparent border-emerald-500/30 text-emerald-600 hover:bg-emerald-500/10"
                onClick={() => window.location.href = "/login"}
              >
                Return to Login
              </Button>
            </motion.div>
          ) : (
            <motion.form
              key="form"
              initial={{ opacity: 0, y: 10 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, scale: 0.95 }}
              onSubmit={handleSubmit(onSubmit)}
              className="space-y-4"
              noValidate
            >
              <AnimatePresence>
                {serverError && (
                  <motion.div 
                    initial={{ opacity: 0, height: 0, marginBottom: 0 }}
                    animate={{ opacity: 1, height: "auto", marginBottom: 16 }}
                    exit={{ opacity: 0, height: 0, marginBottom: 0 }}
                    className="overflow-hidden"
                  >
                    <div className="p-3 text-sm text-rose-600 bg-rose-500/10 border border-rose-500/20 rounded-lg flex items-start gap-2">
                      <AlertTriangle size={16} className="mt-0.5 shrink-0" />
                      <span>{serverError}</span>
                    </div>
                  </motion.div>
                )}
              </AnimatePresence>

              <FloatingInput
                label="Work Email"
                type="email"
                autoComplete="email"
                {...register("email")}
                error={errors.email?.message}
                success={isEmailValid}
              />

              <Button 
                type="submit" 
                className="w-full h-12 font-medium text-base group relative overflow-hidden transition-all duration-300 hover:shadow-[0_0_20px_rgba(var(--primary),0.3)] hover:-translate-y-0.5" 
                disabled={isSubmitting}
              >
                <div className="absolute inset-0 bg-white/20 translate-y-full group-hover:translate-y-0 transition-transform duration-300 ease-out" />
                {isSubmitting ? (
                  <span className="flex items-center gap-2 relative z-10">
                    <svg className="animate-spin -ml-1 mr-3 h-5 w-5 text-current" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                      <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                      <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                    Sending link...
                  </span>
                ) : (
                  <span className="flex items-center gap-2 relative z-10">
                    Send Reset Link
                    <ArrowRight size={16} className="group-hover:translate-x-1 transition-transform" />
                  </span>
                )}
              </Button>
            </motion.form>
          )}
        </AnimatePresence>
      </div>
    </div>
  );
}
