import { useMemo, useState } from "react";
import { Boxes, PackagePlus, Search, Store } from "lucide-react";

import { Button } from "@/components/ui/button";
import {
  Empty,
  EmptyContent,
  EmptyDescription,
  EmptyHeader,
  EmptyMedia,
  EmptyTitle,
} from "@/components/ui/empty";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { AppStoreDialog } from "./components/AppStoreDialog";
import { AppOverviewDrawer } from "./components/AppOverviewDrawer";
import { InstalledAppCard } from "./components/InstalledAppCard";
import type { AppStatus, InstalledApp, StoreApp } from "./types";
import { useAppsData, useCatalogData } from "./queries";
import { useAppMutation } from "./mutations";

const noApps: InstalledApp[] = [];

export function AppsPage() {
  const appsQuery = useAppsData();
  const catalogQuery = useCatalogData();
  const mutation = useAppMutation();
  const apps = appsQuery.data ?? noApps;
  const [storeOpen, setStoreOpen] = useState(false);
  const [selectedAppId, setSelectedAppId] = useState<string | null>(null);
  const selectedApp = apps.find((app) => app.id === selectedAppId) ?? null;
  const [query, setQuery] = useState("");
  const [status, setStatus] = useState<"all" | AppStatus>("all");
  const [notice, setNotice] = useState(
    "Conectado à Labs API.",
  );

  const visibleApps = useMemo(
    () =>
      apps.filter(
        (app) =>
          (status === "all" || app.status === status) &&
          app.name.toLocaleLowerCase().includes(query.toLocaleLowerCase()),
      ),
    [apps, query, status],
  );

  async function act(app: Pick<InstalledApp, "id" | "name">, action: "install" | "start" | "stop" | "restart" | "remove") {
    try {
      const job = await mutation.mutateAsync({ id: app.id, action });
      setNotice(job.message);
      if (action === "install") setStoreOpen(false);
      if (action === "remove") setSelectedAppId(null);
    } catch {
      setNotice("Não foi possível concluir esta ação. Tente novamente.");
    }
  }

  function openApp(app: InstalledApp) {
    if (!app.url) {
      setNotice(`${app.name} não informou uma URL de acesso.`)
      return
    }
    window.open(app.url, "_blank", "noopener,noreferrer")
  }

  if (appsQuery.isPending) return <section className="flex flex-1 items-center justify-center p-8" role="status">Carregando apps…</section>;
  if (appsQuery.isError) return <section className="flex flex-1 flex-col items-center justify-center gap-3 p-8" role="alert"><p>Não foi possível carregar os apps.</p><Button variant="outline" onClick={() => void appsQuery.refetch()}>Tentar novamente</Button></section>;

  return (
    <section className="flex min-w-0 flex-1 flex-col gap-6 p-4 md:p-8">
      <header className="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
        <div>
          <h1 className="text-2xl font-medium md:text-3xl">Apps</h1>
          <p className="mt-1 text-sm text-muted-foreground">
            Seus aplicativos instalados e prontos para usar.
          </p>
        </div>
        <Button onClick={() => setStoreOpen(true)}>
          <Store data-icon="inline-start" />
          Abrir loja
        </Button>
      </header>
      <div className="flex flex-col gap-3 border-y border-border/60 py-4 sm:flex-row sm:items-center">
        <div className="relative min-w-0 flex-1 sm:max-w-md">
          <Search className="pointer-events-none absolute top-1/2 left-3 size-4 -translate-y-1/2 text-muted-foreground" />
          <Input
            value={query}
            onChange={(event) => setQuery(event.target.value)}
            className="pl-9"
            placeholder="Buscar apps instalados..."
          />
        </div>
        <Select
          value={status}
          onValueChange={(value) =>
            value && setStatus(value as "all" | AppStatus)
          }
        >
          <SelectTrigger className="w-full sm:w-44">
            <SelectValue>Todos os estados</SelectValue>
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">Todos os estados</SelectItem>
            <SelectItem value="installing">Instalando</SelectItem>
            <SelectItem value="running">Em execução</SelectItem>
            <SelectItem value="stopped">Parados</SelectItem>
            <SelectItem value="error">Precisa de atenção</SelectItem>
          </SelectContent>
        </Select>
        <p className="text-xs text-muted-foreground sm:ml-auto">
          {apps.length} instalados
        </p>
      </div>
      <p role="status" className="text-sm text-muted-foreground">
        {notice}
      </p>
      {apps.length === 0 ? (
        <Empty className="min-h-72 border-0">
          <EmptyHeader>
            <EmptyMedia variant="icon">
              <Boxes />
            </EmptyMedia>
            <EmptyTitle>Nenhum app instalado</EmptyTitle>
            <EmptyDescription>
              Adicione aplicativos pela App Store para começar.
            </EmptyDescription>
          </EmptyHeader>
          <EmptyContent>
            <Button variant="outline" onClick={() => setStoreOpen(true)}>
              <PackagePlus data-icon="inline-start" />
              Abrir App Store
            </Button>
          </EmptyContent>
        </Empty>
      ) : visibleApps.length === 0 ? (
        <Empty className="min-h-56 border-0">
          <EmptyHeader>
            <EmptyMedia variant="icon">
              <Search />
            </EmptyMedia>
            <EmptyTitle>Nenhum app encontrado</EmptyTitle>
            <EmptyDescription>
              Ajuste a busca ou o filtro de estado.
            </EmptyDescription>
          </EmptyHeader>
        </Empty>
      ) : (
        <div className="grid grid-cols-[repeat(auto-fit,minmax(min(100%,17rem),1fr))] gap-4">
          {visibleApps.map((app) => (
            <InstalledAppCard
              key={app.id}
              app={app}
              onOpen={openApp}
              onOverview={(app) => setSelectedAppId(app.id)}
            />
          ))}
        </div>
      )}
      <AppStoreDialog
        open={storeOpen}
        onOpenChange={setStoreOpen}
        apps={catalogQuery.data ?? []}
        installedApps={apps}
        onInstall={(app: StoreApp) => void act(app, "install")}
      />
      <AppOverviewDrawer
        app={selectedApp}
        onClose={() => setSelectedAppId(null)}
        onOpen={openApp}
        onRestart={(app) => void act(app, "restart")}
        onToggle={(app) => void act(app, app.status === "running" ? "stop" : "start")}
        onRemove={(app) => void act(app, "remove")}
        pending={mutation.isPending}
      />
    </section>
  );
}

export type { AppCategory, AppSource } from "./types";
