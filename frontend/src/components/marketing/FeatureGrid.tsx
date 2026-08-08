"use client";

import { motion } from "framer-motion";
import { Network, Database, ShieldAlert, Cpu, Layers, RefreshCw } from "lucide-react";

const features = [
  {
    title: "Multi-Cloud Routing",
    description: "Write once to the CloudStoreX API, and dynamically route objects to AWS S3, MinIO, or Azure based on cost and latency.",
    icon: Network,
    span: "col-span-1 md:col-span-2 lg:col-span-2",
  },
  {
    title: "Zero-Trust Security",
    description: "Enterprise IAM with fine-grained RBAC and automated key rotation.",
    icon: ShieldAlert,
    span: "col-span-1",
  },
  {
    title: "Global Metadata Engine",
    description: "Search millions of objects across multiple cloud providers in milliseconds with a unified catalog.",
    icon: Database,
    span: "col-span-1",
  },
  {
    title: "AI Storage Intelligence",
    description: "Forecast capacity, detect anomalies, and auto-tier data using operational ML models.",
    icon: Cpu,
    span: "col-span-1 md:col-span-2",
  },
  {
    title: "Disaster Recovery",
    description: "Active-Active replication across regions.",
    icon: RefreshCw,
    span: "col-span-1",
  },
  {
    title: "Provider Abstraction",
    description: "Never write vendor-specific SDK code again.",
    icon: Layers,
    span: "col-span-1 md:col-span-2 lg:col-span-2",
  },
];

export function FeatureGrid() {
  return (
    <section id="features" className="py-24 relative overflow-hidden bg-[#05070A]">
      <div className="container mx-auto px-6 relative z-10">
        <div className="mb-16">
          <motion.h2 
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            className="text-3xl md:text-5xl font-semibold tracking-tight text-white mb-6"
          >
            Everything you need for <br />
            <span className="text-muted-foreground">Enterprise Storage.</span>
          </motion.h2>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-4 gap-6">
          {features.map((feature, i) => (
            <motion.div
              key={i}
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              whileHover={{ y: -4, boxShadow: "0 10px 40px -10px rgba(59,130,246,0.15)" }}
              viewport={{ once: true }}
              transition={{ delay: i * 0.1, duration: 0.3 }}
              className={`group relative overflow-hidden rounded-2xl glass-panel p-8 transition-all hover:bg-white/[0.08] hover:border-primary/30 ${feature.span}`}
            >
              <div className="absolute inset-0 bg-gradient-to-br from-primary/10 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-500" />
              <div className="relative z-10">
                <div className="w-12 h-12 rounded-lg bg-white/5 border border-white/10 flex items-center justify-center mb-6 text-primary group-hover:bg-primary/20 transition-all duration-300">
                  <motion.div whileHover={{ scale: 1.1, rotate: 5 }} transition={{ type: "spring", stiffness: 400 }}>
                    <feature.icon className="w-6 h-6" />
                  </motion.div>
                </div>
                <h3 className="text-xl font-medium text-white mb-3 group-hover:text-primary transition-colors">{feature.title}</h3>
                <p className="text-muted-foreground leading-relaxed">{feature.description}</p>
              </div>
            </motion.div>
          ))}
        </div>
      </div>
    </section>
  );
}
