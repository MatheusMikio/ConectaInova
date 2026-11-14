import { useEffect, useState } from 'react'
import type { FormEvent } from 'react'
import { Field } from '../components/Field'
import { FormDialog } from '../components/FormDialog'
import { Icon } from '../components/Icon'
import { Notice } from '../components/Notice'
import { Pagination } from '../components/Pagination'
import { ScreenHeading } from '../components/ScreenHeading'
import { ProblemBrief } from '../components/ProblemBrief'
import { getDesafios, getDesafiosByNecessidade, updateDesafio } from '../services/desafioService'
import { getNecessidade } from '../services/necessidadeService'
import { errorMessage } from '../services/http'
import type { DesafioResponse, NecessidadeResponse, StatusDesafio } from '../types/domain'

const statusLabels = { aberto: 'Aberto', em_andamento: 'Em andamento', encerrado: 'Encerrado' }

export default function DesafiosPage({ onCreate, problem, onClearProblem, onUpdated }: {
  onCreate: () => void
  problem: NecessidadeResponse | null
  onClearProblem: () => void
  onUpdated: (challenge: DesafioResponse) => void
}) {
  const [items, setItems] = useState<DesafioResponse[]>([])
  const [record, setRecord] = useState<DesafioResponse | null>(null)
  const [titulo, setTitulo] = useState('')
  const [descricao, setDescricao] = useState('')
  const [status, setStatus] = useState<StatusDesafio>('aberto')
  const [criteria, setCriteria] = useState('')
  const [deadline, setDeadline] = useState(7)
  const [dialogOpen, setDialogOpen] = useState(false)
  const [linkedProblem, setLinkedProblem] = useState<NecessidadeResponse | null>(null)
  const [problemLoading, setProblemLoading] = useState(false)
  const [problemError, setProblemError] = useState('')
  const [problemRevision, setProblemRevision] = useState(0)
  const [page, setPage] = useState(1)
  const [revision, setRevision] = useState(0)
  const [loading, setLoading] = useState(true)
  const [busy, setBusy] = useState(false)
  const [formError, setFormError] = useState('')
  const [listError, setListError] = useState('')
  const [message, setMessage] = useState('')

  useEffect(() => {
    let active = true
    const query = problem ? getDesafiosByNecessidade(problem.id, page) : getDesafios(page)
    query.then(result => {
      if (active) { setItems(result); setListError('') }
    }).catch(err => {
      if (active) { setItems([]); setListError(errorMessage(err)) }
    }).finally(() => { if (active) setLoading(false) })
    return () => { active = false }
  }, [page, revision, problem])

  const problemId = record?.necessidade_id
  useEffect(() => {
    if (!dialogOpen || !problemId) return
    let active = true
    getNecessidade(problemId).then(result => {
      if (active) { setLinkedProblem(result); setProblemError('') }
    }).catch(err => {
      if (active) setProblemError(errorMessage(err))
    }).finally(() => { if (active) setProblemLoading(false) })
    return () => { active = false }
  }, [dialogOpen, problemId, problemRevision])

  function reload() { setLoading(true); setRevision(value => value + 1) }

  function edit(item: DesafioResponse) {
    setRecord(item); setTitulo(item.titulo); setDescricao(item.descricao); setStatus(item.status)
    setCriteria((item.criterios_sucesso || []).join('\n')); setDeadline(item.prazo_dias || 7)
    setLinkedProblem(null); setProblemLoading(true); setProblemError('')
    setFormError(''); setMessage(''); setDialogOpen(true)
  }

  async function save(event: FormEvent) {
    event.preventDefault()
    if (!record) return
    setBusy(true); setFormError('')
    try {
      const result = await updateDesafio(record.id, { titulo, descricao, status, criterios_sucesso: criteria.split('\n').map(value => value.trim()).filter(Boolean), prazo_dias: deadline })
      onUpdated(result)
      setDialogOpen(false)
      setMessage(`${result.titulo}: desafio atualizado com sucesso.`)
      reload()
    } catch (err) { setFormError(errorMessage(err)) }
    finally { setBusy(false) }
  }

  return <div className="directory-page">
    <ScreenHeading icon="target" title="Desafios de inovação" subtitle="Explore as oportunidades de inovação e acompanhe os desafios publicados." />
    <Notice message={message} />
    {problem && <div className="problem-summary"><Icon name="clipboard" size={21} /><span><strong>Problema:</strong> {problem.titulo}</span><button type="button" className="text-button" onClick={onClearProblem}>Ver todos os desafios</button></div>}
    <div className="directory-toolbar"><div><h2>Desafios publicados</h2><p>Os desafios nascem de problemas identificados pelos lojistas.</p></div><div className="directory-actions"><button type="button" className="text-button" disabled={loading} onClick={reload}>Atualizar lista</button><button type="button" onClick={onCreate}><Icon name="plus" size={18} />Novo desafio</button></div></div>
    <Notice error={listError} />
    {loading ? <div className="directory-loading" role="status">Carregando desafios…</div>
      : listError ? <div className="empty-state"><h2>Não foi possível carregar os desafios.</h2><button type="button" className="secondary" onClick={reload}>Tentar novamente</button></div>
      : items.length === 0 ? <div className="empty-state"><Icon name="target" size={38} /><h2>{page === 1 ? (problem ? 'Este problema ainda não tem desafios.' : 'Nenhum desafio publicado ainda.') : 'Nenhum desafio nesta página.'}</h2><p>{problem ? 'Busque soluções para este problema. Se nenhuma atender, gere um desafio de inovação.' : 'Cadastre um problema e busque soluções. Se nenhuma atender, gere um desafio a partir dessa necessidade.'}</p>{page === 1 && <button type="button" className="secondary" onClick={onCreate}>{problem ? 'Buscar soluções para este problema' : 'Começar pelo problema'}</button>}</div>
      : <div className="directory-grid">{items.map(item => <article className="directory-card" key={item.id}><div className="directory-card-top"><span className="solution-icon variant-0"><Icon name="target" size={26} /></span><span className={`status-badge ${item.status}`}>{statusLabels[item.status]}</span><button type="button" className="text-button" aria-label={`Editar ${item.titulo}`} onClick={() => edit(item)}>Editar</button></div><h2>{item.titulo}</h2><p className="directory-description">{item.descricao}</p></article>)}</div>}
    <Pagination page={page} count={items.length} loading={loading} onChange={value => { setLoading(true); setPage(value) }} />
    <FormDialog open={dialogOpen} title="Editar desafio" busy={busy} onClose={() => setDialogOpen(false)}>
      <Notice error={formError} />
      <section className="linked-problem" aria-label="Problema vinculado"><h3><Icon name="clipboard" size={17} />Problema de origem</h3>{problemLoading ? <p role="status">Carregando problema…</p> : problemError ? <><Notice error={problemError} /><button type="button" className="text-button" onClick={() => { setProblemLoading(true); setProblemRevision(value => value + 1) }}>Tentar novamente</button></> : <p>{linkedProblem?.titulo}</p>}<small>O vínculo com o problema é definido ao gerar o desafio e é mantido na edição.</small></section>
      {linkedProblem && <details className="brief-details"><summary>Ver briefing completo do problema</summary><ProblemBrief problem={linkedProblem} /></details>}
      <form onSubmit={save}><fieldset disabled={busy}>
        <Field label="Título do desafio" required><input required maxLength={200} value={titulo} onChange={event => setTitulo(event.target.value)} /></Field>
        <Field label="Descrição" required><textarea required maxLength={5000} rows={5} value={descricao} onChange={event => setDescricao(event.target.value)} /></Field>
        <Field label="Status"><select value={status} onChange={event => setStatus(event.target.value as StatusDesafio)}><option value="aberto">Aberto</option><option value="em_andamento">Em andamento</option><option value="encerrado">Encerrado</option></select></Field>
        <Field label="Critérios de sucesso — um por linha"><textarea rows={3} maxLength={6000} value={criteria} onChange={event => setCriteria(event.target.value)} /></Field>
        <Field label="Prazo para propostas (dias)"><input type="number" required min={1} max={365} value={deadline} onChange={event => setDeadline(Number(event.target.value))} /></Field>
        <div className="dialog-actions"><button type="button" className="secondary" onClick={() => setDialogOpen(false)}>Cancelar</button><button type="submit">{busy ? 'Salvando…' : 'Salvar alterações'}</button></div>
      </fieldset></form>
    </FormDialog>
  </div>
}
