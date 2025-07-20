import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import "@/styles/globals.css";
import { AppProviders } from "@/providers/app-providers";
import { ClientAuthWrapper } from "@/components/auth";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "RoutrApp - Route Management System",
  description:
    "Multi-tenant route optimization system for utility and trade companies",
};

// Define route-specific permissions
const routePermissions: Record<string, string[]> = {
  "/admin": ["organizations.*"],
  "/admin/users": ["users.*"],
  "/admin/technicians": ["technicians.*"],
  "/admin/routes": ["routes.*"],
  "/technician": ["routes.read"],
  "/technician/routes": ["routes.read"],
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased`}
      >
        <AppProviders>
          <ClientAuthWrapper routePermissions={routePermissions}>
            {children}
          </ClientAuthWrapper>
        </AppProviders>
      </body>
    </html>
  );
}
