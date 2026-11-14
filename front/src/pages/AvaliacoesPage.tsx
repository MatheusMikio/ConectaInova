import { useEffect, useState } from 'react'
import type { FormEvent } from 'react'
import { Field } from '../components/Field'
import { FormDialog } from '../components/FormDialog'
import { Notice } from '../components/Notice'
import { Pagination } from '../components/Pagination'
import { ProblemBrief } from '../components/ProblemBrief'
import { ScreenHeading } from '../components/ScreenHeading'
import { getAvaliacoes, updateAvaliacao } from '../services/avaliacaoService'
import { getStartupOptions } from '../services/startupService'
import { errorMessage } from '../services/http'
import type { AvaliacaoResponse, StartupResponse, StatusAvaliacao } from '../types/domain'

const labels: Record<StatusAvaliacao, string> = {
  solicitada: 'Solicitada', em_analise: 'Em análise', piloto_proposto: 'Piloto proposto',
  nao_atendida: 'Não atendida', concluida: 'Concluída',
}

export default function AvaliacoesPage() {
  const [items, setItems] = useState<AvaliacaoResponse[]>([])
  const [startups, setStartups] = useState<StartupResponse[]>([])
  const [startupID, setStartupID] = useState(0)
  const [page, setPage] = useState(1)
  const [revision, setRevision] = useState(0)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [message, setMessage] = useState('')
  const [selected, setSelected] = useState<AvaliacaoResponse | null>(null)
  const [status, setStatus] = useState<StatusAvaliacao>('solicitada')
  const [response, setResponse] = useState('')
  const [formError, setFormError] = useState('')
  const [busy, setBusy] = useState(false)

  useEffect(() => {
    let active = true
    getAvaliacoes(page, startupID).then(result => {
      if (active) { setItems(result); setError('') }
    }).catch(err => { if (active) setError(errorMessage(err)) })
      .finally(() => { if (active) setLoading(false) })
    return () => { active = false }
  }, [page, startupID, revision])

  useEffect(() => {
    let active = true
    getStartupOptions().then(result => { if (active) setStartups(result) }).catch(() => {
      // The unfiltered list remains available if startup options cannot be loaded.
    })
    return () => { active = false }
  }, [])

  function open(item: AvaliacaoResponse) {
    setSelected(item); setStatus(item.status); setResponse(item.resposta); setFormError('')
  }

  async function save(event: FormEvent) {
    event.preventDefault()
    if (!selected) return
    setBusy(true); setFormError('')
    try {
      await updateAvaliacao(selected.id, status, response)
      setSelected(null); setMessage('Avaliação atualizada.')
      setLoading(true); setRevision(value => value + 1)
    } catch (err) { setFormError(errorMessage(err)) }
    finally { setBusy(false) }
  }

  return <div className="directory-page">
    <ScreenHeading icon="clipboard" title="Avaliações" subtitle="Briefings encaminhados às startups e próximos passos para avaliar uma solução." />
    <Notice message={message} />
    <div className="directory-toolbar"><Field label="Startup"><select value={startupID} onChange={event => { setStartupID(Number(event.target.value)); setPage(1); setLoading(true) }}><option value={0}>Todas as startups</option>{startups.map(startup => <option key={startup.id} value={startup.id}>{startup.nome}</option>)}</select></Field><button className="text-button" disabled={loading} onClick={() => { setLoading(true); setRevision(value => value + 1) }}>Atualizar lista</button></div>
    <Notice error={error} />
    {loading ? <div className="directory-loading" role="status">Carregando avaliações…</div> : error ? <div className="empty-state"><h2>Não foi possível carregar as avaliações.</h2></div>
      : !items.length ? <div className="empty-state"><h2>Nenhuma avaliação nesta página.</h2><p>Abra uma solução no matching e solicite uma avaliação para encaminhar o briefing.</p></div>
        : <div className="directory-grid">{items.map(item => <article className="directory-card" key={item.id}>
          <div className="card-top"><span className={`status-badge ${item.status}`}>{labels[item.status]}</span><span className="badge">{item.startup.nome}</span></div>
          <h2>{item.necessidade.titulo}</h2><span className="directory-company">{item.necessidade.lojista_nome} · {item.solucao.nome}</span>
          <p className="directory-description">{item.necessidade.resultado_esperado}</p><small>{item.associacao}{item.regiao ? ` · ${item.regiao}` : ''}</small>
          <button type="button" className="secondary evaluation-open" onClick={() => open(item)}>Analisar demanda</button>
        </article>)}</div>}
    <Pagination page={page} count={items.length} loading={loading} onChange={value => { setPage(value); setLoading(true) }} />
    <FormDialog open={!!selected} title="Analisar demanda" busy={busy} onClose={() => setSelected(null)}>
      {selected && <>
        <p className="evaluation-intro"><strong>{selected.startup.nome}</strong> · {selected.solucao.nome}</p>
        <ProblemBrief problem={selected.necessidade} loja={selected.necessidade.lojista_nome} />
        {selected.mensagem && <div className="info-strip"><span><strong>Mensagem do lojista:</strong> {selected.mensagem}</span></div>}
        <Notice error={formError} />
        <form onSubmit={save}><fieldset disabled={busy}>
          <Field label="Etapa da avaliação"><select value={status} onChange={event => setStatus(event.target.value as StatusAvaliacao)}>{Object.entries(labels).map(([value, label]) => <option key={value} value={value}>{label}</option>)}</select></Field>
          <Field label="Parecer e próximo passo"><textarea rows={4} maxLength={4000} required={['piloto_proposto', 'nao_atendida', 'concluida'].includes(status)} value={response} onChange={event => setResponse(event.target.value)} placeholder="Registre condições, dúvidas, proposta de teste ou motivo de não atendimento." /></Field>
          <div className="dialog-actions"><button type="button" className="secondary" onClick={() => setSelected(null)}>Cancelar</button><button type="submit">{busy ? 'Salvando…' : 'Salvar avaliação'}</button></div>
        </fieldset></form>
      </>}
    </FormDialog>
  </div>
}
