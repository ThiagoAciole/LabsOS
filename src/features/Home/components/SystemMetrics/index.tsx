import { MetricCard } from "@/components/shared/MetricCard"
import { systemMetrics } from "../../data"
export function SystemMetrics() { return <>{systemMetrics.map((metric) => <MetricCard key={metric.id} {...metric} />)}</> }
