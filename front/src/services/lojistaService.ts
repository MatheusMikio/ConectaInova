import type { CreateLojistaRequest, LojistaResponse, UpdateLojistaRequest } from '../types/domain'
import { request } from './http'

export function createLojista(data: CreateLojistaRequest) {
  return request<LojistaResponse>('/lojistas', { method: 'POST', body: JSON.stringify(data) })
}

export function getLojista(id: number) {
  return request<LojistaResponse>(`/lojistas/${id}`)
}

export function updateLojista(id: number, data: UpdateLojistaRequest) {
  return request<LojistaResponse>(`/lojistas/${id}`, { method: 'PATCH', body: JSON.stringify(data) })
}
