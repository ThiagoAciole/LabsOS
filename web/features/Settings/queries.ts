import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { applySystemUpdate, getSystemHealth, getSystemSummary, getSystemUpdate, rollbackSystemUpdate, runPowerAction } from "@/api/system"
import { changePassword, getSessions, getSettings, getNetworkSnapshot, getSSHStatus, getStorageSnapshot, getWiFiSnapshot, revokeSession, updateNetwork, updateNetworkSettings, updateSSH, updateSystemSettings, updateWiFi } from "@/api/settings"
import { getDiagnostics } from "@/api/system"

const key = ["settings"] as const
export function useSettingsData() { return useQuery({ queryKey: key, queryFn: async () => { const [settings, summary, health] = await Promise.all([getSettings(), getSystemSummary(), getSystemHealth()]); return { settings, summary, health } } }) }
export function useSystemSettingsMutation() { const client = useQueryClient(); return useMutation({ mutationFn: updateSystemSettings, onSuccess: () => client.invalidateQueries({ queryKey: key }) }) }
export function useNetworkSettingsMutation() { const client = useQueryClient(); return useMutation({ mutationFn: updateNetworkSettings, onSuccess: () => client.invalidateQueries({ queryKey: key }) }) }
export function usePowerMutation() { return useMutation({ mutationFn: runPowerAction }) }
export function useSystemUpdateData() { return useQuery({ queryKey: ["system-update"], queryFn: getSystemUpdate }) }
export function useSystemUpdateMutation() { const client = useQueryClient(); return useMutation({ mutationFn: applySystemUpdate, onSuccess: () => client.invalidateQueries({ queryKey: ["system-update"] }) }) }
export function useSystemRollbackMutation() { const client = useQueryClient(); return useMutation({ mutationFn: rollbackSystemUpdate, onSuccess: () => client.invalidateQueries({ queryKey: ["system-update"] }) }) }
export function useChangePasswordMutation() { return useMutation({ mutationFn: ({ current, next }: { current: string; next: string }) => changePassword(current, next) }) }
export function useSessions() { return useQuery({ queryKey: ["auth-sessions"], queryFn: getSessions }) }
export function useRevokeSessionMutation() { const client = useQueryClient(); return useMutation({ mutationFn: revokeSession, onSuccess: () => client.invalidateQueries({ queryKey: ["auth-sessions"] }) }) }
export function useNetworkSnapshot() { return useQuery({ queryKey: ["network-snapshot"], queryFn: getNetworkSnapshot }) }
export function useNetworkMutation() { const client = useQueryClient(); return useMutation({ mutationFn: updateNetwork, onSuccess: () => client.invalidateQueries({ queryKey: ["network-snapshot"] }) }) }
export function useWiFiSnapshot() { return useQuery({ queryKey: ["wifi-snapshot"], queryFn: getWiFiSnapshot }) }
export function useWiFiMutation() { const client = useQueryClient(); return useMutation({ mutationFn: updateWiFi, onSuccess: () => client.invalidateQueries({ queryKey: ["wifi-snapshot"] }) }) }
export function useStorageSnapshot() { return useQuery({ queryKey: ["storage-snapshot"], queryFn: getStorageSnapshot }) }
export function useSSHStatus() { return useQuery({ queryKey: ["ssh-status"], queryFn: getSSHStatus }) }
export function useSSHMutation() { const client = useQueryClient(); return useMutation({ mutationFn: updateSSH, onSuccess: () => client.invalidateQueries({ queryKey: ["ssh-status"] }) }) }
export function useDiagnostics() { return useQuery({ queryKey: ["diagnostics"], queryFn: getDiagnostics, enabled: false }) }
