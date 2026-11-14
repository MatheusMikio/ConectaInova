import { Field } from '../components/Field'
import { Icon } from '../components/Icon'
import { ScreenHeading } from '../components/ScreenHeading'
import type { InnovationFlow } from '../hooks/useInnovationFlow'
import type { Categoria, Prioridade } from '../types/domain'

const categories: { value: Categoria; label: string }[] = [
  { value: 'estoque', label: 'Estoque' }, { value: 'vendas', label: 'Vendas' },
  { value: 'atendimento', label: 'Atendimento' }, { value: 'financeiro', label: 'Financeiro' },
  { value: 'marketing', label: 'Marketing' }, { value: 'logistica', label: 'Logística' },
]

export default function ProblemPage({ flow }: { flow: InnovationFlow }) {
  const { form, setForm, busy } = flow
  return <div className="flow-page">
    <ScreenHeading icon="clipboard" title={flow.need ? 'Editar problema' : 'Cadastro de problema'} subtitle="Conte qual é o desafio do seu negócio e encontre soluções inovadoras." />
    {flow.need && <div className="saved-problem"><span>Problema #{flow.need.id} · {flow.need.status.replace('_', ' ')}</span><button type="button" className="text-button" disabled={busy} onClick={flow.newProblem}>Cadastrar outro problema</button></div>}
    <form onSubmit={flow.submitProblem} className="problem-form">
      <fieldset disabled={busy}>
        <div className="form-row shop-row">
          <Field label="Nome da loja" required><input required maxLength={150} value={form.nomeLoja} onChange={event => setForm({ ...form, nomeLoja: event.target.value })} placeholder="Nome do seu negócio" autoComplete="organization" /></Field>
          <Field label="Segmento" required><select value={form.segmento} onChange={event => setForm({ ...form, segmento: event.target.value })}>
            {!['Varejo de moda', 'Alimentação', 'Mercado e conveniência', 'Saúde e beleza', 'Casa e decoração', 'Serviços', 'Outros'].includes(form.segmento) && <option>{form.segmento}</option>}
            <option>Varejo de moda</option><option>Alimentação</option><option>Mercado e conveniência</option><option>Saúde e beleza</option><option>Casa e decoração</option><option>Serviços</option><option>Outros</option>
          </select></Field>
        </div>
        <Field label="Título do problema" required><input required maxLength={200} value={form.titulo} onChange={event => setForm({ ...form, titulo: event.target.value })} placeholder="Ex.: Dificuldade para controlar o estoque" /></Field>
        <Field label="Descreva sua necessidade" required><textarea required rows={5} maxLength={5000} value={form.descricao} onChange={event => setForm({ ...form, descricao: event.target.value })} placeholder="Explique o problema, como ele afeta sua rotina e o que você precisa melhorar." /></Field>
        <Field label="Como sua loja vende?" required><select required value={form.canalVenda} onChange={event => setForm({ ...form, canalVenda: event.target.value })}><option value="" disabled>Selecione o perfil da operação</option><option value="fisica">Loja física</option><option value="digital">Loja digital</option><option value="hibrida">Loja física e digital</option></select></Field>
        <div className="form-row">
          <Field label="Impacto no negócio" required><select required value={form.impacto} onChange={event => setForm({ ...form, impacto: event.target.value })}><option value="" disabled>Selecione o impacto</option><option>Perda de vendas e clientes</option><option>Aumento de custos</option><option>Perda de produtividade</option><option>Experiência do cliente</option><option>Outros</option></select></Field>
          <Field label="Prioridade" required><div className={`priority-select ${form.prioridade}`}><span aria-hidden="true" /><select value={form.prioridade} onChange={event => setForm({ ...form, prioridade: event.target.value as Prioridade })}><option value="alta">Alta</option><option value="media">Média</option><option value="baixa">Baixa</option></select></div></Field>
        </div>
        <Field label="Resultado esperado" required><textarea required maxLength={2000} rows={2} value={form.resultadoEsperado} onChange={event => setForm({ ...form, resultadoEsperado: event.target.value })} placeholder="O que precisa melhorar para a solução valer a pena?" /></Field>
        <Field label="Critérios de sucesso — um por linha" required><textarea required maxLength={6000} rows={3} value={form.criteriosSucesso} onChange={event => setForm({ ...form, criteriosSucesso: event.target.value })} placeholder={'Ex.: Receber alertas antes de um produto acabar\nConseguir operar sem treinamento longo'} /></Field>
        <Field label="Restrições e condições para testar"><textarea maxLength={2000} rows={2} value={form.restricoes} onChange={event => setForm({ ...form, restricoes: event.target.value })} placeholder="Informe somente condições que possam ser compartilhadas, como ferramentas existentes e disponibilidade da equipe." /></Field>
        <div className="category-field"><span className="field-caption">Categoria do problema <span className="required">*</span></span>
          <div className="category-options" role="radiogroup" aria-label="Categoria do problema">{categories.map(category => <label className={form.categoria === category.value ? 'category-option selected' : 'category-option'} key={category.value}>
            <input type="radio" name="categoria" value={category.value} checked={form.categoria === category.value} onChange={() => setForm({ ...form, categoria: category.value })} /><span className="category-check"><Icon name="check" size={12} /></span>{category.label}
          </label>)}</div>
          <small>Uma categoria por problema. Não inclua dados de clientes, senhas ou informações comerciais confidenciais.</small>
        </div>
        <button type="submit" className="primary full-width"><Icon name="search" size={19} />{busy ? 'Buscando soluções…' : 'Buscar soluções'}</button>
      </fieldset>
    </form>
  </div>
}
