import { useEffect, useState, type ReactNode } from "react"

import { AppsPage } from "@/features/Apps"
import { HomePage } from "@/features/Home"
import { SettingsPage } from "@/features/Settings"
import { FilesPage } from "@/features/Files"

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
  return <HomePage />
}

export function Router({ children }: { children: ReactNode }) {
  return children
}
