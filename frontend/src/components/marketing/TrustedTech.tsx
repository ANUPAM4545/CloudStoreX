"use client";

import { motion } from "framer-motion";

const technologies = [
  "AWS S3",
  "MinIO",
  "Kubernetes",
  "Docker",
  "Terraform",
  "PostgreSQL",
  "Redis",
  "Prometheus",
  "Grafana",
  "OpenTelemetry",
  "Next.js",
  "Go",
];

export function TrustedTech() {
  return (
    <section className="py-20 border-y border-white/5 bg-white/[0.02] overflow-hidden">
      <div className="container mx-auto px-6 mb-10 text-center">
        <p className="text-sm font-medium text-muted-foreground uppercase tracking-widest">
          Trusted enterprise technology stack
        </p>
      </div>
      
      <div className="relative flex overflow-hidden w-full">
        <motion.div 
          className="flex items-center whitespace-nowrap py-4"
          animate={{ x: ["0%", "-50%"] }}
          transition={{ ease: "linear", duration: 50, repeat: Infinity }}
        >
          {[...technologies, ...technologies, ...technologies, ...technologies].map((tech, i) => (
            <div
              key={i}
              className="mx-12 text-xl md:text-2xl font-bold tracking-tight text-white/30 grayscale hover:grayscale-0 hover:text-white hover:drop-shadow-[0_0_15px_rgba(255,255,255,0.3)] transition-all duration-300"
            >
              {tech}
            </div>
          ))}
        </motion.div>
        
        {/* Gradient Masks */}
        <div className="pointer-events-none absolute inset-y-0 left-0 w-1/4 bg-gradient-to-r from-[#05070A] via-[#05070A]/80 to-transparent" />
        <div className="pointer-events-none absolute inset-y-0 right-0 w-1/4 bg-gradient-to-l from-[#05070A] via-[#05070A]/80 to-transparent" />
      </div>
    </section>
  );
}
