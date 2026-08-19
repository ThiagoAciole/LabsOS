import { AppWindow, ArrowRight } from "lucide-react";
import { Button } from "@/components/ui/button";
import { navigate } from "@/app/navigation";
import type { AppDTO } from "@/api/apps";
import { AppCard } from "./AppCard";
import { useAppMutation } from "@/features/Apps/mutations";
export function ActiveApps({ apps }: { apps: AppDTO[] }) {
  const mutation = useAppMutation();
  const activeApps = apps
    .filter((app) => app.installed && app.status === "running")
    .slice(0, 4);
  return (
    <section className="min-w-0 space-y-3">
      <header className="flex items-center gap-3">
        <AppWindow className="size-4 text-primary" />
        <div>
          <h2 className="text-sm font-medium">Apps principais</h2>
          <p className="text-xs text-muted-foreground">
            {activeApps.length} em execução
          </p>
        </div>
        <Button
          variant="ghost"
          size="sm"
          className="ml-auto min-h-11"
          onClick={() => navigate("/apps")}
        >
          Ver todos <ArrowRight />
        </Button>
      </header>
      <div className="grid gap-3 sm:grid-cols-2">
        {activeApps.map((app) => (
          <AppCard key={app.id} app={app} onAction={(action) => { void mutation.mutateAsync({ id: app.id, action }).catch(() => undefined) }} />
        ))}
      </div>
    </section>
  );
}
