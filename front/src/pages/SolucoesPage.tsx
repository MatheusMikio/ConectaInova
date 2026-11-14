import { useEffect, useState } from 'react'
import type { FormEvent } from 'react'
import { CategorySelect } from '../components/CategorySelect'
import { Field } from '../components/Field'
import { FormDialog } from '../components/FormDialog'
import { Icon } from '../components/Icon'
import { Notice } from '../components/Notice'
import { Pagination } from '../components/Pagination'
import { ScreenHeading } from '../components/ScreenHeading'
import { createSolucao, getSolucoes, updateSolucao } from '../services/solucaoService'
import { getStartupOptions } from '../services/startupService'
import { errorMessage } from '../services/http'
import type { CreateSolucaoRequest, SolucaoResponse, StartupResponse } from '../types/domain'

const initial: CreateSolucaoRequest = { startup_id: 0, nome: '', descricao: '', categoria: 'estoque', palavras_chave: [] }

export default function SolucoesPage({ onOpenStartups }: { onOpenStartups: () => void }) {
  const [items, setItems] = useState<SolucaoResponse[]>([])
  const [startups, setStartups] = useState<StartupResponse[]>([])
  const [record, setRecord] = useState<SolucaoResponse | null>(null)
  const [form, setForm] = useState(initial)
  const [keywords, setKeywords] = useState('')
  const [dialogOpen, setDialogOpen] = useState(false)
  const [page, setPage] = useState(1)
  const [revision, setRevision] = useState(0)
  const [startupRevision, setStartupRevision] = useState(0)
  const [loading, setLoading] = useState(true)
  const [startupsLoading, setStartupsLoading] = useState(true)
  const [busy, setBusy] = useState(false)
  const [formError, setFormError] = useState('')
  const [listError, setListError] = useState('')
  const [startupError, setStartupError] = useState('')
  const [message, setMessage] = useState('')

  useEffect(() => {
    let active = true
    getSolucoes(page).then(result => {
      if (active) { setItems(result); setListError('') }
    }).catch(err => {
      if (active) { setItems([]); setListError(errorMessage(err)) }
    }).finally(() => { if (active) setLoading(false) })
    return () => { active = false }
  }, [page, revision])

  useEffect(() => {
    let active = true
    const controller = new AbortController()
    getStartupOptions(controller.signal).then(result => {
      if (active) { setStartups(result); setStartupError('') }
    }).catch(err => {
      if (active) setStartupError(errorMessage(err))
    }).finally(() => { if (active) setStartupsLoading(false) })
    return () => { active = false; controller.abort() }
  }, [startupRevision])

  function reload() { setLoading(true); setRevision(value => value + 1) }
  function reloadStartups() { setStartupsLoading(true); setStartupRevision(value => value + 1) }

  function openForm(item?: SolucaoResponse) {
    setRecord(item ?? null)
    setForm(item ? { startup_id: item.startup_id, nome: item.nome, descricao: item.descricao, categoria: item.categoria, palavras_chave: item.palavras_chave } : { ...initial })
    setKeywords(item?.palavras_chave.join(', ') ?? '')
    setFormError(''); setMessage(''); setDialogOpen(true)
  }

  async function save(event: FormEvent) {
    event.preventDefault()
    setBusy(true); setFormError('')
    try {
      if (!startups.some(startup => startup.id === form.startup_id)) throw new Error('Selecione a startup responsável pela solução.')
      const data = { ...form, palavras_chave: keywords.split(',').map(value => value.trim()).filter(Boolean) }
      const result = record ? await updateSolucao(record.id, data) : await createSolucao(data)
      setDialogOpen(false)
      setMessage(`${result.nome}: solução ${record ? 'atualizada' : 'cadastrada'} com sucesso.`)
      if (!record) setPage(1)
      reload()
    } catch (err) { setFormError(errorMessage(err)) }
    finally { setBusy(false) }
  }

  const startupNames = new Map(startups.map(startup => [startup.id, startup.nome]))

  return <div className="directory-page">
    <ScreenHeading icon="cube" title="Catálogo de soluções" subtitle="Explore as tecnologias disponíveis para resolver os desafios do seu negócio." />
    <Notice message={message} />
    <div className="directory-toolbar"><div><h2>Soluções disponíveis</h2><p>Descubra o que as startups da rede oferecem.</p></div><div className="directory-actions"><button type="button" className="text-button" disabled={loading} onClick={reload}>Atualizar lista</button><button type="button" onClick={() => openForm()}><Icon name="plus" size={18} />Cadastrar solução</button></div></div>
    <Notice error={listError} />
    {loading ? <div className="directory-loading" role="status">Carregando soluções…</div>
      : listError ? <div className="empty-state"><h2>Não foi possível carregar as soluções.</h2><button className="secondary" type="button" onClick={reload}>Tentar novamente</button></div>
      : items.length === 0 ? <div className="empty-state"><Icon name="cube" size={38} /><h2>{page === 1 ? 'Nenhuma solução cadastrada ainda.' : 'Nenhuma solução nesta página.'}</h2><p>As tecnologias cadastradas pelas startups aparecerão aqui.</p>{page === 1 && <button type="button" className="secondary" onClick={() => openForm()}>Cadastrar primeira solução</button>}</div>
      : <div className="directory-grid">{items.map(item => <article className="directory-card" key={item.id}><div className="directory-card-top"><span className="solution-icon variant-0"><Icon name="cube" size={26} /></span><span className="badge">{item.categoria}</span><button type="button" className="text-button" aria-label={`Editar ${item.nome}`} onClick={() => openForm(item)}>Editar</button></div><h2>{item.nome}</h2><span className="directory-company">{startupNames.get(item.startup_id) ?? (startupsLoading ? 'Carregando startup…' : 'Startup não disponível')}</span><p className="directory-description">{item.descricao}</p><div className="tags">{item.palavras_chave.map((word, index) => <span key={`${index}-${word}`}>{word}</span>)}</div></article>)}</div>}
    <Pagination page={page} count={items.length} loading={loading} onChange={value => { setLoading(true); setPage(value) }} />
    <FormDialog open={dialogOpen} title={record ? 'Editar solução' : 'Cadastrar solução'} busy={busy} onClose={() => setDialogOpen(false)}>
      <Notice error={formError} />
      <Notice error={startupError} />
      {startupError && <button type="button" className="text-button" disabled={startupsLoading} onClick={reloadStartups}>Tentar carregar startups novamente</button>}
      {!startupsLoading && !startupError && startups.length === 0 && <div className="info-strip"><Icon name="info" size={18} /><div>Cadastre uma startup antes de adicionar uma solução.<button type="button" className="text-button" onClick={onOpenStartups}>Ir para startups</button></div></div>}
      <form onSubmit={save}><fieldset disabled={busy}>
        <Field label="Startup responsável" required><select required disabled={startupsLoading || !!startupError || startups.length === 0} value={form.startup_id || ''} onChange={event => setForm({ ...form, startup_id: Number(event.target.value) })}><option value="" disabled>{startupsLoading ? 'Carregando startups…' : 'Selecione uma startup'}</option>{record && !startupNames.has(form.startup_id) && form.startup_id !== 0 && <option value={form.startup_id} disabled>Startup não disponível — selecione outra</option>}{startups.map(startup => <option key={startup.id} value={startup.id}>{startup.nome}</option>)}</select></Field>
        <Field label="Nome da solução" required><input required maxLength={150} value={form.nome} onChange={event => setForm({ ...form, nome: event.target.value })} placeholder="Ex.: Gestão inteligente de estoque" /></Field>
        <Field label="Descrição" required><textarea required maxLength={5000} rows={4} value={form.descricao} onChange={event => setForm({ ...form, descricao: event.target.value })} placeholder="Qual problema sua tecnologia resolve?" /></Field>
        <Field label="Categoria"><CategorySelect value={form.categoria} onChange={categoria => setForm({ ...form, categoria })} /></Field>
        <Field label="Palavras-chave (separadas por vírgula)"><input value={keywords} onChange={event => setKeywords(event.target.value)} placeholder="estoque, reposição, automação" /></Field>
        <div className="dialog-actions"><button type="button" className="secondary" onClick={() => setDialogOpen(false)}>Cancelar</button><button type="submit" disabled={startupsLoading || !!startupError || startups.length === 0}>{busy ? 'Salvando…' : record ? 'Salvar alterações' : 'Cadastrar solução'}</button></div>
      </fieldset></form>
    </FormDialog>
  </div>
}
