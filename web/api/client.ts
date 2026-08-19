const API_BASE_URL = "/api/v1"

type ErrorEnvelope = { error?: { code?: string; message?: string; details?: unknown } }

export class ApiError extends Error {
  code: string
  status: number
  details?: unknown

  constructor(code: string, message: string, status: number, details?: unknown) {
    super(message)
    this.name = "ApiError"
    this.code = code
    this.status = status
    this.details = details
  }
}

export async function api<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(`${API_BASE_URL}${path}`, {
    ...init,
    headers: init?.body ? { "Content-Type": "application/json", ...init.headers } : init?.headers,
  })
  if (!response.ok) {
    const body = await response.json().catch(() => ({})) as ErrorEnvelope
    throw new ApiError(body.error?.code ?? "REQUEST_FAILED", body.error?.message ?? "request failed", response.status, body.error?.details)
  }
  return response.json() as Promise<T>
}

export const json = (method: "POST" | "PUT", body?: unknown): RequestInit => ({
  method,
  ...(body === undefined ? {} : { body: JSON.stringify(body) }),
})
