"use client";

import { motion, useInView } from "framer-motion";
import { useRef, useEffect, useState } from "react";

const Counter = ({ value, suffix = "", prefix = "", label, decimals = 0 }: { value: number; suffix?: string; prefix?: string; label: string, decimals?: number }) => {
  const ref = useRef(null);
  const isInView = useInView(ref, { once: true, margin: "-100px" });
  const [count, setCount] = useState(0);

  useEffect(() => {
    if (isInView) {
      let start = 0;
      const end = value;
      const duration = 2000;
      let startTime: number | null = null;

      const step = (timestamp: number) => {
        if (!startTime) startTime = timestamp;
        const progress = Math.min((timestamp - startTime) / duration, 1);
        
        // Easing out cubic
        const easeOut = 1 - Math.pow(1 - progress, 3);
        setCount(easeOut * end);

        if (progress < 1) {
          window.requestAnimationFrame(step);
        } else {
          setCount(end); // Ensure exact final value
        }
      };
      
      window.requestAnimationFrame(step);
    }
  }, [isInView, value]);

  return (
    <div ref={ref} className="flex flex-col items-center justify-center p-8 rounded-2xl border border-white/5 bg-white/[0.02]">
      <div className="text-4xl md:text-6xl font-bold text-white mb-2 tracking-tighter">
        {prefix}{count.toFixed(decimals)}{suffix}
      </div>
      <div className="text-sm md:text-base text-muted-foreground uppercase tracking-widest font-medium text-center">
        {label}
      </div>
    </div>
  );
};

export function Performance() {
  return (
    <section className="py-24 relative overflow-hidden bg-background">
      <div className="container mx-auto px-6">
        <div className="text-center max-w-3xl mx-auto mb-16">
          <motion.h2 
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            className="text-3xl md:text-5xl font-semibold tracking-tight text-white mb-6"
          >
            Engineered for <span className="text-gradient">Scale.</span>
          </motion.h2>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          <Counter value={10} prefix="<" suffix="ms" label="Policy Evaluation" />
          <Counter value={99.999} decimals={3} suffix="%" label="Availability SLA" />
          <Counter value={50} suffix="B+" label="Objects Indexed" />
          <Counter value={4.2} decimals={1} suffix="TB/s" label="Sustained Throughput" />
        </div>
      </div>
    </section>
  );
}
