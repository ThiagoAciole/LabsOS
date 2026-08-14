import { CloudRain, MapPin } from "lucide-react";
import { clockData } from "../../data";
export function ClockCard() {
  return (
    <aside
      className="hidden flex-col items-end gap-1.5 text-right text-sm text-muted-foreground md:flex"
      aria-label="Horário e clima local"
    >
      <div>
        <time
          dateTime="2026-09-16T10:31:00-03:00"
          className="block text-base font-semibold text-foreground tabular-nums"
        >
          {clockData.time}
        </time>
        <span className="block text-sm">{clockData.date}</span>
      </div>
      <span className="flex items-center gap-2 text-sm">
        <MapPin className="size-4" />
        {clockData.city}
        <CloudRain className="ml-1 size-4 text-primary" />
        {clockData.temperature}°
      </span>
    </aside>
  );
}
