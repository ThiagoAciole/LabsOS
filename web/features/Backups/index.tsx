import { useState } from "react"
import { Archive, CheckCircle2, CircleAlert, RefreshCw, ShieldCheck, Trash2 } from "lucide-react"
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { createBackup, deleteBackup, getBackups, restoreBackup, verifyBackup } from "@/api/settings"

function formatBytes(value: number) {
  if (!Number.isFinite(value) || value < 1) return "0 B"
  const units = ["B", "KB", "MB", "GB", "TB"]
  const index = Math.min(Math.floor(Math.log(value) / Math.log(1024)), units.length - 1)
  return `${(value / 1024 ** index).toFixed(index ? 1 : 0)} ${units[index]}`
}

export function BackupsPage() {
  const client = useQueryClient()
  const backups = useQuery({ queryKey: ["backups"], queryFn: getBackups })
  const [apps, setApps] = useState("")
  const [message, setMessage] = useState<string | null>(null)
  const refresh = () => void client.invalidateQueries({ queryKey: ["backups"] })
  const create = useMutation({ mutationFn: () => createBackup(apps.split(",").map((item) => item.trim()).filter(Boolean)), onSuccess: (result) => { setMessage(result.message ?? (result.simulated ? "Backup planejado; operações reais estão desativadas." : "Backup criado.")); refresh() } })
  const verify = useMutation({ mutationFn: verifyBackup, onSuccess: (result) => setMessage(`Integridade ${result.integrity}: ${result.sha256}`) })
  const restore = useMutation({ mutationFn: restoreBackup, onSuccess: (result) => setMessage(result.message ?? (result.simulated ? "Restauração planejada; nenhum arquivo foi alterado." : "Restauração concluída.")) })
  const remove = useMutation({ mutationFn: deleteBackup, onSuccess: (result) => { setMessage(result.simulated ? "Exclusão planejada; o arquivo foi preservado." : "Backup excluído."); refresh() } })
  const busy = create.isPending || verify.isPending || restore.isPending || remove.isPending

  return <section className="flex min-w-0 flex-1 flex-col gap-6 p-4 md:p-8">
    <header><h1 className="text-2xl font-medium md:text-3xl">Backups</h1><p className="mt-1 text-sm text-muted-foreground">Arquivos de dados e configuração do LabsOS, com operações protegidas.</p></header>
    {message && <div role="status" className="rounded-lg border border-primary/30 bg-primary/5 p-3 text-sm">{message}</div>}
    {backups.isError && <div role="alert" className="flex items-center gap-2 rounded-lg border border-destructive/40 p-4 text-sm"><CircleAlert className="size-4" />Não foi possível carregar os backups.</div>}
    <Card><CardHeader><CardTitle className="flex items-center gap-2 text-base"><Archive className="size-4" />Criar backup</CardTitle></CardHeader><CardContent className="space-y-3"><p className="text-sm text-muted-foreground">Por padrão inclui DATA. Informe Apps opcionais separados por vírgula.</p><div className="flex flex-col gap-2 sm:flex-row"><Input value={apps} onChange={(event) => setApps(event.target.value)} placeholder="jellyfin, nextcloud" aria-label="Apps opcionais" /><Button onClick={() => create.mutate()} disabled={busy}>{create.isPending ? "Criando…" : "Criar backup"}</Button></div></CardContent></Card>
    {backups.isPending ? <p role="status">Carregando backups…</p> : (backups.data ?? []).length === 0 ? <Card><CardContent className="flex min-h-36 flex-col items-center justify-center gap-2 text-center"><ShieldCheck className="size-7 text-muted-foreground" /><p className="font-medium">Nenhum backup encontrado</p><p className="text-sm text-muted-foreground">Crie o primeiro backup acima.</p></CardContent></Card> : <div className="grid gap-3">{(backups.data ?? []).map((backup) => <Card key={backup.id}><CardContent className="flex flex-col gap-3 p-4 sm:flex-row sm:items-center sm:justify-between"><div className="min-w-0"><p className="truncate font-medium">{backup.id}</p><p className="text-sm text-muted-foreground">{formatBytes(backup.size)} · {new Date(backup.createdAt).toLocaleString("pt-BR")}</p>{backup.sha256 && <p className="truncate text-xs text-muted-foreground">SHA-256: {backup.sha256}</p>}</div><div className="flex flex-wrap gap-2"><Button variant="outline" size="sm" onClick={() => verify.mutate(backup.id)} disabled={busy}><CheckCircle2 className="mr-1 size-4" />Verificar</Button><Button variant="outline" size="sm" onClick={() => restore.mutate(backup.id)} disabled={busy}><RefreshCw className="mr-1 size-4" />Restaurar</Button><Button variant="ghost" size="sm" onClick={() => remove.mutate(backup.id)} disabled={busy}><Trash2 className="mr-1 size-4" />Excluir</Button></div></CardContent></Card>)}</div>}
  </section>
}
