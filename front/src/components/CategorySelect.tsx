import type { Categoria } from '../types/domain'

export function CategorySelect({ value, onChange }: { value: Categoria; onChange: (value: Categoria) => void }) {
  return <select value={value} onChange={event => onChange(event.target.value as Categoria)}>
    <option value="estoque">Estoque</option>
    <option value="vendas">Vendas</option>
    <option value="atendimento">Atendimento</option>
    <option value="financeiro">Financeiro</option>
    <option value="marketing">Marketing</option>
    <option value="logistica">Logística</option>
  </select>
}
