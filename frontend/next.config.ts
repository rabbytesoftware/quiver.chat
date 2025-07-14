import type { NextConfig } from "next";

const nextConfig: NextConfig = {
	output: 'export',
	distDir: 'out',
	trailingSlash: true,
	skipTrailingSlashRedirect: true,
	images: {
		unoptimized: true
	}
};

export default nextConfig;
