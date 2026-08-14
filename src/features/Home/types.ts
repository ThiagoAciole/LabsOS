import type { LucideIcon } from "lucide-react"
export type AppStatus = "running" | "stopped"
export interface DashboardApp { id: string; name: string; url: string; icon: string; installed: boolean; status: AppStatus }
export interface SystemMetric { id: string; title: string; value: number; description: string; icon: LucideIcon }
export interface NetworkData { download: number; upload: number; ip: string; online: boolean }
export interface Activity { id: string; title: string; description: string; time: string; icon: LucideIcon }
export interface ClockData { city: string; temperature: number; condition: string; time: string; date: string }
export interface UptimeData { days: number; hours: number; minutes: number }
