import type { AppDTO } from "../../api/apps.ts"
import type { InstalledApp } from "./types.ts"

export function toInstalledApp(app: AppDTO): InstalledApp {
  return { id: app.id, name: app.name, description: app.description, icon: app.icon, version: app.version, status: app.status, url: app.url }
}
