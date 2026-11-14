import { useEffect, useState } from 'react'
import { Icon, Logo } from './components/Icon'
import type { IconName } from './components/Icon'
import { Notice } from './components/Notice'
import { useInnovationFlow } from './hooks/useInnovationFlow'
import DesafiosSection from './pages/DesafiosSection'
import ProblemasPage from './pages/ProblemasPage'
import InicioPage from './pages/InicioPage'
import MatchingPage from './pages/MatchingPage'
import ProblemPage from './pages/ProblemPage'
import SolucoesPage from './pages/SolucoesPage'
import StartupsPage from './pages/StartupsPage'
import AvaliacoesPage from './pages/AvaliacoesPage'
import { getPlataforma } from './services/avaliacaoService'
import type { Screen } from './types/workflow'
import './App.css'

const navigation: { key: Screen; label: string; icon: IconName }[] = [
  { key: 'inicio', label: 'Início', icon: 'home' },
  { key: 'cadastro', label: 'Cadastrar problema', icon: 'plus' },
  { key: 'matching', label: 'Matching (soluções)', icon: 'search' },
  { key: 'solucoes', label: 'Soluções', icon: 'cube' },
  { key: 'startups', label: 'Startups', icon: 'rocket' },
  { key: 'desafio', label: 'Desafios de inovação', icon: 'target' },
  { key: 'problemas', label: 'Problemas cadastrados', icon: 'chart' },
  { key: 'avaliacoes', label: 'Avaliações', icon: 'clipboard' },
]

function App() {
  const flow = useInnovationFlow()
  const [menu, setMenu] = useState<'profile' | 'notifications' | null>(null)
  const activeScreen = flow.screen === 'desafio-draft' ? 'desafio' : flow.screen
  const [platform, setPlatform] = useState<{ associacao: string; regiao: string } | null>(null)
  useEffect(() => {
    let active = true
    getPlataforma().then(result => { if (active) setPlatform(result) }).catch(() => { /* Branding does not block navigation. */ })
    return () => { active = false }
  }, [])

  return <div className="prototype">
    <a className="skip-link" href="#conteudo">Pular para o conteúdo</a>
    <header className="presentation-header"><h1>ConectaInova</h1></header>
    <div className="product-window">
      <header className="app-header">
        <div className="brand-group"><button type="button" className="brand-button" aria-label="ConectaInova — início" disabled={flow.busy} onClick={() => flow.navigate('inicio')}><Logo /></button>{platform && <span className="association-caption">{platform.associacao}{platform.regiao ? ` · ${platform.regiao}` : ''}</span>}</div>
        <div className="header-actions"><div className="popover-anchor"><button type="button" className="icon-button" aria-label="Notificações" aria-expanded={menu === 'notifications'} onClick={() => setMenu(menu === 'notifications' ? null : 'notifications')}><Icon name="bell" size={21} /></button>{menu === 'notifications' && <div className="header-popover"><strong>Notificações</strong><p>Não há notificações nesta sessão.</p></div>}</div><div className="popover-anchor"><button type="button" className="profile-button" aria-label="Menu do varejo" aria-expanded={menu === 'profile'} onClick={() => setMenu(menu === 'profile' ? null : 'profile')}><span className="avatar"><Icon name="user" size={19} /></span><span>Varejo</span><Icon name="chevron" size={14} /></button>{menu === 'profile' && <div className="header-popover profile-popover"><strong>ConectaInova</strong><button type="button" onClick={() => { flow.navigate('problemas'); setMenu(null) }}>Problemas cadastrados</button></div>}</div></div>
      </header>
      <div className="app-body">
        <aside className="sidebar"><nav aria-label="Navegação principal">{navigation.map(item => <button key={item.key} type="button" disabled={flow.busy} className={activeScreen === item.key ? 'nav-item active' : 'nav-item'} aria-current={activeScreen === item.key ? 'page' : undefined} onClick={() => { flow.navigate(item.key); setMenu(null) }}><Icon name={item.icon} size={20} /><span>{item.label}</span></button>)}</nav><div className="sidebar-caption">Inovação para o pequeno varejo</div></aside>
        <main id="conteudo" aria-busy={flow.busy}>
          <Notice error={flow.error} message={flow.message} />
          {flow.screen === 'inicio' && <InicioPage flow={flow} />}
          {flow.screen === 'cadastro' && <ProblemPage flow={flow} />}
          {flow.screen === 'matching' && <MatchingPage flow={flow} />}
          {flow.screen === 'solucoes' && <SolucoesPage onOpenStartups={() => flow.navigate('startups')} />}
          {flow.screen === 'startups' && <StartupsPage />}
          {(flow.screen === 'desafio' || flow.screen === 'desafio-draft') && <DesafiosSection flow={flow} />}
          {flow.screen === 'problemas' && <ProblemasPage flow={flow} />}
          {flow.screen === 'avaliacoes' && <AvaliacoesPage />}
        </main>
      </div>
    </div>
    <footer className="prototype-footer"><Logo /><span>Conectando necessidades reais a novas possibilidades.</span></footer>
  </div>
}

export default App
