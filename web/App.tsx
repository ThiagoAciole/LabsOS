import { Outlet, Router } from "@/app/Router"
import { RootRoute } from "@/routes/Root"
import { InstallerPage } from "@/features/Installer"
import { AuthGate } from "@/features/AuthGate"
import { EventInvalidationBridge } from "@/components/EventInvalidationBridge"

export function App() {
  if (window.location.pathname.startsWith("/installer")) return <InstallerPage />
  return <AuthGate><EventInvalidationBridge /><Router><RootRoute><Outlet /></RootRoute></Router></AuthGate>
}

export default App
