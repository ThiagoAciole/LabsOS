import { api } from "./client"
export type EventDTO = { id: string; type: string; message: string }
export const getEvents = () => api<EventDTO[]>("/events")
export type NotificationDTO = { id: string; type: string; title: string; message: string; source: string; severity: string; createdAt: string; read: boolean }
export const getNotifications = () => api<NotificationDTO[]>("/notifications")
export const markNotificationRead = (id: string) => api<{ read: boolean }>(`/notifications/${id}/read`, { method: "POST" })
export const deleteNotification = (id: string) => api<{ deleted: boolean }>(`/notifications/${id}`, { method: "DELETE" })
export const getSystemLogs = (source: "system" | "kernel", lines = 100) => api<{ source: string; logs: string }>(`/logs/${source}?lines=${lines}`)
export const getAppLogs = (id: string, lines = 100) => api<{ source: string; app: string; logs: string }>(`/apps/${encodeURIComponent(id)}/logs?lines=${lines}`)
export type AuditEntryDTO = { id: string; at: string; actor: string; action: string; target?: string; status: string; details?: string }
export const getAudit = () => api<AuditEntryDTO[]>("/audit")
export function subscribeToEvents(onEvents: (events: EventDTO[]) => void) {
  const source = new EventSource("/api/v1/events/stream")
  source.addEventListener("snapshot", (event) => onEvents(JSON.parse((event as MessageEvent).data) as EventDTO[]))
  source.addEventListener("events", (event) => onEvents(JSON.parse((event as MessageEvent).data) as EventDTO[]))
  return () => source.close()
}
