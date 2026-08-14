import { ArrowRight, Clock3 } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Separator } from "@/components/ui/separator"
import { activities } from "../../data"
import { ActivityItem } from "./ActivityItem"
export function RecentActivity() { return <Card><CardHeader><CardTitle className="flex items-center gap-3"><Clock3 className="size-4 text-primary" /><span><p className="text-sm font-medium">Atividade recente</p><p className="text-xs font-normal text-muted-foreground">Atualizações e eventos do sistema</p></span></CardTitle></CardHeader><CardContent>{activities.map((activity, index) => <div key={activity.id}>{index > 0 && <Separator />}<ActivityItem activity={activity} /></div>)}<Button variant="secondary" className="mt-4 w-full justify-between">Ver todas as atividades <ArrowRight /></Button></CardContent></Card> }
