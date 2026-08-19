import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import { useState } from "react"
import { CheckCircle2, CircleAlert, CircleDashed, CircleX, ListTodo, X } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Progress } from "@/components/ui/progress"
import { Input } from "@/components/ui/input"
import { cancelJob, createScheduledTask, deleteScheduledTask, getJobs, getScheduledTasks, type JobDTO } from "@/api/system"

const statusIcon = { queued: CircleDashed, running: CircleDashed, success: CheckCircle2, error: CircleAlert, cancelled: CircleX } as const
const statusLabel: Record<string, string> = { queued: "Na fila", running: "Executando", success: "Concluída", error: "Erro", cancelled: "Cancelada" }

export function JobsPage() {
  const client = useQueryClient()
  const query = useQuery({ queryKey: ["jobs"], queryFn: getJobs })
  const schedules = useQuery({ queryKey: ["scheduler"], queryFn: getScheduledTasks })
  const cancellation = useMutation({ mutationFn: cancelJob, onSuccess: () => client.invalidateQueries({ queryKey: ["jobs"] }) })
  const createSchedule = useMutation({ mutationFn: createScheduledTask, onSuccess: () => { void client.invalidateQueries({ queryKey: ["scheduler"] }) } })
  const removeSchedule = useMutation({ mutationFn: deleteScheduledTask, onSuccess: () => { void client.invalidateQueries({ queryKey: ["scheduler"] }) } })
  const [target, setTarget] = useState("jellyfin")
  const [action, setAction] = useState("restart")
  const [interval, setInterval] = useState("3600")
  const jobs = [...(query.data ?? [])].sort((a, b) => b.id.localeCompare(a.id, undefined, { numeric: true }))

  function addSchedule() {
    const seconds = Number(interval)
    if (!Number.isFinite(seconds) || seconds < 5) return
    createSchedule.mutate({ name: `${action} ${target}`, action, target, intervalSeconds: seconds, enabled: true })
  }
  return <section className="flex min-w-0 flex-1 flex-col gap-6 p-4 md:p-8">
    <header><h1 className="text-2xl font-medium md:text-3xl">Tarefas</h1><p className="mt-1 text-sm text-muted-foreground">Operações em andamento, progresso e histórico recente.</p></header>
    {query.isError && <p role="alert" className="rounded-lg border border-destructive/40 p-4 text-sm">Não foi possível carregar as tarefas.</p>}
    <Card><CardHeader><CardTitle className="text-base">Tarefas agendadas</CardTitle></CardHeader><CardContent className="space-y-4"><div className="grid gap-2 sm:grid-cols-[1fr_1fr_8rem_auto]"><Input value={target} onChange={(e) => setTarget(e.target.value)} placeholder="App alvo" aria-label="App alvo" /><select className="h-9 rounded-md border bg-background px-3 text-sm" value={action} onChange={(e) => setAction(e.target.value)} aria-label="Ação"><option value="restart">Reiniciar</option><option value="start">Iniciar</option><option value="stop">Parar</option></select><Input type="number" min={5} value={interval} onChange={(e) => setInterval(e.target.value)} aria-label="Intervalo em segundos" /><Button onClick={addSchedule} disabled={createSchedule.isPending}>Agendar</Button></div>{schedules.isError && <p className="text-sm text-destructive">Scheduler indisponível.</p>}{(schedules.data ?? []).map((task) => <div key={task.id} className="flex flex-wrap items-center justify-between gap-2 rounded-md border p-3 text-sm"><span><strong>{task.name || `${task.action} ${task.target}`}</strong><span className="ml-2 text-muted-foreground">a cada {task.intervalSeconds}s · {task.lastState ?? "agendada"}</span></span><Button variant="ghost" size="sm" onClick={() => removeSchedule.mutate(task.id)} disabled={removeSchedule.isPending}>Excluir</Button></div>)}</CardContent></Card>
    {query.isPending ? <p role="status">Carregando tarefas…</p> : jobs.length === 0 ? <Card><CardContent className="flex min-h-48 flex-col items-center justify-center gap-3 text-center"><ListTodo className="size-8 text-muted-foreground" /><p className="font-medium">Nenhuma tarefa registrada</p><p className="text-sm text-muted-foreground">As operações de apps e sistema aparecerão aqui.</p></CardContent></Card> : <div className="grid gap-3">{jobs.map((job) => <JobCard key={job.id} job={job} onCancel={() => cancellation.mutate(job.id)} cancelling={cancellation.isPending} />)}</div>}
  </section>
}

function JobCard({ job, onCancel, cancelling }: { job: JobDTO; onCancel: () => void; cancelling: boolean }) {
  const Icon = statusIcon[job.status as keyof typeof statusIcon] ?? CircleDashed
  const progress = Math.max(0, Math.min(100, job.progress ?? (job.status === "success" ? 100 : 0)))
  return <Card><CardHeader className="flex flex-row items-start justify-between gap-3 pb-3"><CardTitle className="flex items-center gap-2 text-base"><Icon className={job.status === "running" ? "animate-spin" : ""} />{job.kind ?? "Operação"}</CardTitle>{job.status === "running" || job.status === "queued" ? <Button variant="ghost" size="icon" aria-label="Cancelar tarefa" disabled={cancelling} onClick={onCancel}><X /></Button> : null}</CardHeader><CardContent className="space-y-3"><div className="flex justify-between gap-3 text-sm"><span className="text-muted-foreground">{job.target ?? job.id}</span><span>{statusLabel[job.status] ?? job.status}</span></div><Progress value={progress} /><p className="text-sm text-muted-foreground">{job.error ?? job.message}</p>{(job.logs?.length ?? 0) > 0 && <details className="rounded-md border bg-muted/30 p-3"><summary className="cursor-pointer text-sm font-medium">Logs da tarefa ({job.logs?.length})</summary><pre className="mt-2 max-h-40 overflow-auto whitespace-pre-wrap text-xs text-muted-foreground">{job.logs?.join("\n")}</pre></details>}</CardContent></Card>
}
