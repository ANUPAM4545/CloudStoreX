"use client";

import React, { useState } from "react";
import { Shield, CheckCircle2, ChevronRight, Settings } from "lucide-react";
import { Button } from "@/components/ui/button";

export default function SSOPage() {
  const [ssoEnabled, setSsoEnabled] = useState(true);

  return (
    <div className="space-y-8">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Enterprise SSO</h1>
          <p className="text-sm text-muted-foreground mt-1">
            Manage single sign-on for your organization via SAML, OIDC, or Workspace providers.
          </p>
        </div>
        <Button variant={ssoEnabled ? "outline" : "default"} onClick={() => setSsoEnabled(!ssoEnabled)}>
          {ssoEnabled ? "Disable SSO" : "Enable SSO"}
        </Button>
      </div>

      <div className="border rounded-xl bg-card overflow-hidden">
        <div className="p-6 border-b flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
          <div className="flex items-center gap-4">
            <div className="p-3 bg-muted rounded-xl">
              <Shield size={24} className={ssoEnabled ? "text-emerald-500" : "text-muted-foreground"} />
            </div>
            <div>
              <h3 className="font-semibold text-lg flex items-center gap-2">
                SAML Configuration
                {ssoEnabled && <CheckCircle2 size={16} className="text-emerald-500" />}
              </h3>
              <p className="text-sm text-muted-foreground">Authenticate users with your Identity Provider.</p>
            </div>
          </div>
          {ssoEnabled && (
            <Button variant="outline">
              <Settings size={16} className="mr-2" />
              Configure
            </Button>
          )}
        </div>

        {ssoEnabled ? (
          <div className="p-6 bg-muted/20">
            <div className="space-y-4 max-w-xl">
              <div className="grid grid-cols-3 gap-4 border-b pb-4">
                <span className="text-sm font-medium text-muted-foreground col-span-1">Identity Provider</span>
                <span className="text-sm font-semibold col-span-2">Okta Workforce Identity</span>
              </div>
              <div className="grid grid-cols-3 gap-4 border-b pb-4">
                <span className="text-sm font-medium text-muted-foreground col-span-1">Entity ID</span>
                <span className="text-sm font-mono col-span-2 break-all">urn:cloudstorex:enterprise:12345</span>
              </div>
              <div className="grid grid-cols-3 gap-4 border-b pb-4">
                <span className="text-sm font-medium text-muted-foreground col-span-1">ACS URL</span>
                <span className="text-sm font-mono col-span-2 break-all">https://api.cloudstorex.com/sso/saml/acs</span>
              </div>
              <div className="grid grid-cols-3 gap-4">
                <span className="text-sm font-medium text-muted-foreground col-span-1">Status</span>
                <span className="text-sm col-span-2 text-emerald-500 font-semibold">Active & Enforced</span>
              </div>
            </div>
          </div>
        ) : (
          <div className="p-12 text-center text-muted-foreground">
            <p>Single Sign-On is currently disabled.</p>
            <p className="text-sm mt-1">Enable it to configure SAML or OIDC.</p>
          </div>
        )}
      </div>

      {/* Other providers */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
        {["Google Workspace", "Microsoft Entra ID", "Auth0"].map((provider) => (
          <div key={provider} className="border rounded-xl p-4 flex items-center justify-between hover:bg-muted/50 cursor-pointer transition-colors">
            <span className="font-medium">{provider}</span>
            <ChevronRight size={16} className="text-muted-foreground" />
          </div>
        ))}
      </div>
    </div>
  );
}
