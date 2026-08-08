import Link from "next/link";
import { Cloud, GitBranch, MessageCircle, DiscIcon as Discord } from "lucide-react";

export function Footer() {
  return (
    <footer className="border-t border-white/10 bg-[#05070A] py-16 px-6 md:px-12 lg:px-24">
      <div className="max-w-7xl mx-auto grid grid-cols-2 md:grid-cols-4 lg:grid-cols-5 gap-12">
        <div className="col-span-2 lg:col-span-1 flex flex-col gap-4">
          <div className="flex items-center gap-2">
            <div className="flex items-center justify-center w-8 h-8 rounded-lg bg-primary/20 text-primary">
              <Cloud className="w-5 h-5" />
            </div>
            <span className="text-xl font-medium tracking-tight text-white">
              CloudStoreX
            </span>
          </div>
          <p className="text-sm text-muted-foreground mt-2">
            The enterprise AI-native multi-cloud storage control plane.
          </p>
          <div className="flex items-center gap-4 mt-4 text-muted-foreground">
            <Link href="#" className="hover:text-white transition-colors">
              <GitBranch className="w-5 h-5" />
            </Link>
            <Link href="#" className="hover:text-white transition-colors">
              <MessageCircle className="w-5 h-5" />
            </Link>
            <Link href="#" className="hover:text-white transition-colors">
              <Discord className="w-5 h-5" />
            </Link>
          </div>
        </div>

        <div className="flex flex-col gap-4">
          <h3 className="font-medium text-white">Product</h3>
          <Link href="#" className="text-sm text-muted-foreground hover:text-white transition-colors">Storage Routing</Link>
          <Link href="#" className="text-sm text-muted-foreground hover:text-white transition-colors">Policy Engine</Link>
          <Link href="#" className="text-sm text-muted-foreground hover:text-white transition-colors">AI Intelligence</Link>
          <Link href="#" className="text-sm text-muted-foreground hover:text-white transition-colors">Disaster Recovery</Link>
          <Link href="#" className="text-sm text-muted-foreground hover:text-white transition-colors">Enterprise Security</Link>
        </div>

        <div className="flex flex-col gap-4">
          <h3 className="font-medium text-white">Resources</h3>
          <Link href="#" className="text-sm text-muted-foreground hover:text-white transition-colors">Documentation</Link>
          <Link href="#" className="text-sm text-muted-foreground hover:text-white transition-colors">API Reference</Link>
          <Link href="#" className="text-sm text-muted-foreground hover:text-white transition-colors">Go SDK</Link>
          <Link href="#" className="text-sm text-muted-foreground hover:text-white transition-colors">Terraform Provider</Link>
          <Link href="#" className="text-sm text-muted-foreground hover:text-white transition-colors">Open Source</Link>
        </div>

        <div className="flex flex-col gap-4">
          <h3 className="font-medium text-white">Company</h3>
          <Link href="#" className="text-sm text-muted-foreground hover:text-white transition-colors">About</Link>
          <Link href="#" className="text-sm text-muted-foreground hover:text-white transition-colors">Blog</Link>
          <Link href="#" className="text-sm text-muted-foreground hover:text-white transition-colors">Careers</Link>
          <Link href="#" className="text-sm text-muted-foreground hover:text-white transition-colors">Contact</Link>
        </div>

        <div className="flex flex-col gap-4">
          <h3 className="font-medium text-white">Legal</h3>
          <Link href="#" className="text-sm text-muted-foreground hover:text-white transition-colors">Privacy Policy</Link>
          <Link href="#" className="text-sm text-muted-foreground hover:text-white transition-colors">Terms of Service</Link>
          <Link href="#" className="text-sm text-muted-foreground hover:text-white transition-colors">Security</Link>
          <Link href="#" className="text-sm text-muted-foreground hover:text-white transition-colors">Status</Link>
        </div>
      </div>
      
      <div className="max-w-7xl mx-auto mt-16 pt-8 border-t border-white/5 flex flex-col md:flex-row items-center justify-between gap-4 text-sm text-muted-foreground">
        <p>© 2026 CloudStoreX Inc. All rights reserved.</p>
        <div className="flex items-center gap-2">
          <span className="relative flex h-2 w-2">
            <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
            <span className="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"></span>
          </span>
          All systems operational
        </div>
      </div>
    </footer>
  );
}
