"use client";

import React, { useState } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { LoginForm } from "@/features/auth/LoginForm";
import { RegisterForm } from "@/features/auth/RegisterForm";
import { Cloud, ShieldCheck } from "lucide-react";

export default function LoginPage() {
  const [isRegistering, setIsRegistering] = useState(false);

  return (
    <div className="flex min-h-screen items-center justify-center bg-muted/30 p-4">
      <div className="w-full max-w-md space-y-6">
        {/* Brand Header */}
        <div className="flex flex-col items-center text-center space-y-2">
          <div className="flex items-center justify-center h-12 w-12 rounded-xl bg-primary text-primary-foreground shadow-md">
            <Cloud size={24} />
          </div>
          <h1 className="text-2xl font-bold tracking-tight">CloudStoreX</h1>
          <p className="text-sm text-muted-foreground">
            Enterprise Cloud Storage Platform & Provider Router
          </p>
        </div>

        {/* Auth Card */}
        <Card className="shadow-lg border bg-card">
          <CardHeader className="space-y-1 pb-4">
            <CardTitle className="text-xl font-bold tracking-tight">
              {isRegistering ? "Create an account" : "Sign in to CloudStoreX"}
            </CardTitle>
            <CardDescription>
              {isRegistering
                ? "Enter your email below to get started with CloudStoreX"
                : "Enter your credentials to access your storage dashboard"}
            </CardDescription>
          </CardHeader>
          <CardContent>
            {isRegistering ? (
              <RegisterForm onSwitchToLogin={() => setIsRegistering(false)} />
            ) : (
              <LoginForm onSwitchToRegister={() => setIsRegistering(true)} />
            )}
          </CardContent>
        </Card>

        {/* Security / Enterprise Footer */}
        <div className="flex items-center justify-center gap-2 text-xs text-muted-foreground">
          <ShieldCheck size={14} className="text-emerald-500" />
          <span>Zero-trust provider isolation • End-to-end audit ready</span>
        </div>
      </div>
    </div>
  );
}
