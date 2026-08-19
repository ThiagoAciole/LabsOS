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
import { toSystemDashboard } from "@/api/system";
import { useHomeData } from "./queries";

export function HomePage() {
  const home = useHomeData();
  if (home.isPending) return <div className="flex flex-1 items-center justify-center p-8" role="status">Carregando dados do servidor…</div>;
  const dashboard = home.summary.data ? toSystemDashboard(home.summary.data) : null;
  const health = home.health.data;
  const apps = home.apps.data ?? [];
  const events = home.events.data ?? [];
  return (
    <div className="flex min-w-0 flex-1 flex-col gap-6 p-4 md:p-8">
      <header className="flex items-start justify-between gap-4">
        <Welcome />
        <ClockCard />
      </header>
      {health && <SystemHealth health={health} runningApps={apps.filter((app) => app.status === "running").length} />}
      {dashboard && <section aria-labelledby="server-metrics-title">
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
        {home.metrics.data && <p className="mt-2 text-xs text-muted-foreground">Load average: {home.metrics.data.load.map((value) => value.toFixed(2)).join(" · ")} · coleta {new Date(home.metrics.data.collectedAt).toLocaleTimeString("pt-BR")}</p>}
      </section>}
      <section className="grid items-start gap-6 2xl:grid-cols-[minmax(0,1.45fr)_minmax(20rem,0.85fr)]">
        <ActiveApps apps={apps} />
        <RecentActivity events={events} />
      </section>
    </div>
  );
}
