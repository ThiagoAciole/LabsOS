import { useState, type FormEvent, type ReactNode } from "react"
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { LockKeyhole } from "lucide-react"
import { getAuthStatus, login } from "@/api/settings"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"

export function AuthGate({ children }: { children: ReactNode }) {
  const client = useQueryClient()
  const status = useQuery({ queryKey: ["auth-status"], queryFn: getAuthStatus, retry: false })
  const [password, setPassword] = useState("")
  const auth = useMutation({ mutationFn: login, onSuccess: () => { setPassword(""); void client.invalidateQueries({ queryKey: ["auth-status"] }) } })
  if (status.isPending) return <div className="flex min-h-screen items-center justify-center text-sm text-muted-foreground">Verificando sessão…</div>
  if (status.isError) return <div className="flex min-h-screen items-center justify-center p-6"><p role="alert">API indisponível. Verifique se o Labs API está em execução.</p></div>
  if (!status.data.configured || status.data.authenticated) return <>{children}</>
  function submit(event: FormEvent) { event.preventDefault(); if (password) auth.mutate(password) }
  return <main className="flex min-h-screen items-center justify-center bg-background p-6"><Card className="w-full max-w-sm"><CardHeader><div className="mb-2 flex size-10 items-center justify-center rounded-lg bg-primary/10 text-primary"><LockKeyhole /></div><CardTitle>Entrar no LabsOS</CardTitle><p className="text-sm text-muted-foreground">Informe a senha administrativa para acessar o painel.</p></CardHeader><CardContent><form className="space-y-4" onSubmit={submit}><Input type="password" value={password} onChange={(event) => setPassword(event.target.value)} placeholder="Senha" autoFocus aria-label="Senha administrativa" />{auth.isError && <p className="text-sm text-destructive" role="alert">Senha inválida.</p>}<Button className="w-full" type="submit" disabled={!password || auth.isPending}>{auth.isPending ? "Entrando…" : "Entrar"}</Button></form></CardContent></Card></main>
}
