import { useState } from 'react'
import type { FormEvent } from 'react'
import { CategorySelect } from '../components/CategorySelect'
import { Field } from '../components/Field'
import { Notice } from '../components/Notice'
import { createNecessidade, getMatches, getNecessidade, updateNecessidade } from '../services/necessidadeService'
import { createDesafio } from '../services/desafioService'
import { errorMessage } from '../services/http'
import type { CreateNecessidadeRequest, MatchingResponse, NecessidadeResponse, Prioridade, StatusNecessidade } from '../types/domain'

const initial: CreateNecessidadeRequest = { lojista_id: 0, titulo: '', descricao: '', categoria: 'estoque', prioridade: 'media', status: 'aberta' }

export default function NecessidadesPage() {
  const [record, setRecord] = useState<NecessidadeResponse | null>(null)
  const [form, setForm] = useState(initial)
  const [lookup, setLookup] = useState('')
  const [matches, setMatches] = useState<MatchingResponse | null>(null)
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState('')
  const [message, setMessage] = useState('')
  const [challengeCreated, setChallengeCreated] = useState(false)

  function show(result: NecessidadeResponse) {
    setRecord(result)
    setForm({ lojista_id: result.lojista_id, titulo: result.titulo, descricao: result.descricao, categoria: result.categoria, prioridade: result.prioridade, status: result.status })
    setMatches(null)
    setChallengeCreated(false)
  }

  async function search(event: FormEvent) {
    event.preventDefault()
    setBusy(true); setError(''); setMessage('')
    try { show(await getNecessidade(Number(lookup))) }
    catch (err) { setError(errorMessage(err)) }
    finally { setBusy(false) }
  }

  async function save(event: FormEvent) {
    event.preventDefault()
    setBusy(true); setError(''); setMessage('')
    try {
      const result = record ? await updateNecessidade(record.id, form) : await createNecessidade(form)
      show(result)
      setMessage(`Necessidade #${result.id} salva. Agora você pode buscar soluções compatíveis.`)
    } catch (err) { setError(errorMessage(err)) }
    finally { setBusy(false) }
  }

  async function match() {
    if (!record) return
    setBusy(true); setError(''); setMessage(''); setMatches(null)
    try { setMatches(await getMatches(record.id)) }
    catch (err) { setError(errorMessage(err)) }
    finally { setBusy(false) }
  }

  async function challenge() {
    if (!record) return
    setBusy(true); setError(''); setMessage('')
    try {
      const result = await createDesafio({
        necessidade_id: record.id,
        titulo: `Inovação para: ${record.titulo}`.slice(0, 200),
        descricao: record.descricao,
      })
      setChallengeCreated(true)
      setMessage(`Desafio #${result.id} criado. Acompanhe e edite na área de desafios.`)
    } catch (err) { setError(errorMessage(err)) }
    finally { setBusy(false) }
  }

  return <>
    <div className="page-heading"><div><span className="eyebrow">02 / Necessidades</span><h1>Transforme desafios em conexões.</h1><p>Conte o que seu negócio precisa. Encontre quem pode ajudar.</p></div></div>
    <Notice error={error} message={message} />
    <section className="panel search-panel"><form className="lookup" onSubmit={search}><Field label="Consultar necessidade por ID"><input type="number" min="1" step="1" required value={lookup} onChange={event => setLookup(event.target.value)} /></Field><button className="secondary" disabled={busy}>Consultar</button></form></section>
    <div className="two-columns">
      <section className="panel">
        <div className="panel-heading"><h2>{record ? `Necessidade #${record.id}` : 'Nova necessidade'}</h2>{record && <button className="text-button" disabled={busy} onClick={() => { setRecord(null); setForm(initial); setMatches(null); setMessage(''); setError('') }}>Novo cadastro</button>}</div>
        <form onSubmit={save}><fieldset disabled={busy}>
          <Field label="ID do lojista"><input required type="number" min="1" step="1" value={form.lojista_id || ''} onChange={event => setForm({ ...form, lojista_id: Number(event.target.value) })} /></Field>
          <Field label="Título"><input required maxLength={200} placeholder="Ex.: Quero automatizar meu estoque" value={form.titulo} onChange={event => setForm({ ...form, titulo: event.target.value })} /></Field>
          <Field label="Descreva sua necessidade"><textarea required maxLength={5000} rows={4} value={form.descricao} onChange={event => setForm({ ...form, descricao: event.target.value })} placeholder="O que acontece hoje? O que você gostaria de melhorar?" /></Field>
          <div className="form-row"><Field label="Categoria"><CategorySelect value={form.categoria} onChange={categoria => setForm({ ...form, categoria })} /></Field><Field label="Prioridade"><select value={form.prioridade} onChange={event => setForm({ ...form, prioridade: event.target.value as Prioridade })}><option value="baixa">Baixa</option><option value="media">Média</option><option value="alta">Alta</option></select></Field></div>
          <Field label="Status"><select value={form.status} onChange={event => setForm({ ...form, status: event.target.value as StatusNecessidade })}><option value="aberta">Aberta</option><option value="em_analise">Em análise</option><option value="resolvida">Resolvida</option></select></Field>
          <button>{busy ? 'Aguarde…' : record ? 'Salvar alterações' : 'Cadastrar necessidade'}</button>
        </fieldset></form>
      </section>
      <section className="panel matching-panel">
        <span className="eyebrow">Matching / ConectaInova</span><h2>A solução pode estar mais perto do que você imagina.</h2>
        <p>A busca usa a categoria e as palavras-chave do título e da descrição salvos.</p>
        <button type="button" disabled={!record || busy} onClick={match}>{busy ? 'Aguarde…' : 'Buscar soluções compatíveis'} <span aria-hidden="true">↗</span></button>
        {!record && <p className="muted">Cadastre ou consulte uma necessidade para começar.</p>}
        {matches && <div className="results" aria-live="polite">
          <h3>{matches.matches.length} solução(ões) encontrada(s)</h3>
          {matches.matches.map(item => <article className="match-card" key={item.solucao.id}>
            <div className="card-top"><span className="badge">{item.solucao.categoria}</span><strong className="score">{item.pontuacao}%</strong></div>
            <h3>{item.solucao.nome}</h3><p>{item.solucao.descricao}</p>
            <small>{item.startup.nome} · Solução #{item.solucao.id}</small>
          </article>)}
          {matches.sugerir_desafio && <div className="empty-state"><h3>Uma oportunidade para inovar.</h3><p>Ainda não há soluções compatíveis. Crie um desafio para essa necessidade.</p><button type="button" disabled={busy || challengeCreated} onClick={challenge}>{challengeCreated ? 'Desafio criado' : 'Criar desafio de inovação'}</button></div>}
        </div>}
      </section>
    </div>
  </>
}
