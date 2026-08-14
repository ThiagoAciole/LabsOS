import { AlertCircle, MoreHorizontal, Play } from "lucide-react";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Progress } from "@/components/ui/progress";
import { AppIcon } from "./AppIcon";
import type { InstalledApp } from "../types";

const statusLabel = {
  installing: "Instalando",
  running: "Em execução",
  stopped: "Parado",
  error: "Precisa de atenção",
} as const;

export function InstalledAppCard({
  app,
  onOpen,
  onOverview,
}: {
  app: InstalledApp;
  onOpen: (app: InstalledApp) => void;
  onOverview: (app: InstalledApp) => void;
}) {
  const running = app.status === "running";

  return (
    <Card
      className={`flex h-full ${
        app.status === "installing"
          ? "border-primary/30"
          : app.status === "error"
            ? "border-destructive/30"
            : ""
      }`}
      data-status={app.status}
    >
      <CardContent className="flex min-h-44 w-full flex-col gap-4 p-4">
        <div className="flex min-w-0 items-start gap-3">
          <AppIcon
            icon={app.icon}
            className="size-11 shrink-0 object-contain"
          />
          <div className="min-w-0 flex-1">
            <div className="flex items-start justify-between gap-2">
              <h2 className="truncate font-medium">{app.name}</h2>
              <Badge
                variant={
                  app.status === "error"
                    ? "destructive"
                    : app.status === "installing"
                      ? "secondary"
                      : "outline"
                }
                className="shrink-0"
              >
                <span
                  className={
                    app.status === "running"
                      ? "mr-1.5 size-1.5 rounded-full bg-emerald-500"
                      : "mr-1.5"
                  }
                >
                  {app.status === "error" ? (
                    <AlertCircle className="size-3" />
                  ) : null}
                </span>
                {statusLabel[app.status]}
              </Badge>
            </div>
            <p className="mt-1 line-clamp-2 text-sm text-muted-foreground">
              {app.description}
            </p>
          </div>
        </div>
        {app.status === "installing" ? (
          <div className="space-y-1.5">
            <div className="flex items-center justify-between text-xs text-muted-foreground">
              <span>Preparando o app</span>
              <span className="tabular-nums">{app.progress}%</span>
            </div>
            <Progress value={app.progress ?? null} />
          </div>
        ) : (
          <div className="mt-auto flex items-center gap-2">
            <Button
              variant={running ? "secondary" : "outline"}
              size="sm"
              className="min-h-11 flex-1"
              disabled={!running}
              onClick={() => onOpen(app)}
            >
              <Play data-icon="inline-start" />
              Abrir
            </Button>
            <Button
              size="icon"
              variant="secondary"
              className="size-11"
              aria-label={`Visão geral de ${app.name}`}
              onClick={() => onOverview(app)}
            >
              <MoreHorizontal />
            </Button>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
