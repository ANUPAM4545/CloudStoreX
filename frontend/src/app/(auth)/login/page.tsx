"use client";

import React from "react";
import { LoginForm } from "@/features/auth/LoginForm";
import { ShieldCheck } from "lucide-react";

export default function LoginPage() {
  return (
    <div className="w-full">
      <LoginForm />
      
      {/* Security Footer */}
      <div className="flex items-center justify-center gap-2 text-xs text-muted-foreground mt-8">
        <ShieldCheck size={14} className="text-emerald-500" />
        <span>Enterprise single sign-on enabled</span>
      </div>
    </div>
  );
}
