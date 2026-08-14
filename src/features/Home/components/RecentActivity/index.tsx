import { ArrowRight, Clock3 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import { Settings } from "lucide-react";
import type { EventDTO } from "@/api/events";
import { ActivityItem } from "./ActivityItem";
export function RecentActivity({ events }: { events: EventDTO[] }) {
  const activities = events.map((event) => ({ id: event.id, icon: Settings, title: event.message, description: event.type === "job" ? "Operação concluída" : "Evento do sistema", time: "recente", dateTime: "" }));
  return (
    <section aria-labelledby="recent-activity-title">
      <Card>
        <CardHeader>
          <div className="flex items-start gap-3">
            <Clock3 className="mt-0.5 size-4 text-primary" />
            <div>
              <h2 id="recent-activity-title" className="text-sm font-medium">
                Atividade recente
              </h2>
              <p className="mt-1 text-xs text-muted-foreground">
                Atualizações e eventos do sistema
              </p>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          {activities.map((activity, index) => (
            <div key={activity.id}>
              {index > 0 && <Separator />}
              <ActivityItem activity={activity} />
            </div>
          ))}
          <Button
            variant="secondary"
            className="mt-4 min-h-11 w-full justify-between"
            disabled
          >
            Ver todas as atividades <ArrowRight />
          </Button>
        </CardContent>
      </Card>
    </section>
  );
}
