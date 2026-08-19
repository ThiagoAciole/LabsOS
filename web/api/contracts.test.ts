/// <reference types="node" />
import assert from "node:assert/strict"
import test from "node:test"

import { ApiError, api } from "./client.ts"
import { toSystemDashboard } from "./system.ts"
import { toInstalledApp } from "../features/Apps/api.ts"

test("converte o envelope de erro da Labs API", async () => {
  const originalFetch = globalThis.fetch
  globalThis.fetch = async () => new Response(JSON.stringify({ error: { code: "PROVIDER_UNAVAILABLE", message: "active provider is unavailable", details: {} } }), { status: 503, headers: { "Content-Type": "application/json" } })
  try {
    await assert.rejects(api("/system/summary"), (error: unknown) => error instanceof ApiError && error.code === "PROVIDER_UNAVAILABLE" && error.status === 503)
  } finally {
    globalThis.fetch = originalFetch
  }
})

test("preserva estado da API e complementa somente metadata editorial", () => {
  const app = toInstalledApp({ id: "jellyfin", name: "Jellyfin", icon: "/jellyfin.svg", description: "Media", status: "stopped", version: "1", url: "http://labsos.local:8096", updateAvailable: false, installed: true })
  assert.equal(app.status, "stopped")
  assert.equal(app.url, "http://labsos.local:8096")
})

test("adapta bytes e uptime para a Home", () => {
  const dashboard = toSystemDashboard({ hostname: "labsos-dev", status: "healthy", uptimeSeconds: 223380, version: "0.1.0", cpuUsage: 18.2, memoryUsedBytes: 3_221_225_472, memoryTotalBytes: 8_589_934_592, temperatureCelsius: 43, storageUsedBytes: 459_561_500_672, storageTotalBytes: 999_653_638_144, ipAddress: "172.20.0.2", networkOnline: true, networkDownloadBytesPerSecond: 2_000_000, networkUploadBytesPerSecond: 1_000_000 })
  assert.equal(dashboard.metrics[0].value, 38)
  assert.equal(dashboard.uptime.days, 2)
  assert.equal(dashboard.uptime.hours, 14)
  assert.equal(dashboard.uptime.minutes, 3)
  assert.deepEqual(dashboard.network, { download: 2, upload: 1, ip: "172.20.0.2", online: true })
})
