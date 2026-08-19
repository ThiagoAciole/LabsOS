"use client"

import * as React from "react"
import { navigation } from "@/config/navigation"
import { navigate } from "@/app/navigation"
import { Sidebar, SidebarContent, SidebarHeader, SidebarMenu, SidebarMenuButton, SidebarMenuItem, useSidebar } from "@/components/ui/sidebar"

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  const { toggleSidebar } = useSidebar()
  const [pathname, setPathname] = React.useState(window.location.pathname)

  React.useEffect(() => {
    const handlePopState = () => setPathname(window.location.pathname)
    window.addEventListener("popstate", handlePopState)
    return () => window.removeEventListener("popstate", handlePopState)
  }, [])

  return (
    <Sidebar className="items-center" variant="inset" collapsible="icon" {...props}>
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton size="lg" onClick={toggleSidebar}>
              <div className="flex aspect-square size-8 items-center justify-center"><img src="/Icon.svg" alt="LabsOS" className="size-8" /></div>
              <span className="truncate font-medium">LabsOS</span>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>
      <SidebarContent >
        <SidebarMenu>
          {navigation.map(({ title, href, icon: Icon }) => (
            <SidebarMenuItem key={href}>
              <SidebarMenuButton className="p-2 py-6" isActive={pathname === href} tooltip={title} render={<a href={href} onClick={(event) => { event.preventDefault(); navigate(href) }}  />}>
                <Icon />
                <span>{title}</span>
              </SidebarMenuButton>
            </SidebarMenuItem>
          ))}
        </SidebarMenu>
      </SidebarContent>
    </Sidebar>
  )
}
