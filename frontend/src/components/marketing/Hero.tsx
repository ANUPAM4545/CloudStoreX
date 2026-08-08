"use client";

import { motion, useMotionValue, useSpring, useTransform } from "framer-motion";
import Link from "next/link";
import { ArrowRight, GitBranch, BookOpen } from "lucide-react";
import { useEffect } from "react";

const AbstractTopology = () => {
  const mouseX = useMotionValue(0);
  const mouseY = useMotionValue(0);

  const springX = useSpring(mouseX, { stiffness: 50, damping: 20 });
  const springY = useSpring(mouseY, { stiffness: 50, damping: 20 });

  useEffect(() => {
    const handleMouseMove = (e: MouseEvent) => {
      mouseX.set(e.clientX / window.innerWidth - 0.5);
      mouseY.set(e.clientY / window.innerHeight - 0.5);
    };
    window.addEventListener("mousemove", handleMouseMove);
    return () => window.removeEventListener("mousemove", handleMouseMove);
  }, [mouseX, mouseY]);

  const translateX = useTransform(springX, [-0.5, 0.5], [-20, 20]);
  const translateY = useTransform(springY, [-0.5, 0.5], [-20, 20]);

  return (
    <div className="absolute inset-0 overflow-hidden pointer-events-none flex items-center justify-center opacity-25">
      <motion.svg 
        style={{ x: translateX, y: translateY }}
        width="100%" height="100%" viewBox="0 0 1200 800" fill="none" xmlns="http://www.w3.org/2000/svg"
      >
        <defs>
          <filter id="glow" x="-20%" y="-20%" width="140%" height="140%">
            <feGaussianBlur stdDeviation="8" result="blur" />
            <feComposite in="SourceGraphic" in2="blur" operator="over" />
          </filter>
        </defs>
        <motion.path
          d="M 200 400 L 400 200 L 800 200 L 1000 400 L 800 600 L 400 600 Z"
          stroke="currentColor"
          strokeWidth="1.5"
          className="text-primary/40"
          initial={{ pathLength: 0, opacity: 0 }}
          animate={{ pathLength: 1, opacity: 1 }}
          transition={{ duration: 4, ease: "easeInOut", repeat: Infinity, repeatType: "mirror" }}
          filter="url(#glow)"
        />
        <motion.path
          d="M 400 200 L 400 600 M 800 200 L 800 600 M 200 400 L 1000 400"
          stroke="currentColor"
          strokeWidth="1"
          className="text-white/10"
          initial={{ pathLength: 0, opacity: 0 }}
          animate={{ pathLength: 1, opacity: 1 }}
          transition={{ duration: 5, ease: "easeInOut", repeat: Infinity, repeatType: "mirror" }}
        />
        
        {/* Animated Data Packets */}
        <motion.circle
          r="4"
          className="fill-primary"
          filter="url(#glow)"
          initial={{ x: 200, y: 400, opacity: 0 }}
          animate={{ x: 400, y: 200, opacity: 1 }}
          transition={{ duration: 2, ease: "linear", repeat: Infinity }}
        />
        <motion.circle
          r="4"
          className="fill-emerald-400"
          filter="url(#glow)"
          initial={{ x: 800, y: 200, opacity: 0 }}
          animate={{ x: 1000, y: 400, opacity: 1 }}
          transition={{ duration: 2.5, ease: "linear", repeat: Infinity, delay: 1 }}
        />
        <motion.circle
          r="4"
          className="fill-primary"
          filter="url(#glow)"
          initial={{ x: 400, y: 600, opacity: 0 }}
          animate={{ x: 200, y: 400, opacity: 1 }}
          transition={{ duration: 2, ease: "linear", repeat: Infinity, delay: 0.5 }}
        />
        
        {/* Nodes */}
        {[
          [200, 400], [400, 200], [800, 200], [1000, 400], [800, 600], [400, 600]
        ].map(([cx, cy], i) => (
          <motion.circle
            key={i}
            cx={cx}
            cy={cy}
            r="5"
            className="fill-white/30"
            filter="url(#glow)"
            initial={{ scale: 0 }}
            animate={{ scale: [1, 1.3, 1] }}
            transition={{ delay: 0.5 + i * 0.1, duration: 3, repeat: Infinity }}
          />
        ))}
      </motion.svg>
      
      {/* Radial Gradient for fade out on edges */}
      <div className="absolute inset-0 bg-[radial-gradient(ellipse_at_center,transparent_20%,#05070A_70%)]" />
    </div>
  );
};

export function Hero() {
  return (
    <section className="relative pt-32 pb-20 md:pt-48 md:pb-32 overflow-hidden flex flex-col items-center justify-center min-h-[90vh]">
      <AbstractTopology />
      
      <div className="container relative z-10 mx-auto px-6 flex flex-col items-center text-center">
        <motion.div
          initial={{ y: 20, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          transition={{ duration: 0.8, ease: [0.16, 1, 0.3, 1] }}
          className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-white/5 border border-white/10 text-sm text-muted-foreground mb-8 backdrop-blur-sm"
        >
          <span className="flex h-2 w-2 rounded-full bg-primary animate-pulse" />
          Enterprise AI Storage Control Plane
        </motion.div>

        <motion.h1
          initial={{ y: 20, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          transition={{ duration: 0.8, delay: 0.1, ease: [0.16, 1, 0.3, 1] }}
          className="text-5xl md:text-7xl lg:text-8xl font-semibold tracking-tighter text-white max-w-5xl leading-[1.1]"
        >
          AI-Native Multi-Cloud <br className="hidden md:block" />
          <span className="text-gradient">Storage Infrastructure.</span>
        </motion.h1>

        <motion.p
          initial={{ y: 20, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          transition={{ duration: 0.8, delay: 0.2, ease: [0.16, 1, 0.3, 1] }}
          className="mt-6 text-lg md:text-xl text-muted-foreground max-w-2xl font-light"
        >
          The enterprise control plane that unifies AWS, MinIO, and Kubernetes with dynamic policy routing, zero-trust security, and operational AI.
        </motion.p>

        <motion.div
          initial={{ y: 20, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          transition={{ duration: 0.8, delay: 0.3, ease: [0.16, 1, 0.3, 1] }}
          className="flex flex-col sm:flex-row items-center gap-4 mt-10"
        >
          <Link
            href="/login"
            className="flex items-center justify-center gap-2 px-8 py-4 text-sm font-medium text-white transition-all rounded-full bg-primary hover:bg-primary/90 hover:scale-105 active:scale-95 w-full sm:w-auto shadow-[0_0_40px_-10px_rgba(59,130,246,0.5)]"
          >
            Start Free <ArrowRight className="w-4 h-4" />
          </Link>
          <Link
            href="#"
            className="flex items-center justify-center gap-2 px-8 py-4 text-sm font-medium text-white transition-all rounded-full bg-white/5 border border-white/10 hover:bg-white/10 w-full sm:w-auto"
          >
            <GitBranch className="w-4 h-4" /> GitHub
          </Link>
          <Link
            href="#"
            className="flex items-center justify-center gap-2 px-8 py-4 text-sm font-medium text-white transition-all rounded-full bg-white/5 border border-white/10 hover:bg-white/10 w-full sm:w-auto"
          >
            <BookOpen className="w-4 h-4" /> Documentation
          </Link>
        </motion.div>
      </div>
    </section>
  );
}
