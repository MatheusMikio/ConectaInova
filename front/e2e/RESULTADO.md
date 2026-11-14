# Resultado dos testes reais

URL: `http://localhost:5173/`

Execução final após correções e seleção direta no Matching: `2026-11-14T17-48-48-814Z`

**71 cenários: 71 aprovados, 0 reprovados, nenhum pulado.**

Foram usados navegador, frontend, API e PostgreSQL reais, sem interceptação ou simulação de respostas. Os testes criam e editam apenas registros próprios com prefixo `Teste E2E`. Não foram executadas exclusões nem resets.

## Defeitos corrigidos e validados

| Casos originais | Defeito | Resultado após correção |
|---|---|---|
| 2 | Briefing aceita somente espaços em campos obrigatórios | O formulário bloqueia a submissão antes de qualquer POST/PATCH. A API também rejeita o apagamento explícito do resultado esperado e dos critérios, sem alterar dados persistidos. |
| 1 | PATCH com JSON `null` é aceito | A API retorna HTTP 400, preservando os dados e `updated_at`. A proteção foi aplicada a todos os handlers de criação e edição. |
| 3 | Overflow na paginação | A verificação ocorre antes da multiplicação. Páginas extremas retornam listas vazias em startups, soluções e desafios; necessidades e avaliações também passaram. |

Os 64 cenários originais foram preservados e cinco regressões de API foram acrescentadas para resultado esperado em branco/vazio, critérios vazios, critério somente com espaços e apagamento simultâneo do briefing. Todos passaram; a suíte retornou código de saída 0. A compatibilidade com briefings legados omitidos foi mantida.

Também passaram `go test ./...`, `go vet ./...`, `bun run build` e `bun run lint`.

## Matching acessível pelo menu

A seção agora lista problemas reais quando não há seleção e permite analisá-los sem sair do Matching. Os resultados oferecem “Trocar problema” e a seleção permite voltar aos resultados atuais. A lista reutiliza paginação, atualização e tratamento de falha de carregamento existentes.

Dois novos cenários reais validam entrada direta após recarregar, seleção, troca para um problema sem soluções sem reaproveitar resultados anteriores e navegação mobile sem overflow. Build e lint passaram novamente.

## Cenários aprovados

- Cadastros, edições e consultas pelo navegador, conferindo a persistência pela API.
- Seleção de startup pelo nome e vínculo automático entre problema e desafio.
- Matching com palavras-chave, pontuação, evidências e ordenação.
- Matching somente por categoria, sem inventar palavras em comum.
- Caminho sem solução, publicação de desafio e persistência de critérios e prazo.
- Sincronização da aba do desafio publicado depois da edição.
- Briefing incompleto bloqueando uma solicitação de avaliação.
- Encaminhamento da demanda, parecer e proposta de piloto.
- Bloqueio de piloto sem parecer pela interface.
- Oito solicitações simultâneas para o mesmo par: **um HTTP 201, sete HTTP 200 e um único registro persistido**.
- Filtros de avaliação sem misturar startups.
- Rejeição de solução incompatível e referências inexistentes.
- Rejeição de JSON malformado, IDs inválidos e recursos inexistentes.
- Limites de textos Unicode, palavras-chave, critérios, prazos e status.
- Verificação de que PATCHes rejeitados não alteram os dados nem `updated_at`.
- Conteúdo HTML armazenado exibido como texto, sem execução de script.
- Falha de rede real usando modo offline do navegador, seguida de recuperação.
- Layout mobile, cancelamento de modais e persistência após recarregar.
- Nenhum erro JavaScript durante a execução.

## Evidências locais

Ignoradas pelo Git:

- `front/test-results/real-2026-11-14T17-48-48-814Z/report.json`
- `front/test-results/real-2026-11-14T17-48-48-814Z/trace.zip`
- Screenshots dos fluxos na mesma pasta.

A execução anterior, `real-2026-11-14T17-29-05-603Z`, foi preservada com relatório, trace e screenshots: 64 cenários, 58 aprovados e seis reprovados, correspondentes aos três defeitos acima.

A rodada `real-2026-11-14T17-37-13-660Z` preserva os 69 cenários aprovados após essas correções. A rodada `real-2026-11-14T17-48-14-484Z` teve 70 aprovados e uma falha no novo seletor textual do teste, corrigido para conferir o título dentro do resumo do problema; não houve alteração de comportamento da aplicação entre essa rodada e a final.

A primeira rodada também foi preservada. Nela houve uma ambiguidade no seletor da automação, corrigida antes da rodada de 64 cenários; essa ambiguidade não foi contada como defeito da aplicação.

Consulte `README.md` nesta pasta para repetir a suíte.
