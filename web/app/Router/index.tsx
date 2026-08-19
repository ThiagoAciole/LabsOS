import { useEffect, useState, type ReactNode } from "react"

import { AppsPage } from "@/features/Apps"
import { HomePage } from "@/features/Home"
import { SettingsPage } from "@/features/Settings"
import { FilesPage } from "@/features/Files"
import { NotificationsPage } from "@/features/Notifications"
import { LogsPage } from "@/features/Logs"
import { JobsPage } from "@/features/Jobs"
import { AuditPage } from "@/features/Audit"
import { BackupsPage } from "@/features/Backups"
import { ServicesPage } from "@/features/Services"
import { SecretsPage } from "@/features/Secrets"

export function Outlet() {
  const [path, setPath] = useState(window.location.pathname)

  useEffect(() => {
    const handlePopState = () => setPath(window.location.pathname)
    window.addEventListener("popstate", handlePopState)
    return () => window.removeEventListener("popstate", handlePopState)
  }, [])

  if (path === "/apps") return <AppsPage />
  if (path === "/settings") return <SettingsPage />
  if (path === "/files") return <FilesPage />
  if (path === "/notifications") return <NotificationsPage />
  if (path === "/logs") return <LogsPage />
  if (path === "/jobs") return <JobsPage />
  if (path === "/audit") return <AuditPage />
  if (path === "/backups") return <BackupsPage />
  if (path === "/services") return <ServicesPage />
  if (path === "/secrets") return <SecretsPage />
  return <HomePage />
}

export function Router({ children }: { children: ReactNode }) {
  return children
}
