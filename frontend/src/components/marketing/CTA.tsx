"use client";

import { motion } from "framer-motion";
import Link from "next/link";
import { ArrowRight, Terminal } from "lucide-react";

export function CTA() {
  return (
    <section className="py-32 relative bg-[#05070A] overflow-hidden">
      {/* Background Glow */}
      <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[800px] h-[400px] bg-primary/20 blur-[120px] rounded-[100%] pointer-events-none" />
      
      <div className="container mx-auto px-6 relative z-10 text-center">
        <motion.h2 
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          className="text-4xl md:text-6xl lg:text-7xl font-semibold tracking-tight text-white mb-8"
        >
          Build Enterprise <br />
          <span className="text-gradient">Storage Infrastructure.</span>
        </motion.h2>
        
        <motion.p
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ delay: 0.1 }}
          className="text-lg text-muted-foreground max-w-2xl mx-auto mb-10"
        >
          Deploy CloudStoreX in minutes and start routing your data intelligently across AWS and MinIO with zero-trust security.
        </motion.p>
        
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ delay: 0.2 }}
          className="flex flex-col sm:flex-row items-center justify-center gap-4 mb-16"
        >
          <Link
            href="/login"
            className="flex items-center justify-center gap-2 px-8 py-4 text-sm font-medium text-white transition-all rounded-full bg-primary hover:bg-primary/90 hover:scale-105 active:scale-95 w-full sm:w-auto shadow-[0_0_40px_-10px_rgba(59,130,246,0.5)]"
          >
            Start Building <ArrowRight className="w-4 h-4" />
          </Link>
          <Link
            href="#"
            className="flex items-center justify-center gap-2 px-8 py-4 text-sm font-medium text-white transition-all rounded-full bg-white/5 border border-white/10 hover:bg-white/10 w-full sm:w-auto"
          >
            Read Documentation
          </Link>
        </motion.div>
        
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ delay: 0.3 }}
          className="inline-flex items-center gap-4 bg-white/5 border border-white/10 rounded-full px-6 py-3 font-mono text-sm text-muted-foreground"
        >
          <Terminal className="w-4 h-4" />
          <span>npm create cloudstorex@latest</span>
        </motion.div>
      </div>
    </section>
  );
}
