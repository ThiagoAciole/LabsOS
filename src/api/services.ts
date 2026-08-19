import { api, json } from "./client"

export type ServiceExposureDTO = { id: string; name: string; ports?: number[]; url?: string; exposed: boolean; provider?: string; conflict?: boolean }
export const getServices = () => api<ServiceExposureDTO[]>("/services")
export const updateServiceExposure = (id: string, exposed: boolean) => api<{ id: string; exposed: boolean; applied: boolean; planned: boolean; message: string }>(`/services/${encodeURIComponent(id)}/exposure`, json("PUT", { exposed }))
