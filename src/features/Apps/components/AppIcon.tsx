import { useState } from "react";

export function AppIcon({
  icon,
  name,
  className,
}: {
  icon?: string;
  name?: string;
  className: string;
}) {
  const [failed, setFailed] = useState(false);
  if (!icon || failed) {
    return <span className={`${className} flex items-center justify-center rounded-lg bg-muted text-sm font-semibold text-muted-foreground`} aria-hidden="true">{name?.slice(0, 2).toUpperCase() ?? "APP"}</span>;
  }
  return <img src={icon} alt="" className={className} onError={() => setFailed(true)} />;
}
