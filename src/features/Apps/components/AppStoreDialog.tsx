import { useMemo, useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { ArrowLeft, Check, PackagePlus, Search } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { AppIcon } from "./AppIcon";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import type { AppCategory, InstalledApp, StoreApp } from "../types";
import { importDeclarativeApp } from "@/api/apps";

export function AppStoreDialog({
  open,
  onOpenChange,
  apps,
  installedApps,
  onInstall,
}: {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  apps: StoreApp[];
  installedApps: InstalledApp[];
  onInstall: (app: StoreApp) => void;
}) {
  const [query, setQuery] = useState("");
  const [category, setCategory] = useState<AppCategory | "all">("all");
  const [source, setSource] = useState("all");
  const [selectedApp, setSelectedApp] = useState<StoreApp | null>(null);
  const [importOpen, setImportOpen] = useState(false);
  const [importID, setImportID] = useState("");
  const [importName, setImportName] = useState("");
  const [importCompose, setImportCompose] = useState("services:\n  app:\n    image: example/app");
  const client = useQueryClient();
  const importMutation = useMutation({ mutationFn: importDeclarativeApp, onSuccess: () => { setImportOpen(false); setImportID(""); setImportName(""); void client.invalidateQueries({ queryKey: ["catalog"] }); } });
  const installedById = new Map(installedApps.map((app) => [app.id, app]));
  const visibleApps = useMemo(
    () =>
      apps.filter(
        (app) =>
          (category === "all" || app.category === category) &&
          (source === "all" || app.source === source) &&
          app.name.toLocaleLowerCase().includes(query.toLocaleLowerCase()),
      ),
    [apps, category, query, source],
  );

  function closeStore(nextOpen: boolean) {
    if (!nextOpen) setSelectedApp(null);
    onOpenChange(nextOpen);
  }

  function installSelectedApp() {
    if (!selectedApp) return;
    onInstall(selectedApp);
    setSelectedApp(null);
  }

  return (
    <Dialog open={open} onOpenChange={closeStore}>
      <DialogContent className="flex h-[min(88vh,56rem)] w-[calc(100vw-2rem)] max-w-6xl flex-col gap-0 overflow-hidden p-0 sm:max-w-6xl">
        <DialogHeader className="px-6 pt-5 sm:px-8 sm:pt-6">
          <DialogTitle className="text-xl">
            {selectedApp ? selectedApp.name : "App Store"}
          </DialogTitle>
          <DialogDescription>
            {selectedApp
              ? "Detalhes para decidir antes de instalar."
              : "Adicione novos aplicativos ao seu servidor."}
          </DialogDescription>
        </DialogHeader>
        {selectedApp ? (
          <AppDetails
            app={selectedApp}
            installed={installedById.get(selectedApp.id)}
            onBack={() => setSelectedApp(null)}
            onInstall={installSelectedApp}
          />
        ) : (
          <>
            <div className="border-b border-border/60 px-6 py-5 sm:px-8 sm:py-6">
              <div className="flex flex-col gap-3 sm:flex-row sm:items-center">
                <Select
                  value={category}
                  onValueChange={(value) =>
                    value && setCategory(value as AppCategory | "all")
                  }
                >
                  <SelectTrigger className="w-full sm:w-36">
                    <SelectValue>Todas</SelectValue>
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">Todas</SelectItem>
                    <SelectItem value="media">Mídia</SelectItem>
                    <SelectItem value="storage">Arquivos</SelectItem>
                    <SelectItem value="automation">Automação</SelectItem>
                    <SelectItem value="network">Rede</SelectItem>
                  </SelectContent>
                </Select>
                <Select value={source} onValueChange={(value) => value && setSource(value)}>
                  <SelectTrigger className="w-full sm:w-36">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">Todas as fontes</SelectItem>
                    <SelectItem value="BigBearCasaOS">BigBear CasaOS</SelectItem>
                    <SelectItem value="CasaOS/ZimaOS AppStore">CasaOS / ZimaOS</SelectItem>
                    <SelectItem value="labs">LabsOS</SelectItem>
                  </SelectContent>
                </Select>
                <div className="relative min-w-0 flex-1">
                  <Search className="pointer-events-none absolute top-1/2 left-3 size-4 -translate-y-1/2 text-muted-foreground" />
                  <Input
                    value={query}
                    onChange={(event) => setQuery(event.target.value)}
                    className="pl-9"
                    placeholder="Buscar um app..."
                  />
                </div>
                <Button variant="secondary" onClick={() => setImportOpen((value) => !value)}>
                  <PackagePlus data-icon="inline-start" />
                  Importar App
                </Button>
              </div>
              {importOpen && <form className="mt-4 grid gap-2 rounded-lg border bg-muted/20 p-3 sm:grid-cols-2" onSubmit={(event) => { event.preventDefault(); importMutation.mutate({ id: importID, name: importName, compose: importCompose }); }}>
                <Input required value={importID} onChange={(event) => setImportID(event.target.value)} placeholder="ID (ex.: minha-app)" aria-label="ID do App" />
                <Input required value={importName} onChange={(event) => setImportName(event.target.value)} placeholder="Nome do App" aria-label="Nome do App" />
                <textarea required value={importCompose} onChange={(event) => setImportCompose(event.target.value)} className="min-h-32 rounded-md border bg-background p-2 font-mono text-xs sm:col-span-2" aria-label="Manifesto Docker Compose" />
                <div className="flex items-center justify-between gap-3 sm:col-span-2"><span className="text-xs text-muted-foreground">O manifesto será apenas registrado. A instalação exige uma ação protegida.</span><Button type="submit" disabled={importMutation.isPending}>{importMutation.isPending ? "Importando…" : "Registrar manifesto"}</Button></div>
                {importMutation.isError && <p className="text-sm text-destructive sm:col-span-2">Não foi possível validar o manifesto.</p>}
              </form>}
            </div>
            <div className="min-h-0 flex-1 overflow-y-auto p-5 sm:p-8">
              <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
                {visibleApps.map((app) => {
                  const installed = installedById.get(app.id);
                  return (
                    <Card
                      key={app.id}
                      className="transition-colors hover:border-primary/40"
                    >
                      <CardContent className="flex min-h-56 flex-col gap-4 p-4">
                        <div className="flex items-start gap-3">
                          <AppIcon
                            icon={app.icon}
                            name={app.name}
                            className="size-11 object-contain"
                          />
                          <div className="min-w-0">
                            <h2 className="font-medium">{app.name}</h2>
                            <p className="mt-1 line-clamp-3 text-sm text-muted-foreground">
                              {app.description}
                            </p>
                          </div>
                        </div>
                        <div className="mt-auto flex gap-2">
                          <Button
                            variant="outline"
                            size="sm"
                            className="min-h-10 flex-1"
                            onClick={() => setSelectedApp(app)}
                          >
                            Detalhes
                          </Button>
                          <Button
                            size="sm"
                            className={`min-h-10 flex-1 ${
                              installed
                                ? "bg-primary/15 text-primary opacity-100 disabled:opacity-100"
                                : ""
                            }`}
                            disabled={Boolean(installed) || !app.installable}
                            onClick={() => {
                              setSelectedApp(app);
                            }}
                          >
                            {installed
                              ? installed.status === "installing"
                                ? "Instalando"
                                : "Instalado"
                              : app.installable ? "Instalar" : "Em breve"}
                          </Button>
                        </div>
                      </CardContent>
                    </Card>
                  );
                })}
              </div>
            </div>
          </>
        )}
      </DialogContent>
    </Dialog>
  );
}

function AppDetails({
  app,
  installed,
  onBack,
  onInstall,
}: {
  app: StoreApp;
  installed?: InstalledApp;
  onBack: () => void;
  onInstall: () => void;
}) {
  return (
    <div className="flex min-h-0 flex-1 flex-col overflow-y-auto">
      <div className="flex items-center border-b border-border/60 px-5 py-3 sm:px-6">
        <Button variant="ghost" size="sm" onClick={onBack}>
          <ArrowLeft data-icon="inline-start" />
          Voltar
        </Button>
      </div>
      <div className="grid gap-6 p-5 sm:p-6 lg:grid-cols-[minmax(0,1fr)_minmax(16rem,0.7fr)]">
        <div>
          <div className="flex items-start gap-4">
            <AppIcon icon={app.icon} name={app.name} className="size-16 object-contain" />
            <div>
              <h2 className="text-xl font-medium">{app.name}</h2>
              <p className="mt-1 text-muted-foreground">{app.description}</p>
            </div>
          </div>
          <section className="mt-7">
            <h3 className="text-sm font-medium">Recursos principais</h3>
            <ul className="mt-3 space-y-2 text-sm text-muted-foreground">
              {app.highlights.map((highlight) => (
                <li key={highlight} className="flex gap-2">
                  <Check className="mt-0.5 size-4 shrink-0 text-primary" />
                  {highlight}
                </li>
              ))}
            </ul>
          </section>
        </div>
        <aside className="space-y-4 rounded-xl border bg-muted/20 p-4">
          <div className="grid grid-cols-2 gap-3 text-sm">
            <div>
              <span className="block text-xs text-muted-foreground">
                Categoria
              </span>
              <span className="mt-1 block capitalize">{app.category}</span>
            </div>
            <div>
              <span className="block text-xs text-muted-foreground">
                Versão
              </span>
              <span className="mt-1 block">{app.version}</span>
            </div>
            <div>
              <span className="block text-xs text-muted-foreground">
                Tamanho
              </span>
              <span className="mt-1 block">{app.size}</span>
            </div>
            <div>
              <span className="block text-xs text-muted-foreground">
                Origem
              </span>
              <span className="mt-1 block">Oficial</span>
            </div>
          </div>
          <Button
            className="w-full"
            disabled={Boolean(installed) || !app.installable}
            onClick={onInstall}
          >
            <PackagePlus data-icon="inline-start" />
            {installed
              ? installed.status === "installing"
                ? "Instalando"
                : "Já instalado"
              : app.installable ? "Instalar app" : "Em breve"}
          </Button>
          {installed ? (
            <p className="text-xs text-muted-foreground">
              Este app já está na sua página Apps.
            </p>
          ) : (
            <p className="text-xs text-muted-foreground">
              Este app ainda não possui instalação segura no LabsOS.
            </p>
          )}
        </aside>
      </div>
    </div>
  );
}
