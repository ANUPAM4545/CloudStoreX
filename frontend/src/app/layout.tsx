import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import { Providers } from "@/components/providers";
import { SmoothScroll } from "@/components/smooth-scroll";
import { TooltipProvider } from "@/components/ui/tooltip";

import { Toaster } from "sonner";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "CloudStoreX",
  description: "One API. Any Cloud. Intelligent Storage Platform.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body className={`${inter.className} antialiased`}>
        <Providers
          attribute="class"
          defaultTheme="dark"
          forcedTheme="dark"
          disableTransitionOnChange
        >
          <TooltipProvider>
            <SmoothScroll>
              {children}
              <Toaster position="bottom-right" richColors />
            </SmoothScroll>
          </TooltipProvider>
        </Providers>
      </body>
    </html>
  );
}
