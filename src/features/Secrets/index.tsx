import { useState } from "react"
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { KeyRound, Trash2 } from "lucide-react"
import { getSecrets, putSecret, deleteSecret } from "@/api/secrets"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"

export function SecretsPage() {
  const client = useQueryClient()
  const query = useQuery({ queryKey: ["secrets"], queryFn: getSecrets })
  const [name, setName] = useState("")
  const [value, setValue] = useState("")
  const [notice, setNotice] = useState<string | null>(null)
  const refresh = () => void client.invalidateQueries({ queryKey: ["secrets"] })
  const save = useMutation({ mutationFn: () => putSecret(name.trim(), value), onSuccess: (result) => { setNotice(result.simulated ? "Gravação planejada; nenhum segredo foi alterado." : "Segredo armazenado sem exibir o valor."); setValue(""); refresh() } })
  const remove = useMutation({ mutationFn: deleteSecret, onSuccess: (result) => { setNotice(result.simulated ? "Remoção planejada." : "Segredo removido."); refresh() } })
  return <section className="flex min-w-0 flex-1 flex-col gap-6 p-4 md:p-8"><header><h1 className="text-2xl font-medium md:text-3xl">Segredos</h1><p className="mt-1 text-sm text-muted-foreground">Credenciais armazenadas localmente. O valor nunca é retornado pela API.</p></header>{notice && <p role="status" className="rounded-lg border border-primary/30 p-3 text-sm">{notice}</p>}<Card><CardHeader><CardTitle className="flex items-center gap-2 text-base"><KeyRound className="size-4" />Adicionar ou atualizar</CardTitle></CardHeader><CardContent className="flex flex-col gap-2 sm:flex-row"><Input value={name} onChange={(event) => setName(event.target.value)} placeholder="Nome (ex.: registry-token)" aria-label="Nome do segredo" /><Input type="password" value={value} onChange={(event) => setValue(event.target.value)} placeholder="Valor" aria-label="Valor do segredo" /><Button disabled={!name.trim() || !value || save.isPending} onClick={() => save.mutate()}>Salvar</Button></CardContent></Card>{query.isError && <p role="alert" className="rounded-lg border border-destructive/40 p-4 text-sm">Não foi possível carregar os segredos.</p>}{query.isPending ? <p role="status">Carregando segredos…</p> : <div className="grid gap-3">{(query.data ?? []).map((secret) => <Card key={secret.name}><CardContent className="flex items-center justify-between gap-3 p-4"><div><p className="font-medium">{secret.name}</p><p className="text-xs text-muted-foreground">Atualizado em {new Date(secret.updatedAt).toLocaleString("pt-BR")}</p></div><Button variant="ghost" size="icon" aria-label={`Remover ${secret.name}`} disabled={remove.isPending} onClick={() => remove.mutate(secret.name)}><Trash2 /></Button></CardContent></Card>)}</div>}</section>
}
