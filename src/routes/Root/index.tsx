import type { ReactNode } from "react"
import { AppShell } from "@/components/layout/AppShell"

export function RootRoute({ children }: { children: ReactNode }) { return <AppShell>{children}</AppShell> }
