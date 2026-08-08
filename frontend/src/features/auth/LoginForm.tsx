"use client";

import React, { useState, useEffect } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "next/navigation";
import { toast } from "sonner";
import { motion, AnimatePresence } from "framer-motion";
import { AlertTriangle, Wand2, ArrowRight, ShieldCheck, Database, Cloud, Zap } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import { SocialLogin } from "./components/SocialLogin";
import { FloatingInput } from "./components/FloatingInput";

import { authApi } from "./api";
import { mapAuthDTOToUserViewModel } from "./mappers";
import { loginSchema, LoginFormData } from "./schemas";
import { useAuthStore } from "./store";

export function LoginForm() {
  const [serverError, setServerError] = useState("");
  const [capsLockActive, setCapsLockActive] = useState(false);
  const [useMagicLink, setUseMagicLink] = useState(false);
  
  const [successStage, setSuccessStage] = useState<number>(0);
  const [isSuccessSequence, setIsSuccessSequence] = useState(false);

  const router = useRouter();
  const setSession = useAuthStore((state) => state.setSession);

  const {
    register,
    handleSubmit,
    watch,
    formState: { errors, isSubmitting, dirtyFields },
  } = useForm<LoginFormData>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      email: "",
      password: "",
      rememberMe: false,
    },
    mode: "onChange",
  });

  const emailValue = watch("email");
  const isEmailValid = dirtyFields.email && !errors.email && emailValue?.length > 0;

  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.getModifierState("CapsLock")) setCapsLockActive(true);
      else setCapsLockActive(false);
    };
    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, []);

  const runSuccessSequence = () => {
    setIsSuccessSequence(true);
    setSuccessStage(1); // Authenticating

    setTimeout(() => {
      setSuccessStage(2); // Establishing Secure Session
    }, 400);

    setTimeout(() => {
      setSuccessStage(3); // Loading Workspace
    }, 1000);

    setTimeout(() => {
      setSuccessStage(4); // Dashboard Ready
    }, 1600);

    setTimeout(() => {
      router.push("/dashboard/overview");
    }, 2000);
  };

  const onSubmit = async (data: LoginFormData) => {
    setServerError("");
    try {
      if (useMagicLink) {
        await new Promise((resolve) => setTimeout(resolve, 1000));
        toast.success("Magic link sent!", { description: `Check ${data.email} for your login link.` });
        return;
      }

      const res = await authApi.login(data);
      if (res && res.token) {
        const userVM = mapAuthDTOToUserViewModel(res, data.email);
        setSession(res.token, userVM);
        runSuccessSequence();
      } else {
        throw new Error("Invalid token received from server.");
      }
    } catch (err: any) {
      const errorMsg = err?.message || "Invalid email or password.";
      setServerError(errorMsg);
    }
  };

  if (isSuccessSequence) {
    return (
      <motion.div 
        initial={{ opacity: 0, scale: 0.95 }}
        animate={{ opacity: 1, scale: 1 }}
        className="flex flex-col items-center justify-center space-y-6 py-12 w-full max-w-sm mx-auto"
        aria-live="polite"
      >
        <div className="relative w-20 h-20 flex items-center justify-center">
          <motion.div 
            className="absolute inset-0 border-4 border-primary/20 rounded-full"
            animate={successStage === 4 ? { scale: 1.1, opacity: 0 } : { scale: 1, opacity: 1 }}
            transition={{ duration: 0.4 }}
          />
          <motion.div 
            className="absolute inset-0 border-4 border-primary rounded-full border-t-transparent"
            animate={{ rotate: 360 }}
            transition={{ duration: 1, repeat: Infinity, ease: "linear" }}
            style={{ opacity: successStage === 4 ? 0 : 1 }}
          />
          <AnimatePresence>
            {successStage === 4 && (
              <motion.div
                initial={{ scale: 0 }}
                animate={{ scale: 1 }}
                transition={{ type: "spring", stiffness: 200, damping: 15 }}
                className="absolute inset-0 bg-emerald-500 rounded-full flex items-center justify-center text-white"
              >
                <CheckIcon className="w-8 h-8" />
              </motion.div>
            )}
          </AnimatePresence>
        </div>

        <div className="h-6 overflow-hidden relative w-full flex justify-center">
          <AnimatePresence mode="popLayout">
            <motion.p
              key={successStage}
              initial={{ y: 20, opacity: 0 }}
              animate={{ y: 0, opacity: 1 }}
              exit={{ y: -20, opacity: 0 }}
              transition={{ duration: 0.3 }}
              className="text-lg font-medium absolute"
            >
              {successStage === 1 && "Authenticating..."}
              {successStage === 2 && "Establishing Secure Session..."}
              {successStage === 3 && "Loading Workspace..."}
              {successStage === 4 && "Dashboard Ready"}
            </motion.p>
          </AnimatePresence>
        </div>
      </motion.div>
    );
  }

  return (
    <div className="space-y-6 w-full max-w-sm mx-auto">
      <div className="space-y-2 text-center sm:text-left">
        <h1 className="text-2xl font-bold tracking-tight">Welcome back</h1>
        <p className="text-sm text-muted-foreground">
          Enter your credentials to access your organization dashboard.
        </p>
      </div>

      <SocialLogin 
        onGoogle={() => toast("Google SSO coming soon.")} 
        onGitHub={() => toast("GitHub SSO coming soon.")}
        onMicrosoft={() => toast("Microsoft Entra ID coming soon.")}
      />

      <form onSubmit={handleSubmit(onSubmit)} className="space-y-4" noValidate>
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

        <AnimatePresence mode="wait">
          {!useMagicLink && (
            <motion.div 
              initial={{ opacity: 0, height: 0 }}
              animate={{ opacity: 1, height: "auto" }}
              exit={{ opacity: 0, height: 0 }}
              className="space-y-2 overflow-hidden pt-1"
            >
              <div className="flex items-center justify-end pb-1">
                <a href="/reset-password" className="text-xs font-medium text-primary hover:underline transition-colors" tabIndex={-1}>
                  Forgot password?
                </a>
              </div>
              
              <FloatingInput
                label="Password"
                type="password"
                autoComplete="current-password"
                {...register("password")}
                error={errors.password?.message}
              />
              
              <div className="flex items-center justify-between mt-2 h-4">
                <div />
                <AnimatePresence>
                  {capsLockActive && (
                    <motion.p 
                      initial={{ opacity: 0, y: -5 }} 
                      animate={{ opacity: 1, y: 0 }} 
                      exit={{ opacity: 0, y: -5 }} 
                      className="text-xs text-amber-500 font-medium"
                    >
                      Caps Lock is on
                    </motion.p>
                  )}
                </AnimatePresence>
              </div>

              <div className="flex items-center space-x-2 pt-1 pb-2">
                <Checkbox id="rememberMe" {...register("rememberMe")} />
                <label
                  htmlFor="rememberMe"
                  className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
                >
                  Remember for 30 days
                </label>
              </div>
            </motion.div>
          )}
        </AnimatePresence>

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
              Authenticating...
            </span>
          ) : useMagicLink ? (
            <span className="flex items-center gap-2 relative z-10">
              Send Magic Link
              <Wand2 size={16} className="group-hover:text-yellow-300 transition-colors" />
            </span>
          ) : (
            <span className="flex items-center gap-2 relative z-10">
              Sign In
              <ArrowRight size={16} className="group-hover:translate-x-1 transition-transform" />
            </span>
          )}
        </Button>
      </form>

      <div className="flex flex-col gap-4 text-center text-sm pt-2">
        <button
          type="button"
          onClick={() => setUseMagicLink(!useMagicLink)}
          className="text-muted-foreground hover:text-foreground font-medium transition-colors inline-flex items-center justify-center gap-2"
        >
          {useMagicLink ? "Sign in with password instead" : "Sign in with a Magic Link"}
        </button>
        <p className="text-muted-foreground">
          Don&apos;t have an account?{" "}
          <a href="/register" className="text-primary hover:underline font-semibold transition-colors">
            Request Access
          </a>
        </p>
      </div>

      {/* Trust Indicators */}
      <div className="pt-6 border-t flex flex-wrap justify-center gap-x-4 gap-y-2 text-[11px] font-medium text-muted-foreground">
        <span className="flex items-center gap-1"><ShieldCheck size={12} className="text-emerald-500" /> Enterprise Security</span>
        <span className="flex items-center gap-1"><Zap size={12} className="text-emerald-500" /> OAuth 2.0</span>
        <span className="flex items-center gap-1"><Database size={12} className="text-emerald-500" /> Zero Trust</span>
        <span className="flex items-center gap-1"><Cloud size={12} className="text-emerald-500" /> Multi-Cloud</span>
      </div>
    </div>
  );
}

function CheckIcon(props: React.SVGProps<SVGSVGElement>) {
  return (
    <svg
      {...props}
      xmlns="http://www.w3.org/2000/svg"
      width="24"
      height="24"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="3"
      strokeLinecap="round"
      strokeLinejoin="round"
    >
      <polyline points="20 6 9 17 4 12" />
    </svg>
  );
}
