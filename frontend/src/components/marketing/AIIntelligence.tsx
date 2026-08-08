"use client";

import { motion, useInView } from "framer-motion";
import { Sparkles, BarChart3, ArrowUpRight, CheckCircle2 } from "lucide-react";
import { useState, useEffect, useRef } from "react";

const TypewriterText = ({ text, delay = 0, onComplete, cursor = true }: { text: string; delay?: number, onComplete?: () => void, cursor?: boolean }) => {
  const [displayText, setDisplayText] = useState("");
  const [isTyping, setIsTyping] = useState(false);
  const [hasStarted, setHasStarted] = useState(false);

  useEffect(() => {
    let i = 0;
    let interval: NodeJS.Timeout;
    const timer = setTimeout(() => {
      setHasStarted(true);
      setIsTyping(true);
      interval = setInterval(() => {
        if (i < text.length) {
          setDisplayText(text.substring(0, i + 1));
          i++;
        } else {
          clearInterval(interval);
          setIsTyping(false);
          if (onComplete) onComplete();
        }
      }, 30);
    }, delay * 1000);
    return () => {
      clearTimeout(timer);
      if (interval) clearInterval(interval);
    };
  }, [text, delay, onComplete]);

  return (
    <span>
      {displayText}
      {cursor && (isTyping || !hasStarted) && (
        <motion.span 
          initial={{ opacity: 0 }} 
          animate={{ opacity: 1 }} 
          transition={{ repeat: Infinity, duration: 0.4, repeatType: "reverse" }} 
          className="inline-block w-1.5 h-4 ml-0.5 bg-white/70 align-middle" 
        />
      )}
    </span>
  );
};

const CountUp = ({ to, prefix = "", suffix = "", delay = 0 }: { to: number, prefix?: string, suffix?: string, delay?: number }) => {
  const [value, setValue] = useState(0);
  const ref = useRef(null);
  const isInView = useInView(ref, { once: true });

  useEffect(() => {
    if (!isInView) return;
    let startTime: number;
    let animationFrame: number;
    const duration = 1500;
    
    const timer = setTimeout(() => {
      const step = (timestamp: number) => {
        if (!startTime) startTime = timestamp;
        const progress = Math.min((timestamp - startTime) / duration, 1);
        const ease = 1 - Math.pow(1 - progress, 4); // easeOutQuart
        setValue(Math.floor(ease * to));
        if (progress < 1) {
          animationFrame = requestAnimationFrame(step);
        }
      };
      animationFrame = requestAnimationFrame(step);
    }, delay * 1000);
    
    return () => {
      clearTimeout(timer);
      if (animationFrame) cancelAnimationFrame(animationFrame);
    };
  }, [to, delay, isInView]);

  return <span ref={ref}>{prefix}{value.toLocaleString()}{suffix}</span>;
};

export function AIIntelligence() {
  const [questionDone, setQuestionDone] = useState(false);

  return (
    <section className="py-24 relative overflow-hidden bg-background">
      <div className="container mx-auto px-6">
        <div className="flex flex-col lg:flex-row items-center gap-16">
          <div className="flex-1 w-full max-w-xl">
            <motion.div
              initial={{ opacity: 0, x: -20 }}
              whileInView={{ opacity: 1, x: 0 }}
              viewport={{ once: true }}
              className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-primary/10 text-sm font-medium text-primary mb-6"
            >
              <Sparkles className="w-4 h-4" />
              Operational Intelligence
            </motion.div>
            <motion.h2 
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ delay: 0.1 }}
              className="text-3xl md:text-5xl font-semibold tracking-tight text-white mb-6"
            >
              Talk to your infrastructure.
            </motion.h2>
            <motion.p 
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ delay: 0.2 }}
              className="text-lg text-muted-foreground mb-8"
            >
              CloudStoreX AI doesn't just answer questions—it analyzes metadata, forecasts capacity, and automatically optimizes your routing policies across clouds.
            </motion.p>
            
            <ul className="space-y-4">
              {["Identify cost anomalies instantly", "Forecast storage capacity for the next quarter", "Automatically generate optimal routing policies", "Understand latency bottlenecks"].map((item, i) => (
                <motion.li 
                  key={i}
                  initial={{ opacity: 0, x: -10 }}
                  whileInView={{ opacity: 1, x: 0 }}
                  viewport={{ once: true }}
                  transition={{ delay: 0.3 + (i * 0.1) }}
                  className="flex items-center gap-3 text-muted-foreground"
                >
                  <CheckCircle2 className="w-5 h-5 text-emerald-500 flex-shrink-0" />
                  <span>{item}</span>
                </motion.li>
              ))}
            </ul>
          </div>

          <div className="flex-1 w-full relative">
            {/* Chat Interface Mockup */}
            <motion.div 
              initial={{ opacity: 0, scale: 0.95 }}
              whileInView={{ opacity: 1, scale: 1 }}
              viewport={{ once: true }}
              transition={{ duration: 0.5, delay: 0.2 }}
              className="glass-panel rounded-2xl overflow-hidden border border-white/10 shadow-2xl"
            >
              <div className="flex items-center gap-2 px-4 py-3 border-b border-white/5 bg-white/[0.02]">
                <div className="w-3 h-3 rounded-full bg-red-500/80" />
                <div className="w-3 h-3 rounded-full bg-amber-500/80" />
                <div className="w-3 h-3 rounded-full bg-emerald-500/80" />
                <span className="ml-2 text-xs font-medium text-muted-foreground">CloudStoreX Intelligence</span>
              </div>
              
              <div className="p-6 flex flex-col gap-6 bg-[#05070A]/50 min-h-[400px]">
                {/* User Message */}
                <div className="flex gap-4 items-start">
                  <div className="w-8 h-8 rounded-full bg-white/10 flex items-center justify-center text-xs font-bold text-white shrink-0">
                    JD
                  </div>
                  <div className="bg-white/5 rounded-2xl rounded-tl-sm px-5 py-3 text-sm text-white">
                    <TypewriterText 
                      text="Which bucket is causing the sudden spike in our AWS bill this week?" 
                      delay={0.5} 
                      onComplete={() => setQuestionDone(true)}
                    />
                  </div>
                </div>

                {/* AI Response */}
                {questionDone && (
                  <motion.div 
                    initial="hidden"
                    animate="visible"
                    variants={{
                      visible: { transition: { staggerChildren: 0.8, delayChildren: 0.5 } }
                    }}
                    className="flex gap-4 items-start"
                  >
                    <motion.div variants={{ hidden: { opacity: 0, scale: 0.8 }, visible: { opacity: 1, scale: 1 } }} className="w-8 h-8 rounded-full bg-primary/20 flex items-center justify-center text-primary shrink-0">
                      <Sparkles className="w-4 h-4" />
                    </motion.div>
                    <div className="flex flex-col gap-4 w-full">
                      <motion.div variants={{ hidden: { opacity: 0, y: 10 }, visible: { opacity: 1, y: 0 } }} className="text-sm text-muted-foreground leading-relaxed">
                        <TypewriterText text="I analyzed your billing metadata. The spike is originating from " delay={0.6} cursor={false} />
                        <motion.span initial={{ opacity: 0 }} animate={{ opacity: 1 }} transition={{ delay: 2.5 }} className="text-white font-mono bg-white/10 px-1 rounded">prod-analytics-raw</motion.span>
                        <TypewriterText text=" in AWS us-east-1." delay={2.6} cursor={false} />
                        <br /><br />
                        <TypewriterText text="This bucket received 4.2TB of new objects in the last 72 hours with a high retrieval rate." delay={3.5} cursor={false} />
                      </motion.div>
                      
                      {/* Embedded Card */}
                      <motion.div variants={{ hidden: { opacity: 0, y: 20 }, visible: { opacity: 1, y: 0 } }} className="bg-white/5 border border-white/10 rounded-xl p-4">
                        <div className="flex items-center justify-between mb-4">
                          <div className="flex items-center gap-2">
                            <BarChart3 className="w-4 h-4 text-emerald-400" />
                            <span className="text-sm font-medium text-white">Cost Anomaly Detected</span>
                          </div>
                          <span className="text-xs text-red-400 font-mono">
                            +<CountUp to={842} prefix="$" delay={4} />.<CountUp to={50} delay={4.5} /> (+<CountUp to={412} suffix="%" delay={4} />)
                          </span>
                        </div>
                        
                        <div className="space-y-2 mb-4">
                          <div className="flex justify-between text-xs text-muted-foreground">
                            <span>Normal baseline</span>
                            <span>$204.00</span>
                          </div>
                          <div className="flex justify-between text-xs text-muted-foreground">
                            <span>Current spend (72h)</span>
                            <span className="text-white">
                              $<CountUp to={1046} delay={4} />.50
                            </span>
                          </div>
                        </div>

                        <button className="w-full py-2 bg-primary/20 hover:bg-primary/30 text-primary text-xs font-medium rounded-lg transition-colors flex items-center justify-center gap-2 group">
                          Apply Auto-Tiering Policy <ArrowUpRight className="w-3 h-3 group-hover:translate-x-0.5 group-hover:-translate-y-0.5 transition-transform" />
                        </button>
                      </motion.div>
                    </div>
                  </motion.div>
                )}
              </div>
            </motion.div>
            
            {/* Background Glow */}
            <div className="absolute -inset-10 bg-primary/20 blur-[100px] rounded-full pointer-events-none opacity-50 z-[-1]" />
          </div>
        </div>
      </div>
    </section>
  );
}
