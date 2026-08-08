"use client";

import * as React from "react";
import { motion, AnimatePresence } from "framer-motion";
import { cn } from "@/lib/utils";
import { Eye, EyeOff, CheckCircle2, AlertCircle } from "lucide-react";

export interface FloatingInputProps
  extends React.InputHTMLAttributes<HTMLInputElement> {
  label: string;
  error?: string;
  success?: boolean;
}

const FloatingInput = React.forwardRef<HTMLInputElement, FloatingInputProps>(
  ({ className, type, label, error, success, id, ...props }, ref) => {
    const [isFocused, setIsFocused] = React.useState(false);
    const [showPassword, setShowPassword] = React.useState(false);
    const [hasValue, setHasValue] = React.useState(!!props.value || !!props.defaultValue);

    const inputType = type === "password" && showPassword ? "text" : type;
    const inputId = id || `floating-input-${label.replace(/\s+/g, '-').toLowerCase()}`;

    // Handle detecting autofill (a bit tricky in React, but listening to change helps)
    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
      setHasValue(e.target.value.length > 0);
      if (props.onChange) props.onChange(e);
    };

    const isFloating = isFocused || hasValue || props.value;

    return (
      <div className="relative w-full group">
        <div className="relative flex items-center">
          <input
            id={inputId}
            type={inputType}
            className={cn(
              "peer flex h-14 w-full rounded-xl border border-input bg-background/50 px-4 pt-4 pb-1 text-sm transition-all outline-none",
              "hover:border-primary/50 hover:bg-background/80",
              "focus:border-primary focus:bg-background focus:ring-1 focus:ring-primary/20",
              "disabled:cursor-not-allowed disabled:opacity-50",
              error && "border-rose-500/50 hover:border-rose-500 focus:border-rose-500 focus:ring-rose-500/20",
              success && !error && "border-emerald-500/50 hover:border-emerald-500 focus:border-emerald-500 focus:ring-emerald-500/20",
              type === "password" && "pr-10",
              className
            )}
            ref={ref}
            onFocus={(e) => {
              setIsFocused(true);
              if (props.onFocus) props.onFocus(e);
            }}
            onBlur={(e) => {
              setIsFocused(false);
              setHasValue(e.target.value.length > 0);
              if (props.onBlur) props.onBlur(e);
            }}
            onChange={handleChange}
            {...props}
          />
          
          <label
            htmlFor={inputId}
            className={cn(
              "absolute left-4 cursor-text transition-all duration-200 pointer-events-none select-none text-muted-foreground",
              isFloating 
                ? "top-1.5 text-[11px] font-medium opacity-70" 
                : "top-4 text-sm"
            )}
          >
            {label}
          </label>

          {/* Status Icons */}
          <div className="absolute right-3 flex items-center gap-2">
            <AnimatePresence>
              {error && (
                <motion.div
                  initial={{ opacity: 0, scale: 0.5 }}
                  animate={{ opacity: 1, scale: 1 }}
                  exit={{ opacity: 0, scale: 0.5 }}
                >
                  <AlertCircle className="h-4 w-4 text-rose-500" />
                </motion.div>
              )}
              {success && !error && (
                <motion.div
                  initial={{ opacity: 0, scale: 0.5 }}
                  animate={{ opacity: 1, scale: 1 }}
                  exit={{ opacity: 0, scale: 0.5 }}
                >
                  <CheckCircle2 className="h-4 w-4 text-emerald-500" />
                </motion.div>
              )}
            </AnimatePresence>

            {type === "password" && (
              <button
                type="button"
                onClick={() => setShowPassword(!showPassword)}
                className="text-muted-foreground hover:text-foreground transition-colors p-1 rounded-md focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary/50"
                tabIndex={-1}
                aria-label={showPassword ? "Hide password" : "Show password"}
              >
                <AnimatePresence mode="wait">
                  {showPassword ? (
                    <motion.div
                      key="hide"
                      initial={{ opacity: 0, rotate: -45 }}
                      animate={{ opacity: 1, rotate: 0 }}
                      exit={{ opacity: 0, rotate: 45 }}
                      transition={{ duration: 0.15 }}
                    >
                      <EyeOff size={16} />
                    </motion.div>
                  ) : (
                    <motion.div
                      key="show"
                      initial={{ opacity: 0, rotate: 45 }}
                      animate={{ opacity: 1, rotate: 0 }}
                      exit={{ opacity: 0, rotate: -45 }}
                      transition={{ duration: 0.15 }}
                    >
                      <Eye size={16} />
                    </motion.div>
                  )}
                </AnimatePresence>
              </button>
            )}
          </div>
        </div>

        {/* Focus Glow Effect */}
        <AnimatePresence>
          {isFocused && !error && (
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              className="absolute -inset-[1px] -z-10 rounded-xl bg-gradient-to-r from-primary/30 to-blue-500/30 blur-sm transition-all"
            />
          )}
          {isFocused && error && (
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              className="absolute -inset-[1px] -z-10 rounded-xl bg-rose-500/30 blur-sm transition-all"
            />
          )}
        </AnimatePresence>
      </div>
    );
  }
);

FloatingInput.displayName = "FloatingInput";

export { FloatingInput };
