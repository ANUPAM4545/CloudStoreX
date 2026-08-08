"use client";

import React from "react";
import { motion, AnimatePresence } from "framer-motion";
import { Check, X } from "lucide-react";

interface PasswordStrengthProps {
  password?: string;
  showChecklist?: boolean;
}

export function PasswordStrength({ password = "", showChecklist = true }: PasswordStrengthProps) {
  const reqs = {
    length: password.length >= 12,
    upper: /[A-Z]/.test(password),
    lower: /[a-z]/.test(password),
    number: /[0-9]/.test(password),
    special: /[^A-Za-z0-9]/.test(password),
  };

  const strengthScore = Object.values(reqs).filter(Boolean).length;
  
  let strengthLabel = "Weak";
  let strengthColor = "bg-muted";
  
  if (password.length > 0) {
    if (strengthScore <= 2) {
      strengthLabel = "Weak";
      strengthColor = "bg-rose-500";
    } else if (strengthScore === 3) {
      strengthLabel = "Fair";
      strengthColor = "bg-amber-500";
    } else if (strengthScore === 4) {
      strengthLabel = "Strong";
      strengthColor = "bg-emerald-400";
    } else if (strengthScore === 5) {
      strengthLabel = "Excellent";
      strengthColor = "bg-emerald-500";
    }
  } else {
    strengthLabel = "";
  }

  const checklist = [
    { key: "length", label: "At least 12 characters", met: reqs.length },
    { key: "upper", label: "Contains uppercase letter", met: reqs.upper },
    { key: "lower", label: "Contains lowercase letter", met: reqs.lower },
    { key: "number", label: "Contains number", met: reqs.number },
    { key: "special", label: "Contains special character", met: reqs.special },
  ];

  return (
    <div className="space-y-3 mt-3">
      {/* Strength Meter */}
      <div className="space-y-1">
        <div className="flex justify-between items-center text-xs font-medium">
          <span className="text-muted-foreground">Password strength</span>
          <span 
            className={
              strengthLabel === "Weak" ? "text-rose-500" :
              strengthLabel === "Fair" ? "text-amber-500" :
              strengthLabel === "Strong" ? "text-emerald-400" :
              strengthLabel === "Excellent" ? "text-emerald-500" : "text-muted-foreground"
            }
          >
            {strengthLabel}
          </span>
        </div>
        <div className="flex gap-1 h-1.5 w-full">
          {[1, 2, 3, 4].map((step) => {
            // Determine segment color dynamically
            let isActive = false;
            if (strengthScore <= 2 && step === 1 && password.length > 0) isActive = true;
            if (strengthScore === 3 && step <= 2) isActive = true;
            if (strengthScore === 4 && step <= 3) isActive = true;
            if (strengthScore === 5 && step <= 4) isActive = true;

            return (
              <motion.div
                key={step}
                className={`flex-1 rounded-full ${isActive ? strengthColor : "bg-muted"}`}
                initial={false}
                animate={{ 
                  backgroundColor: isActive ? `var(--${strengthColor.replace('bg-', '')})` : "hsl(var(--muted))" 
                }}
                transition={{ duration: 0.3 }}
              />
            );
          })}
        </div>
      </div>

      {/* Live Checklist */}
      <AnimatePresence>
        {showChecklist && (
          <motion.div 
            initial={{ opacity: 0, height: 0 }}
            animate={{ opacity: 1, height: "auto" }}
            exit={{ opacity: 0, height: 0 }}
            className="space-y-1.5 overflow-hidden"
            aria-live="polite"
          >
            {checklist.map((item) => (
              <motion.div 
                key={item.key} 
                className={`flex items-center gap-2 text-xs transition-colors duration-300 ${item.met ? "text-emerald-500" : "text-muted-foreground"}`}
                layout
              >
                <div className="relative w-3.5 h-3.5 flex items-center justify-center">
                  <AnimatePresence mode="wait">
                    {item.met ? (
                      <motion.div
                        key="met"
                        initial={{ scale: 0, opacity: 0 }}
                        animate={{ scale: 1, opacity: 1 }}
                        exit={{ scale: 0, opacity: 0 }}
                        transition={{ duration: 0.2 }}
                      >
                        <Check size={14} />
                      </motion.div>
                    ) : (
                      <motion.div
                        key="unmet"
                        initial={{ scale: 0, opacity: 0 }}
                        animate={{ scale: 1, opacity: 1 }}
                        exit={{ scale: 0, opacity: 0 }}
                        transition={{ duration: 0.2 }}
                        className="w-1 h-1 rounded-full bg-muted-foreground/40"
                      />
                    )}
                  </AnimatePresence>
                </div>
                <span>{item.label}</span>
              </motion.div>
            ))}
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}
