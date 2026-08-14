import { ArrowDown, ArrowUp, Wifi } from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import type { NetworkData } from "../../types";
export function NetworkCard({ network: networkData }: { network: NetworkData }) {
  return (
    <Card className="flex min-h-28 flex-col">
      <CardHeader className="p-3 pb-2">
        <CardTitle className="flex items-center gap-2 text-sm">
          <span className="flex size-8 items-center justify-center rounded-md bg-primary/10 text-primary">
            <Wifi className="size-4" />
          </span>
          Rede
        </CardTitle>
      </CardHeader>
      <CardContent className="flex flex-1 items-end p-3 pt-0">
        <div className="space-y-2">
          <p className="text-xs text-muted-foreground">{networkData.online ? networkData.ip : "Rede indisponível"}</p>
          <div className="flex gap-3">
          <span className="flex items-center gap-2">
            <ArrowDown className="size-4 text-muted-foreground" />
            <span>
              <small className="block text-muted-foreground">Download</small>
              <strong className="whitespace-nowrap text-sm">
                {networkData.download} MB/s
              </strong>
            </span>
          </span>
          <span className="flex items-center gap-2">
            <ArrowUp className="size-4 text-muted-foreground" />
            <span>
              <small className="block text-muted-foreground">Upload</small>
              <strong className="whitespace-nowrap text-sm">
                {networkData.upload} MB/s
              </strong>
            </span>
          </span>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
