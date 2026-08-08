"use client";

import { motion } from "framer-motion";
import { Network, Database, Cloud, Zap, ArrowRight, CornerDownRight } from "lucide-react";

export function PolicyEngine() {
  return (
    <section className="py-24 relative overflow-hidden bg-[#05070A] border-y border-white/5">
      <div className="container mx-auto px-6">
        <div className="flex flex-col lg:flex-row items-center gap-16">
          <div className="flex-1 w-full order-2 lg:order-1 relative">
            {/* Visual Decision Tree */}
            <div className="glass-panel p-8 rounded-2xl border border-white/10 relative">
              <div className="absolute inset-0 bg-[radial-gradient(ellipse_at_center,rgba(16,185,129,0.1)_0%,transparent_70%)]" />
              
              <div className="relative z-10 font-mono text-sm">
                {/* Incoming Request */}
                <motion.div 
                  initial={{ opacity: 0, x: -20 }}
                  whileInView={{ opacity: 1, x: 0 }}
                  viewport={{ once: true }}
                  className="flex items-center gap-3 bg-white/5 p-3 rounded-lg border border-white/10 w-max mb-6"
                >
                  <Zap className="w-4 h-4 text-emerald-400" />
                  <span className="text-white">PUT /dataset/model-weights.h5</span>
                  <span className="text-muted-foreground ml-2">Size: 4.2GB</span>
                </motion.div>

                <div className="ml-6 border-l border-white/10 pl-6 space-y-6">
                  {/* IF Condition */}
                  <motion.div 
                    initial={{ opacity: 0 }}
                    whileInView={{ opacity: 1 }}
                    viewport={{ once: true }}
                    transition={{ delay: 0.3 }}
                    className="relative"
                  >
                    <CornerDownRight className="w-5 h-5 text-white/20 absolute -left-11 -top-1" />
                    <div className="flex items-center gap-2 mb-2">
                      <span className="text-primary font-bold">IF</span>
                      <span className="text-white">object.size {">"} 1GB</span>
                    </div>
                    
                    {/* Routing Action 1 */}
                    <motion.div 
                      initial={{ opacity: 0, x: -10 }}
                      whileInView={{ opacity: 1, x: 0 }}
                      viewport={{ once: true }}
                      transition={{ delay: 0.6 }}
                      className="ml-4 mt-2 flex items-center gap-3 bg-primary/20 border border-primary/30 p-3 rounded-lg w-max"
                    >
                      <ArrowRight className="w-4 h-4 text-primary" />
                      <span className="text-white">Route to</span>
                      <div className="flex items-center gap-2 bg-white/10 px-2 py-1 rounded text-white">
                        <Cloud className="w-4 h-4" /> AWS S3 (Infrequent Access)
                      </div>
                    </motion.div>
                  </motion.div>

                  {/* ELSE Condition */}
                  <motion.div 
                    initial={{ opacity: 0 }}
                    whileInView={{ opacity: 1 }}
                    viewport={{ once: true }}
                    transition={{ delay: 0.9 }}
                    className="relative"
                  >
                    <CornerDownRight className="w-5 h-5 text-white/20 absolute -left-11 -top-1" />
                    <div className="flex items-center gap-2 mb-2">
                      <span className="text-muted-foreground font-bold">ELSE IF</span>
                      <span className="text-muted-foreground">request.latency_requirement {"<"} 10ms</span>
                    </div>
                    
                    {/* Routing Action 2 */}
                    <div className="ml-4 mt-2 flex items-center gap-3 bg-white/5 border border-white/10 p-3 rounded-lg w-max opacity-50">
                      <ArrowRight className="w-4 h-4 text-muted-foreground" />
                      <span className="text-muted-foreground">Route to</span>
                      <div className="flex items-center gap-2 bg-white/5 px-2 py-1 rounded text-muted-foreground">
                        <Database className="w-4 h-4" /> MinIO (Local NVMe)
                      </div>
                    </div>
                  </motion.div>
                </div>
              </div>
            </div>
          </div>

          <div className="flex-1 w-full order-1 lg:order-2">
            <motion.div
              initial={{ opacity: 0, x: 20 }}
              whileInView={{ opacity: 1, x: 0 }}
              viewport={{ once: true }}
              className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-emerald-500/10 text-sm font-medium text-emerald-400 mb-6"
            >
              <Network className="w-4 h-4" />
              Dynamic Policy Engine
            </motion.div>
            <motion.h2 
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ delay: 0.1 }}
              className="text-3xl md:text-5xl font-semibold tracking-tight text-white mb-6"
            >
              Code as infrastructure. <br />
              <span className="text-muted-foreground">Rules as code.</span>
            </motion.h2>
            <motion.p 
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ delay: 0.2 }}
              className="text-lg text-muted-foreground mb-8"
            >
              Control exactly where your data goes based on size, extension, region, or custom metadata tags. Eliminate vendor lock-in by writing one unified API request and letting the policy engine determine the optimal storage provider.
            </motion.p>
          </div>
        </div>
      </div>
    </section>
  );
}
