import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { applySystemUpdate, getSystemHealth, getSystemSummary, getSystemUpdate, runPowerAction } from "@/api/system"
import { getSettings, updateNetworkSettings, updateSystemSettings } from "@/api/settings"

const key = ["settings"] as const
export function useSettingsData() { return useQuery({ queryKey: key, queryFn: async () => { const [settings, summary, health] = await Promise.all([getSettings(), getSystemSummary(), getSystemHealth()]); return { settings, summary, health } } }) }
export function useSystemSettingsMutation() { const client = useQueryClient(); return useMutation({ mutationFn: updateSystemSettings, onSuccess: () => client.invalidateQueries({ queryKey: key }) }) }
export function useNetworkSettingsMutation() { const client = useQueryClient(); return useMutation({ mutationFn: updateNetworkSettings, onSuccess: () => client.invalidateQueries({ queryKey: key }) }) }
export function usePowerMutation() { return useMutation({ mutationFn: runPowerAction }) }
export function useSystemUpdateData() { return useQuery({ queryKey: ["system-update"], queryFn: getSystemUpdate }) }
export function useSystemUpdateMutation() { const client = useQueryClient(); return useMutation({ mutationFn: applySystemUpdate, onSuccess: () => client.invalidateQueries({ queryKey: ["system-update"] }) }) }
