import type { InnovationFlow } from '../hooks/useInnovationFlow'
import ChallengePage from './ChallengePage'
import DesafiosPage from './DesafiosPage'

export default function DesafiosSection({ flow }: { flow: InnovationFlow }) {
  const draftOpen = flow.screen === 'desafio-draft'
  return <div className="challenge-section">
    {flow.draft && <nav className="section-navigation" aria-label="Desafios de inovação">
      <button type="button" className={!draftOpen ? 'selected' : ''} aria-current={!draftOpen ? 'page' : undefined} disabled={flow.busy} onClick={() => flow.navigate('desafio')}>Consultar e gerenciar</button>
      <button type="button" className={draftOpen ? 'selected' : ''} aria-current={draftOpen ? 'page' : undefined} disabled={flow.busy} onClick={() => flow.navigate('desafio-draft')}>{flow.published ? 'Desafio publicado' : 'Desafio gerado'}</button>
    </nav>}
    {draftOpen ? <ChallengePage flow={flow} /> : <DesafiosPage
      key={flow.challengeProblem?.id ?? 'all'}
      problem={flow.challengeProblem}
      onUpdated={flow.challengeUpdated}
      onClearProblem={() => flow.navigate('desafio')}
      onCreate={() => {
        if (flow.challengeProblem) void flow.openProblem(flow.challengeProblem.id, 'matching')
        else flow.navigate(flow.need ? 'matching' : 'cadastro')
      }}
    />}
  </div>
}
