import { api } from "./client"
export type EventDTO = { id: string; type: string; message: string }
export const getEvents = () => api<EventDTO[]>("/events")
