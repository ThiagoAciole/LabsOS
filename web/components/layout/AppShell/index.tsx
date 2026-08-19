import type { ReactNode } from "react"

import { AppHeader } from "@/components/layout/AppHeader"
import { AppSidebar } from "@/components/layout/AppSidebar"
import { SidebarInset, SidebarProvider } from "@/components/ui/sidebar"

export function AppShell({ children }: { children: ReactNode }) {
  return (
    <SidebarProvider>
      <AppSidebar />
      <SidebarInset className="min-w-0 overflow-x-hidden">
        <AppHeader />
        <main className="m-0 flex min-w-0 flex-1 overflow-x-hidden p-0">{children}</main>
      </SidebarInset>
    </SidebarProvider>
  )
}
