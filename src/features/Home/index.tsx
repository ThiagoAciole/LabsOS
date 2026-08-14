import {
  ActiveApps,
  ClockCard,
  NetworkCard,
  RecentActivity,
  SystemHealth,
  SystemMetrics,
  UptimeCard,
  Welcome,
} from "./components";
import { Button } from "@/components/ui/button";
import { useHomeData } from "./queries";

export function HomePage() {
  const home = useHomeData();
  if (home.isPending) return <div className="flex flex-1 items-center justify-center p-8" role="status">Carregando dados do servidor…</div>;
  if (home.isError) return <div className="flex flex-1 flex-col items-center justify-center gap-3 p-8 text-center" role="alert"><p>Não foi possível conectar ao servidor.</p><Button variant="outline" onClick={() => void home.refetch()}>Tentar novamente</Button></div>;
  const { dashboard, health, apps, events } = home.data;
  return (
    <div className="flex min-w-0 flex-1 flex-col gap-6 p-4 md:p-8">
      <header className="flex items-start justify-between gap-4">
        <Welcome />
        <ClockCard />
      </header>
      <SystemHealth health={health} runningApps={apps.filter((app) => app.status === "running").length} />
      <section aria-labelledby="server-metrics-title">
        <header className="mb-3">
          <h2 id="server-metrics-title" className="text-sm font-medium">
            Visão do servidor
          </h2>
          <p className="mt-1 text-xs text-muted-foreground">
            Recursos e conectividade
          </p>
        </header>
        <div className="grid grid-cols-[repeat(auto-fit,minmax(16rem,1fr))] gap-3">
          <SystemMetrics metrics={dashboard.metrics} />
          <UptimeCard uptime={dashboard.uptime} />
          <NetworkCard network={dashboard.network} />
        </div>
      </section>
      <section className="grid items-start gap-6 2xl:grid-cols-[minmax(0,1.45fr)_minmax(20rem,0.85fr)]">
        <ActiveApps apps={apps} />
        <RecentActivity events={events} />
      </section>
    </div>
  );
}
