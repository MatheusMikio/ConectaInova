import { useState } from 'react'
import { Icon } from '../components/Icon'
import type { IconName } from '../components/Icon'
import { Field } from '../components/Field'
import { FormDialog } from '../components/FormDialog'
import { Notice } from '../components/Notice'
import { ProblemBrief } from '../components/ProblemBrief'
import { ScreenHeading } from '../components/ScreenHeading'
import type { InnovationFlow } from '../hooks/useInnovationFlow'
import { createAvaliacao } from '../services/avaliacaoService'
import { errorMessage } from '../services/http'
import type { AvaliacaoResponse, MatchResponse } from '../types/domain'
import ProblemasPage from './ProblemasPage'

function adherence(score: number): { label: string; tone: string; icon: IconName } {
  if (score >= 90) return { label: 'Alta aderência', tone: 'high', icon: 'trophy' }
  if (score >= 80) return { label: 'Boa aderência', tone: 'good', icon: 'check' }
  return { label: 'Aderência moderada', tone: 'moderate', icon: 'tag' }
}

export default function MatchingPage({ flow }: { flow: InnovationFlow }) {
  const [selected, setSelected] = useState<MatchResponse | null>(null)
  const [message, setMessage] = useState('')
  const [error, setError] = useState('')
  const [busy, setBusy] = useState(false)
  const [created, setCreated] = useState<AvaliacaoResponse | null>(null)
  const [choosingProblem, setChoosingProblem] = useState(false)
  const qualified = !!flow.need?.resultado_esperado?.trim() && !!flow.need?.criterios_sucesso?.length

  function openDetails(item: MatchResponse) {
    setSelected(item); setMessage(''); setError(''); setCreated(null)
  }

  async function requestEvaluation() {
    if (!flow.need || !selected) return
    setBusy(true); setError('')
    try {
      setCreated(await createAvaliacao(flow.need.id, selected.solucao.id, message))
    } catch (err) { setError(errorMessage(err)) }
    finally { setBusy(false) }
  }

  async function selectProblem(id: number) {
    if (await flow.openProblem(id, 'matching')) setChoosingProblem(false)
  }

  return <div className="flow-page">
    <ScreenHeading icon="search" title={flow.need && !choosingProblem ? 'Soluções encontradas' : 'Matching de soluções'} subtitle="Selecione um problema, entenda as recomendações e avance para uma avaliação com a startup." />
    {!flow.need || choosingProblem ? <>
      {flow.need && <button type="button" className="text-button" disabled={flow.busy} onClick={() => setChoosingProblem(false)}>Voltar aos resultados</button>}
      <ProblemasPage flow={flow} matchingSelection onSelect={selectProblem} />
    </> : <>
      <div className="problem-summary"><Icon name="cube" size={23} /><span><strong>Problema analisado:</strong> {flow.need.titulo}</span><button type="button" className="text-button" disabled={flow.busy} onClick={() => setChoosingProblem(true)}>Trocar problema</button></div>
      {flow.busy ? <div className="empty-state" role="status">Buscando soluções compatíveis…</div> : !flow.matches ? <div className="empty-state"><p>A busca ainda não foi concluída. Seu problema já está salvo.</p><button onClick={flow.refreshMatches}><Icon name="search" />Buscar novamente</button></div> : <>
        <div className="solution-list">{flow.matches.matches.map((item, index) => {
          const badge = adherence(item.pontuacao)
          return <article className="solution-card" key={item.solucao.id}>
            <div className={`solution-icon variant-${index % 3}`}><Icon name={index % 3 === 0 ? 'cube' : index % 3 === 1 ? 'chart' : 'tag'} size={29} /></div>
            <div className="solution-description"><h2>{item.solucao.nome}</h2><p>{item.solucao.descricao}</p></div>
            <div className="solution-compatibility"><span className={`adherence ${badge.tone}`}><Icon name={badge.icon} size={13} />{badge.label}</span><span>Índice de aderência: <strong>{item.pontuacao.toLocaleString('pt-BR')}/100</strong></span><progress max="100" value={item.pontuacao} aria-label={`Aderência de ${item.solucao.nome}`} /></div>
            <button className="view-solution" type="button" onClick={() => openDetails(item)}>Avaliar solução</button>
            <div className="match-context"><strong>Por que foi recomendada?</strong>{(item.motivos || []).map(reason => <p key={reason}>{reason}</p>)}</div>
          </article>
        })}</div>
        {flow.matches.matches.length === 0 && <div className="empty-state"><Icon name="search" size={36} /><h2>Nenhuma solução compatível por enquanto.</h2><p>Transforme essa necessidade em uma oportunidade para o ecossistema.</p></div>}
        <div className="info-strip"><Icon name="info" size={19} /><span>O índice combina categoria e palavras-chave. Preço, implantação e integrações precisam ser avaliados com a startup.</span></div>
        <div className="matching-actions"><button className="secondary" type="button" onClick={flow.prepareChallenge}><Icon name="plus" size={19} />{flow.matches.matches.length ? 'Nenhuma solução atendeu? Gerar desafio' : 'Gerar desafio de inovação'}</button></div>
      </>}
    </>}
    <FormDialog open={!!selected} title={selected?.solucao.nome || 'Avaliar solução'} busy={busy} onClose={() => setSelected(null)}>
      {selected && flow.need && <>
        <p className="evaluation-intro">{selected.solucao.descricao}</p>
        <div className="detail-grid"><span>Startup<strong>{selected.startup.nome}</strong></span><span>Categoria<strong>{selected.solucao.categoria}</strong></span><span>Índice de aderência<strong>{selected.pontuacao}/100</strong></span></div>
        <h3 className="brief-title">Por que essa solução é pertinente?</h3><ul className="reason-list">{(selected.motivos || []).map(reason => <li key={reason}>{reason}</li>)}</ul>
        <h3 className="brief-title">O que ainda precisa ser confirmado</h3><ul className="reason-list">{(selected.pontos_a_validar || []).map(point => <li key={point}>{point}</li>)}</ul>
        <details className="brief-details"><summary>Briefing que a startup poderá avaliar</summary><ProblemBrief problem={flow.need} loja={flow.shop?.nome} /></details>
        <Notice error={error} />
        {created ? <><Notice message="Solicitação registrada. O briefing está disponível na área de avaliações." /><button type="button" className="full-width" onClick={() => flow.navigate('avaliacoes')}>Acompanhar avaliação</button></>
          : !qualified ? <div className="info-strip"><Icon name="info" /><div>Complete o resultado esperado e os critérios de sucesso para a startup avaliar a demanda.<button type="button" className="text-button" onClick={() => flow.navigate('cadastro')}>Completar briefing do problema</button></div></div>
            : <><Field label="Mensagem para a avaliação (opcional)"><textarea rows={3} maxLength={2000} value={message} disabled={busy} onChange={event => setMessage(event.target.value)} placeholder="Quais pontos você gostaria de esclarecer antes de testar?" /></Field><button type="button" className="full-width" disabled={busy} onClick={requestEvaluation}><Icon name="rocket" size={19} />{busy ? 'Registrando…' : 'Solicitar avaliação'}</button><p className="contract-note">O pedido fica disponível na plataforma. Não há envio automático de e-mail nem agendamento de piloto.</p></>}
      </>}
    </FormDialog>
  </div>
}
