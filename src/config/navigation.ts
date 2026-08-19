import { Archive, Bell, Boxes, FileText, Folder, Globe2, House, KeyRound, ListTodo, ScrollText, Settings } from "lucide-react"

export const navigation = [
  { title: "Home", href: "/", icon: House },
  { title: "Apps", href: "/apps", icon: Boxes },
  { title: "Files", href: "/files", icon: Folder },
  { title: "Settings", href: "/settings", icon: Settings },
  { title: "Segredos", href: "/secrets", icon: KeyRound },
  { title: "Notificações", href: "/notifications", icon: Bell },
  { title: "Logs", href: "/logs", icon: FileText },
  { title: "Tarefas", href: "/jobs", icon: ListTodo },
  { title: "Auditoria", href: "/audit", icon: ScrollText },
  { title: "Backups", href: "/backups", icon: Archive },
  { title: "Serviços", href: "/services", icon: Globe2 },
] as const
