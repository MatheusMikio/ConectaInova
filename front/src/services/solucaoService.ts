import type { CreateSolucaoRequest, SolucaoResponse, UpdateSolucaoRequest } from '../types/domain'
import { request } from './http'

export function createSolucao(data: CreateSolucaoRequest) {
  return request<SolucaoResponse>('/solucoes', { method: 'POST', body: JSON.stringify(data) })
}

export function getSolucoes(page = 1) {
  return request<SolucaoResponse[]>(`/solucoes?page=${page}&size=12`)
}

export function updateSolucao(id: number, data: UpdateSolucaoRequest) {
  return request<SolucaoResponse>(`/solucoes/${id}`, { method: 'PATCH', body: JSON.stringify(data) })
}
