import { api, json } from "./client"

export type SettingsDTO = { hostname?: string; timezone?: string; language?: string; dhcp?: boolean; remoteAccess?: boolean; [key: string]: unknown }
export const getSettings = () => api<SettingsDTO>("/settings")
export const updateSystemSettings = (patch: Pick<SettingsDTO, "hostname" | "language">) => api<SettingsDTO>("/settings/system", json("PUT", patch))
export const updateNetworkSettings = (patch: Pick<SettingsDTO, "remoteAccess">) => api<SettingsDTO>("/settings/network", json("PUT", patch))
export type NetworkInterface = { ifname?: string; operstate?: string; addr_info?: Array<{ local?: string }> }
export type NetworkSnapshot = { interfaces: NetworkInterface[]; routes: string; dns: string[]; simulated: boolean }
export const getNetworkSnapshot = () => api<NetworkSnapshot>("/network")
export type WiFiSnapshot = { available: boolean; devices: string[]; networks: string[]; current?: string; simulated: boolean }
export const getWiFiSnapshot = () => api<WiFiSnapshot>("/network/wifi")
export const updateWiFi = (patch: { action: "connect" | "reconnect" | "forget"; ssid: string; password?: string; device?: string }) => api<{ applied: boolean; simulated: boolean; message: string }>("/network/wifi", json("PUT", patch))
export type StorageDevice = { name: string; path?: string; size: number; fstype?: string; uuid?: string; mountpoints?: string[]; type: string; readOnly?: boolean; transport?: string; model?: string; children?: StorageDevice[] }
export type StorageSnapshot = { devices: StorageDevice[]; usage: Array<{ path: string; usedBytes: number; availableBytes: number; totalBytes: number; usePercent: number }>; dataMounted?: boolean; readOnly: boolean; message: string }
export const getStorageSnapshot = () => api<StorageSnapshot>("/storage")
export type FileEntry = { name: string; path: string; type: "file" | "directory"; size: number; modifiedAt: string }
export type FilesSnapshot = { root: string; path: string; entries: FileEntry[]; readOnly: boolean }
export const getFiles = (path = ".") => api<FilesSnapshot>(`/files?path=${encodeURIComponent(path)}`)
export const updateNetwork = (patch: Record<string, unknown>) => api<{ applied: boolean; simulated: boolean; message: string }>("/network", json("PUT", patch))
export type SSHStatus = { enabled: boolean; active: boolean; port: number; keys: Array<{ fingerprint: string; comment?: string }>; simulated: boolean }
export const getSSHStatus = () => api<SSHStatus>("/access/ssh")
export const updateSSH = (patch: Record<string, unknown>) => api<{ applied: boolean; simulated: boolean; message: string }>("/access/ssh", json("PUT", patch))
export type AuthStatus = { authenticated: boolean; configured: boolean; localDevelopmentBypass: boolean }
export const getAuthStatus = () => api<AuthStatus>("/auth/status")
export const login = (password: string) => api<{ authenticated: boolean }>("/auth/login", json("POST", { password }))
export const logout = () => api<{ authenticated: boolean }>("/auth/logout", json("POST"))
export type SessionDTO = { id: string; createdAt: string; lastSeen: string; userAgent?: string; current: boolean }
export const getSessions = () => api<SessionDTO[]>("/auth/sessions")
export const revokeSession = (id: string) => api<{ revoked: boolean }>(`/auth/sessions/${encodeURIComponent(id)}`, { method: "DELETE" })
export const changePassword = (current: string, next: string) => api<{ changed: boolean; simulated: boolean; message?: string }>("/auth/password", json("POST", { current, new: next }))
export type BackupRecord = { id: string; path: string; size: number; sha256?: string; createdAt: string; simulated?: boolean }
export const getBackups = () => api<BackupRecord[]>("/backups")
export const createBackup = (apps: string[]) => api<{ id: string; status?: string; simulated?: boolean; message?: string }>("/backups", json("POST", { apps }))
export const verifyBackup = (id: string) => api<{ id: string; size: number; sha256: string; integrity: string }>(`/backups/${encodeURIComponent(id)}/verify`)
export const restoreBackup = (id: string) => api<{ id: string; status: string; simulated: boolean; entries?: string[]; message?: string }>(`/backups/${encodeURIComponent(id)}/restore`, json("POST"))
export const deleteBackup = (id: string) => api<{ id: string; status: string; simulated: boolean }>(`/backups/${encodeURIComponent(id)}`, { method: "DELETE" })
