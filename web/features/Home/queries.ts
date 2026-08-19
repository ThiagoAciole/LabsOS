import { useQuery } from "@tanstack/react-query"
import { getApps } from "@/api/apps"
import { getEvents } from "@/api/events"
import { getSystemHealth, getSystemMetrics, getSystemSummary } from "@/api/system"

export function useHomeData() {
  const summary = useQuery({ queryKey: ["system", "summary"], queryFn: getSystemSummary })
  const health = useQuery({ queryKey: ["system", "health"], queryFn: getSystemHealth })
  const apps = useQuery({ queryKey: ["apps"], queryFn: getApps })
  const events = useQuery({ queryKey: ["events"], queryFn: getEvents })
  const metrics = useQuery({ queryKey: ["system", "metrics"], queryFn: getSystemMetrics })
  return { summary, health, apps, events, metrics, isPending: summary.isPending && health.isPending && apps.isPending && events.isPending }
}
