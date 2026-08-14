import { api, json } from "./client.ts"
import type { JobDTO } from "./system"

export type AppDTO = { id: string; name: string; icon: string; description: string; status: "running" | "stopped" | "installing" | "error"; version?: string; url?: string; updateAvailable: boolean; installed: boolean }
export const getApps = () => api<AppDTO[]>("/apps")
export const getCatalogApps = () => api<AppDTO[]>("/catalog/apps")
export const runAppAction = (id: string, action: "install" | "start" | "stop" | "restart") => api<JobDTO>(`/apps/${id}/${action}`, json("POST"))
export const removeApp = (id: string) => api<JobDTO>(`/apps/${id}`, { method: "DELETE" })
