"use client";

import { motion } from "framer-motion";
import { Server, Database, Cloud, Shield, Search, BrainCircuit, Activity } from "lucide-react";

const Node = ({ icon: Icon, label, x, y, delay = 0, isActive = false }: any) => (
  <motion.div
    initial={{ opacity: 0, scale: 0.8 }}
    whileInView={{ opacity: 1, scale: 1 }}
    whileHover={{ scale: 1.05, borderColor: "rgba(255,255,255,0.3)", backgroundColor: "rgba(255,255,255,0.1)" }}
    viewport={{ once: true, margin: "-100px" }}
    animate={isActive ? { boxShadow: ["0 0 0px 0px rgba(59,130,246,0)", "0 0 30px 5px rgba(59,130,246,0.4)", "0 0 0px 0px rgba(59,130,246,0)"] } : {}}
    transition={{ duration: isActive ? 2 : 0.3, delay: isActive ? 0 : delay, repeat: isActive ? Infinity : 0 }}
    className={`absolute flex flex-col items-center justify-center p-4 rounded-xl border backdrop-blur-md transition-colors w-32 cursor-pointer ${
      isActive 
        ? "bg-primary/20 border-primary shadow-[0_0_30px_-5px_rgba(59,130,246,0.3)] text-white" 
        : "bg-white/5 border-white/10 text-muted-foreground"
    }`}
    style={{ left: `calc(${x}% - 64px)`, top: `calc(${y}% - 40px)` }}
  >
    <Icon className={`w-6 h-6 mb-2 ${isActive ? "text-primary" : ""}`} />
    <span className="text-xs font-medium text-center">{label}</span>
  </motion.div>
);

const Connection = ({ d, duration = 3, delay = 0, reverse = false }: any) => (
  <>
    <path d={d} stroke="rgba(255,255,255,0.05)" strokeWidth="2" fill="none" />
    <motion.path
      d={d}
      stroke="url(#gradient)"
      strokeWidth="2"
      fill="none"
      initial={{ pathLength: 0, opacity: 0 }}
      whileInView={{ pathLength: 1, opacity: 1 }}
      viewport={{ once: true }}
      transition={{ 
        duration, 
        delay, 
        ease: "easeInOut", 
      }}
    />
    <motion.path
      d={d}
      stroke="#3b82f6"
      strokeWidth="3"
      fill="none"
      strokeDasharray="4 40"
      initial={{ strokeDashoffset: reverse ? -100 : 100, opacity: 0 }}
      whileInView={{ opacity: 1 }}
      viewport={{ once: true }}
      animate={{ strokeDashoffset: reverse ? 100 : -100 }}
      transition={{ 
        duration: duration * 0.8, 
        delay: delay + 0.5,
        ease: "linear", 
        repeat: Infinity,
      }}
      className="drop-shadow-[0_0_8px_rgba(59,130,246,0.8)]"
    />
  </>
);

export function InteractiveArchitecture() {
  return (
    <section id="architecture" className="py-24 md:py-32 relative">
      <div className="container mx-auto px-6">
        <div className="flex flex-col items-center text-center mb-20">
          <motion.h2 
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            className="text-3xl md:text-5xl font-semibold tracking-tight text-white mb-6"
          >
            The Orchestration <span className="text-gradient">Control Plane</span>
          </motion.h2>
          <motion.p 
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            transition={{ delay: 0.1 }}
            className="text-lg text-muted-foreground max-w-2xl"
          >
            A unified architecture abstracting complex storage layers into a single intelligent API gateway with real-time policy routing.
          </motion.p>
        </div>

        <div className="relative w-full max-w-5xl mx-auto h-[600px] glass-panel rounded-2xl overflow-hidden hidden md:block">
          <svg className="absolute inset-0 w-full h-full" style={{ zIndex: 0 }}>
            <defs>
              <linearGradient id="gradient" x1="0%" y1="0%" x2="100%" y2="0%">
                <stop offset="0%" stopColor="rgba(59,130,246,0)" />
                <stop offset="50%" stopColor="rgba(59,130,246,1)" />
                <stop offset="100%" stopColor="rgba(59,130,246,0)" />
              </linearGradient>
            </defs>
            
            {/* Connections */}
            <Connection d="M 150 300 L 250 300" delay={0.5} />
            <Connection d="M 350 300 L 450 300" delay={0.8} />
            <Connection d="M 550 300 C 650 300, 650 150, 750 150" delay={1.1} />
            <Connection d="M 550 300 C 650 300, 650 450, 750 450" delay={1.1} />
            <Connection d="M 450 100 L 450 250" duration={2} delay={1.5} reverse />
            <Connection d="M 450 500 L 450 350" duration={2} delay={1.7} reverse />
            <Connection d="M 850 150 C 950 150, 950 300, 750 300" duration={4} delay={2} />
            <Connection d="M 850 450 C 950 450, 950 300, 750 300" duration={4} delay={2.2} />
          </svg>

          {/* Nodes */}
          <div className="absolute inset-0 z-10">
            <Node icon={Server} label="Client Application" x={10} y={50} delay={0.2} />
            <Node icon={Shield} label="API Gateway" x={30} y={50} isActive delay={0.4} />
            
            <Node icon={BrainCircuit} label="AI Intelligence" x={50} y={15} delay={0.6} />
            <Node icon={Search} label="Policy Engine" x={50} y={50} isActive delay={0.8} />
            <Node icon={Activity} label="Observability" x={50} y={85} delay={0.7} />
            
            <Node icon={Cloud} label="AWS S3" x={80} y={25} delay={1.0} />
            <Node icon={Database} label="MinIO Cluster" x={80} y={75} delay={1.0} />
            
            <Node icon={Database} label="Metadata Catalog" x={95} y={50} delay={1.2} />
          </div>
        </div>

        {/* Mobile View */}
        <div className="md:hidden flex flex-col gap-4">
          <div className="p-6 rounded-xl glass-panel flex items-center gap-4">
            <Server className="w-8 h-8 text-white" />
            <div>
              <h3 className="text-white font-medium">1. Client Application</h3>
              <p className="text-sm text-muted-foreground">Standard S3-compatible requests</p>
            </div>
          </div>
          <div className="w-0.5 h-8 bg-gradient-to-b from-white/20 to-primary/50 mx-auto" />
          <div className="p-6 rounded-xl bg-primary/20 border border-primary shadow-[0_0_30px_-5px_rgba(59,130,246,0.3)] flex items-center gap-4">
            <Shield className="w-8 h-8 text-primary" />
            <div>
              <h3 className="text-white font-medium">2. Policy Engine</h3>
              <p className="text-sm text-white/70">Analyzes and routes request</p>
            </div>
          </div>
          <div className="w-0.5 h-8 bg-gradient-to-b from-primary/50 to-white/20 mx-auto" />
          <div className="p-6 rounded-xl glass-panel flex items-center gap-4">
            <Cloud className="w-8 h-8 text-white" />
            <div>
              <h3 className="text-white font-medium">3. Storage Provider</h3>
              <p className="text-sm text-muted-foreground">AWS S3 or MinIO executes</p>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
