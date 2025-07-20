"use client";

import { createContext, useContext, useEffect, useState } from "react";

type Theme = "dark" | "light" | "system";

type ThemeProviderProps = {
  children: React.ReactNode;
  defaultTheme?: Theme;
  storageKey?: string;
};

type ThemeProviderState = {
  theme: Theme;
  systemTheme: "dark" | "light";
  actualTheme: "dark" | "light";
  setTheme: (theme: Theme) => void;
};

const initialState: ThemeProviderState = {
  theme: "system",
  systemTheme: "light",
  actualTheme: "light",
  setTheme: () => null,
};

const ThemeProviderContext = createContext<ThemeProviderState>(initialState);

export function ThemeProvider({
  children,
  defaultTheme = "system",
  storageKey = "routrapp-ui-theme",
  ...props
}: ThemeProviderProps) {
  const [theme, setTheme] = useState<Theme>(() => {
    // Try to get theme from localStorage first, fallback to defaultTheme
    if (typeof window !== "undefined") {
      try {
        const stored = localStorage.getItem(storageKey);
        if (
          stored &&
          (stored === "dark" || stored === "light" || stored === "system")
        ) {
          return stored as Theme;
        }
      } catch {
        // Ignore localStorage errors
      }
    }
    return defaultTheme;
  });
  const [systemTheme, setSystemTheme] = useState<"dark" | "light">("light");
  const [actualTheme, setActualTheme] = useState<"dark" | "light">("light");

  useEffect(() => {
    const mediaQuery = window.matchMedia("(prefers-color-scheme: dark)");

    const updateSystemTheme = () => {
      const isDark = mediaQuery.matches;
      setSystemTheme(isDark ? "dark" : "light");
    };

    // Set initial theme
    updateSystemTheme();

    // Listen for changes
    mediaQuery.addEventListener("change", updateSystemTheme);

    return () => {
      mediaQuery.removeEventListener("change", updateSystemTheme);
    };
  }, []);

  useEffect(() => {
    const newActualTheme = theme === "system" ? systemTheme : theme;
    setActualTheme(newActualTheme);

    const root = window.document.documentElement;
    root.classList.remove("light", "dark");
    root.classList.add(newActualTheme);
  }, [theme, systemTheme]);

  // Persist theme changes to localStorage
  useEffect(() => {
    try {
      localStorage.setItem(storageKey, theme);
    } catch {
      // Ignore localStorage errors
    }
  }, [theme, storageKey]);

  const value = {
    theme,
    systemTheme,
    actualTheme,
    setTheme,
  };

  return (
    <ThemeProviderContext.Provider {...props} value={value}>
      {children}
    </ThemeProviderContext.Provider>
  );
}

export const useTheme = () => {
  const context = useContext(ThemeProviderContext);

  if (context === undefined)
    throw new Error("useTheme must be used within a ThemeProvider");

  return context;
};
