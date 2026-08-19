import { Play, Power, RefreshCw, Trash2 } from "lucide-react";

import { Button } from "@/components/ui/button";
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetFooter,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet";
import { AppIcon } from "./AppIcon";
import type { InstalledApp } from "../types";
import { useAppHealth, useAppMetrics } from "../queries";

const statusLabel = {
  installing: "Instalando",
  updating: "Atualizando",
  running: "Em execução",
  stopped: "Parado",
  error: "Precisa de atenção",
} as const;

export function AppOverviewDrawer({
  app,
  onClose,
  onOpen,
  onRestart,
  onToggle,
  onRemove,
  onUpdate,
  pending,
}: {
  app: InstalledApp | null;
  onClose: () => void;
  onOpen: (app: InstalledApp) => void;
  onRestart: (app: InstalledApp) => void;
  onToggle: (app: InstalledApp) => void;
  onRemove: (app: InstalledApp) => void;
  onUpdate: (app: InstalledApp) => void;
  pending: boolean;
}) {
  const running = app?.status === "running";
  const healthQuery = useAppHealth(app?.id ?? null);
  const metricsQuery = useAppMetrics(app?.id ?? null);

  return (
    <Sheet open={Boolean(app)} onOpenChange={(open) => !open && onClose()}>
      <SheetContent side="right" className="w-full sm:max-w-sm">
        {app && (
          <>
            <SheetHeader className="pr-12">
              <div className="flex items-center gap-3">
                <AppIcon icon={app.icon} name={app.name} className="size-12 object-contain" />
                <div>
                  <SheetTitle>{app.name}</SheetTitle>
                  <SheetDescription>{statusLabel[app.status]}</SheetDescription>
                </div>
              </div>
              <p className="pt-3 text-sm text-muted-foreground">
                {app.description}
              </p>
            </SheetHeader>
            <div className="space-y-4 p-6 pt-5">
              <div className="grid grid-cols-2 gap-3 rounded-lg border p-4 text-sm">
                <div>
                  <span className="block text-xs text-muted-foreground">
                    Status
                  </span>
                  <span className="mt-1 block font-medium">
                    {statusLabel[app.status]}
                  </span>
                </div>
                <div>
                  <span className="block text-xs text-muted-foreground">
                    Versão
                  </span>
                  <span className="mt-1 block font-medium">
                    {app.version ?? "Disponível"}
                  </span>
                </div>
              </div>
              <div className="rounded-lg border p-4 text-sm">
                <span className="block text-xs text-muted-foreground">Health check</span>
                <span className={`mt-1 block font-medium ${healthQuery.data?.healthy ? "text-emerald-600" : "text-muted-foreground"}`}>
                  {healthQuery.isPending ? "Verificando…" : healthQuery.data?.message ?? "Indisponível"}
                </span>
                {healthQuery.data?.dependencies?.length ? <span className="mt-2 block text-xs text-muted-foreground">Dependências: {healthQuery.data.dependencies.join(", ")}</span> : null}
              </div>
              <div className="rounded-lg border p-4 text-sm"><span className="block text-xs text-muted-foreground">Métricas</span>{metricsQuery.data?.available ? <div className="mt-2 grid grid-cols-2 gap-2 text-xs"><span>CPU: {metricsQuery.data.cpuPercent.toFixed(1)}%</span><span>RAM: {formatBytes(metricsQuery.data.memoryUsedBytes)} / {formatBytes(metricsQuery.data.memoryLimitBytes)}</span><span>Rede RX: {formatBytes(metricsQuery.data.networkRXBytes)}</span><span>Rede TX: {formatBytes(metricsQuery.data.networkTXBytes)}</span><span>Leitura: {formatBytes(metricsQuery.data.blockReadBytes)}</span><span>Escrita: {formatBytes(metricsQuery.data.blockWriteBytes)}</span></div> : <span className="mt-1 block text-xs text-muted-foreground">{metricsQuery.isPending ? "Coletando…" : metricsQuery.data?.message ?? "Indisponível"}</span>}</div>
            </div>
            <SheetFooter className="gap-2 border-t">
              <Button disabled={!running} onClick={() => onOpen(app)}>
                <Play data-icon="inline-start" />
                Abrir
              </Button>
              <div className="grid grid-cols-2 gap-2">
                <Button variant="outline" disabled={pending} onClick={() => onRestart(app)}>
                  <RefreshCw data-icon="inline-start" />
                  Reiniciar
                </Button>
                <Button
                  variant="outline"
                  disabled={pending}
                  onClick={() => onToggle(app)}
                >
                  <Power data-icon="inline-start" />
                  {running ? "Parar" : "Iniciar"}
                </Button>
              </div>
              <Button variant="outline" disabled={pending} onClick={() => onRemove(app)}><Trash2 data-icon="inline-start" />Desinstalar</Button>
              {app.updateAvailable && <Button variant="outline" disabled={pending} onClick={() => onUpdate(app)}><RefreshCw data-icon="inline-start" />Atualizar</Button>}
            </SheetFooter>
          </>
        )}
      </SheetContent>
    </Sheet>
  );
}

function formatBytes(value: number) { if (!Number.isFinite(value) || value < 1) return "0 B"; const units = ["B", "KB", "MB", "GB", "TB"]; const index = Math.min(Math.floor(Math.log(value) / Math.log(1024)), units.length - 1); return `${(value / 1024 ** index).toFixed(index ? 1 : 0)} ${units[index]}` }
