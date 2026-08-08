"use client";

import React, { useState, useEffect } from "react";
import { motion } from "framer-motion";
import { CheckCircle2, Loader2, ShieldCheck } from "lucide-react";
import { useRouter } from "next/navigation";

export default function VerifyEmailPage() {
  const [status, setStatus] = useState<"verifying" | "success" | "error">("verifying");
  const router = useRouter();

  useEffect(() => {
    // Mock verification process
    const timer = setTimeout(() => {
      setStatus("success");
      
      // Auto-redirect after success
      setTimeout(() => {
        router.push("/dashboard/overview");
      }, 2000);
      
    }, 2500);
    
    return () => clearTimeout(timer);
  }, [router]);

  return (
    <div className="w-full max-w-sm mx-auto space-y-8 flex flex-col items-center text-center">
      {status === "verifying" && (
        <motion.div
          initial={{ opacity: 0, scale: 0.9 }}
          animate={{ opacity: 1, scale: 1 }}
          className="space-y-6 flex flex-col items-center"
        >
          <div className="w-16 h-16 rounded-full bg-muted flex items-center justify-center animate-pulse">
            <Loader2 size={32} className="text-foreground animate-spin" />
          </div>
          <div className="space-y-2">
            <h1 className="text-2xl font-bold tracking-tight">Verifying Email...</h1>
            <p className="text-sm text-muted-foreground">
              Please wait while we verify your email address securely.
            </p>
          </div>
        </motion.div>
      )}

      {status === "success" && (
        <motion.div
          initial={{ opacity: 0, scale: 0.9 }}
          animate={{ opacity: 1, scale: 1 }}
          className="space-y-6 flex flex-col items-center"
        >
          <div className="w-16 h-16 rounded-full bg-emerald-500/20 flex items-center justify-center">
            <motion.div
              initial={{ scale: 0 }}
              animate={{ scale: 1 }}
              transition={{ type: "spring", stiffness: 200, damping: 15 }}
              className="w-10 h-10 rounded-full bg-emerald-500 flex items-center justify-center text-white"
            >
              <CheckCircle2 size={24} />
            </motion.div>
          </div>
          <div className="space-y-2">
            <h1 className="text-2xl font-bold tracking-tight">Email Verified</h1>
            <p className="text-sm text-muted-foreground">
              Your email has been successfully verified. Redirecting to dashboard...
            </p>
          </div>
        </motion.div>
      )}

      {/* Security Footer */}
      <div className="flex items-center justify-center gap-2 text-xs text-muted-foreground mt-12">
        <ShieldCheck size={14} className="text-emerald-500" />
        <span>Enterprise identity protected</span>
      </div>
    </div>
  );
}
