"use client";

import React, { useState } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { Loader2 } from "lucide-react";

interface SocialLoginProps {
  onGoogle: () => void;
  onGitHub: () => void;
  onMicrosoft: () => void;
}

export function SocialLogin({ onGoogle, onGitHub, onMicrosoft }: SocialLoginProps) {
  const [loadingProvider, setLoadingProvider] = useState<string | null>(null);

  const handleClick = (provider: string, onClick: () => void) => {
    if (loadingProvider) return; // Prevent multiple clicks
    setLoadingProvider(provider);
    onClick();
    // Reset after a timeout for mock purposes (in reality, it redirects)
    setTimeout(() => {
      if (loadingProvider === provider) {
        setLoadingProvider(null);
      }
    }, 2000);
  };

  return (
    <div className="space-y-4">
      <div className="space-y-3">
        <SocialButton
          provider="google"
          label="Continue with Google"
          isLoading={loadingProvider === "google"}
          disabled={loadingProvider !== null && loadingProvider !== "google"}
          onClick={() => handleClick("google", onGoogle)}
          icon={
            <svg viewBox="0 0 24 24" width="20" height="20" xmlns="http://www.w3.org/2000/svg">
              <path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" fill="#4285F4" />
              <path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" fill="#34A853" />
              <path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" fill="#FBBC05" />
              <path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" fill="#EA4335" />
            </svg>
          }
        />
        
        <SocialButton
          provider="github"
          label="Continue with GitHub"
          isLoading={loadingProvider === "github"}
          disabled={loadingProvider !== null && loadingProvider !== "github"}
          onClick={() => handleClick("github", onGitHub)}
          icon={
            <svg viewBox="0 0 24 24" width="20" height="20" fill="currentColor" xmlns="http://www.w3.org/2000/svg">
              <path d="M12 1.27a11 11 0 00-3.48 21.46c.55.09.73-.28.73-.55v-1.84c-3.03.64-3.67-1.46-3.67-1.46-.55-1.29-1.28-1.65-1.28-1.65-.92-.65.1-.65.1-.65 1.1.09 1.73 1.1 1.73 1.1.92 1.65 2.57 1.2 3.21.92a2 2 0 01.64-1.47c-2.47-.27-5.04-1.19-5.04-5.5 0-1.19.46-2.2 1.2-2.94a4 4 0 01.09-2.93s.91-.28 3.02 1.1c.82-.2 1.74-.28 2.66-.28.92 0 1.83.09 2.66.28 2.1-1.38 3.02-1.1 3.02-1.1a4 4 0 01.1 2.93c.73.74 1.2 1.75 1.2 2.94 0 4.32-2.57 5.23-5.04 5.5.55.46 1.01 1.38 1.01 2.76v4.14c0 .28.18.65.73.55A11 11 0 0012 1.27" />
            </svg>
          }
        />

        <SocialButton
          provider="microsoft"
          label="Continue with Microsoft Entra ID"
          isLoading={loadingProvider === "microsoft"}
          disabled={loadingProvider !== null && loadingProvider !== "microsoft"}
          onClick={() => handleClick("microsoft", onMicrosoft)}
          icon={
            <svg viewBox="0 0 21 21" width="20" height="20" xmlns="http://www.w3.org/2000/svg">
              <path fill="#f25022" d="M1 1h9v9H1z"/>
              <path fill="#00a4ef" d="M1 11h9v9H1z"/>
              <path fill="#7fba00" d="M11 1h9v9h-9z"/>
              <path fill="#ffb900" d="M11 11h9v9h-9z"/>
            </svg>
          }
        />
      </div>

      <div className="relative">
        <div className="absolute inset-0 flex items-center">
          <span className="w-full border-t border-muted" />
        </div>
        <div className="relative flex justify-center text-[10px] font-medium uppercase tracking-widest">
          <span className="bg-background px-4 text-muted-foreground">Or continue with email</span>
        </div>
      </div>
    </div>
  );
}

interface SocialButtonProps {
  provider: string;
  label: string;
  icon: React.ReactNode;
  isLoading: boolean;
  disabled: boolean;
  onClick: () => void;
}

function SocialButton({ label, icon, isLoading, disabled, onClick }: SocialButtonProps) {
  return (
    <motion.button
      type="button"
      onClick={onClick}
      disabled={disabled || isLoading}
      whileHover={!disabled && !isLoading ? { y: -1, backgroundColor: "hsl(var(--muted) / 0.8)" } : {}}
      whileTap={!disabled && !isLoading ? { scale: 0.98 } : {}}
      className={`
        relative flex w-full items-center justify-center gap-3 rounded-xl border border-input bg-background px-4 h-11 text-sm font-medium transition-colors
        focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary/50 focus-visible:ring-offset-1 focus-visible:ring-offset-background
        ${disabled ? "opacity-50 cursor-not-allowed" : "hover:border-primary/30"}
      `}
      aria-label={label}
    >
      <div className="absolute left-4 flex h-full items-center justify-center">
        {icon}
      </div>
      
      <span className={isLoading ? "opacity-0" : "opacity-100"}>
        {label}
      </span>

      <AnimatePresence>
        {isLoading && (
          <motion.div
            initial={{ opacity: 0, scale: 0.8 }}
            animate={{ opacity: 1, scale: 1 }}
            exit={{ opacity: 0, scale: 0.8 }}
            className="absolute inset-0 flex items-center justify-center bg-background/50 rounded-xl"
          >
            <Loader2 className="h-5 w-5 animate-spin text-primary" />
          </motion.div>
        )}
      </AnimatePresence>
    </motion.button>
  );
}
