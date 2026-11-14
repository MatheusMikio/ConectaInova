import { useState } from 'react'
import type { FormEvent } from 'react'
import { Field } from '../components/Field'
import { Notice } from '../components/Notice'
import { createLojista, getLojista, updateLojista } from '../services/lojistaService'
import { errorMessage } from '../services/http'
import type { LojistaResponse } from '../types/domain'

export default function LojistasPage() {
  const [record, setRecord] = useState<LojistaResponse | null>(null)
  const [nome, setNome] = useState('')
  const [segmento, setSegmento] = useState('')
  const [lookup, setLookup] = useState('')
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState('')
  const [message, setMessage] = useState('')

  function show(result: LojistaResponse) {
    setRecord(result)
    setNome(result.nome)
    setSegmento(result.segmento)
  }

  async function search(event: FormEvent) {
    event.preventDefault()
    setBusy(true); setError(''); setMessage('')
    try {
      show(await getLojista(Number(lookup)))
    } catch (err) { setError(errorMessage(err)) }
    finally { setBusy(false) }
  }

  async function save(event: FormEvent) {
    event.preventDefault()
    setBusy(true); setError(''); setMessage('')
    try {
      const result = record
        ? await updateLojista(record.id, { nome, segmento })
        : await createLojista({ nome, segmento })
      show(result)
      setMessage(`Lojista #${result.id} salvo. Use esse ID para cadastrar suas necessidades.`)
    } catch (err) { setError(errorMessage(err)) }
    finally { setBusy(false) }
  }

  return <>
    <div className="page-heading"><div><span className="eyebrow">01 / Varejo</span><h1>Quem precisa inovar?</h1><p>Cadastre seu negócio e dê o primeiro passo para encontrar soluções.</p></div></div>
    <Notice error={error} message={message} />
    <div className="two-columns">
      <section className="panel">
        <div className="panel-heading"><h2>{record ? `Editar lojista #${record.id}` : 'Cadastrar lojista'}</h2>{record && <button className="text-button" disabled={busy} onClick={() => { setRecord(null); setNome(''); setSegmento(''); setMessage(''); setError('') }}>Novo cadastro</button>}</div>
        <form onSubmit={save}><fieldset disabled={busy}>
          <Field label="Nome da loja"><input required maxLength={150} value={nome} onChange={event => setNome(event.target.value)} placeholder="Ex.: Boutique da Ana" /></Field>
          <Field label="Segmento"><input required maxLength={150} value={segmento} onChange={event => setSegmento(event.target.value)} placeholder="Ex.: Moda e vestuário" /></Field>
          <button type="submit">{busy ? 'Aguarde…' : record ? 'Salvar alterações' : 'Cadastrar negócio'} <span aria-hidden="true">↗</span></button>
        </fieldset></form>
      </section>
      <aside className="stack">
        <section className="panel">
          <h2>Já tem cadastro?</h2><p>Consulte pelo ID retornado ao cadastrar seu negócio.</p>
          <form className="lookup" onSubmit={search}><Field label="ID do lojista"><input type="number" min="1" step="1" required value={lookup} onChange={event => setLookup(event.target.value)} /></Field><button type="submit" className="secondary" disabled={busy}>Consultar</button></form>
        </section>
        <div className="tip"><span className="eyebrow">Do problema à oportunidade</span><h2>Sua próxima melhoria começa com uma conversa.</h2><p>Estoque, vendas, atendimento ou logística: transforme o que dificulta sua rotina em uma necessidade de inovação.</p></div>
      </aside>
    </div>
  </>
}
