import { RefreshCw } from "lucide-react";

export function AppIcon({
  icon,
  className,
}: {
  icon?: string;
  className: string;
}) {
  return icon ? (
    <img src={icon} alt="" className={className} />
  ) : (
    <RefreshCw className={className} aria-hidden="true" />
  );
}
