import { Timer } from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import type { UptimeData } from "../../types";
export function UptimeCard({ uptime: uptimeData }: { uptime: UptimeData }) {
  return (
    <Card className="flex min-h-28 flex-col">
      <CardHeader className="p-3 pb-2">
        <CardTitle className="flex items-center gap-2 text-sm font-medium">
          <Timer className="size-4 text-primary" />
          Tempo ligado
        </CardTitle>
      </CardHeader>
      <CardContent className="flex flex-1 items-end justify-between gap-3 p-3 pt-0">
        <div>
          <strong className="block whitespace-nowrap text-xl font-medium tabular-nums sm:text-2xl">
            {uptimeData.days} dias e {uptimeData.hours} horas
          </strong>
          <p className="mt-1 text-xs text-muted-foreground">
            {uptimeData.minutes} minutos adicionais
          </p>
        </div>
        <div
          className="hidden h-3 w-1/3 items-end gap-1 sm:flex"
          aria-hidden="true"
        >
          {[3, 5, 4, 7, 6, 9, 8, 11, 10, 12, 9, 13, 11, 14, 12, 15, 13, 16].map(
            (height, index) => (
              <span
                key={index}
                className="min-w-0 flex-1 rounded-sm bg-primary"
                style={{
                  height: `${height}px`,
                  opacity: 0.45 + (index / 18) * 0.4,
                }}
              />
            ),
          )}
        </div>
      </CardContent>
    </Card>
  );
}
