import type { NecessidadeResponse } from '../types/domain'

const channels: Record<string, string> = { fisica: 'Loja física', digital: 'Loja digital', hibrida: 'Operação híbrida' }

export function ProblemBrief({ problem, loja }: { problem: NecessidadeResponse; loja?: string }) {
  return <dl className="problem-brief">
    {loja && <div><dt>Loja</dt><dd>{loja}</dd></div>}
    <div><dt>Problema</dt><dd>{problem.titulo}</dd></div>
    <div><dt>Situação atual</dt><dd>{problem.descricao}</dd></div>
    <div><dt>Perfil da operação</dt><dd>{channels[problem.canal_venda] || 'Não informado'}</dd></div>
    <div><dt>Impacto no negócio</dt><dd>{problem.impacto || 'Não informado'}</dd></div>
    <div><dt>Resultado esperado</dt><dd>{problem.resultado_esperado || 'Não informado'}</dd></div>
    <div><dt>Restrições</dt><dd>{problem.restricoes || 'Não informadas'}</dd></div>
    <div><dt>Critérios de sucesso</dt><dd>{problem.criterios_sucesso?.length ? <ul>{problem.criterios_sucesso.map((criterion, index) => <li key={index}>{criterion}</li>)}</ul> : 'Não informados'}</dd></div>
    <div><dt>Prioridade</dt><dd>{({ baixa: 'Baixa', media: 'Média', alta: 'Alta' })[problem.prioridade]}</dd></div>
  </dl>
}
