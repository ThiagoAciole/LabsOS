import { Cpu, Download, HardDrive, MemoryStick, Settings, ShieldCheck, Upload } from "lucide-react"
import type { Activity, ClockData, DashboardApp, NetworkData, SystemMetric, UptimeData } from "./types"
export const clockData: ClockData = { city: "Campina Grande", temperature: 22, condition: "Chuva", time: "10:31", date: "Sexta-feira, 16 Set 2026" }
export const uptimeData: UptimeData = { days: 2, hours: 14, minutes: 32 }
export const systemMetrics: SystemMetric[] = [
  { id: "memory", title: "RAM Usage", value: 40, description: "6.6 GB / 16 GB", icon: MemoryStick },
  { id: "disk", title: "Disk Usage", value: 29, description: "128 GB / 447 GB", icon: HardDrive },
  { id: "cpu", title: "CPU Usage", value: 17, description: "Temp 42°", icon: Cpu },
]
export const networkData: NetworkData = { download: 12, upload: 38, ip: "192.168.0.10", online: true }
export const apps: DashboardApp[] = [
  ["jellyfin", "Jellyfin", "jellyfin.labs.local", "/app-icons/jellyfin.svg"], ["immich", "Immich", "immich.labs.local", "/app-icons/immich.svg"], ["home-assistant", "Home Assistant", "home-assistant.labs.local", "/app-icons/home-assistant.svg"], ["adguard", "AdGuard Home", "adguard.labs.local", "/app-icons/adguard-home.svg"],
].map(([id, name, url, icon]) => ({ id, name, url, icon, installed: true, status: "running" as const }))
export const activities: Activity[] = [
  { id: "adguard", icon: ShieldCheck, title: "AdGuard atualizado", description: "Versão 0.107.52 instalada", time: "há 2 min" },
  { id: "backup", icon: Upload, title: "Backup concluído", description: "Arquivos salvos no disco USB", time: "há 1 h" },
  { id: "jellyfin", icon: Download, title: "Biblioteca do Jellyfin atualizada", description: "12 novos itens adicionados", time: "há 3 h" },
  { id: "check", icon: Settings, title: "Verificação do sistema concluída", description: "Todos os serviços estão saudáveis", time: "há 5 h" },
]
