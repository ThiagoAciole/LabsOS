import { CheckCircle2, CircleAlert, RefreshCw } from "lucide-react";

import { Card, CardContent } from "@/components/ui/card";
import type { SystemHealthDTO } from "@/api/system";

export function SystemHealth({ health, runningApps }: { health: SystemHealthDTO; runningApps: number }) {
  const ready = health.status === "healthy";
  const Icon = ready ? CheckCircle2 : CircleAlert;

  return (
    <section aria-labelledby="system-health-title">
      <Card className="border-primary/20 bg-primary/5">
        <CardContent className="flex flex-col gap-4 p-4 sm:flex-row sm:items-center sm:justify-between sm:p-5">
          <div className="flex items-start gap-3">
            <span className="flex size-10 shrink-0 items-center justify-center rounded-lg bg-primary/10 text-primary">
              <Icon className="size-5" />
            </span>
            <div>
              <h2 id="system-health-title" className="font-medium">
                {ready ? "Tudo funcionando normalmente" : "O servidor precisa de atenção"}
              </h2>
              <p className="mt-1 text-sm text-muted-foreground">
                {runningApps} {runningApps === 1 ? "app em execução" : "apps em execução"} · modo {health.mode}
              </p>
            </div>
          </div>
          <span className="inline-flex items-center gap-1.5 text-xs text-muted-foreground">
            <RefreshCw className="size-3.5" />
            Atualizado agora
          </span>
        </CardContent>
      </Card>
    </section>
  );
}
