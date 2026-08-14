import type { LucideIcon } from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Progress } from "@/components/ui/progress";
export function MetricCard({
  title,
  value,
  description,
  icon: Icon,
}: {
  title: string;
  value: number;
  description: string;
  icon: LucideIcon;
}) {
  return (
    <Card className="flex min-h-28 flex-col">
      <CardHeader className="p-3 pb-2">
        <CardTitle className="flex items-center gap-2 text-sm font-medium">
          <span className="flex size-8 items-center justify-center rounded-md bg-primary/10 text-primary">
            <Icon className="size-4" />
          </span>
          {title}
        </CardTitle>
      </CardHeader>
      <CardContent className="flex flex-1 items-end gap-3 p-3 pt-0">
        <div className="shrink-0">
          <strong className="block text-2xl font-medium tabular-nums">
            {value}%
          </strong>
          <p className="mt-1 text-xs text-muted-foreground">{description}</p>
        </div>
        <Progress value={value} className="mb-1 h-1 flex-1" />
      </CardContent>
    </Card>
  );
}
