"use client";

import React from "react";
import { motion } from "framer-motion";
import { Server, Activity, Database, Shield, FileText, Bot, ShieldCheck, ArrowRight } from "lucide-react";

export function RightPanel() {
  return (
    <div className="hidden lg:flex w-[50%] bg-zinc-950 relative overflow-hidden items-center justify-center border-l border-white/5">
      
      {/* 1. Animated Background */}
      <div className="absolute inset-0 z-0">
        <div className="absolute inset-0 bg-[linear-gradient(to_right,#80808012_1px,transparent_1px),linear-gradient(to_bottom,#80808012_1px,transparent_1px)] bg-[size:24px_24px] [mask-image:radial-gradient(ellipse_60%_50%_at_50%_50%,#000_70%,transparent_100%)]" />
        
        {/* Soft Ambient Glows */}
        <motion.div 
          animate={{ opacity: [0.3, 0.5, 0.3], scale: [1, 1.05, 1] }}
          transition={{ duration: 8, repeat: Infinity, ease: "easeInOut" }}
          className="absolute top-1/4 left-1/4 w-[400px] h-[400px] bg-blue-600/20 rounded-full blur-[120px] mix-blend-screen" 
        />
        <motion.div 
          animate={{ opacity: [0.2, 0.4, 0.2], scale: [1, 1.1, 1] }}
          transition={{ duration: 10, repeat: Infinity, ease: "easeInOut", delay: 2 }}
          className="absolute bottom-1/4 right-1/4 w-[500px] h-[500px] bg-emerald-600/15 rounded-full blur-[150px] mix-blend-screen" 
        />
        
        <div className="absolute inset-0 bg-[url('/noise.svg')] opacity-[0.03] mix-blend-overlay" />
      </div>

      {/* 2. Top Branding */}
      <div className="absolute top-8 left-8 z-20 flex items-center gap-2">
        <div className="w-8 h-8 bg-white rounded-lg flex items-center justify-center">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M12 2L2 7L12 12L22 7L12 2Z" fill="black"/>
            <path d="M2 17L12 22L22 17" stroke="black" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
            <path d="M2 12L12 17L22 12" stroke="black" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
          </svg>
        </div>
        <span className="font-semibold text-white tracking-tight text-lg">CloudStoreX</span>
      </div>

      {/* 3. Floating Enterprise Widgets */}
      <div className="relative z-10 w-full max-w-lg h-[600px]">
        
        {/* Widget 1: Provider Health */}
        <motion.div 
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.2, duration: 0.8 }}
          className="absolute top-[5%] left-[5%] w-64 bg-black/40 backdrop-blur-md border border-white/10 rounded-xl p-4 shadow-2xl"
        >
          <div className="flex items-center gap-2 mb-3 text-white/90">
            <Server size={14} className="text-blue-400" />
            <span className="text-xs font-semibold uppercase tracking-wider">Provider Health</span>
          </div>
          <div className="space-y-3">
            <div className="flex justify-between items-center text-sm">
              <span className="text-white/70">AWS S3</span>
              <div className="flex items-center gap-1.5"><div className="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse" /><span className="text-emerald-400 text-xs">Healthy</span></div>
            </div>
            <div className="flex justify-between items-center text-sm">
              <span className="text-white/70">MinIO Cluster</span>
              <div className="flex items-center gap-1.5"><div className="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse" /><span className="text-emerald-400 text-xs">Healthy</span></div>
            </div>
            <div className="flex justify-between items-center text-sm">
              <span className="text-white/70">Azure Blob</span>
              <div className="flex items-center gap-1.5"><div className="w-1.5 h-1.5 rounded-full bg-blue-500" /><span className="text-blue-400 text-xs">Ready</span></div>
            </div>
          </div>
        </motion.div>

        {/* Widget 2: Storage Overview */}
        <motion.div 
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.4, duration: 0.8 }}
          className="absolute top-[20%] right-[0%] w-56 bg-black/40 backdrop-blur-md border border-white/10 rounded-xl p-4 shadow-2xl"
        >
          <div className="flex items-center gap-2 mb-4 text-white/90">
            <Database size={14} className="text-emerald-400" />
            <span className="text-xs font-semibold uppercase tracking-wider">Storage Overview</span>
          </div>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <div className="text-2xl font-bold text-white mb-0.5">24</div>
              <div className="text-[10px] text-white/50 uppercase tracking-wide">Buckets</div>
            </div>
            <div>
              <div className="text-2xl font-bold text-white mb-0.5">82<span className="text-sm text-white/50">TB</span></div>
              <div className="text-[10px] text-white/50 uppercase tracking-wide">Total Size</div>
            </div>
          </div>
        </motion.div>

        {/* Widget 3: AI Recommendation */}
        <motion.div 
          initial={{ opacity: 0, scale: 0.95 }}
          animate={{ opacity: 1, scale: 1 }}
          transition={{ delay: 0.6, duration: 0.8 }}
          className="absolute top-[45%] left-[20%] w-72 bg-gradient-to-br from-indigo-500/10 to-purple-500/10 backdrop-blur-xl border border-indigo-500/20 rounded-xl p-5 shadow-[0_0_40px_rgba(99,102,241,0.15)]"
        >
          <div className="flex items-center gap-2 mb-3 text-indigo-300">
            <Bot size={16} />
            <span className="text-xs font-semibold uppercase tracking-wider">AI Intelligence</span>
          </div>
          <p className="text-sm text-white/90 leading-relaxed mb-3">
            Identified <span className="text-white font-medium">1.2M</span> unused objects in <code className="text-xs bg-white/10 px-1 py-0.5 rounded">logs/</code>. Moving to Glacier will reduce costs.
          </p>
          <div className="flex items-center justify-between mt-4 pt-3 border-t border-indigo-500/20">
            <span className="text-xs text-indigo-300">Est. Savings</span>
            <span className="text-sm font-bold text-emerald-400">$2,450 / mo</span>
          </div>
        </motion.div>

        {/* Widget 4: Policy Decision */}
        <motion.div 
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.8, duration: 0.8 }}
          className="absolute bottom-[10%] right-[10%] w-64 bg-black/40 backdrop-blur-md border border-white/10 rounded-xl p-4 shadow-2xl"
        >
          <div className="flex items-center gap-2 mb-3 text-white/90">
            <Shield size={14} className="text-amber-400" />
            <span className="text-xs font-semibold uppercase tracking-wider">Policy Engine</span>
          </div>
          <div className="bg-white/5 border border-white/10 rounded-lg p-3">
            <div className="flex items-center justify-between mb-1">
              <span className="text-[10px] text-white/50 uppercase">Match</span>
              <span className="text-[10px] text-emerald-400">Allowed</span>
            </div>
            <div className="text-sm text-white font-medium mb-1">Large File Policy</div>
            <div className="flex items-center gap-2 text-xs text-white/70">
              <span>Client</span>
              <ArrowRight size={10} className="text-white/30" />
              <span className="text-amber-400">AWS S3</span>
            </div>
          </div>
        </motion.div>

        {/* Connecting SVG Lines (Visualizing Data Flow) */}
        <svg className="absolute inset-0 w-full h-full pointer-events-none z-0" style={{ filter: "drop-shadow(0 0 8px rgba(255,255,255,0.1))" }}>
          <motion.path 
            d="M 150 120 C 250 120, 250 300, 350 300" 
            fill="none" 
            stroke="url(#gradient-line)" 
            strokeWidth="1.5"
            initial={{ pathLength: 0, opacity: 0 }}
            animate={{ pathLength: 1, opacity: 0.3 }}
            transition={{ delay: 1, duration: 1.5, ease: "easeInOut" }}
          />
          <motion.circle 
            cx="350" cy="300" r="3" fill="#60A5FA"
            initial={{ opacity: 0, scale: 0 }}
            animate={{ opacity: [0, 1, 0.5], scale: [0, 1.5, 1] }}
            transition={{ delay: 2.2, duration: 0.5 }}
          />
          <defs>
            <linearGradient id="gradient-line" x1="0%" y1="0%" x2="100%" y2="0%">
              <stop offset="0%" stopColor="#34D399" stopOpacity="0" />
              <stop offset="50%" stopColor="#60A5FA" stopOpacity="1" />
              <stop offset="100%" stopColor="#818CF8" stopOpacity="0.2" />
            </linearGradient>
          </defs>
        </svg>

      </div>

      {/* 4. Tech Stack Badges */}
      <div className="absolute bottom-8 left-0 right-0 z-20 flex justify-center">
        <div className="flex flex-col items-center gap-3">
          <span className="text-[10px] font-semibold text-white/30 uppercase tracking-widest">Enterprise Infrastructure</span>
          <div className="flex flex-wrap justify-center gap-2 max-w-sm px-8">
            {["Go", "Next.js", "PostgreSQL", "Redis", "Kubernetes", "AWS", "Terraform", "OpenTelemetry"].map((tech, i) => (
              <motion.div 
                key={tech}
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: 1 + i * 0.1 }}
                className="px-2.5 py-1 rounded-md bg-white/5 border border-white/10 text-[11px] font-medium text-white/60 hover:text-white hover:bg-white/10 transition-colors cursor-default"
              >
                {tech}
              </motion.div>
            ))}
          </div>
        </div>
      </div>

    </div>
  );
}
