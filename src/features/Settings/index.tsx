import { useState, type ReactNode } from "react";
import {
  CirclePower,
  Copy,
  FileText,
  Globe2,
  Languages,
  MonitorCog,
  Network,
  RefreshCw,
  Server,
  ShieldCheck,
  Wrench,
} from "lucide-react";

import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogMedia,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Switch } from "@/components/ui/switch";
import { useTheme } from "@/components/theme-provider";
import { useNetworkSettingsMutation, usePowerMutation, useSettingsData, useSystemSettingsMutation, useSystemUpdateData, useSystemUpdateMutation } from "./queries";

type SensitiveAction = "restart" | "shutdown";
type Theme = "dark" | "light" | "system";

const actionCopy = {
  restart: {
    title: "Reiniciar o servidor?",
    description:
      "Os apps ficarão indisponíveis por alguns minutos enquanto o LabsOS reinicia.",
    confirm: "Reiniciar servidor",
  },
  shutdown: {
    title: "Desligar o servidor?",
    description:
      "Os apps e o acesso aos arquivos ficarão indisponíveis até o servidor ser ligado novamente.",
    confirm: "Desligar servidor",
  },
} as const;

export function SettingsPage() {
  const { theme, setTheme } = useTheme();
  const query = useSettingsData();
  const updateQuery = useSystemUpdateData();
  const updateMutation = useSystemUpdateMutation();
  const systemMutation = useSystemSettingsMutation();
  const networkMutation = useNetworkSettingsMutation();
  const powerMutation = usePowerMutation();
  const [serverName, setServerName] = useState<string | null>(null);
  const [language, setLanguage] = useState<string | null>(null);
  const [remoteAccess, setRemoteAccess] = useState<boolean | null>(null);
  const [pendingRemoteAccess, setPendingRemoteAccess] = useState<
    boolean | null
  >(null);
  const [pendingAction, setPendingAction] = useState<SensitiveAction | null>(
    null,
  );
  const [diagnosticsOpen, setDiagnosticsOpen] = useState(false);
  const [notice, setNotice] = useState(
    "Conectado à Labs API.",
  );

  async function saveGeneral() {
    try { await systemMutation.mutateAsync({ hostname: serverName ?? query.data?.settings.hostname, language: language ?? query.data?.settings.language }); setNotice("Preferências gerais atualizadas."); }
    catch { setNotice("Não foi possível salvar as preferências. Tente novamente."); }
  }

  async function confirmRemoteAccess() {
    if (pendingRemoteAccess === null) return;
    try { await networkMutation.mutateAsync({ remoteAccess: pendingRemoteAccess }); setRemoteAccess(pendingRemoteAccess); setNotice(pendingRemoteAccess ? "Acesso remoto ativado." : "Acesso remoto desativado."); setPendingRemoteAccess(null); }
    catch { setNotice("Não foi possível alterar o acesso remoto. Tente novamente."); }
  }

  async function confirmMaintenance() {
    if (!pendingAction) return;
    try { const job = await powerMutation.mutateAsync(pendingAction === "restart" ? "reboot" : "shutdown"); setNotice(job.message); setPendingAction(null); }
    catch { setNotice("Não foi possível concluir esta ação. Tente novamente."); }
  }

  function copyLocalAddress() {
    void navigator.clipboard?.writeText(`${query.data?.summary.hostname ?? serverName ?? "labsos"}.local`);
    setNotice("Endereço local copiado.");
  }

  if (query.isPending) return <section className="flex flex-1 items-center justify-center p-8" role="status">Carregando configurações…</section>;
  if (query.isError) return <section className="flex flex-1 flex-col items-center justify-center gap-3 p-8" role="alert"><p>Não foi possível carregar as configurações.</p><Button variant="outline" onClick={() => void query.refetch()}>Tentar novamente</Button></section>;
  const { summary, health } = query.data;
  const uptimeDays = Math.floor(summary.uptimeSeconds / 86400);
  const currentServerName = serverName ?? query.data.settings.hostname ?? "";
  const currentLanguage = language ?? query.data.settings.language ?? "pt-BR";
  const currentRemoteAccess = remoteAccess ?? Boolean(query.data.settings.remoteAccess);

  return (
    <section className="flex min-w-0 flex-1 flex-col gap-6 p-4 md:p-8">
      <header>
        <h1 className="text-2xl font-medium md:text-3xl">Settings</h1>
        <p className="mt-1 text-sm text-muted-foreground">
          Preferências, acesso e cuidados do seu servidor.
        </p>
      </header>

      <p role="status" className="text-sm text-muted-foreground">
        {notice}
      </p>

      <div className="grid items-start gap-6 xl:grid-cols-2">
        <SettingsSection
          icon={MonitorCog}
          title="Geral"
          description="Nome, idioma e aparência."
        >
          <div className="grid gap-4 sm:grid-cols-[minmax(0,1fr)_auto] sm:items-end">
            <label className="grid gap-2 text-sm font-medium">
              Nome do servidor
              <Input
                value={currentServerName}
                onChange={(event) => setServerName(event.target.value)}
              />
            </label>
            <Button disabled={systemMutation.isPending} onClick={() => void saveGeneral()}>Salvar</Button>
          </div>
          <SettingsRow
            icon={Languages}
            label="Idioma"
            description="Usado no painel LabsOS."
          >
            <Select
              value={currentLanguage}
              onValueChange={(value) => value && setLanguage(value)}
            >
              <SelectTrigger className="w-36">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="pt-BR">Português</SelectItem>
                <SelectItem value="en-US">English</SelectItem>
              </SelectContent>
            </Select>
          </SettingsRow>
          <SettingsRow
            icon={MonitorCog}
            label="Aparência"
            description="Escolha como o painel acompanha sua tela."
          >
            <Select
              value={theme}
              onValueChange={(value) => setTheme(value as Theme)}
            >
              <SelectTrigger className="w-36">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="dark">Escuro</SelectItem>
                <SelectItem value="light">Claro</SelectItem>
                <SelectItem value="system">Sistema</SelectItem>
              </SelectContent>
            </Select>
          </SettingsRow>
          <SettingsRow
            icon={Server}
            label="Informações do sistema"
            description={`Hostname: ${summary.hostname} · Versão: ${summary.version} · Uptime: ${uptimeDays} dias`}
          />
        </SettingsSection>

        <SettingsSection
          icon={Network}
          title="Rede e acesso"
          description="Conexão atual e como você chega ao servidor."
        >
          <SettingsRow
            icon={Globe2}
            label="Conexão"
            description={health.status === "healthy" ? "Online · IP local indisponível" : "Conexão indisponível"}
          />
          <SettingsRow
            icon={Network}
            label="Endereço local"
            description={`${summary.hostname}.local`}
          >
            <Button
              variant="ghost"
              size="icon"
              aria-label="Copiar endereço local"
              onClick={copyLocalAddress}
            >
              <Copy />
            </Button>
          </SettingsRow>
          <SettingsRow
            icon={Network}
            label="Acesso local"
            description="Disponível na sua rede doméstica."
          />
          <SettingsRow
            icon={ShieldCheck}
            label="Acesso remoto"
            description={
              currentRemoteAccess
                ? "Ativado. Você pode acessar o LabsOS fora de casa."
                : "Desativado. O LabsOS só está disponível na rede local."
            }
          >
            <Switch
              checked={currentRemoteAccess}
              onCheckedChange={setPendingRemoteAccess}
              aria-label="Alterar acesso remoto"
            />
          </SettingsRow>
        </SettingsSection>

        <SettingsSection
          icon={RefreshCw}
          title="Atualizações"
          description="Versão e novidades do LabsOS."
        >
          <SettingsRow
            icon={Server}
            label="LabsOS"
            description={`Versão ${updateQuery.data?.currentVersion ?? summary.version}${updateQuery.data?.updateAvailable ? ` · Nova versão ${updateQuery.data.latestVersion}` : " · Atualizado"}`}
          />
          <SettingsRow
            icon={RefreshCw}
            label="Canal de atualização"
            description="Receba versões estáveis do LabsOS."
          >
            <Select value="stable">
              <SelectTrigger className="w-36">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="stable">Estável</SelectItem>
              </SelectContent>
            </Select>
          </SettingsRow>
          <SettingsRow
            icon={RefreshCw}
            label="Atualizações automáticas"
            description="Mantenha o LabsOS atualizado automaticamente."
          >
            <span className="text-xs text-muted-foreground">Manual</span>
          </SettingsRow>
          <SettingsRow
            icon={RefreshCw}
            label="Verificar atualizações"
            description={updateQuery.isError ? "Não foi possível consultar atualizações." : updateQuery.data?.updateAvailable ? "Há uma atualização disponível." : "Sistema atualizado."}
          >
            <Button
              variant="outline"
              disabled={updateMutation.isPending || !updateQuery.data?.updateAvailable}
              onClick={() => void updateMutation.mutateAsync().then(() => setNotice("Atualização instalada; serviços serão reiniciados.")).catch(() => setNotice("Não foi possível instalar a atualização."))}
              >
              {updateQuery.data?.updateAvailable ? "Atualizar" : "Verificar"}
            </Button>
          </SettingsRow>
        </SettingsSection>

        <SettingsSection
          icon={Wrench}
          title="Manutenção"
          description="Ações que afetam a disponibilidade do servidor."
          tone="maintenance"
        >
          <SettingsRow
            icon={Wrench}
            label="Diagnóstico"
            description="Todos os serviços estão saudáveis."
          >
            <Button variant="outline" onClick={() => setDiagnosticsOpen(true)}>
              Ver diagnóstico
            </Button>
          </SettingsRow>
          <SettingsRow
            icon={FileText}
            label="Relatório do sistema"
            description="Gere informações úteis para identificar problemas."
          >
            <Button
              variant="outline"
              onClick={() => setNotice("Relatório gerado para demonstração.")}
            >
              Gerar relatório
            </Button>
          </SettingsRow>
          <div className="border-t border-border/60 pt-4">
            <p className="text-sm font-medium">Ações do servidor</p>
            <p className="mt-1 text-xs text-muted-foreground">
              Estas ações interrompem temporariamente o acesso ao servidor.
            </p>
          </div>
          <div className="grid gap-2 sm:grid-cols-2">
            <Button
              variant="outline"
              onClick={() => setPendingAction("restart")}
            >
              <RefreshCw data-icon="inline-start" />
              Reiniciar servidor
            </Button>
            <Button
              variant="outline"
              onClick={() => setPendingAction("shutdown")}
            >
              <CirclePower data-icon="inline-start" />
              Desligar servidor
            </Button>
          </div>
        </SettingsSection>
      </div>

      <AlertDialog
        open={pendingRemoteAccess !== null}
        onOpenChange={(open) => !open && setPendingRemoteAccess(null)}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogMedia>
              <ShieldCheck />
            </AlertDialogMedia>
            <AlertDialogTitle>
              {pendingRemoteAccess
                ? "Ativar acesso remoto?"
                : "Desativar acesso remoto?"}
            </AlertDialogTitle>
            <AlertDialogDescription>
              {pendingRemoteAccess
                ? "Você poderá acessar o LabsOS fora da rede local. Ative apenas se reconhecer esse impacto."
                : "Você só poderá acessar o LabsOS pela rede local até ativar o acesso remoto novamente."}
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancelar</AlertDialogCancel>
            <AlertDialogAction disabled={networkMutation.isPending} onClick={() => void confirmRemoteAccess()}>
              {pendingRemoteAccess ? "Ativar acesso remoto" : "Desativar acesso remoto"}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      <AlertDialog
        open={pendingAction !== null}
        onOpenChange={(open) => !open && setPendingAction(null)}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogMedia>
              <CirclePower />
            </AlertDialogMedia>
            <AlertDialogTitle>
              {pendingAction ? actionCopy[pendingAction].title : ""}
            </AlertDialogTitle>
            <AlertDialogDescription>
              {pendingAction ? actionCopy[pendingAction].description : ""}
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancelar</AlertDialogCancel>
            <AlertDialogAction
              variant="destructive"
              disabled={powerMutation.isPending}
              onClick={() => void confirmMaintenance()}
            >
              {pendingAction ? actionCopy[pendingAction].confirm : "Confirmar"}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      <Dialog open={diagnosticsOpen} onOpenChange={setDiagnosticsOpen}>
        <DialogContent className="max-w-md">
          <DialogHeader>
            <DialogTitle>Diagnóstico do servidor</DialogTitle>
            <DialogDescription>Resumo atual do LabsOS.</DialogDescription>
          </DialogHeader>
          <div className="space-y-3 text-sm">
            <p className="rounded-lg border p-3">
              Serviços do LabsOS: {health.status === "healthy" ? "saudáveis" : "precisam de atenção"}
            </p>
            <p className="rounded-lg border p-3">Última verificação: agora</p>
            <p className="rounded-lg border p-3">
              Tempo ligado: {uptimeDays} dias
            </p>
          </div>
        </DialogContent>
      </Dialog>
    </section>
  );
}

function SettingsSection({
  icon: Icon,
  title,
  description,
  tone,
  children,
}: {
  icon: typeof MonitorCog;
  title: string;
  description: string;
  tone?: "maintenance";
  children: ReactNode;
}) {
  return (
    <Card className={tone ? "border-primary/20" : undefined}>
      <CardHeader>
        <CardTitle className="flex items-center gap-2 text-base">
          <Icon className="size-4 text-primary" />
          {title}
        </CardTitle>
        <p className="text-sm text-muted-foreground">{description}</p>
      </CardHeader>
      <CardContent className="space-y-4">{children}</CardContent>
    </Card>
  );
}

function SettingsRow({
  icon: Icon,
  label,
  description,
  children,
}: {
  icon: typeof MonitorCog;
  label: string;
  description: string;
  children?: ReactNode;
}) {
  return (
    <div className="flex min-h-14 items-center gap-3 rounded-lg border border-border/60 px-3 py-2.5">
      <Icon className="size-4 shrink-0 text-muted-foreground" />
      <div className="min-w-0 flex-1">
        <p className="text-sm font-medium">{label}</p>
        <p className="text-xs text-muted-foreground">{description}</p>
      </div>
      {children}
    </div>
  );
}
