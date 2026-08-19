export type AppStatus = "installing" | "running" | "stopped" | "updating" | "error";
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
  updateAvailable?: boolean;
  progress?: number;
  url?: string;
	 health?: string;
	 dependencies?: string[];
	 volumes?: string[];
	 ports?: number[];
}
export interface StoreApp {
  id: string;
  name: string;
  description: string;
  icon?: string;
  category: AppCategory;
	source: string;
  version: string;
  size: string;
	highlights: string[];
	installable?: boolean;
	architecture?: string[];
	requirements?: string[];
}
