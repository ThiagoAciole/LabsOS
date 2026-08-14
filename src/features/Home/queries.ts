import { useQuery } from "@tanstack/react-query"
import { getApps } from "@/api/apps"
import { getEvents } from "@/api/events"
import { getSystemHealth, getSystemSummary, toSystemDashboard } from "@/api/system"

export function useHomeData() {
  return useQuery({
    queryKey: ["home"],
    queryFn: async () => {
      const [summary, health, apps, events] = await Promise.all([getSystemSummary(), getSystemHealth(), getApps(), getEvents()])
      return { summary, health, apps, events, dashboard: toSystemDashboard(summary) }
    },
  })
}
