import { useQuery } from "@tanstack/react-query"
import { getAppHealth, getAppMetrics, getApps, getCatalogApps } from "@/api/apps"
import { toInstalledApp } from "./api"
import type { StoreApp } from "./types"

const categories = new Set<string>(["media", "network", "storage", "automation", "development", "utilities"])

export const appsQueryKey = ["apps"] as const

export function useAppsData() {
  return useQuery({ queryKey: appsQueryKey, queryFn: async () => (await getApps()).map(toInstalledApp) })
}

export function useCatalogData() {
  return useQuery({
    queryKey: ["catalog"],
    queryFn: async () => (await getCatalogApps()).map((app): StoreApp => {
      return { id: app.id, name: app.name, description: app.description, icon: app.icon, version: app.version ?? "Disponível", category: categories.has(app.category ?? "") ? (app.category as StoreApp["category"]) : "utilities", source: app.source ?? "unknown", size: "Tamanho variável", highlights: [], installable: app.installable ?? false, architecture: app.architecture, requirements: app.requirements }
    }),
  })
}
export function useAppHealth(id: string | null) {
  return useQuery({ queryKey: ["app-health", id], queryFn: () => getAppHealth(id as string), enabled: Boolean(id) })
}
export function useAppMetrics(id: string | null) { return useQuery({ queryKey: ["app-metrics", id], queryFn: () => getAppMetrics(id as string), enabled: Boolean(id) }) }
