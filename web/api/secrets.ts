import { api, json } from "./client"

export type SecretDTO = { name: string; updatedAt: string }
export const getSecrets = () => api<SecretDTO[]>("/secrets")
export const putSecret = (name: string, value: string) => api<{ name: string; stored: boolean; simulated?: boolean; message?: string }>(`/secrets/${encodeURIComponent(name)}`, json("PUT", { value }))
export const deleteSecret = (name: string) => api<{ name: string; deleted: boolean; simulated?: boolean }>(`/secrets/${encodeURIComponent(name)}`, { method: "DELETE" })
