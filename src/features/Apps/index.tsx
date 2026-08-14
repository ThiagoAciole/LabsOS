import { useMemo, useState, type ChangeEvent } from "react"
import { ArrowDownAZ, ArrowUpAZ, SlidersHorizontal } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Separator } from "@/components/ui/separator"
import { apps } from "./data/apps"

type AppType = "all" | "connected" | "notConnected"
const appText = new Map<AppType, string>([["all", "All Apps"], ["connected", "Connected"], ["notConnected", "Not Connected"]])

export function AppsPage() {
  const [searchTerm, setSearchTerm] = useState("")
  const [appType, setAppType] = useState<AppType>(() => {
    const value = new URLSearchParams(window.location.search).get("type")
    return value === "connected" || value === "notConnected" ? value : "all"
  })
  const [sort, setSort] = useState<"asc" | "desc">("asc")
  const filteredApps = useMemo(() => apps.filter((app) => appType === "all" || (appType === "connected" ? app.connected : !app.connected)).filter((app) => app.name.toLowerCase().includes(searchTerm.toLowerCase())).sort((a, b) => sort === "asc" ? a.name.localeCompare(b.name) : b.name.localeCompare(a.name)), [appType, searchTerm, sort])
  const handleSearch = (event: ChangeEvent<HTMLInputElement>) => setSearchTerm(event.target.value)
  const handleTypeChange = (value: AppType | null) => {
    if (!value) return
    setAppType(value)
    const url = new URL(window.location.href)
    if (value === "all") url.searchParams.delete("type")
    else url.searchParams.set("type", value)
    window.history.replaceState(null, "", url)
  }

  return <section className="flex w-full flex-col gap-6 overflow-auto rounded-2xl border bg-background/80 p-4 md:p-8">
    <div><h1 className="text-3xl font-medium tracking-tight">App Integrations</h1><p className="mt-1 text-muted-foreground">Here&apos;s a list of your apps for the integration!</p></div>
    <div className="flex flex-col justify-between gap-4 sm:flex-row sm:items-center">
      <div className="flex flex-col gap-4 sm:flex-row"><Input placeholder="Filter apps..." className="h-9 w-40 lg:w-64" value={searchTerm} onChange={handleSearch} /><Select value={appType} onValueChange={handleTypeChange}><SelectTrigger className="w-36"><SelectValue>{appText.get(appType)}</SelectValue></SelectTrigger><SelectContent><SelectItem value="all">All Apps</SelectItem><SelectItem value="connected">Connected</SelectItem><SelectItem value="notConnected">Not Connected</SelectItem></SelectContent></Select></div>
      <Button variant="outline" size="sm" onClick={() => setSort((current) => current === "asc" ? "desc" : "asc")}>{sort === "asc" ? <ArrowUpAZ /> : <ArrowDownAZ />}<span className="hidden sm:inline">{sort === "asc" ? "Ascending" : "Descending"}</span><SlidersHorizontal className="sm:hidden" /></Button>
    </div><Separator />
    <ul className="grid gap-4 pb-8 md:grid-cols-2 lg:grid-cols-3">{filteredApps.map((app) => { const Icon = app.logo; return <li key={app.name} className="rounded-md border bg-card p-5 transition-colors hover:border-primary/40 hover:bg-primary/5"><div className="mb-8 flex items-center justify-between"><div className="flex size-10 items-center justify-center rounded-md bg-primary/10 p-2 text-primary"><Icon /></div><Button variant="outline" size="sm" className={app.connected ? "border-emerald-400/30 bg-emerald-500/10 text-emerald-400" : "border-red-400/30 bg-red-500/10 text-red-400"}>{app.connected ? "Connected" : "Not Connected"}</Button></div><h2 className="mb-1 font-semibold">{app.name}</h2><p className="line-clamp-2 text-muted-foreground">{app.desc}</p></li> })}</ul>
  </section>
}
