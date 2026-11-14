import { Icon } from '../components/Icon'
import { ScreenHeading } from '../components/ScreenHeading'
import type { InnovationFlow } from '../hooks/useInnovationFlow'

export default function InicioPage({ flow }: { flow: InnovationFlow }) {
  return <div className="flow-page">
    <ScreenHeading icon="home" title="Inovação que começa no seu negócio" subtitle="Conectamos os desafios do varejo às soluções de startups." />
    <section className="welcome-card"><span className="welcome-tag">VAREJO + TECNOLOGIA</span><h2>Seu problema pode encontrar<br />uma nova possibilidade.</h2><p>Conte sua necessidade, descubra soluções compatíveis ou gere um desafio de inovação para o ecossistema.</p><button onClick={() => flow.navigate('cadastro')}><Icon name="plus" />Cadastrar problema</button></section>
    <div className="steps-grid"><article><Icon name="clipboard" size={26} /><span>01</span><h3>Conte seu problema</h3><p>Explique o que você precisa melhorar no seu negócio.</p></article><article><Icon name="search" size={26} /><span>02</span><h3>Encontre soluções</h3><p>Conheça startups compatíveis com a sua necessidade.</p></article><article><Icon name="target" size={26} /><span>03</span><h3>Abra um novo desafio</h3><p>Se nenhuma solução atender, crie uma oportunidade de inovação.</p></article></div>
    <div className="info-strip"><Icon name="info" /><span>Acesse Problemas cadastrados para acompanhar as necessidades dos lojistas, buscar soluções e consultar os desafios vinculados.</span></div>
  </div>
}
