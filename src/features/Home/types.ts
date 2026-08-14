import type { LucideIcon } from "lucide-react";
export type AppStatus = "installing" | "running" | "stopped" | "error";
export type DataState = "loading" | "ready" | "stale" | "unavailable" | "error";
export interface DashboardApp {
  id: string;
  name: string;
  url?: string;
  icon: string;
  installed: boolean;
  status: AppStatus;
}
export interface SystemMetric {
  id: string;
  title: string;
  value: number;
  description: string;
  icon: LucideIcon;
}
export interface NetworkData {
  download: number;
  upload: number;
  ip: string;
  online: boolean;
}
export interface Activity {
  id: string;
  title: string;
  description: string;
  time: string;
  dateTime: string;
  icon: LucideIcon;
}
export interface ClockData {
  city: string;
  temperature: number;
  condition: string;
  time: string;
  date: string;
}
export interface UptimeData {
  days: number;
  hours: number;
  minutes: number;
}
export interface SystemHealth {
  state: DataState;
  title: string;
  description: string;
  updatedLabel: string;
}
