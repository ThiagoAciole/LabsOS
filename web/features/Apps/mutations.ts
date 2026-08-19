import { useMutation, useQueryClient } from "@tanstack/react-query"
import { removeApp, runAppAction } from "@/api/apps"
import { appsQueryKey } from "./queries"
import { getJob } from "@/api/system"

async function waitForJob(id: string) {
  const initial = await getJob(id)
  if (["success", "error", "cancelled"].includes(initial.status)) return initial
  return new Promise<typeof initial>((resolve, reject) => {
    const source = new EventSource("/api/v1/events/stream")
    let settled = false
    const finish = (callback: () => void) => {
      if (settled) return
      settled = true
      window.clearTimeout(timeout)
      source.close()
      callback()
    }
    const timeout = window.setTimeout(() => finish(() => reject(new Error("Tempo limite aguardando o job"))), 120_000)
    source.addEventListener("events", (event) => {
      try {
        const items = JSON.parse((event as MessageEvent<string>).data) as Array<{ type?: string; message?: string }>
        if (!items.some((item) => item.type === "job" && item.message?.startsWith(id))) return
        void getJob(id).then((job) => {
          if (["success", "error", "cancelled"].includes(job.status)) finish(() => resolve(job))
        }).catch((error: unknown) => finish(() => reject(error)))
      } catch {
        // Ignore malformed events; the stream remains the source of truth.
      }
    })
    source.onerror = () => finish(() => reject(new Error("Conexão de eventos encerrada")))
  })
}

export function useAppMutation() {
  const client = useQueryClient()
  return useMutation({
    mutationFn: async ({ id, action }: { id: string; action: "install" | "start" | "stop" | "restart" | "update" | "downgrade" | "remove" }) => { const job = action === "remove" ? await removeApp(id) : await runAppAction(id, action); return waitForJob(job.id) },
    onSuccess: async () => { await Promise.all([client.invalidateQueries({ queryKey: appsQueryKey }), client.invalidateQueries({ queryKey: ["catalog"] }), client.invalidateQueries({ queryKey: ["home"] })]) },
  })
}
