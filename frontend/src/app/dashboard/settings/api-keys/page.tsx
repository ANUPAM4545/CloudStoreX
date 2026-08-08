"use client";

import React, { useState } from "react";
import { Plus, MoreVertical, Copy, Key, Calendar, Trash2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { toast } from "sonner";
import { motion, AnimatePresence } from "framer-motion";

interface ApiKey {
  id: string;
  name: string;
  prefix: string;
  created: string;
  lastUsed: string | null;
  expires: string | null;
}

const MOCK_KEYS: ApiKey[] = [
  {
    id: "key_1",
    name: "Production Worker",
    prefix: "csx_prod_",
    created: "2026-07-01",
    lastUsed: "2 mins ago",
    expires: null,
  },
  {
    id: "key_2",
    name: "Development Env",
    prefix: "csx_dev_",
    created: "2026-07-15",
    lastUsed: "Yesterday",
    expires: "2026-10-15",
  },
];

export default function ApiKeysPage() {
  const [keys, setKeys] = useState(MOCK_KEYS);

  const handleCreate = () => {
    toast("Creating API keys is mocked in this prototype.");
  };

  const handleRevoke = (id: string) => {
    setKeys(keys.filter((k) => k.id !== id));
    toast.success("API Key revoked successfully.");
  };

  const copyToClipboard = () => {
    toast.success("Key ID copied to clipboard.");
  };

  return (
    <div className="space-y-8">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">API Keys</h1>
          <p className="text-sm text-muted-foreground mt-1">
            Manage secret keys used to authenticate with the CloudStoreX API.
          </p>
        </div>
        <Button onClick={handleCreate}>
          <Plus size={16} className="mr-2" />
          Create new key
        </Button>
      </div>

      <div className="border rounded-xl bg-card overflow-hidden">
        {keys.length > 0 ? (
          <table className="w-full text-sm text-left">
            <thead className="bg-muted/50 border-b text-muted-foreground">
              <tr>
                <th className="px-6 py-4 font-medium">Name</th>
                <th className="px-6 py-4 font-medium">Key Prefix</th>
                <th className="px-6 py-4 font-medium">Last Used</th>
                <th className="px-6 py-4 font-medium">Created</th>
                <th className="px-6 py-4 font-medium text-right">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y">
              <AnimatePresence>
                {keys.map((key) => (
                  <motion.tr 
                    key={key.id}
                    initial={{ opacity: 1 }}
                    exit={{ opacity: 0, backgroundColor: "rgba(225, 29, 72, 0.1)" }}
                    className="hover:bg-muted/30 transition-colors group"
                  >
                    <td className="px-6 py-4 font-medium">{key.name}</td>
                    <td className="px-6 py-4">
                      <div className="flex items-center gap-2">
                        <code className="px-2 py-1 bg-muted rounded font-mono text-xs">
                          {key.prefix}••••••••
                        </code>
                        <button 
                          onClick={copyToClipboard}
                          className="text-muted-foreground hover:text-foreground opacity-0 group-hover:opacity-100 transition-opacity"
                        >
                          <Copy size={14} />
                        </button>
                      </div>
                    </td>
                    <td className="px-6 py-4 text-muted-foreground">{key.lastUsed || "Never"}</td>
                    <td className="px-6 py-4 text-muted-foreground">{key.created}</td>
                    <td className="px-6 py-4 text-right">
                      <Button 
                        variant="ghost" 
                        size="icon" 
                        className="text-rose-500 hover:text-rose-600 hover:bg-rose-500/10"
                        onClick={() => handleRevoke(key.id)}
                      >
                        <Trash2 size={16} />
                      </Button>
                    </td>
                  </motion.tr>
                ))}
              </AnimatePresence>
            </tbody>
          </table>
        ) : (
          <div className="p-12 text-center flex flex-col items-center">
            <div className="w-12 h-12 bg-muted rounded-full flex items-center justify-center mb-4">
              <Key size={24} className="text-muted-foreground" />
            </div>
            <h3 className="text-lg font-medium">No API Keys</h3>
            <p className="text-sm text-muted-foreground max-w-sm mt-1 mb-4">
              You haven't created any API keys yet. Create one to start using the CloudStoreX API.
            </p>
            <Button onClick={handleCreate} variant="outline">Create your first key</Button>
          </div>
        )}
      </div>
    </div>
  );
}
