import React from "react";
import { RightPanel } from "@/features/auth/components/RightPanel";

export default function AuthLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="flex min-h-screen w-full bg-background">
      {/* Left Pane - Forms */}
      <div className="flex-1 flex flex-col justify-center items-center p-6 sm:p-12 lg:p-24 overflow-y-auto">
        <div className="w-full max-w-sm mx-auto">
          {children}
        </div>
      </div>
      
      {/* Right Pane - Visuals */}
      <RightPanel />
    </div>
  );
}
