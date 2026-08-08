"use client";

import { motion } from "framer-motion";
import { GitBranch, Star, GitFork, GitPullRequest, CircleDot } from "lucide-react";
import Link from "next/link";

export function OpenSource() {
  return (
    <section className="py-24 relative bg-background overflow-hidden">
      <div className="container mx-auto px-6 max-w-5xl">
        <div className="flex flex-col lg:flex-row items-center gap-16">
          <div className="flex-1 w-full">
            <motion.h2 
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              className="text-3xl md:text-5xl font-semibold tracking-tight text-white mb-6"
            >
              Open by <span className="text-gradient">design.</span>
            </motion.h2>
            <motion.p 
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ delay: 0.1 }}
              className="text-lg text-muted-foreground mb-8"
            >
              CloudStoreX is built in the open. We believe enterprise infrastructure should be transparent, extensible, and community-driven. Check out our architecture diagrams, RFCs, and source code on GitHub.
            </motion.p>

            <motion.div
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ delay: 0.2 }}
              className="flex flex-wrap gap-4"
            >
              <Link
                href="#"
                className="inline-flex items-center gap-2 px-6 py-3 text-sm font-medium text-white transition-all rounded-full bg-white/5 border border-white/10 hover:bg-white/10"
              >
                <GitBranch className="w-5 h-5" /> View on GitHub
              </Link>
              <Link
                href="#"
                className="inline-flex items-center gap-2 px-6 py-3 text-sm font-medium text-muted-foreground transition-all rounded-full hover:text-white"
              >
                Architecture Docs
              </Link>
              <Link
                href="#"
                className="inline-flex items-center gap-2 px-6 py-3 text-sm font-medium text-muted-foreground transition-all rounded-full hover:text-white"
              >
                Public Roadmap
              </Link>
            </motion.div>
          </div>

          <div className="flex-1 w-full">
            <motion.div
              initial={{ opacity: 0, scale: 0.95 }}
              whileInView={{ opacity: 1, scale: 1 }}
              viewport={{ once: true }}
              transition={{ delay: 0.2 }}
              className="rounded-2xl border border-white/10 bg-[#0D1117] p-8 shadow-2xl relative overflow-hidden"
            >
              <div className="absolute inset-0 bg-gradient-to-br from-primary/10 to-transparent opacity-50" />
              
              <div className="relative z-10">
                <div className="flex items-center justify-between mb-8">
                  <div className="flex items-center gap-3">
                    <GitBranch className="w-8 h-8 text-white" />
                    <div>
                      <div className="text-white font-medium">cloudstorex/cloudstorex</div>
                      <div className="text-sm text-muted-foreground">The enterprise storage control plane</div>
                    </div>
                  </div>
                </div>

                <div className="grid grid-cols-2 gap-4 mb-6">
                  <div className="p-4 rounded-xl bg-[#05070A] border border-white/5 flex flex-col gap-1">
                    <div className="flex items-center gap-2 text-muted-foreground text-sm">
                      <Star className="w-4 h-4 text-amber-400" /> Stars
                    </div>
                    <div className="text-2xl font-semibold text-white">14.2k</div>
                  </div>
                  <div className="p-4 rounded-xl bg-[#05070A] border border-white/5 flex flex-col gap-1">
                    <div className="flex items-center gap-2 text-muted-foreground text-sm">
                      <GitFork className="w-4 h-4" /> Forks
                    </div>
                    <div className="text-2xl font-semibold text-white">1.8k</div>
                  </div>
                  <div className="p-4 rounded-xl bg-[#05070A] border border-white/5 flex flex-col gap-1">
                    <div className="flex items-center gap-2 text-muted-foreground text-sm">
                      <CircleDot className="w-4 h-4 text-emerald-400" /> Issues
                    </div>
                    <div className="text-2xl font-semibold text-white">42</div>
                  </div>
                  <div className="p-4 rounded-xl bg-[#05070A] border border-white/5 flex flex-col gap-1">
                    <div className="flex items-center gap-2 text-muted-foreground text-sm">
                      <GitPullRequest className="w-4 h-4 text-primary" /> PRs
                    </div>
                    <div className="text-2xl font-semibold text-white">18</div>
                  </div>
                </div>

                <div className="flex items-center gap-6 mb-6 pt-6 border-t border-white/5">
                  <div className="flex items-center gap-2">
                    <div className="px-2 py-1 rounded bg-emerald-500/10 text-emerald-400 text-xs font-mono font-medium border border-emerald-500/20">v2.0.0</div>
                    <span className="text-xs text-muted-foreground">Latest Release</span>
                  </div>
                  <div className="flex items-center gap-2">
                    <div className="px-2 py-1 rounded bg-white/5 text-white/70 text-xs font-mono font-medium border border-white/10">MIT</div>
                    <span className="text-xs text-muted-foreground">License</span>
                  </div>
                </div>

                <div className="flex flex-wrap gap-2">
                  <span className="px-2 py-1 rounded bg-white/5 text-xs text-muted-foreground font-mono">go</span>
                  <span className="px-2 py-1 rounded bg-white/5 text-xs text-muted-foreground font-mono">typescript</span>
                  <span className="px-2 py-1 rounded bg-white/5 text-xs text-muted-foreground font-mono">kubernetes</span>
                  <span className="px-2 py-1 rounded bg-white/5 text-xs text-muted-foreground font-mono">terraform</span>
                </div>
              </div>
            </motion.div>
          </div>
        </div>
      </div>
    </section>
  );
}
