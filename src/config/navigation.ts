import { Boxes, Folder, House, Settings } from "lucide-react"

export const navigation = [
  { title: "Home", href: "/", icon: House },
  { title: "Apps", href: "/apps", icon: Boxes },
  { title: "Files", href: "/files", icon: Folder },
  { title: "Settings", href: "/settings", icon: Settings },
] as const
