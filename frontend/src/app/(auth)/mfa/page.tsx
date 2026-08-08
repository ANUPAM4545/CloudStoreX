"use client";

import React, { useState } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { ShieldAlert, Fingerprint, Lock, ArrowRight, Smartphone } from "lucide-react";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { FloatingInput } from "@/features/auth/components/FloatingInput";
import { useAuthStore } from "@/features/auth/store";
import { authApi } from "@/features/auth/api";

export default function MFAPage() {
  const [code, setCode] = useState("");
  const [error, setError] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const router = useRouter();
  
  const token = useAuthStore((state) => state.token);
  const user = useAuthStore((state) => state.user);

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (code.length !== 6) {
      setError("Please enter a valid 6-digit code.");
      return;
    }

    setIsSubmitting(true);
    setError("");

    try {
      if (!token) throw new Error("No active session found. Please log in again.");
      
      await authApi.verifyMFA(token, code);
      // Wait for success animation
      await new Promise(res => setTimeout(res, 1000));
      router.push("/dashboard/overview");
    } catch (err: any) {
      setError(err?.message || "Invalid authentication code.");
      setCode("");
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="w-full">
      <div className="max-w-sm mx-auto space-y-6">
        <div className="w-12 h-12 bg-amber-500/10 rounded-xl flex items-center justify-center mb-6">
          <ShieldAlert className="text-amber-500" size={24} />
        </div>
        
        <div className="space-y-2">
          <h1 className="text-2xl font-bold tracking-tight">Two-Factor Authentication</h1>
          <p className="text-sm text-muted-foreground leading-relaxed">
            Your organization requires MFA. Please enter the 6-digit code from your authenticator app.
          </p>
        </div>

        <form onSubmit={onSubmit} className="space-y-6 pt-4" noValidate>
          <div className="space-y-4">
            <FloatingInput
              label="Authentication Code"
              type="text"
              autoComplete="one-time-code"
              inputMode="numeric"
              maxLength={6}
              value={code}
              onChange={(e) => {
                setError("");
                setCode(e.target.value.replace(/[^0-9]/g, ""));
              }}
              error={error}
              className="text-center tracking-[0.5em] font-mono text-xl h-16"
              placeholder="000000"
            />
            
            {/* Alternative Methods */}
            <div className="flex items-center justify-between gap-4">
              <button type="button" className="flex items-center gap-2 text-xs font-medium text-muted-foreground hover:text-foreground transition-colors p-2 rounded-md hover:bg-muted">
                <Fingerprint size={14} />
                Passkey
              </button>
              <button type="button" className="flex items-center gap-2 text-xs font-medium text-muted-foreground hover:text-foreground transition-colors p-2 rounded-md hover:bg-muted">
                <Smartphone size={14} />
                SMS / Push
              </button>
            </div>
          </div>

          <Button 
            type="submit" 
            className="w-full h-12 font-medium text-base group relative overflow-hidden transition-all duration-300 hover:shadow-[0_0_20px_rgba(var(--primary),0.3)] hover:-translate-y-0.5" 
            disabled={isSubmitting || code.length !== 6}
          >
            <div className="absolute inset-0 bg-white/20 translate-y-full group-hover:translate-y-0 transition-transform duration-300 ease-out" />
            {isSubmitting ? (
              <span className="flex items-center gap-2 relative z-10">
                <svg className="animate-spin -ml-1 mr-3 h-5 w-5 text-current" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                Verifying...
              </span>
            ) : (
              <span className="flex items-center gap-2 relative z-10">
                Verify and Continue
                <ArrowRight size={16} className="group-hover:translate-x-1 transition-transform" />
              </span>
            )}
          </Button>
        </form>
      </div>
    </div>
  );
}
