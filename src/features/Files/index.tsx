import { FolderOpen } from "lucide-react"

import { Card, CardContent } from "@/components/ui/card"
import { env } from "@/lib/env"

export function FilesPage() {
  return (
    <div className="flex min-w-0 flex-1 flex-col gap-6 rounded-2xl bg-background/80 p-4 md:p-8">
      <div>
        <h1 className="text-3xl font-medium tracking-tight">Files</h1>
        <p className="mt-1 text-muted-foreground">Arquivos do LabsOS</p>
      </div>
      {env.filesUrl ? (
        <Card className="min-h-[calc(100dvh-10rem)] min-w-0 overflow-hidden">
          <iframe className="h-[calc(100dvh-10rem)] min-h-96 w-full border-0" src={env.filesUrl} title="Gerenciador de arquivos" />
        </Card>
      ) : (
        <Card>
          <CardContent className="flex min-h-48 flex-col items-center justify-center gap-3 text-center">
            <FolderOpen className="size-8 text-primary" />
            <p className="font-medium">Gerenciador de arquivos indisponível</p>
            <p className="text-sm text-muted-foreground">O System File App ainda não foi provisionado neste servidor.</p>
          </CardContent>
        </Card>
      )}
    </div>
  )
}
