import { useQuery } from "@tanstack/react-query"
import { getApps, getCatalogApps } from "@/api/apps"
import { storeApps as editorialApps } from "./data"
import { toInstalledApp } from "./api"
import type { StoreApp } from "./types"

export const appsQueryKey = ["apps"] as const

export function useAppsData() {
  return useQuery({ queryKey: appsQueryKey, queryFn: async () => (await getApps()).map(toInstalledApp) })
}

export function useCatalogData() {
  return useQuery({
    queryKey: ["catalog"],
    queryFn: async () => (await getCatalogApps()).map((app): StoreApp => {
      const metadata = editorialApps.find((item) => item.id === app.id)
      return { id: app.id, name: app.name, description: app.description, icon: app.icon, version: app.version ?? "Disponível", category: metadata?.category ?? "utilities", source: metadata?.source ?? "official", size: metadata?.size ?? "Tamanho variável", highlights: metadata?.highlights ?? [] }
    }),
  })
}
