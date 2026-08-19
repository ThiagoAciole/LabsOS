import { api, json } from "./client"

export type InstallerStatus = {
  status: "needs-install" | "running" | "complete" | "cancelled" | "cancelling"
  job?: InstallerJob
}

export type InstallerDisk = {
  id: string
  name: string
  path: string
  sizeBytes: number
  size: string
  kind: string
  ready: boolean
  hasLabsOSData: boolean
}

export type InstallerRequest = {
  operation: "fresh" | "restore" | "transfer"
  disk: string
  preserve: boolean
  serverName: string
  password: string
}

export type InstallerJob = {
  id: string
  status: "running" | "complete" | "cancelled" | "cancelling" | "error"
  progress: number
  message: string
  error?: string
  simulated: boolean
}

export const installerApi = {
  status: () => api<InstallerStatus>("/installer/status"),
  disks: () => api<InstallerDisk[]>("/installer/disks"),
  validate: (request: InstallerRequest) => api<{ valid: boolean; destructive: boolean; disk: string }>("/installer/validate", json("POST", request)),
  start: (request: InstallerRequest) => api<InstallerJob>("/installer/start", json("POST", request)),
  cancel: (id: string) => api<{ cancelled: boolean }>(`/installer/jobs/${encodeURIComponent(id)}/cancel`, json("POST")),
  job: (id: string) => api<InstallerJob>(`/installer/jobs/${id}`),
  events: () => api<Array<{ type: string; jobId: string; progress: number; message: string }>>("/installer/events"),
  reboot: () => api<{ accepted: boolean; message: string }>("/installer/reboot", json("POST")),
}
