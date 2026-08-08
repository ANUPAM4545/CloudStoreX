"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import { useAuthStore } from "@/store/auth";
import { Cloud, ArrowRight } from "lucide-react";
import { motion } from "framer-motion";

export function Header() {
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated);
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
  }, []);

  return (
    <motion.header
      initial={{ y: -20, opacity: 0 }}
      animate={{ y: 0, opacity: 1 }}
      transition={{ duration: 0.8, ease: [0.16, 1, 0.3, 1] }}
      className="fixed top-0 left-0 right-0 z-50 flex items-center justify-between px-6 py-4 backdrop-blur-md border-b border-white/5 bg-background/50"
    >
      <div className="flex items-center gap-2">
        <div className="flex items-center justify-center w-8 h-8 rounded-lg bg-primary/20 text-primary">
          <Cloud className="w-5 h-5" />
        </div>
        <span className="text-xl font-medium tracking-tight text-white">
          CloudStoreX
        </span>
      </div>

      <nav className="hidden md:flex items-center gap-8 text-sm font-medium">
        {['Architecture', 'Features', 'Developers', 'Enterprise'].map((item) => (
          <Link key={item} href={`#${item.toLowerCase()}`} className="relative group text-muted-foreground hover:text-white transition-colors">
            {item}
            <span className="absolute -bottom-1 left-0 w-0 h-0.5 bg-primary transition-all duration-300 group-hover:w-full" />
          </Link>
        ))}
      </nav>

      <div className="flex items-center gap-4">
        {mounted && isAuthenticated ? (
          <Link
            href="/dashboard/overview"
            className="flex items-center gap-2 px-4 py-2 text-sm font-medium text-white transition-colors rounded-full bg-white/10 hover:bg-white/15"
          >
            Dashboard
          </Link>
        ) : (
          <div className={`flex items-center gap-4 transition-opacity duration-300 ${mounted ? "opacity-100" : "opacity-0"}`}>
            <Link
              href="/login"
              className="text-sm font-medium text-muted-foreground hover:text-white transition-colors"
            >
              Sign In
            </Link>
            <Link
              href="/login"
              className="group flex items-center gap-2 px-4 py-2 text-sm font-medium text-white transition-all rounded-full bg-primary hover:bg-primary/90 hover:shadow-[0_0_20px_-5px_rgba(59,130,246,0.5)] hover:scale-105 active:scale-95"
            >
              Start Free <ArrowRight className="w-4 h-4 group-hover:translate-x-1 transition-transform" />
            </Link>
          </div>
        )}
      </div>
    </motion.header>
  );
}
