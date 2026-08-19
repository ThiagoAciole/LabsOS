import { useQuery } from "@tanstack/react-query"
import { ScrollText } from "lucide-react"
import { getAudit } from "@/api/events"
import { Card, CardContent } from "@/components/ui/card"

export function AuditPage() {
  const query = useQuery({ queryKey: ["audit"], queryFn: getAudit })
  return <section className="flex min-w-0 flex-1 flex-col gap-6 p-4 md:p-8"><header><h1 className="text-2xl font-medium md:text-3xl">Auditoria</h1><p className="mt-1 text-sm text-muted-foreground">Histórico persistente de ações administrativas.</p></header>{query.isError && <p role="alert" className="rounded-lg border border-destructive/40 p-4 text-sm">Não foi possível carregar a auditoria.</p>}{query.isPending ? <p role="status">Carregando registros…</p> : (query.data ?? []).length === 0 ? <Card><CardContent className="flex min-h-48 flex-col items-center justify-center gap-3 text-center"><ScrollText className="size-8 text-muted-foreground" /><p>Nenhum registro ainda.</p></CardContent></Card> : <div className="space-y-2">{query.data?.map((entry) => <Card key={entry.id}><CardContent className="flex flex-wrap items-center justify-between gap-3 p-4 text-sm"><div><p className="font-medium">{entry.action}{entry.target ? ` · ${entry.target}` : ""}</p><p className="text-muted-foreground">{entry.details || "Sem detalhes"}</p></div><div className="text-right text-xs text-muted-foreground"><p>{entry.status}</p><p>{new Date(entry.at).toLocaleString("pt-BR")}</p></div></CardContent></Card>)}</div>}</section>
}
