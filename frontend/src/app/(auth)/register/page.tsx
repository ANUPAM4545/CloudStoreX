"use client";

import React from "react";
import { RegisterForm } from "@/features/auth/RegisterForm";
import { ShieldCheck } from "lucide-react";

export default function RegisterPage() {
  return (
    <div className="w-full">
      <RegisterForm />
      
      {/* Security Footer */}
      <div className="flex items-center justify-center gap-2 text-xs text-muted-foreground mt-8">
        <ShieldCheck size={14} className="text-emerald-500" />
        <span>End-to-end encrypted storage</span>
      </div>
    </div>
  );
}
