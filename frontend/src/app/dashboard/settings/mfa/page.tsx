"use client";

import React, { useState } from "react";
import { ShieldCheck, Smartphone, Fingerprint, Lock, ChevronRight } from "lucide-react";
import { Button } from "@/components/ui/button";

export default function MFASettingsPage() {
  const [totpEnabled, setTotpEnabled] = useState(true);

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-2xl font-bold tracking-tight">Two-Factor Authentication</h1>
        <p className="text-sm text-muted-foreground mt-1">
          Add an extra layer of security to your account to ensure only you can access your data.
        </p>
      </div>

      <div className="grid gap-6">
        {/* Authenticator App */}
        <div className="border rounded-xl p-6 flex flex-col sm:flex-row gap-6 items-start sm:items-center justify-between bg-card">
          <div className="flex gap-4">
            <div className="p-3 bg-muted rounded-xl shrink-0 h-fit">
              <Smartphone size={24} className="text-foreground" />
            </div>
            <div>
              <div className="flex items-center gap-2 mb-1">
                <h3 className="font-semibold text-lg">Authenticator App</h3>
                {totpEnabled && (
                  <span className="px-2 py-0.5 text-[10px] uppercase tracking-wider font-semibold bg-emerald-500/15 text-emerald-500 rounded-full">
                    Enabled
                  </span>
                )}
              </div>
              <p className="text-sm text-muted-foreground max-w-md">
                Use an app like 1Password or Google Authenticator to generate one-time codes when you sign in.
              </p>
            </div>
          </div>
          <Button variant={totpEnabled ? "outline" : "default"} onClick={() => setTotpEnabled(!totpEnabled)}>
            {totpEnabled ? "Disable" : "Set Up"}
          </Button>
        </div>

        {/* Passkeys */}
        <div className="border rounded-xl p-6 flex flex-col sm:flex-row gap-6 items-start sm:items-center justify-between bg-card">
          <div className="flex gap-4">
            <div className="p-3 bg-muted rounded-xl shrink-0 h-fit">
              <Fingerprint size={24} className="text-foreground" />
            </div>
            <div>
              <div className="flex items-center gap-2 mb-1">
                <h3 className="font-semibold text-lg">Passkeys & Security Keys</h3>
                <span className="px-2 py-0.5 text-[10px] uppercase tracking-wider font-semibold bg-blue-500/15 text-blue-500 rounded-full">
                  Recommended
                </span>
              </div>
              <p className="text-sm text-muted-foreground max-w-md">
                Sign in safely and easily with Touch ID, Face ID, or a hardware security key.
              </p>
            </div>
          </div>
          <Button variant="outline">
            Manage <ChevronRight size={16} className="ml-1" />
          </Button>
        </div>

        {/* Recovery Codes */}
        <div className="border rounded-xl p-6 flex flex-col sm:flex-row gap-6 items-start sm:items-center justify-between bg-card">
          <div className="flex gap-4">
            <div className="p-3 bg-muted rounded-xl shrink-0 h-fit">
              <Lock size={24} className="text-foreground" />
            </div>
            <div>
              <h3 className="font-semibold text-lg mb-1">Recovery Codes</h3>
              <p className="text-sm text-muted-foreground max-w-md">
                Generate codes you can use to access your account if you lose your device.
              </p>
            </div>
          </div>
          <Button variant="outline">
            View Codes
          </Button>
        </div>
      </div>
    </div>
  );
}
