import { FolderOpen } from "lucide-react"

import { Card, CardContent } from "@/components/ui/card"

export function FilesPage() {
  return (
    <div className="flex flex-1 flex-col gap-6 rounded-2xl bg-background/80 p-4 md:p-8">
      <div>
        <h1 className="text-3xl font-medium tracking-tight">Files</h1>
        <p className="mt-1 text-muted-foreground">Arquivos do LabsOS</p>
      </div>
      <Card>
        <CardContent className="flex min-h-48 flex-col items-center justify-center gap-3 text-center">
          <FolderOpen className="size-8 text-primary" />
          <p className="font-medium">Nenhum arquivo disponível</p>
          <p className="text-sm text-muted-foreground">O explorador de arquivos será exibido aqui.</p>
        </CardContent>
      </Card>
    </div>
  )
}
