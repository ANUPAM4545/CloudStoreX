import { Header } from "@/components/marketing/Header";
import { Footer } from "@/components/marketing/Footer";
import { Hero } from "@/components/marketing/Hero";
import { TrustedTech } from "@/components/marketing/TrustedTech";
import { InteractiveArchitecture } from "@/components/marketing/InteractiveArchitecture";
import { FeatureGrid } from "@/components/marketing/FeatureGrid";
import { AIIntelligence } from "@/components/marketing/AIIntelligence";
import { PolicyEngine } from "@/components/marketing/PolicyEngine";
import { DashboardShowcase } from "@/components/marketing/DashboardShowcase";
import { DeveloperExperience } from "@/components/marketing/DeveloperExperience";
import { Performance } from "@/components/marketing/Performance";
import { EnterpriseTable } from "@/components/marketing/EnterpriseTable";
import { OpenSource } from "@/components/marketing/OpenSource";
import { CTA } from "@/components/marketing/CTA";

export default function Home() {
  return (
    <main className="min-h-screen bg-background text-foreground flex flex-col selection:bg-primary/30 selection:text-white">
      <Header />
      
      <Hero />
      <TrustedTech />
      <InteractiveArchitecture />
      <FeatureGrid />
      <AIIntelligence />
      <PolicyEngine />
      <DashboardShowcase />
      <DeveloperExperience />
      <Performance />
      <EnterpriseTable />
      <OpenSource />
      <CTA />

      <Footer />
    </main>
  );
}
