export const env = {
  apiUrl: import.meta.env.VITE_API_URL ?? "",
  filesUrl: import.meta.env.VITE_FILES_URL ?? "/file-manager/",
} as const
