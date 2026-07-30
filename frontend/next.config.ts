import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: "standalone",
  eslint: {
    ignoreDuringBuilds: true, // We handle linting in a separate CI step
  }
};

export default nextConfig;
