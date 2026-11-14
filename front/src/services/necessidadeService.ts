import type { CreateNecessidadeRequest, MatchingResponse, NecessidadeListResponse, NecessidadeResponse, UpdateNecessidadeRequest } from '../types/domain'
import { request } from './http'

export function createNecessidade(data: CreateNecessidadeRequest) {
  return request<NecessidadeResponse>('/necessidades', { method: 'POST', body: JSON.stringify(data) })
}

export function getNecessidade(id: number) {
  return request<NecessidadeResponse>(`/necessidades/${id}`)
}

export function getNecessidades(page = 1) {
  return request<NecessidadeListResponse[]>(`/necessidades?page=${page}&size=12`)
}

export function updateNecessidade(id: number, data: UpdateNecessidadeRequest) {
  return request<NecessidadeResponse>(`/necessidades/${id}`, { method: 'PATCH', body: JSON.stringify(data) })
}

export function getMatches(id: number) {
  return request<MatchingResponse>(`/necessidades/${id}/matches`)
}
