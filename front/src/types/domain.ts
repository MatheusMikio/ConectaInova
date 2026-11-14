export type Categoria = 'estoque' | 'vendas' | 'atendimento' | 'financeiro' | 'marketing' | 'logistica'
export type Prioridade = 'baixa' | 'media' | 'alta'
export type StatusNecessidade = 'aberta' | 'em_analise' | 'resolvida'
export type StatusDesafio = 'aberto' | 'em_andamento' | 'encerrado'

export interface ResourceMetadata {
  id: number
  created_at: string
  updated_at: string
}

export interface CreateLojistaRequest {
  nome: string
  segmento: string
}
export interface UpdateLojistaRequest {
  nome?: string
  segmento?: string
}
export interface LojistaResponse extends ResourceMetadata, CreateLojistaRequest {}

export interface CreateNecessidadeRequest {
  lojista_id: number
  titulo: string
  descricao: string
  categoria: Categoria
  prioridade: Prioridade
  status?: StatusNecessidade
  impacto?: string
  resultado_esperado?: string
  restricoes?: string
  canal_venda?: string
  criterios_sucesso?: string[]
}
export interface UpdateNecessidadeRequest {
  lojista_id?: number
  titulo?: string
  descricao?: string
  categoria?: Categoria
  prioridade?: Prioridade
  status?: StatusNecessidade
  impacto?: string
  resultado_esperado?: string
  restricoes?: string
  canal_venda?: string
  criterios_sucesso?: string[]
}
export interface NecessidadeResponse extends ResourceMetadata, CreateNecessidadeRequest {
  status: StatusNecessidade
  impacto: string
  resultado_esperado: string
  restricoes: string
  canal_venda: string
  criterios_sucesso: string[]
}

export interface NecessidadeListResponse extends NecessidadeResponse {
  lojista_nome: string
}

export interface CreateStartupRequest {
  nome: string
  descricao: string
}
export interface UpdateStartupRequest {
  nome?: string
  descricao?: string
}
export interface StartupResponse extends ResourceMetadata, CreateStartupRequest {}

export interface CreateSolucaoRequest {
  startup_id: number
  nome: string
  descricao: string
  categoria: Categoria
  palavras_chave: string[]
}
export interface UpdateSolucaoRequest {
  startup_id?: number
  nome?: string
  descricao?: string
  categoria?: Categoria
  palavras_chave?: string[]
}
export interface SolucaoResponse extends ResourceMetadata, CreateSolucaoRequest {}

export interface CreateDesafioRequest {
  necessidade_id: number
  titulo: string
  descricao: string
  status?: StatusDesafio
  criterios_sucesso?: string[]
  prazo_dias?: number
}
export interface UpdateDesafioRequest {
  necessidade_id?: number
  titulo?: string
  descricao?: string
  status?: StatusDesafio
  criterios_sucesso?: string[]
  prazo_dias?: number
}
export interface DesafioResponse extends ResourceMetadata, CreateDesafioRequest {
  status: StatusDesafio
  criterios_sucesso: string[]
  prazo_dias: number
}

export interface MatchResponse {
  solucao: SolucaoResponse
  startup: StartupResponse
  pontuacao: number
  motivos: string[]
  palavras_encontradas: string[]
  pontos_a_validar: string[]
}

export type StatusAvaliacao = 'solicitada' | 'em_analise' | 'piloto_proposto' | 'nao_atendida' | 'concluida'

export interface AvaliacaoResponse extends ResourceMetadata {
  necessidade: NecessidadeListResponse
  solucao: SolucaoResponse
  startup: StartupResponse
  mensagem: string
  resposta: string
  associacao: string
  regiao: string
  status: StatusAvaliacao
}
export interface MatchingResponse {
  necessidade_id: number
  matches: MatchResponse[]
  sugerir_desafio: boolean
}
