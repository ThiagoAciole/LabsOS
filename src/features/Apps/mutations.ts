import { useMutation, useQueryClient } from "@tanstack/react-query"
import { removeApp, runAppAction } from "@/api/apps"
import { appsQueryKey } from "./queries"

export function useAppMutation() {
  const client = useQueryClient()
  return useMutation({
    mutationFn: ({ id, action }: { id: string; action: "install" | "start" | "stop" | "restart" | "remove" }) => action === "remove" ? removeApp(id) : runAppAction(id, action),
    onSuccess: async () => { await Promise.all([client.invalidateQueries({ queryKey: appsQueryKey }), client.invalidateQueries({ queryKey: ["catalog"] }), client.invalidateQueries({ queryKey: ["home"] })]) },
  })
}
