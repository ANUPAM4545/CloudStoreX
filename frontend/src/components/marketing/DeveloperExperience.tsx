"use client";

import { useState } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { Code2, Terminal, CloudLightning, Copy, Check, Globe } from "lucide-react";

const codeExamples = {
  go: {
    label: "Go SDK",
    icon: Code2,
    code: `import "github.com/cloudstorex/sdk-go/v2/cx"

client := cx.NewClient(
    cx.WithAPIKey(os.Getenv("CX_API_KEY")),
    cx.WithRegion("global"),
)

// The Policy Engine determines if this goes to AWS or MinIO
res, err := client.Storage.PutObject(ctx, &cx.PutObjectInput{
    Bucket:      "enterprise-data-lake",
    Key:         "models/weights-v2.h5",
    Body:        fileData,
    ContentType: "application/x-hdf5",
    Tags: map[string]string{
        "compliance": "hipaa",
        "tier":       "intelligent",
    },
})`,
  },
  terraform: {
    label: "Terraform",
    icon: CloudLightning,
    code: `terraform {
  required_providers {
    cloudstorex = {
      source = "cloudstorex/cloudstorex"
    }
  }
}

resource "cloudstorex_routing_policy" "ml_data" {
  name        = "ml-weights-routing"
  description = "Route large ML models to cost-effective storage"

  condition {
    object_size_gt = "1GB"
    tags           = { tier = "intelligent" }
  }

  action {
    primary_provider = "aws-us-east-1"
    fallback_provider = "minio-onprem"
    replication      = true
  }
}`,
  },
  rest: {
    label: "REST API",
    icon: Globe,
    code: `POST /v1/storage/objects HTTP/1.1
Host: api.cloudstorex.corp
Authorization: Bearer cx_prod_99823...
Content-Type: application/json

{
  "bucket": "enterprise-data-lake",
  "key": "models/weights-v2.h5",
  "routing_policy": "intelligent-tier",
  "metadata": {
    "compliance": "hipaa"
  }
}`
  },
  cli: {
    label: "CLI",
    icon: Terminal,
    code: `# Upload a file using the global routing engine
$ cx storage cp ./dataset.csv cx://enterprise-data-lake/

# Ask the AI to optimize routing policies
$ cx ai optimize --target=cost --dry-run
> Analyzing 4.2TB of metadata...
> Recommendation: Route objects > 30 days to AWS Glacier
> Estimated monthly savings: $4,200

# Apply the recommendation automatically
$ cx ai optimize --target=cost --apply`,
  }
};

const highlightCode = (line: string) => {
  if (line.trim().startsWith('//') || line.trim().startsWith('#') || line.trim().startsWith('>')) {
    return <span className="text-muted-foreground">{line}</span>;
  }
  
  const parts = line.split(/("[^"]*")/g);
  
  return (
    <>
      {parts.map((part, i) => {
        if (part.startsWith('"') && part.endsWith('"')) {
          return <span key={i} className="text-emerald-400">{part}</span>;
        }
        
        const keywords = ['import', 'terraform', 'resource', 'POST', 'Host:', 'Authorization:', 'Content-Type:'];
        let html = part.replace(/</g, "&lt;").replace(/>/g, "&gt;");
        
        keywords.forEach(keyword => {
          if (keyword.endsWith(':')) {
            html = html.replace(new RegExp(`${keyword}`, 'g'), `<span class="text-primary">${keyword}</span>`);
          } else {
            html = html.replace(new RegExp(`\\b${keyword}\\b`, 'g'), `<span class="text-primary">${keyword}</span>`);
          }
        });
        
        html = html.replace(/([a-zA-Z0-9_]+)(?=\()/g, '<span class="text-blue-400">$1</span>');

        return <span key={i} dangerouslySetInnerHTML={{ __html: html }} />;
      })}
    </>
  );
};

export function DeveloperExperience() {
  const [activeTab, setActiveTab] = useState<keyof typeof codeExamples>("go");
  const [copied, setCopied] = useState(false);

  const handleCopy = () => {
    navigator.clipboard.writeText(codeExamples[activeTab].code);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <section id="developers" className="py-24 relative bg-[#05070A] border-t border-white/5">
      <div className="container mx-auto px-6">
        <div className="flex flex-col lg:flex-row gap-16 items-center">
          <div className="flex-1 w-full">
            <motion.h2 
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              className="text-3xl md:text-5xl font-semibold tracking-tight text-white mb-6"
            >
              Built for <span className="text-gradient">engineers.</span>
            </motion.h2>
            <motion.p 
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ delay: 0.1 }}
              className="text-lg text-muted-foreground mb-8"
            >
              CloudStoreX provides first-class Terraform providers, a powerful CLI, and strongly typed SDKs for Go, TypeScript, and Python. We abstract the complexity so you can focus on building.
            </motion.p>

            <ul className="space-y-4 text-muted-foreground">
              <li className="flex items-center gap-3">
                <div className="w-1.5 h-1.5 rounded-full bg-primary" /> Terraform & OpenTofu support
              </li>
              <li className="flex items-center gap-3">
                <div className="w-1.5 h-1.5 rounded-full bg-emerald-400" /> Declarative GitOps ready
              </li>
              <li className="flex items-center gap-3">
                <div className="w-1.5 h-1.5 rounded-full bg-white/50" /> Kubernetes Operator
              </li>
            </ul>
          </div>

          <div className="flex-1 w-full max-w-2xl">
            <motion.div 
              initial={{ opacity: 0, scale: 0.95 }}
              whileInView={{ opacity: 1, scale: 1 }}
              viewport={{ once: true }}
              transition={{ delay: 0.2 }}
              className="rounded-xl overflow-hidden border border-white/10 bg-[#0D1117] shadow-2xl"
            >
              {/* Tabs */}
              <div className="flex items-center gap-1 border-b border-white/5 p-2 bg-[#05070A]/50">
                {(Object.entries(codeExamples) as [keyof typeof codeExamples, typeof codeExamples[keyof typeof codeExamples]][]).map(([key, item]) => (
                  <button
                    key={key}
                    onClick={() => setActiveTab(key)}
                    className={`relative flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
                      activeTab === key ? "text-white" : "text-muted-foreground hover:text-white"
                    }`}
                  >
                    {activeTab === key && (
                      <motion.div
                        layoutId="activeTab"
                        className="absolute inset-0 bg-white/10 rounded-lg"
                        transition={{ type: "spring", bounce: 0.2, duration: 0.6 }}
                      />
                    )}
                    <item.icon className="w-4 h-4 relative z-10" />
                    <span className="relative z-10">{item.label}</span>
                  </button>
                ))}
                
                <button 
                  onClick={handleCopy}
                  className="ml-auto p-2 text-muted-foreground hover:text-white transition-colors relative"
                  title="Copy code"
                >
                  <AnimatePresence mode="wait">
                    {copied ? (
                      <motion.div key="check" initial={{ scale: 0 }} animate={{ scale: 1 }} exit={{ scale: 0 }}>
                        <Check className="w-4 h-4 text-emerald-400" />
                      </motion.div>
                    ) : (
                      <motion.div key="copy" initial={{ scale: 0 }} animate={{ scale: 1 }} exit={{ scale: 0 }}>
                        <Copy className="w-4 h-4" />
                      </motion.div>
                    )}
                  </AnimatePresence>
                </button>
              </div>

              {/* Code Area */}
              <div className="p-6 overflow-x-auto min-h-[350px]">
                <AnimatePresence mode="wait">
                  <motion.pre
                    key={activeTab}
                    initial={{ opacity: 0, y: 5 }}
                    animate={{ opacity: 1, y: 0 }}
                    exit={{ opacity: 0, y: -5 }}
                    transition={{ duration: 0.2 }}
                    className="text-sm font-mono text-white/80 leading-relaxed"
                  >
                    <code>
                      {codeExamples[activeTab].code.split('\n').map((line, i) => (
                        <div key={i} className="table-row">
                          <span className="table-cell pr-6 text-white/20 select-none text-right">{i + 1}</span>
                          <span className="table-cell">
                            {highlightCode(line)}
                          </span>
                        </div>
                      ))}
                    </code>
                  </motion.pre>
                </AnimatePresence>
              </div>
            </motion.div>
          </div>
        </div>
      </div>
    </section>
  );
}
