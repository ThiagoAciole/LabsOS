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
  KeyRound,
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
import { useChangePasswordMutation, useDiagnostics, useNetworkMutation, useNetworkSettingsMutation, useNetworkSnapshot, usePowerMutation, useRevokeSessionMutation, useSessions, useSettingsData, useSSHMutation, useSSHStatus, useStorageSnapshot, useSystemRollbackMutation, useSystemSettingsMutation, useSystemUpdateData, useSystemUpdateMutation, useWiFiMutation, useWiFiSnapshot } from "./queries";

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
  const rollbackMutation = useSystemRollbackMutation();
  const systemMutation = useSystemSettingsMutation();
  const networkMutation = useNetworkSettingsMutation();
  const networkQuery = useNetworkSnapshot();
  const networkChangeMutation = useNetworkMutation();
  const sshQuery = useSSHStatus();
  const sshMutation = useSSHMutation();
  const wifiQuery = useWiFiSnapshot();
  const wifiMutation = useWiFiMutation();
  const storageQuery = useStorageSnapshot();
  const passwordMutation = useChangePasswordMutation();
  const sessionsQuery = useSessions();
  const revokeSessionMutation = useRevokeSessionMutation();
  const diagnosticsQuery = useDiagnostics();
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
  const [currentPassword, setCurrentPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [selectedInterface, setSelectedInterface] = useState("");
  const [staticAddress, setStaticAddress] = useState("");
  const [staticGateway, setStaticGateway] = useState("");
  const [staticDNS, setStaticDNS] = useState("");
  const [sshKey, setSSHKey] = useState("");
  const [wifiSSID, setWifiSSID] = useState("");
  const [wifiPassword, setWifiPassword] = useState("");

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

  function openDiagnostics() {
    setDiagnosticsOpen(true);
    void diagnosticsQuery.refetch();
  }

  async function savePassword() {
    if (newPassword.length < 8) { setNotice("A nova senha precisa ter pelo menos 8 caracteres."); return }
    if (newPassword !== confirmPassword) { setNotice("A confirmação da senha não confere."); return }
    try { const result = await passwordMutation.mutateAsync({ current: currentPassword, next: newPassword }); setNotice(result.simulated ? "Troca de senha planejada; operações reais estão desativadas." : "Senha alterada com sucesso."); setCurrentPassword(""); setNewPassword(""); setConfirmPassword(""); } catch { setNotice("Não foi possível alterar a senha."); }
  }

  async function requestDHCP() {
    if (!selectedInterface) { setNotice("Selecione uma interface de rede."); return }
    try { const result = await networkChangeMutation.mutateAsync({ interface: selectedInterface, dhcp: true }); setNotice(result.simulated ? "DHCP planejado; nenhuma alteração de rede foi aplicada." : result.message); }
    catch { setNotice("Não foi possível configurar DHCP nessa interface."); }
  }

  async function requestStaticNetwork() {
    if (!selectedInterface || !staticAddress || !staticGateway || !staticDNS) { setNotice("Preencha interface, endereço CIDR, gateway e DNS."); return }
    try { const result = await networkChangeMutation.mutateAsync({ interface: selectedInterface, dhcp: false, address: staticAddress, gateway: staticGateway, dns: staticDNS }); setNotice(result.simulated ? "Rede estática planejada; nenhuma alteração foi aplicada." : result.message); }
    catch { setNotice("Não foi possível configurar o endereço estático."); }
  }

  async function addSSHKey() {
    if (!sshKey.trim()) return;
    try { const result = await sshMutation.mutateAsync({ action: "add-key", key: sshKey.trim() }); setNotice(result.simulated ? "Chave SSH planejada; nenhuma alteração foi aplicada." : result.message); if (!result.simulated) setSSHKey(""); } catch { setNotice("Não foi possível adicionar a chave SSH."); }
  }

  async function removeSSHKey(fingerprint: string) {
    try { const result = await sshMutation.mutateAsync({ action: "remove-key", fingerprint }); setNotice(result.simulated ? "Remoção da chave planejada; nenhuma alteração foi aplicada." : result.message); } catch { setNotice("Não foi possível remover a chave SSH."); }
  }

  async function connectWiFi() {
    if (!wifiSSID) { setNotice("Selecione uma rede Wi-Fi."); return }
    try { const result = await wifiMutation.mutateAsync({ action: "connect", ssid: wifiSSID, password: wifiPassword, device: wifiQuery.data?.devices[0] }); setNotice(result.simulated ? "Conexão Wi-Fi planejada; nenhuma alteração foi aplicada." : result.message); } catch { setNotice("Não foi possível conectar ao Wi-Fi."); }
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

        <SettingsSection
          icon={KeyRound}
          title="Segurança"
          description="Altere a senha administrativa do painel."
        >
          <div className="grid gap-3 sm:grid-cols-3">
            <Input type="password" placeholder="Senha atual" value={currentPassword} onChange={(event) => setCurrentPassword(event.target.value)} aria-label="Senha atual" />
            <Input type="password" placeholder="Nova senha" value={newPassword} onChange={(event) => setNewPassword(event.target.value)} aria-label="Nova senha" />
            <Input type="password" placeholder="Confirmar nova senha" value={confirmPassword} onChange={(event) => setConfirmPassword(event.target.value)} aria-label="Confirmar nova senha" />
          </div>
          <Button className="mt-3" variant="outline" disabled={passwordMutation.isPending || !currentPassword || !newPassword || !confirmPassword} onClick={() => void savePassword()}>{passwordMutation.isPending ? "Salvando…" : "Alterar senha"}</Button>
          {(sessionsQuery.data ?? []).length > 0 && <div className="mt-4 space-y-2 border-t pt-3"><p className="text-sm font-medium">Sessões ativas</p>{(sessionsQuery.data ?? []).map((session) => <div key={session.id} className="flex items-center justify-between gap-2 text-xs"><span className="truncate">{session.current ? "Esta sessão" : session.userAgent || "Dispositivo"} · último acesso {new Date(session.lastSeen).toLocaleString("pt-BR")}</span>{!session.current && <Button variant="ghost" size="sm" disabled={revokeSessionMutation.isPending} onClick={() => revokeSessionMutation.mutate(session.id)}>Revogar</Button>}</div>)}</div>}
        </SettingsSection>
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
          <SettingsRow
            icon={Network}
            label="Interfaces"
            description={networkQuery.data ? `${networkQuery.data.interfaces.length} interface(s) detectada(s)` : "Dados de rede indisponíveis"}
          />
          <div className="grid gap-2 rounded-lg border p-3 sm:grid-cols-[minmax(0,1fr)_auto]">
            <select className="h-9 rounded-md border bg-background px-3 text-sm" value={selectedInterface} onChange={(event) => setSelectedInterface(event.target.value)} aria-label="Interface para DHCP">
              <option value="">Selecionar interface</option>
              {(networkQuery.data?.interfaces ?? []).map((item) => item.ifname ? <option key={item.ifname} value={item.ifname}>{item.ifname} · {item.operstate ?? "estado desconhecido"}</option> : null)}
            </select>
            <Button variant="outline" disabled={networkChangeMutation.isPending || !selectedInterface} onClick={() => void requestDHCP()}>{networkChangeMutation.isPending ? "Solicitando…" : "Usar DHCP"}</Button>
          </div>
          <div className="grid gap-2 rounded-lg border p-3 sm:grid-cols-4">
            <Input value={staticAddress} onChange={(event) => setStaticAddress(event.target.value)} placeholder="Endereço CIDR" aria-label="Endereço IP CIDR" />
            <Input value={staticGateway} onChange={(event) => setStaticGateway(event.target.value)} placeholder="Gateway" aria-label="Gateway" />
            <Input value={staticDNS} onChange={(event) => setStaticDNS(event.target.value)} placeholder="DNS (vírgulas)" aria-label="DNS" />
            <Button variant="outline" disabled={networkChangeMutation.isPending || !selectedInterface} onClick={() => void requestStaticNetwork()}>Usar IP estático</Button>
          </div>
          <SettingsRow
            icon={ShieldCheck}
            label="SSH"
            description={sshQuery.data ? (sshQuery.data.active ? `Ativo na porta ${sshQuery.data.port}` : "Serviço SSH inativo") : "Status SSH indisponível"}
          >
            <Button variant="outline" size="sm" disabled={sshMutation.isPending || !sshQuery.data} onClick={() => void sshMutation.mutateAsync({ enabled: !sshQuery.data?.active }).then((result) => setNotice(result.simulated ? "Alteração do SSH planejada." : result.message)).catch(() => setNotice("Não foi possível alterar o SSH."))}>{sshQuery.data?.active ? "Desativar" : "Ativar"}</Button>
          </SettingsRow>
          <div className="space-y-3 rounded-lg border p-3">
            <div className="flex flex-wrap items-center justify-between gap-2 text-sm"><span>Comando de acesso</span><code className="text-xs text-muted-foreground">ssh labs@{summary.hostname}.local</code></div>
            <div className="flex flex-col gap-2 sm:flex-row"><Input value={sshKey} onChange={(event) => setSSHKey(event.target.value)} placeholder="ssh-ed25519 AAAA... comentário" aria-label="Chave pública SSH" /><Button variant="outline" disabled={sshMutation.isPending || !sshKey.trim()} onClick={() => void addSSHKey()}>Adicionar chave</Button></div>
            {(sshQuery.data?.keys ?? []).map((key) => <div key={`${key.fingerprint}-${key.comment ?? ""}`} className="flex items-center justify-between gap-2 text-xs"><span className="truncate">{key.fingerprint}{key.comment ? ` · ${key.comment}` : ""}</span><Button variant="ghost" size="sm" disabled={sshMutation.isPending} onClick={() => void removeSSHKey(key.fingerprint)}>Remover</Button></div>)}
          </div>
          <SettingsRow
            icon={Network}
            label="Wi-Fi"
            description={wifiQuery.data?.available ? (wifiQuery.data.current ? `Conectado a ${wifiQuery.data.current}` : `${wifiQuery.data.networks.length} rede(s) visível(is)`) : "Wi-Fi não disponível"}
          />
          <div className="space-y-2 rounded-lg border p-3"><div className="flex flex-col gap-2 sm:flex-row"><select className="h-9 min-w-0 flex-1 rounded-md border bg-background px-3 text-sm" value={wifiSSID} onChange={(event) => setWifiSSID(event.target.value)} aria-label="Rede Wi-Fi"><option value="">Selecionar rede</option>{(wifiQuery.data?.networks ?? []).map((ssid) => <option key={ssid} value={ssid}>{ssid}</option>)}</select><Input type="password" value={wifiPassword} onChange={(event) => setWifiPassword(event.target.value)} placeholder="Senha Wi-Fi" aria-label="Senha Wi-Fi" /><Button variant="outline" disabled={wifiMutation.isPending || !wifiSSID} onClick={() => void connectWiFi()}>{wifiMutation.isPending ? "Conectando…" : "Conectar"}</Button></div>{wifiQuery.data?.current && <div className="flex justify-end gap-2"><Button variant="ghost" size="sm" disabled={wifiMutation.isPending} onClick={() => void wifiMutation.mutateAsync({ action: "reconnect", ssid: wifiQuery.data?.current ?? "" }).then((result) => setNotice(result.simulated ? "Reconexão Wi-Fi planejada." : result.message)).catch(() => setNotice("Não foi possível reconectar ao Wi-Fi."))}>Reconectar</Button><Button variant="ghost" size="sm" disabled={wifiMutation.isPending} onClick={() => void wifiMutation.mutateAsync({ action: "forget", ssid: wifiQuery.data?.current ?? "" }).then((result) => setNotice(result.simulated ? "Remoção da rede planejada." : result.message)).catch(() => setNotice("Não foi possível esquecer a rede."))}>Esquecer rede</Button></div>}</div>
          <SettingsRow
            icon={Server}
            label="Armazenamento"
            description={storageQuery.data?.usage[0] ? `${storageQuery.data.usage[0].usePercent}% usado em ${storageQuery.data.usage[0].path}` : "Uso de armazenamento indisponível"}
          />
          <div className="space-y-2 rounded-lg border p-3 text-xs"><div className="flex justify-between gap-2"><span className="font-medium">Discos e partições</span><span className="text-muted-foreground">{storageQuery.data?.dataMounted ? "/DATA montado" : "/DATA não detectado"}</span></div>{(storageQuery.data?.devices ?? []).map((device) => <div key={device.path ?? device.name}><div className="flex justify-between gap-2"><span>{device.name} · {device.type}{device.transport ? ` · ${device.transport}` : ""}</span><span className="text-muted-foreground">{device.fstype ?? "sem filesystem"}{device.uuid ? ` · ${device.uuid.slice(0, 12)}` : ""}</span></div>{device.children?.map((child) => <div key={child.path ?? child.name} className="ml-4 text-muted-foreground">└ {child.name} · {child.mountpoints?.filter(Boolean).join(", ") || "não montada"}</div>)}</div>)}</div>
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
          {updateQuery.data?.updateAvailable && <div className="rounded-lg border p-3 text-sm"><p className="font-medium">Novidades {updateQuery.data.latestVersion}{updateQuery.data.channel ? ` · canal ${updateQuery.data.channel}` : ""}</p><p className="mt-1 whitespace-pre-wrap text-xs text-muted-foreground">{updateQuery.data.changelog || "Nenhum changelog informado."}</p>{updateQuery.data.sha256 && <p className="mt-2 break-all font-mono text-[10px] text-muted-foreground">SHA-256: {updateQuery.data.sha256}</p>}</div>}
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
              onClick={() => void updateMutation.mutateAsync().then((result) => setNotice(result.simulated ? "Atualização planejada; nenhuma alteração foi aplicada." : result.message ?? "Atualização instalada; serviços serão reiniciados.")).catch(() => setNotice("Não foi possível instalar a atualização."))}
              >
              {updateQuery.data?.updateAvailable ? "Atualizar" : "Verificar"}
            </Button>
          </SettingsRow>
          <SettingsRow
            icon={RefreshCw}
            label="Voltar para a versão anterior"
            description="Usa a release preservada pelo último update."
          >
            <Button
              variant="outline"
              disabled={rollbackMutation.isPending}
              onClick={() => void rollbackMutation.mutateAsync().then((result) => setNotice(result.simulated ? "Rollback planejado; nenhuma alteração foi aplicada." : result.message ?? "Rollback concluído.")).catch(() => setNotice("Não foi possível fazer rollback."))}
            >
              {rollbackMutation.isPending ? "Restaurando…" : "Rollback"}
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
            description={diagnosticsQuery.data ? (diagnosticsQuery.data.healthy ? "Todos os checks estão saudáveis." : `${diagnosticsQuery.data.checks.filter((check) => check.status !== "healthy").length} check(s) precisam de atenção.`) : "Execute o diagnóstico para verificar o servidor."}
          >
            <Button variant="outline" onClick={openDiagnostics}>
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
            {diagnosticsQuery.isPending && <p className="rounded-lg border p-3">Executando verificações…</p>}
            {diagnosticsQuery.data?.checks.map((check) => <p key={check.id} className="rounded-lg border p-3">{check.id}: {check.message}</p>)}
            {diagnosticsQuery.isError && <p className="rounded-lg border p-3 text-destructive">Não foi possível executar o diagnóstico.</p>}
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
