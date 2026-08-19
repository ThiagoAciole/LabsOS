import { api, json } from "./client.ts"
import type { JobDTO } from "./system"

export type AppDTO = { id: string; name: string; icon: string; description: string; category?: string; source?: string; status: "running" | "stopped" | "installing" | "error"; health?: string; dependencies?: string[]; volumes?: string[]; ports?: number[]; actions?: string[]; architecture?: string[]; requirements?: string[]; version?: string; url?: string; updateAvailable: boolean; installed: boolean; installable?: boolean }
export const getApps = () => api<AppDTO[]>("/apps")
export const getCatalogApps = () => api<AppDTO[]>("/catalog/apps")
export const importDeclarativeApp = (input: { id: string; name: string; description?: string; category?: string; version?: string; compose: string }) => api<{ app: AppDTO; installed: boolean; message: string }>("/catalog/apps/import", json("POST", input))
export const runAppAction = (id: string, action: "install" | "start" | "stop" | "restart" | "update" | "downgrade") => api<JobDTO>(`/apps/${id}/${action}`, json("POST"))
export const removeApp = (id: string) => api<JobDTO>(`/apps/${id}`, { method: "DELETE" })
export type AppHealthDTO = { id: string; status: string; healthy: boolean; dependencies?: string[]; message: string }
export const getAppHealth = (id: string) => api<AppHealthDTO>(`/apps/${encodeURIComponent(id)}/health`)
export type AppMetricsDTO = { id: string; cpuPercent: number; memoryUsedBytes: number; memoryLimitBytes: number; memoryPercent: number; networkRXBytes: number; networkTXBytes: number; blockReadBytes: number; blockWriteBytes: number; available: boolean; message: string }
export const getAppMetrics = (id: string) => api<AppMetricsDTO>(`/apps/${encodeURIComponent(id)}/metrics`)
