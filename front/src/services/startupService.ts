import type { CreateStartupRequest, StartupResponse, UpdateStartupRequest } from '../types/domain'
import { request } from './http'

export function createStartup(data: CreateStartupRequest) {
  return request<StartupResponse>('/startups', { method: 'POST', body: JSON.stringify(data) })
}

export function getStartups(page = 1) {
  return request<StartupResponse[]>(`/startups?page=${page}&size=12`)
}

export async function getStartupOptions(signal?: AbortSignal) {
  const startups: StartupResponse[] = []
  for (let page = 1; ; page++) {
    const items = await request<StartupResponse[]>(`/startups?page=${page}&size=30`, { signal })
    startups.push(...items)
    if (items.length < 30) return startups
  }
}

export function updateStartup(id: number, data: UpdateStartupRequest) {
  return request<StartupResponse>(`/startups/${id}`, { method: 'PATCH', body: JSON.stringify(data) })
}
