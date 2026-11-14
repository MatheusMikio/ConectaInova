import { useEffect, useState } from 'react'
import { Icon } from '../components/Icon'
import { Notice } from '../components/Notice'
import { Pagination } from '../components/Pagination'
import { ScreenHeading } from '../components/ScreenHeading'
import { ProblemBrief } from '../components/ProblemBrief'
import type { InnovationFlow } from '../hooks/useInnovationFlow'
import { getNecessidades } from '../services/necessidadeService'
import { errorMessage } from '../services/http'
import type { Categoria, NecessidadeListResponse } from '../types/domain'

const statusLabels = { aberta: 'Aberta', em_analise: 'Em análise', resolvida: 'Resolvida' }
const priorityLabels = { baixa: 'Baixa', media: 'Média', alta: 'Alta' }
const categoryLabels: Record<Categoria, string> = {
  estoque: 'Estoque', vendas: 'Vendas', atendimento: 'Atendimento',
  financeiro: 'Financeiro', marketing: 'Marketing', logistica: 'Logística',
}

export default function ProblemasPage({ flow, matchingSelection = false, onSelect }: { flow: InnovationFlow; matchingSelection?: boolean; onSelect?: (id: number) => void }) {
  const [items, setItems] = useState<NecessidadeListResponse[]>([])
  const [page, setPage] = useState(1)
  const [revision, setRevision] = useState(0)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    let active = true
    getNecessidades(page).then(result => {
      if (active) { setItems(result); setError('') }
    }).catch(err => {
      if (active) { setItems([]); setError(errorMessage(err)) }
    }).finally(() => { if (active) setLoading(false) })
    return () => { active = false }
  }, [page, revision])

  function reload() { setLoading(true); setRevision(value => value + 1) }

  return <div className="directory-page">
    {!matchingSelection && <ScreenHeading icon="chart" title="Problemas cadastrados" subtitle="Conheça os desafios dos lojistas e encontre oportunidades de inovação." />}
    <div className="directory-toolbar"><div><h2>Problemas do varejo</h2><p>Necessidades cadastradas por todos os lojistas na plataforma.</p></div><div className="directory-actions"><button type="button" className="text-button" disabled={loading || flow.busy} onClick={reload}>Atualizar lista</button><button type="button" disabled={flow.busy} onClick={flow.newProblem}><Icon name="plus" size={18} />Cadastrar problema</button></div></div>
    <Notice error={error} />
    {loading ? <div className="directory-loading" role="status">Carregando problemas…</div>
      : error ? <div className="empty-state"><h2>Não foi possível carregar os problemas.</h2><button type="button" className="secondary" onClick={reload}>Tentar novamente</button></div>
      : items.length === 0 ? <div className="empty-state"><Icon name="clipboard" size={38} /><h2>{page === 1 ? 'Nenhum problema cadastrado ainda.' : 'Nenhum problema nesta página.'}</h2><p>Cadastre uma necessidade para buscar soluções e gerar oportunidades de inovação.</p></div>
      : <div className="directory-grid">{items.map(problem => <article className="directory-card problem-card" key={problem.id}>
        <div className="directory-card-top"><span className="solution-icon variant-0"><Icon name="clipboard" size={26} /></span><span className={`status-badge ${problem.status}`}>{statusLabels[problem.status]}</span><span className="badge">{categoryLabels[problem.categoria]}</span></div>
        <h2>{problem.titulo}</h2><span className="directory-company">{problem.lojista_nome || 'Loja indisponível'}</span>
        <p className="directory-description">{problem.descricao}</p>
        <details className="brief-details"><summary>Ver briefing</summary><ProblemBrief problem={problem} loja={problem.lojista_nome} /></details>
        <div className="problem-metadata"><span>Prioridade: <strong>{priorityLabels[problem.prioridade]}</strong></span><span>{new Date(problem.created_at).toLocaleDateString('pt-BR')}</span></div>
        <div className="problem-actions"><button type="button" disabled={flow.busy} onClick={() => onSelect ? onSelect(problem.id) : flow.openProblem(problem.id, 'matching')}><Icon name="search" size={16} />{matchingSelection ? 'Analisar este problema' : 'Ver soluções'}</button>{!matchingSelection && <><button type="button" className="secondary" disabled={flow.busy} onClick={() => flow.openProblem(problem.id, 'desafio')}><Icon name="target" size={16} />Ver desafios</button><button type="button" className="text-button" aria-label={`Editar ${problem.titulo}`} disabled={flow.busy} onClick={() => flow.openProblem(problem.id)}>Editar</button></>}</div>
      </article>)}</div>}
    <Pagination page={page} count={items.length} loading={loading || flow.busy} onChange={value => { setLoading(true); setPage(value) }} />
  </div>
}
