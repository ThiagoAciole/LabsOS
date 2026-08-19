import { useEffect } from "react"
import { useQueryClient } from "@tanstack/react-query"
import { subscribeToEvents } from "@/api/events"

const liveQueries = ["jobs", "scheduler", "apps", "app-health", "app-metrics", "services", "audit", "notifications", "events", "settings", "network-snapshot", "wifi-snapshot", "storage-snapshot", "ssh-status", "system"]

export function EventInvalidationBridge() {
  const client = useQueryClient()
  useEffect(() => subscribeToEvents(() => {
    for (const queryKey of liveQueries) void client.invalidateQueries({ queryKey: [queryKey] })
  }), [client])
  return null
}
