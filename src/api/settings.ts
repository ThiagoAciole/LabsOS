import { api, json } from "./client"

export type SettingsDTO = { hostname?: string; timezone?: string; language?: string; dhcp?: boolean; remoteAccess?: boolean; [key: string]: unknown }
export const getSettings = () => api<SettingsDTO>("/settings")
export const updateSystemSettings = (patch: Pick<SettingsDTO, "hostname" | "language">) => api<SettingsDTO>("/settings/system", json("PUT", patch))
export const updateNetworkSettings = (patch: Pick<SettingsDTO, "remoteAccess">) => api<SettingsDTO>("/settings/network", json("PUT", patch))
