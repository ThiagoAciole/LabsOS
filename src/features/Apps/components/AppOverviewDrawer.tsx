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

const statusLabel = {
  installing: "Instalando",
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
  pending,
}: {
  app: InstalledApp | null;
  onClose: () => void;
  onOpen: (app: InstalledApp) => void;
  onRestart: (app: InstalledApp) => void;
  onToggle: (app: InstalledApp) => void;
  onRemove: (app: InstalledApp) => void;
  pending: boolean;
}) {
  const running = app?.status === "running";

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
            </SheetFooter>
          </>
        )}
      </SheetContent>
    </Sheet>
  );
}
