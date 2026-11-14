import { useEffect, useState } from 'react'
import type { FormEvent } from 'react'
import { Field } from '../components/Field'
import { FormDialog } from '../components/FormDialog'
import { Icon } from '../components/Icon'
import { Notice } from '../components/Notice'
import { Pagination } from '../components/Pagination'
import { ScreenHeading } from '../components/ScreenHeading'
import { createStartup, getStartups, updateStartup } from '../services/startupService'
import { errorMessage } from '../services/http'
import type { StartupResponse } from '../types/domain'

export default function StartupsPage() {
  const [items, setItems] = useState<StartupResponse[]>([])
  const [record, setRecord] = useState<StartupResponse | null>(null)
  const [nome, setNome] = useState('')
  const [descricao, setDescricao] = useState('')
  const [dialogOpen, setDialogOpen] = useState(false)
  const [page, setPage] = useState(1)
  const [revision, setRevision] = useState(0)
  const [loading, setLoading] = useState(true)
  const [busy, setBusy] = useState(false)
  const [formError, setFormError] = useState('')
  const [listError, setListError] = useState('')
  const [message, setMessage] = useState('')

  useEffect(() => {
    let active = true
    getStartups(page).then(result => {
      if (active) { setItems(result); setListError('') }
    }).catch(err => {
      if (active) { setItems([]); setListError(errorMessage(err)) }
    }).finally(() => { if (active) setLoading(false) })
    return () => { active = false }
  }, [page, revision])

  function reload() { setLoading(true); setRevision(value => value + 1) }

  function openForm(item?: StartupResponse) {
    setRecord(item ?? null)
    setNome(item?.nome ?? '')
    setDescricao(item?.descricao ?? '')
    setFormError(''); setMessage(''); setDialogOpen(true)
  }

  async function save(event: FormEvent) {
    event.preventDefault()
    setBusy(true); setFormError('')
    try {
      const result = record
        ? await updateStartup(record.id, { nome, descricao })
        : await createStartup({ nome, descricao })
      setDialogOpen(false)
      setMessage(`${result.nome}: ${record ? 'cadastro atualizado' : 'startup cadastrada'} com sucesso.`)
      if (!record) setPage(1)
      reload()
    } catch (err) { setFormError(errorMessage(err)) }
    finally { setBusy(false) }
  }

  return <div className="directory-page">
    <ScreenHeading icon="rocket" title="Startups" subtitle="Conheça as empresas que estão transformando o varejo." />
    <Notice message={message} />
    <div className="directory-toolbar"><div><h2>Startups da rede</h2><p>Explore o ecossistema de inovação.</p></div><div className="directory-actions"><button type="button" className="text-button" disabled={loading} onClick={reload}>Atualizar lista</button><button type="button" onClick={() => openForm()}><Icon name="plus" size={18} />Cadastrar startup</button></div></div>
    <Notice error={listError} />
    {loading ? <div className="directory-loading" role="status">Carregando startups…</div>
      : listError ? <div className="empty-state"><h2>Não foi possível carregar as startups.</h2><button className="secondary" type="button" onClick={reload}>Tentar novamente</button></div>
      : items.length === 0 ? <div className="empty-state"><Icon name="rocket" size={38} /><h2>{page === 1 ? 'Nenhuma startup cadastrada ainda.' : 'Nenhuma startup nesta página.'}</h2><p>As startups cadastradas na plataforma aparecerão aqui.</p>{page === 1 && <button type="button" className="secondary" onClick={() => openForm()}>Cadastrar primeira startup</button>}</div>
      : <div className="directory-grid">{items.map(item => <article className="directory-card" key={item.id}><div className="directory-card-top"><span className="solution-icon variant-0"><Icon name="rocket" size={26} /></span><button type="button" className="text-button" onClick={() => openForm(item)} aria-label={`Editar ${item.nome}`}>Editar</button></div><h2>{item.nome}</h2><p className="directory-description">{item.descricao}</p></article>)}</div>}
    <Pagination page={page} count={items.length} loading={loading} onChange={value => { setLoading(true); setPage(value) }} />
    <FormDialog open={dialogOpen} title={record ? 'Editar startup' : 'Cadastrar startup'} busy={busy} onClose={() => setDialogOpen(false)}>
      <Notice error={formError} />
      <form onSubmit={save}><fieldset disabled={busy}>
        <Field label="Nome da startup" required><input required maxLength={150} value={nome} onChange={event => setNome(event.target.value)} placeholder="Nome da empresa" /></Field>
        <Field label="Sobre a startup" required><textarea required rows={5} maxLength={5000} value={descricao} onChange={event => setDescricao(event.target.value)} placeholder="O que a startup faz e como ajuda o varejo?" /></Field>
        <div className="dialog-actions"><button type="button" className="secondary" onClick={() => setDialogOpen(false)}>Cancelar</button><button type="submit">{busy ? 'Salvando…' : record ? 'Salvar alterações' : 'Cadastrar startup'}</button></div>
      </fieldset></form>
    </FormDialog>
  </div>
}
