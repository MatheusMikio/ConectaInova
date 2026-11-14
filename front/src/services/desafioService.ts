import type { CreateDesafioRequest, DesafioResponse, UpdateDesafioRequest } from '../types/domain'
import { request } from './http'

export function createDesafio(data: CreateDesafioRequest) {
  return request<DesafioResponse>('/desafios', { method: 'POST', body: JSON.stringify(data) })
}

export function getDesafios(page = 1) {
  return request<DesafioResponse[]>(`/desafios?page=${page}&size=12`)
}

export function getDesafiosByNecessidade(necessidadeID: number, page = 1) {
  return request<DesafioResponse[]>(`/necessidades/${necessidadeID}/desafios?page=${page}&size=12`)
}

export function updateDesafio(id: number, data: UpdateDesafioRequest) {
  return request<DesafioResponse>(`/desafios/${id}`, { method: 'PATCH', body: JSON.stringify(data) })
}
