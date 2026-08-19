import { Clock3 } from "lucide-react";
export function ClockCard() {
  const now = new Date();
  const time = now.toLocaleTimeString("pt-BR", { hour: "2-digit", minute: "2-digit" });
  const date = now.toLocaleDateString("pt-BR", { weekday: "long", day: "2-digit", month: "short" });
  return (
    <aside
      className="hidden flex-col items-end gap-1.5 text-right text-sm text-muted-foreground md:flex"
      aria-label="Horário e clima local"
    >
      <div>
        <time
          dateTime={now.toISOString()}
          className="block text-base font-semibold text-foreground tabular-nums"
        >
          {time}
        </time>
        <span className="block text-sm capitalize">{date}</span>
      </div>
      <span className="flex items-center gap-2 text-sm">
        <Clock3 className="size-4" />
        Horário local
      </span>
    </aside>
  );
}
