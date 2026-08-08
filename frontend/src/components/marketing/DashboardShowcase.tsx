"use client";

import { motion, useScroll, useTransform } from "framer-motion";
import { HardDrive, Server, Shield, BrainCircuit, Activity } from "lucide-react";
import { useRef } from "react";

export function DashboardShowcase() {
  const containerRef = useRef(null);
  const { scrollYProgress } = useScroll({
    target: containerRef,
    offset: ["start end", "center center"]
  });
  
  const scale = useTransform(scrollYProgress, [0, 1], [0.85, 1]);
  const opacity = useTransform(scrollYProgress, [0, 1], [0.4, 1]);

  return (
    <section className="py-24 relative overflow-hidden bg-background">
      <div className="container mx-auto px-6 text-center max-w-4xl mb-16">
        <motion.h2 
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          className="text-3xl md:text-5xl font-semibold tracking-tight text-white mb-6"
        >
          A control plane built for <span className="text-gradient">production.</span>
        </motion.h2>
        <motion.p 
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ delay: 0.1 }}
          className="text-lg text-muted-foreground"
        >
          No more jumping between AWS consoles and MinIO terminals. Manage your entire global storage infrastructure from a single, unified enterprise dashboard.
        </motion.p>
      </div>

      <div className="container mx-auto px-6" ref={containerRef}>
        {/* CSS Mockup of the CloudStoreX Dashboard */}
        <motion.div 
          style={{ scale, opacity }}
          className="relative max-w-6xl mx-auto rounded-xl overflow-hidden border border-white/10 shadow-2xl bg-[#0D1117] transform-gpu"
        >
          {/* Topbar */}
          <div className="h-12 border-b border-white/5 flex items-center px-4 gap-4 bg-[#05070A]">
            <div className="flex gap-1.5">
              <div className="w-3 h-3 rounded-full bg-white/10" />
              <div className="w-3 h-3 rounded-full bg-white/10" />
              <div className="w-3 h-3 rounded-full bg-white/10" />
            </div>
            <div className="flex-1 max-w-md mx-auto h-6 rounded bg-white/5 border border-white/10 flex items-center px-3 text-xs text-muted-foreground font-mono">
              cloudstorex.corp.internal / production
            </div>
          </div>

          <div className="flex h-[600px]">
            {/* Sidebar */}
            <div className="w-64 border-r border-white/5 p-4 flex flex-col gap-2 bg-[#05070A]/50 hidden md:flex">
              <div className="text-xs font-medium text-muted-foreground mb-2 mt-2 px-2 uppercase tracking-wider">Overview</div>
              <div className="flex items-center gap-3 px-2 py-1.5 rounded-lg bg-primary/10 text-primary text-sm font-medium">
                <Activity className="w-4 h-4" /> Global Metrics
              </div>
              <div className="flex items-center gap-3 px-2 py-1.5 text-muted-foreground text-sm font-medium hover:text-white">
                <HardDrive className="w-4 h-4" /> Storage Buckets
              </div>
              <div className="flex items-center gap-3 px-2 py-1.5 text-muted-foreground text-sm font-medium hover:text-white">
                <Server className="w-4 h-4" /> Providers
              </div>
              <div className="text-xs font-medium text-muted-foreground mb-2 mt-6 px-2 uppercase tracking-wider">Intelligence</div>
              <div className="flex items-center gap-3 px-2 py-1.5 text-muted-foreground text-sm font-medium hover:text-white">
                <BrainCircuit className="w-4 h-4" /> AI Operations
              </div>
              <div className="flex items-center gap-3 px-2 py-1.5 text-muted-foreground text-sm font-medium hover:text-white">
                <Shield className="w-4 h-4" /> Security & IAM
              </div>
            </div>

            {/* Main Content Area */}
            <div className="flex-1 p-8 bg-[#0D1117] overflow-hidden flex flex-col gap-6">
              <div className="flex items-center justify-between">
                <h3 className="text-2xl font-semibold text-white">Global Metrics</h3>
                <div className="flex gap-2">
                  <div className="h-8 w-24 rounded border border-white/10 bg-white/5" />
                  <div className="h-8 w-32 rounded bg-primary" />
                </div>
              </div>

              {/* Stat Cards */}
              <div className="grid grid-cols-3 gap-4">
                {[
                  { label: "Total Storage", value: "4.2 PB", trend: "+12%" },
                  { label: "Active Objects", value: "1.2B", trend: "+5%" },
                  { label: "Network Egress", value: "850 TB/mo", trend: "-2%" },
                ].map((stat, i) => (
                  <div key={i} className="p-4 rounded-xl border border-white/10 bg-white/5 flex flex-col gap-2">
                    <span className="text-sm text-muted-foreground">{stat.label}</span>
                    <div className="flex items-end justify-between">
                      <span className="text-3xl font-semibold text-white">{stat.value}</span>
                      <span className={`text-xs font-medium ${stat.trend.startsWith("+") ? "text-emerald-400" : "text-primary"}`}>{stat.trend}</span>
                    </div>
                  </div>
                ))}
              </div>

              {/* Chart Mockup */}
              <div className="flex-1 rounded-xl border border-white/10 bg-white/5 p-4 flex flex-col">
                <div className="text-sm text-muted-foreground mb-6">Storage Growth vs Capacity</div>
                <div className="flex-1 flex items-end gap-2 px-4 pb-4">
                  {[40, 50, 45, 60, 75, 65, 80, 95, 85, 100, 110, 105].map((height, i) => (
                    <motion.div
                      key={i}
                      initial={{ height: 0 }}
                      whileInView={{ height: `${height}%` }}
                      viewport={{ once: true }}
                      transition={{ delay: 0.2 + i * 0.05, duration: 0.8, ease: "easeOut" }}
                      className="flex-1 bg-gradient-to-t from-primary/20 to-primary/80 rounded-t-sm"
                    />
                  ))}
                </div>
              </div>
            </div>
          </div>
          
          {/* Subtle Reflection */}
          <div className="absolute inset-0 bg-gradient-to-b from-white/5 to-transparent pointer-events-none rounded-xl" />
        </motion.div>
      </div>
    </section>
  );
}
