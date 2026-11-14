import type { AvaliacaoResponse, StatusAvaliacao } from '../types/domain'
import { request } from './http'

export function createAvaliacao(necessidadeID: number, solucaoID: number, mensagem: string) {
  return request<AvaliacaoResponse>('/avaliacoes', {
    method: 'POST', body: JSON.stringify({ necessidade_id: necessidadeID, solucao_id: solucaoID, mensagem }),
  })
}

export function getAvaliacoes(page = 1, startupID = 0) {
  return request<AvaliacaoResponse[]>(`/avaliacoes?page=${page}&size=12${startupID ? `&startup_id=${startupID}` : ''}`)
}

export function updateAvaliacao(id: number, status: StatusAvaliacao, resposta: string) {
  return request<AvaliacaoResponse>(`/avaliacoes/${id}`, { method: 'PATCH', body: JSON.stringify({ status, resposta }) })
}

export function getPlataforma() {
  return request<{ associacao: string; regiao: string }>('/plataforma')
}
