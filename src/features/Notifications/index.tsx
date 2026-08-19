import { useEffect } from "react"
import { Bell, CheckCircle2, CircleAlert, Info, Trash2, TriangleAlert } from "lucide-react"
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { deleteNotification, getNotifications, markNotificationRead, subscribeToEvents } from "@/api/events"

const icons = { success: CheckCircle2, error: CircleAlert, warning: TriangleAlert, info: Info }
export function NotificationsPage() {
  const client = useQueryClient()
  const query = useQuery({ queryKey: ["notifications"], queryFn: getNotifications })
  const read = useMutation({ mutationFn: markNotificationRead, onSuccess: () => client.invalidateQueries({ queryKey: ["notifications"] }) })
  const remove = useMutation({ mutationFn: deleteNotification, onSuccess: () => client.invalidateQueries({ queryKey: ["notifications"] }) })
  useEffect(() => subscribeToEvents(() => void client.invalidateQueries({ queryKey: ["notifications"] })), [client])
  if (query.isPending) return <section className="flex flex-1 items-center justify-center p-8">Carregando notificações…</section>
  return <section className="flex min-w-0 flex-1 flex-col gap-6 p-4 md:p-8"><header><h1 className="text-2xl font-semibold">Notificações</h1><p className="mt-1 text-sm text-muted-foreground">Atividade persistente do sistema e dos Apps.</p></header><div className="space-y-3">{query.data?.length ? query.data.map(item => { const Icon = icons[item.severity as keyof typeof icons] ?? Bell; return <article key={item.id} className={`flex gap-3 rounded-lg border p-4 ${item.read ? "opacity-60" : ""}`}><Icon className="mt-0.5 size-5 text-primary" /><div className="min-w-0 flex-1"><h2 className="font-medium">{item.title}</h2><p className="text-sm text-muted-foreground">{item.message}</p><p className="mt-1 text-xs text-muted-foreground">{item.source} · {new Date(item.createdAt).toLocaleString("pt-BR")}</p></div><div className="flex items-center gap-3">{!item.read && <button className="text-xs text-primary" onClick={() => read.mutate(item.id)}>Marcar como lida</button>}<button className="text-muted-foreground hover:text-destructive" aria-label={`Excluir ${item.title}`} onClick={() => remove.mutate(item.id)}><Trash2 className="size-4" /></button></div></article> }) : <div className="rounded-lg border border-dashed p-8 text-center text-sm text-muted-foreground">Nenhuma notificação.</div>}</div></section>
}
