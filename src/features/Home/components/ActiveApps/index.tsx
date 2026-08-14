import { AppWindow, ArrowRight } from "lucide-react"
import { Button } from "@/components/ui/button"
import { apps } from "../../data"
import { AppCard } from "./AppCard"
export function ActiveApps() { const activeApps = apps.filter((app) => app.installed && app.status === "running").slice(0, 4); return <section className="min-w-0 space-y-3"><header className="flex items-center gap-3"><AppWindow className="size-4 text-primary" /><div><h2 className="text-sm font-medium">Apps ativos</h2><p className="text-xs text-muted-foreground">{activeApps.length} em execução</p></div><Button variant="ghost" size="sm" className="ml-auto">Ver todos <ArrowRight /></Button></header><div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-4">{activeApps.map((app) => <AppCard key={app.id} app={app} />)}</div></section> }
