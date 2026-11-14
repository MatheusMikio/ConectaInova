import { Icon } from '../components/Icon'
import { ScreenHeading } from '../components/ScreenHeading'
import type { InnovationFlow } from '../hooks/useInnovationFlow'

export default function ChallengePage({ flow }: { flow: InnovationFlow }) {
  const { draft, need, published, busy } = flow
  if (!need || !draft) return <div className="flow-page"><ScreenHeading icon="target" title="Desafios de inovação" subtitle="Novas oportunidades para solucionar problemas do varejo." /><div className="empty-state"><Icon name="target" size={40} /><h2>Um problema pode ser o começo de uma nova solução.</h2><p>Busque soluções para uma necessidade e gere um desafio a partir do matching.</p><button onClick={() => flow.navigate(need ? 'matching' : 'cadastro')}>Começar</button></div></div>

  return <div className="flow-page">
    <ScreenHeading icon="target" title={published ? 'Desafio de inovação publicado' : 'Desafio de inovação gerado'} subtitle="Uma oportunidade de desenvolver uma nova solução para o ecossistema." />
    <div className={`info-strip ${flow.matches?.sugerir_desafio ? 'warning-strip' : ''}`}><Icon name={flow.matches?.sugerir_desafio ? 'warning' : 'info'} size={20} /><span>{flow.matches?.sugerir_desafio ? 'Nenhuma solução cadastrada é compatível com a categoria deste problema.' : 'As soluções encontradas não atenderam? Proponha um novo desafio de inovação.'}</span></div>
    <section className="challenge-card">
      <div className="challenge-card-heading"><div className="solution-icon variant-0"><Icon name="cube" size={29} /></div><div><label className="sr-only" htmlFor="challenge-title">Título do desafio</label><input id="challenge-title" className="editable-title" maxLength={200} required disabled={busy || !!published} value={draft.titulo} onChange={event => flow.setDraft({ ...draft, titulo: event.target.value })} /><p>Uma solução acessível para pequenos lojistas, com implantação simples e foco nas necessidades reais do negócio.</p></div><span className="open-badge"><span />{published ? ({ aberto: 'Aberto para startups', em_andamento: 'Em andamento', encerrado: 'Encerrado' })[published.status] : 'Rascunho do desafio'}</span></div>
      <dl className="challenge-details">
        <div><dt><Icon name="clipboard" />Problema original</dt><dd>{need.titulo}.</dd></div>
        <div><dt><Icon name="target" />Objetivo do desafio</dt><dd><label className="sr-only" htmlFor="challenge-description">Objetivo do desafio</label><textarea id="challenge-description" className="editable-description" maxLength={5000} required rows={3} disabled={busy || !!published} value={draft.descricao} onChange={event => flow.setDraft({ ...draft, descricao: event.target.value })} /></dd></div>
        <div><dt><Icon name="list" />Critérios de sucesso</dt><dd><label className="sr-only" htmlFor="challenge-criteria">Critérios de sucesso</label><textarea id="challenge-criteria" rows={4} maxLength={6000} disabled={busy || !!published} value={draft.criterios} onChange={event => flow.setDraft({ ...draft, criterios: event.target.value })} placeholder="Um critério por linha" /></dd></div>
        <div><dt><Icon name="calendar" />Prazo para propostas</dt><dd><label className="sr-only" htmlFor="challenge-deadline">Prazo para propostas</label><select id="challenge-deadline" className="deadline-select" value={draft.prazo} disabled={busy || !!published} onChange={event => flow.setDraft({ ...draft, prazo: Number(event.target.value) })}>{![7, 14, 30].includes(draft.prazo) && <option value={draft.prazo}>{draft.prazo} dias após a publicação</option>}<option value={7}>7 dias após a publicação</option><option value={14}>14 dias após a publicação</option><option value={30}>30 dias após a publicação</option></select></dd></div>
      </dl>
    </section>
    <p className="contract-note">Revise o objetivo, os critérios de sucesso e o prazo antes de publicar. O desafio fica vinculado ao briefing do problema.</p>
    <div className="challenge-actions"><button type="button" disabled={busy || !!published} onClick={flow.publishChallenge}><Icon name={published ? 'check' : 'rocket'} size={21} />{busy ? 'Publicando…' : published ? 'Desafio publicado' : 'Publicar desafio'}</button><button type="button" className="secondary" disabled={busy} onClick={() => flow.navigate('matching')}><Icon name="arrowLeft" size={22} />Voltar ao matching</button></div>
  </div>
}
