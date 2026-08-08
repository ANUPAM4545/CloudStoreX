"use client";

import React from "react";
import { motion } from "framer-motion";
import { Hammer } from "lucide-react";
import { PageHeader } from "./PageHeader";

interface PlaceholderPageProps {
  title: string;
  description: string;
}

export function PlaceholderPage({ title, description }: PlaceholderPageProps) {
  return (
    <div className="flex flex-col h-full space-y-6">
      <PageHeader title={title} description={description} />
      <motion.div 
        initial={{ opacity: 0, y: 10 }}
        animate={{ opacity: 1, y: 0 }}
        className="flex-1 rounded-2xl border border-dashed border-border/60 bg-muted/10 flex flex-col items-center justify-center text-center p-8 min-h-[400px]"
      >
        <div className="w-16 h-16 bg-muted/30 rounded-2xl flex items-center justify-center mb-6 ring-1 ring-border">
          <Hammer className="text-muted-foreground w-8 h-8" />
        </div>
        <h3 className="text-xl font-semibold tracking-tight mb-2">Coming Soon</h3>
        <p className="text-muted-foreground max-w-md mx-auto">
          The {title} module is currently under active development. This feature will be available in an upcoming release.
        </p>
      </motion.div>
    </div>
  );
}
