"use client";

import { Link2, Moon, Sun } from "lucide-react";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { useTheme } from "@/components/theme-provider";

export function AppHeader() {
  const { theme, setTheme } = useTheme();
  const isDark = theme === "dark";

  return (
    <header className="flex h-16 shrink-0 items-center justify-end gap-2 px-4 md:px-6">
      <Button
        variant="outline"
        size="icon-sm"
        className="rounded-md bg-card/60"
        aria-label="Alternar tema"
        onClick={() => setTheme(isDark ? "light" : "dark")}
      >
        {isDark ? <Sun /> : <Moon />}
      </Button>
      <Button
        variant="outline"
        size="sm"
        className="gap-2 rounded-md bg-card/60"
      >
        <Link2 />
        <span className="hidden sm:inline">Acesso local</span>
      </Button>
      <Badge
        variant="outline"
        className="h-9 gap-2 rounded-md px-3 text-muted-foreground"
      >
        <span className="size-2 rounded-full bg-emerald-500" />
        <span className="hidden sm:inline">Servidor conectado</span>
      </Badge>
    </header>
  );
}
