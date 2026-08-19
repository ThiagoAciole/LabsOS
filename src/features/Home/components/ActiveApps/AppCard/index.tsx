import { MoreHorizontal } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import type { DashboardApp } from "../../../types";
import { AppIcon } from "@/features/Apps/components/AppIcon";
export function AppCard({ app, onAction }: { app: DashboardApp; onAction: (action: "restart" | "stop") => void }) {
  return (
    <Card>
      <CardContent className="flex items-center gap-3 p-3">
        <AppIcon icon={app.icon} name={app.name} className="size-10 shrink-0 object-contain" />
        <div className="min-w-0 flex-1">
          <h3 className="truncate text-sm font-medium">{app.name}</h3>
          <span className="mt-1 flex items-center gap-1.5 text-xs text-muted-foreground">
            <span className="size-1.5 rounded-full bg-emerald-500" />
            Em execução
          </span>
        </div>
        <div className="flex shrink-0 items-center gap-1">
          <Button
            variant="secondary"
            size="sm"
            className="min-h-11"
            disabled={!app.url}
            onClick={() => app.url && window.open(app.url, "_blank", "noopener,noreferrer")}
            aria-label={`Abrir ${app.name}`}
          >
            Abrir
          </Button>
          <DropdownMenu>
            <DropdownMenuTrigger
              render={
                <Button
                  size="icon"
                  variant="secondary"
                  className="size-11"
                  aria-label={`Mais opções de ${app.name}`}
                />
              }
            >
              <MoreHorizontal />
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuItem onClick={() => onAction("restart")}>Reiniciar</DropdownMenuItem>
              <DropdownMenuItem variant="destructive" onClick={() => onAction("stop")}>
                Parar
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </CardContent>
    </Card>
  );
}
