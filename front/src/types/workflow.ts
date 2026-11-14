import type { Categoria, Prioridade } from './domain'

export type Screen = 'inicio' | 'cadastro' | 'matching' | 'solucoes' | 'startups' | 'desafio' | 'desafio-draft' | 'problemas' | 'avaliacoes'

export interface ProblemForm {
  nomeLoja: string
  segmento: string
  titulo: string
  descricao: string
  impacto: string
  prioridade: Prioridade
  categoria: Categoria
  resultadoEsperado: string
  restricoes: string
  canalVenda: string
  criteriosSucesso: string
}

export interface ChallengeDraft {
  titulo: string
  descricao: string
  criterios: string
  prazo: number
}
