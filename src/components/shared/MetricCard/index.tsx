import type { LucideIcon } from "lucide-react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Progress } from "@/components/ui/progress"
export function MetricCard({ title, value, description, icon: Icon }: { title: string; value: number; description: string; icon: LucideIcon }) {
  return <Card className="flex min-h-36 flex-col"><CardHeader className="pb-3"><CardTitle className="flex items-center gap-2 text-sm font-medium"><span className="flex size-8 items-center justify-center rounded-md bg-primary/10 text-primary"><Icon className="size-4" /></span>{title}</CardTitle></CardHeader><CardContent className="flex flex-1 flex-col justify-end gap-3"><strong className="text-2xl font-medium tabular-nums">{value}%</strong><Progress value={value} className="h-1" /><p className="text-xs text-muted-foreground">{description}</p></CardContent></Card>
}
