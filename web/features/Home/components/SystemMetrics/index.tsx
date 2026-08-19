import { MetricCard } from "@/components/shared/MetricCard"
import { Cpu, HardDrive, MemoryStick } from "lucide-react"
import type { SystemMetric } from "../../types"
const icons = { memory: MemoryStick, disk: HardDrive, cpu: Cpu }
export function SystemMetrics({ metrics }: { metrics: Omit<SystemMetric, "icon">[] }) { return <>{metrics.map((metric) => <MetricCard key={metric.id} {...metric} icon={icons[metric.id as keyof typeof icons]} />)}</> }
