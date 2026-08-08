"use client";

import { motion } from "framer-motion";
import { Check, X } from "lucide-react";

const comparisonFeatures = [
  { name: "Multi-Cloud Abstraction", cloudstorex: true, traditional: false },
  { name: "Dynamic Policy Routing", cloudstorex: true, traditional: false },
  { name: "AI Storage Intelligence", cloudstorex: true, traditional: false },
  { name: "Global Metadata Catalog", cloudstorex: true, traditional: false },
  { name: "Active-Active Replication", cloudstorex: true, traditional: "Add-on" },
  { name: "Zero-Trust IAM", cloudstorex: true, traditional: true },
  { name: "Kubernetes Operator", cloudstorex: true, traditional: false },
  { name: "Vendor Lock-in", cloudstorex: false, traditional: true },
];

export function EnterpriseTable() {
  return (
    <section className="py-24 relative bg-[#05070A] border-y border-white/5">
      <div className="container mx-auto px-6 max-w-5xl">
        <div className="text-center mb-16">
          <motion.h2 
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            className="text-3xl md:text-5xl font-semibold tracking-tight text-white mb-6"
          >
            The next generation of <br />
            <span className="text-muted-foreground">storage infrastructure.</span>
          </motion.h2>
        </div>

        <motion.div
          initial={{ opacity: 0, y: 40 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ duration: 0.8 }}
          className="rounded-2xl border border-white/10 bg-[#0D1117] overflow-hidden shadow-2xl"
        >
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="border-b border-white/10 bg-[#05070A]/50">
                <th className="py-6 px-6 font-medium text-muted-foreground w-1/2">Capability</th>
                <th className="py-6 px-6 font-semibold text-white w-1/4">CloudStoreX</th>
                <th className="py-6 px-6 font-medium text-muted-foreground w-1/4">Traditional Storage</th>
              </tr>
            </thead>
            <tbody>
              {comparisonFeatures.map((feature, i) => (
                <tr key={i} className="border-b border-white/5 hover:bg-white/[0.02] transition-colors">
                  <td className="py-5 px-6 text-sm text-white">{feature.name}</td>
                  <td className="py-5 px-6">
                    {feature.cloudstorex === true ? (
                      <Check className="w-5 h-5 text-emerald-400" />
                    ) : feature.cloudstorex === false ? (
                      <X className="w-5 h-5 text-muted-foreground" />
                    ) : (
                      <span className="text-sm text-muted-foreground">{feature.cloudstorex}</span>
                    )}
                  </td>
                  <td className="py-5 px-6">
                    {feature.traditional === true ? (
                      <Check className="w-5 h-5 text-emerald-400/50" />
                    ) : feature.traditional === false ? (
                      <X className="w-5 h-5 text-muted-foreground" />
                    ) : (
                      <span className="text-sm text-muted-foreground">{feature.traditional}</span>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </motion.div>
      </div>
    </section>
  );
}
