import { api, json } from "./client.ts"

export type SystemSummaryDTO = { hostname: string; status: string; uptimeSeconds: number; version: string; cpuUsage: number; memoryUsedBytes: number; memoryTotalBytes: number; temperatureCelsius: number; storageUsedBytes: number; storageTotalBytes: number; ipAddress?: string; networkOnline?: boolean; networkDownloadBytesPerSecond?: number; networkUploadBytesPerSecond?: number }
export type SystemHealthDTO = { status: string; components: Record<string, string> }
export type JobDTO = { id: string; kind?: string; target?: string; status: "queued" | "running" | "success" | "error" | "cancelled" | string; progress?: number; message: string; error?: string; logs?: string[] }
export type UpdateStatusDTO = { currentVersion: string; latestVersion: string; updateAvailable: boolean; channel?: string; changelog?: string; sha256?: string; applied?: boolean; simulated?: boolean; message?: string }

const percent = (used: number, total: number) => total > 0 ? Math.round((used / total) * 100) : 0
const gib = (bytes: number) => `${(bytes / 1024 ** 3).toLocaleString("pt-BR", { maximumFractionDigits: 1 })} GB`

export function toSystemDashboard(dto: SystemSummaryDTO) {
  return {
    metrics: [
      { id: "memory", title: "Uso de memória", value: percent(dto.memoryUsedBytes, dto.memoryTotalBytes), description: `${gib(dto.memoryUsedBytes)} de ${gib(dto.memoryTotalBytes)}` },
      { id: "disk", title: "Uso do disco", value: percent(dto.storageUsedBytes, dto.storageTotalBytes), description: `${gib(dto.storageUsedBytes)} de ${gib(dto.storageTotalBytes)}` },
      { id: "cpu", title: "Uso de CPU", value: Math.round(dto.cpuUsage), description: dto.temperatureCelsius > 0 ? `Temperatura: ${dto.temperatureCelsius} °C` : "Temperatura não disponível" },
    ],
    uptime: { days: Math.floor(dto.uptimeSeconds / 86400), hours: Math.floor(dto.uptimeSeconds % 86400 / 3600), minutes: Math.floor(dto.uptimeSeconds % 3600 / 60) },
    network: { download: Math.round((dto.networkDownloadBytesPerSecond ?? 0) / 100_000) / 10, upload: Math.round((dto.networkUploadBytesPerSecond ?? 0) / 100_000) / 10, ip: dto.ipAddress ?? "Não disponível", online: dto.networkOnline ?? false },
  }
}

export const getSystemSummary = () => api<SystemSummaryDTO>("/system/summary")
export const getSystemHealth = () => api<SystemHealthDTO>("/system/health")
export type MetricsDTO = { collectedAt: string; cpuPercent: number; memory: { used: number; total: number; percent: number }; storage: { used: number; total: number; percent: number }; load: number[]; networkBytes: Record<string, number> }
export const getSystemMetrics = () => api<MetricsDTO>("/system/metrics")
export const runPowerAction = (action: "reboot" | "shutdown") => api<JobDTO>(`/system/${action}`, json("POST"))
export const getSystemUpdate = () => api<UpdateStatusDTO>("/system/update")
export const applySystemUpdate = () => api<UpdateStatusDTO>("/system/update", json("POST"))
export const rollbackSystemUpdate = () => api<{ applied: boolean; simulated?: boolean; status?: UpdateStatusDTO; message?: string }>("/system/rollback", json("POST"))
export const getJob = (id: string) => api<JobDTO>(`/jobs/${id}`)
export const getJobs = () => api<JobDTO[]>("/jobs")
export type ScheduledTaskDTO = { id: string; name?: string; action: string; target: string; intervalSeconds: number; enabled: boolean; lastRun?: string; nextRun: string; lastState?: string }
export const getScheduledTasks = () => api<ScheduledTaskDTO[]>("/scheduler/tasks")
export const createScheduledTask = (task: Pick<ScheduledTaskDTO, "name" | "action" | "target" | "intervalSeconds" | "enabled">) => api<ScheduledTaskDTO>("/scheduler/tasks", json("POST", task))
export const deleteScheduledTask = (id: string) => api<{ deleted: boolean; id: string }>(`/scheduler/tasks/${encodeURIComponent(id)}`, { method: "DELETE" })
export const cancelJob = (id: string) => api<{ cancelled: boolean; id: string }>(`/jobs/${id}/cancel`, json("POST"))
export type DiagnosticCheck = { id: string; status: string; message: string }
export const getDiagnostics = () => api<{ checks: DiagnosticCheck[]; healthy: boolean }>("/diagnostics")
