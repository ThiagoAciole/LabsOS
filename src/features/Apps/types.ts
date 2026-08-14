export type AppStatus = "installing" | "running" | "stopped" | "error";
export type AppCategory =
  | "all"
  | "media"
  | "network"
  | "storage"
  | "automation"
  | "development"
  | "utilities";
export type AppSource = "all" | "official" | "community";

export interface InstalledApp {
  id: string;
  name: string;
  description: string;
  icon?: string;
  version?: string;
  status: AppStatus;
  progress?: number;
  url?: string;
}
export interface StoreApp {
  id: string;
  name: string;
  description: string;
  icon?: string;
  category: AppCategory;
  source: AppSource;
  version: string;
  size: string;
	highlights: string[];
	installable?: boolean;
}
