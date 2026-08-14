import { ChevronDown, Globe2, Link2 } from "lucide-react"

import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "@/components/ui/dropdown-menu"

export function AppHeader() {
  return (
    <header className="flex h-16 shrink-0 items-center justify-end gap-2 px-4 md:px-6">
      <Button variant="outline" size="sm" className="gap-2 rounded-md bg-card/60">
        <Link2 />
        <span className="hidden sm:inline">Acesso local</span>
      </Button>
      <DropdownMenu>
        <DropdownMenuTrigger render={<Button variant="outline" size="sm" className="gap-2 rounded-md bg-card/60" />}>
          <Globe2 />
          <span>PT-BR</span>
          <ChevronDown />
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end">
          <DropdownMenuItem>PT-BR</DropdownMenuItem>
          <DropdownMenuItem>EN-US</DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
      <Badge variant="outline" className="h-9 gap-2 rounded-md px-3 text-muted-foreground">
        <span className="size-2 rounded-full bg-emerald-500" />
        <span className="hidden sm:inline">Servidor conectado</span>
      </Badge>
    </header>
  )
}
