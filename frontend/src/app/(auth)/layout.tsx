import { ReactNode } from "react";
import { ThemeProvider } from "@/components/theme-provider";

interface AuthLayoutProps {
  children: ReactNode;
}

export default function AuthLayout({ children }: AuthLayoutProps) {
  return (
    <ThemeProvider defaultTheme="system">
      <div className="min-h-screen bg-background">{children}</div>
    </ThemeProvider>
  );
}
