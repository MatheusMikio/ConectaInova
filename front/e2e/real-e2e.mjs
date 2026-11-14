import { chromium } from 'playwright'
import assert from 'node:assert/strict'
import { mkdir, writeFile } from 'node:fs/promises'
import { fileURLToPath } from 'node:url'
import path from 'node:path'

// This suite writes only records it creates, against a real running API.
// No page.route(), response mocks, database resets or existing-data cleanup.
if (process.env.E2E_ALLOW_WRITES !== '1') {
  throw new Error('This suite creates persistent Teste E2E records. Set E2E_ALLOW_WRITES=1 to authorize the run.')
}
const base = process.env.E2E_BASE_URL || 'http://localhost:5173'
const runID = new Date().toISOString().replace(/[:.]/g, '-')
const prefix = `Teste E2E ${runID}`
const output = fileURLToPath(new URL(`../test-results/real-${runID}/`, import.meta.url))
await mkdir(output, { recursive: true })
const report = {
  baseURL: base, runID, prefix, simulatedResponses: false,
  startedAt: new Date().toISOString(), cases: [], records: {}, requests: [], runtimeErrors: [],
}
const browser = await chromium.launch({
  headless: true,
  ...(process.env.PLAYWRIGHT_EXECUTABLE_PATH
    ? { executablePath: process.env.PLAYWRIGHT_EXECUTABLE_PATH }
    : process.platform === 'win32' ? { channel: 'msedge' } : {}),
})
const context = await browser.newContext({ viewport: { width: 1440, height: 1080 } })
await context.tracing.start({ screenshots: true, snapshots: true, sources: false })
const page = await context.newPage()
page.setDefaultTimeout(12000)
page.on('pageerror', error => report.runtimeErrors.push(error.message))
page.on('response', response => {
  const url = new URL(response.url())
  if (url.pathname.startsWith('/api/v1/')) report.requests.push({
    method: response.request().method(), path: url.pathname + url.search, status: response.status(),
  })
})
const nav = () => page.getByRole('navigation', { name: 'Navegação principal', exact: true })
const saved = {}
const persisted = () => writeFile(path.join(output, 'report.json'), JSON.stringify(report, null, 2))

async function check(name, fn, dependencies = []) {
  const missing = dependencies.filter(key => !saved[key])
  if (missing.length) {
    report.cases.push({ name, status: 'skipped', reason: `Prerequisites failed: ${missing.join(', ')}` })
    console.log(`SKIP: ${name}`)
    return
  }
  const start = Date.now()
  try {
    const evidence = await fn()
    report.cases.push({ name, status: 'passed', durationMs: Date.now() - start, evidence })
    console.log(`PASS: ${name}`)
  } catch (error) {
    const item = { name, status: 'failed', durationMs: Date.now() - start, error: error.message }
    report.cases.push(item)
    console.error(`FAIL: ${name}\n${error.message}`)
    const screenshot = `failure-${report.cases.length}.png`
    await page.screenshot({ path: path.join(output, screenshot), fullPage: true, timeout: 5000 }).then(() => { item.screenshot = screenshot }).catch(() => {})
    await context.setOffline(false)
    await page.keyboard.press('Escape').catch(() => {})
  } finally { await persisted() }
}
async function api(method, route, data) {
  return context.request.fetch(`${base}/api/v1/${route}`, {
    method,
    ...(data === undefined ? {} : { data, headers: { 'Content-Type': 'application/json' } }),
  })
}
async function payload(response, expected = 200) {
  const body = await response.json()
  assert.equal(response.status(), expected, `${response.url()}: ${JSON.stringify(body)}`)
  assert.ok(Object.hasOwn(body, 'data'), 'Missing data envelope')
  return body.data
}
async function get(route) { return payload(await api('GET', route)) }
async function all(resource) {
  const rows = []
  for (let p = 1; p <= 1000; p++) {
    const batch = await get(`${resource}${resource.includes('?') ? '&' : '?'}page=${p}&size=30`)
    rows.push(...batch)
    if (batch.length < 30) return rows
  }
  throw new Error('Pagination did not terminate')
}
function waitAPI(method, route) {
  return page.waitForResponse(response => response.request().method() === method &&
    new URL(response.url()).pathname === `/api/v1/${route}`)
}
async function act(method, route, action, expected = 200) {
  const pending = waitAPI(method, route)
  await action()
  return payload(await pending, expected)
}
async function track(key, entity) {
  saved[key] = entity
  report.records[key] = { id: entity.id, nome: entity.nome, titulo: entity.titulo }
  await persisted()
  return entity
}
async function go(name) { await nav().getByRole('button', { name, exact: true }).click() }
async function entity(resource, id) {
  if (['lojistas', 'necessidades'].includes(resource)) return get(`${resource}/${id}`)
  const result = (await all(resource)).find(row => row.id === id)
  assert.ok(result, `Missing persisted ${resource}/${id}`)
  return result
}
async function rejectWithoutChanging(name, resource, id, data, expected = 400) {
  await check(name, async () => {
    const before = await entity(resource, id)
    const response = await api('PATCH', `${resource}/${id}`, data)
    const error = await response.json()
    const after = await entity(resource, id)
    assert.equal(response.status(), expected, `Expected rejection; received ${response.status()}: ${JSON.stringify(error)}`)
    assert.ok(error.error, 'Missing error envelope')
    assert.deepEqual(after, before, 'Rejected PATCH changed the persisted entity')
    return { expectedStatus: expected, actualStatus: response.status(), unchanged: true }
  })
}
async function fillProblem(name, category = 'estoque') {
  await page.getByRole('textbox', { name: /Nome da loja/ }).fill(`${prefix} Loja`)
  await page.getByRole('textbox', { name: /Título do problema/ }).fill(name)
  await page.getByRole('textbox', { name: /Descreva sua necessidade/ }).fill(`${prefix}: controlar estoque e reposição sem perder vendas.`)
  await page.getByRole('combobox', { name: /Como sua loja vende/ }).selectOption('hibrida')
  await page.getByRole('combobox', { name: /Impacto no negócio/ }).selectOption('Perda de vendas e clientes')
  await page.getByRole('combobox', { name: /Prioridade/ }).selectOption('alta')
  await page.getByRole('textbox', { name: /Resultado esperado/ }).fill(`${prefix}: reduzir rupturas durante a operação.`)
  await page.getByRole('textbox', { name: /Critérios de sucesso/ }).fill('Receber alertas de estoque\nSimplificar a reposição')
  await page.getByRole('textbox', { name: /Restrições e condições/ }).fill(`${prefix}: manter as planilhas atuais até validar a integração.`)
  await page.locator(`input[name="categoria"][value="${category}"]`).check()
}
async function problemCard(problem) {
  await page.getByRole('heading', { name: problem.titulo, exact: true }).waitFor()
  return page.locator('.problem-card').filter({ has: page.getByRole('heading', { name: problem.titulo, exact: true }) })
}

try {
  await page.goto(base, { waitUntil: 'domcontentloaded', timeout: 30000 })
  await check('UI: formulário vazio bloqueia gravação e não cria loja nem necessidade', async () => {
    const before = report.requests.filter(r => r.method === 'POST').length
    await page.getByRole('button', { name: 'Buscar soluções', exact: true }).click()
    assert.equal(await page.locator('.problem-form').evaluate(form => form.checkValidity()), false)
    assert.equal(report.requests.filter(r => r.method === 'POST').length, before)
  })

  await check('UI: startup criada no modal e confirmada no banco', async () => {
    await go('Startups')
    assert.equal(await page.getByRole('dialog').count(), 0)
    await page.getByRole('button', { name: 'Cadastrar startup', exact: true }).click()
    const dialog = page.getByRole('dialog', { name: 'Cadastrar startup', exact: true })
    await dialog.getByRole('textbox', { name: /Nome da startup/ }).fill(`${prefix} Startup`)
    await dialog.getByRole('textbox', { name: /Sobre a startup/ }).fill(`${prefix}: tecnologia para estoque e reposição.`)
    const result = await act('POST', 'startups', () => dialog.getByRole('button', { name: 'Cadastrar startup', exact: true }).click(), 201)
    await track('startup', result)
    await page.getByRole('heading', { name: result.nome, exact: true }).waitFor()
    assert.equal((await entity('startups', result.id)).descricao, result.descricao)
  })

  await check('UI: erro de nome em branco mantém modal e não altera startup', async () => {
    await page.getByRole('button', { name: `Editar ${saved.startup.nome}`, exact: true }).click()
    const dialog = page.getByRole('dialog', { name: 'Editar startup', exact: true })
    await dialog.getByRole('textbox', { name: /Nome da startup/ }).fill('   ')
    const pending = waitAPI('PATCH', `startups/${saved.startup.id}`)
    await dialog.getByRole('button', { name: 'Salvar alterações', exact: true }).click()
    assert.equal((await pending).status(), 400)
    await dialog.getByRole('alert').waitFor()
    assert.equal((await entity('startups', saved.startup.id)).nome, saved.startup.nome)
    await dialog.getByRole('textbox', { name: /Nome da startup/ }).fill(`${prefix} Startup Editada`)
    saved.startup = await act('PATCH', `startups/${saved.startup.id}`, () => dialog.getByRole('button', { name: 'Salvar alterações', exact: true }).click())
    await page.getByRole('heading', { name: saved.startup.nome, exact: true }).waitFor()
  }, ['startup'])

  await check('UI: cancelar e Escape não gravam cadastros', async () => {
    const before = (await all('startups')).length
    await page.getByRole('button', { name: 'Cadastrar startup', exact: true }).click()
    const dialog = page.getByRole('dialog', { name: 'Cadastrar startup', exact: true })
    await dialog.getByRole('textbox', { name: /Nome da startup/ }).fill(`${prefix} CANCELADA`)
    await dialog.getByRole('button', { name: 'Cancelar', exact: true }).click()
    await dialog.waitFor({ state: 'hidden' })
    await page.getByRole('button', { name: 'Cadastrar startup', exact: true }).click()
    await page.getByRole('dialog').waitFor({ state: 'visible' })
    await page.keyboard.press('Escape')
    await page.getByRole('dialog').waitFor({ state: 'hidden' })
    assert.equal((await all('startups')).length, before)
  })

  await check('UI: solução criada por seleção de startup pelo nome e persistida', async () => {
    await go('Soluções')
    await page.getByRole('button', { name: 'Cadastrar solução', exact: true }).click()
    const dialog = page.getByRole('dialog', { name: 'Cadastrar solução', exact: true })
    const select = dialog.getByRole('combobox', { name: /Startup responsável/ })
    await select.locator('option').filter({ hasText: saved.startup.nome }).waitFor({ state: 'attached' })
    assert.equal(await dialog.locator('input[type="number"]').count(), 0)
    await select.selectOption({ label: saved.startup.nome })
    await dialog.getByRole('textbox', { name: /Nome da solução/ }).fill(`${prefix} Estoque`)
    await dialog.getByRole('textbox', { name: /Descrição/ }).fill(`${prefix}: controle de estoque e reposição.`)
    await dialog.getByRole('textbox', { name: /Palavras-chave/ }).fill('estoque, reposição')
    const result = await act('POST', 'solucoes', () => dialog.getByRole('button', { name: 'Cadastrar solução', exact: true }).click(), 201)
    await track('solution', result)
    assert.equal(result.startup_id, saved.startup.id)
    await page.getByRole('heading', { name: result.nome, exact: true }).waitFor()
  }, ['startup'])

  await check('Dados de teste reais para vínculos incompatíveis e concorrência', async () => {
    const other = await payload(await api('POST', 'startups', { nome: `${prefix} Outra Startup`, descricao: `${prefix}: empresa para testar isolamento de filtros.` }), 201)
    await track('otherStartup', other)
    await track('secondSolution', await payload(await api('POST', 'solucoes', {
      startup_id: saved.startup.id, nome: `${prefix} Estoque Categoria`, descricao: `${prefix}: solução somente por categoria.`,
      categoria: 'estoque', palavras_chave: [],
    }), 201))
    await track('incompatibleSolution', await payload(await api('POST', 'solucoes', {
      startup_id: other.id, nome: `${prefix} Financeiro`, descricao: `${prefix}: conciliação financeira.`,
      categoria: 'financeiro', palavras_chave: ['caixa'],
    }), 201))
  }, ['startup'])

  await check('UI: lojista + necessidade qualificada + matching real de 100 pontos', async () => {
    await go('Cadastrar problema')
    await fillProblem(`${prefix} Problema Estoque`)
    const storePending = waitAPI('POST', 'lojistas')
    const needPending = waitAPI('POST', 'necessidades')
    await page.getByRole('button', { name: 'Buscar soluções', exact: true }).click()
    await track('shop', await payload(await storePending, 201))
    const need = await track('need', await payload(await needPending, 201))
    assert.equal(need.lojista_id, saved.shop.id)
    assert.equal(need.criterios_sucesso.length, 2)
    assert.equal(need.canal_venda, 'hibrida')
    assert.ok(need.restricoes.includes('planilhas'))
    await page.getByRole('heading', { name: 'Soluções encontradas', exact: true }).waitFor()
    const match = (await get(`necessidades/${need.id}/matches`)).matches.find(item => item.solucao.id === saved.solution.id)
    assert.equal(match.pontuacao, 100)
    assert.equal(match.startup.id, saved.startup.id)
    assert.equal(match.palavras_encontradas.length, 2)
    assert.ok(match.motivos[1].includes('estoque'))
  }, ['solution'])

  await check('Matching: categoria sem palavras em comum é explicada, não inventa evidências', async () => {
    const result = await get(`necessidades/${saved.need.id}/matches`)
    const match = result.matches.find(item => item.solucao.id === saved.secondSolution.id)
    assert.equal(match.pontuacao, 70)
    assert.deepEqual(match.palavras_encontradas, [])
    assert.ok(match.motivos.some(reason => reason.includes('apenas a categoria')))
    assert.ok(result.matches.every((item, index) => index === 0 || result.matches[index - 1].pontuacao >= item.pontuacao))
  }, ['need', 'secondSolution'])

  await check('UI: encaminha briefing à startup e avaliação persiste com associação correta', async () => {
    const card = page.locator('.solution-card').filter({ has: page.getByRole('heading', { name: saved.solution.nome, exact: true }) })
    await card.getByRole('button', { name: 'Avaliar solução', exact: true }).click()
    const dialog = page.getByRole('dialog', { name: saved.solution.nome, exact: true })
    await dialog.getByRole('heading', { name: 'O que ainda precisa ser confirmado', exact: true }).waitFor()
    await dialog.getByRole('textbox', { name: /Mensagem para a avaliação/ }).fill(`${prefix}: confirmar implantação, suporte e integração.`)
    const pending = waitAPI('POST', 'avaliacoes')
    await dialog.getByRole('button', { name: 'Solicitar avaliação', exact: true }).click()
    const evaluation = await track('evaluation', await payload(await pending, 201))
    assert.equal(evaluation.necessidade.id, saved.need.id)
    assert.equal(evaluation.necessidade.criterios_sucesso.length, 2)
    assert.equal(evaluation.necessidade.lojista_nome, saved.shop.nome)
    assert.equal(evaluation.solucao.id, saved.solution.id)
    assert.equal(evaluation.associacao, (await get('plataforma')).associacao)
    await dialog.getByRole('button', { name: 'Acompanhar avaliação', exact: true }).click()
    await page.getByRole('heading', { name: 'Avaliações', exact: true, level: 1 }).waitFor()
  }, ['need', 'solution'])

  await check('UI: piloto sem parecer é bloqueado; parecer válido permite avançar', async () => {
    const card = page.locator('.directory-card').filter({ has: page.getByRole('heading', { name: saved.need.titulo, exact: true }) }).first()
    await card.getByRole('button', { name: 'Analisar demanda', exact: true }).click()
    const dialog = page.getByRole('dialog', { name: 'Analisar demanda', exact: true })
    await dialog.getByRole('combobox', { name: 'Etapa da avaliação', exact: true }).selectOption('piloto_proposto')
    await dialog.getByRole('button', { name: 'Salvar avaliação', exact: true }).click()
    assert.equal(await dialog.locator('form').evaluate(form => form.checkValidity()), false)
    assert.equal((await entity('avaliacoes', saved.evaluation.id)).status, 'solicitada')
    await dialog.getByRole('textbox', { name: 'Parecer e próximo passo', exact: true }).fill(`${prefix}: confirmar escopo e condições com a loja antes de propor o teste.`)
    saved.evaluation = await act('PATCH', `avaliacoes/${saved.evaluation.id}`, () => dialog.getByRole('button', { name: 'Salvar avaliação', exact: true }).click())
    assert.equal(saved.evaluation.status, 'piloto_proposto')
    assert.ok(saved.evaluation.resposta.includes('confirmar escopo'))
  }, ['evaluation'])

  await check('Concorrência real: 8 solicitações simultâneas geram um único registro', async () => {
    const data = { necessidade_id: saved.need.id, solucao_id: saved.secondSolution.id, mensagem: `${prefix}: concorrência real.` }
    const responses = await Promise.all(Array.from({ length: 8 }, () => api('POST', 'avaliacoes', data)))
    const statuses = responses.map(response => response.status())
    const results = await Promise.all(responses.map(async response => (await response.json()).data))
    assert.equal(statuses.filter(status => status === 201).length, 1)
    assert.equal(statuses.filter(status => status === 200).length, 7)
    assert.equal(new Set(results.map(item => item.id)).size, 1)
    await track('concurrentEvaluation', results[0])
    const matches = (await all('avaliacoes')).filter(row => row.necessidade.id === saved.need.id && row.solucao.id === saved.secondSolution.id)
    assert.equal(matches.length, 1)
    return { statuses, persistedCount: matches.length }
  }, ['need', 'secondSolution'])

  await check('Filtros por startup não misturam demandas de outra startup', async () => {
    const own = await all(`avaliacoes?startup_id=${saved.startup.id}`)
    assert.ok(own.length >= 2)
    assert.ok(own.every(row => row.startup.id === saved.startup.id))
    const other = await all(`avaliacoes?startup_id=${saved.otherStartup.id}`)
    assert.deepEqual(other, [])
  }, ['concurrentEvaluation', 'otherStartup'])

  await check('Avaliação de solução de categoria incompatível é rejeitada e não grava', async () => {
    const before = (await all('avaliacoes')).filter(row => row.necessidade.id === saved.need.id).length
    const response = await api('POST', 'avaliacoes', { necessidade_id: saved.need.id, solucao_id: saved.incompatibleSolution.id })
    assert.equal(response.status(), 400)
    assert.equal((await all('avaliacoes')).filter(row => row.necessidade.id === saved.need.id).length, before)
  }, ['need', 'incompatibleSolution'])

  await check('Briefing incompleto não permite solicitar avaliação; UI oferece completar dados', async () => {
    const legacy = await payload(await api('POST', 'necessidades', {
      lojista_id: saved.shop.id, titulo: `${prefix} Briefing incompleto`, descricao: `${prefix}: estoque.`,
      categoria: 'estoque', prioridade: 'media',
    }), 201)
    await track('legacyNeed', legacy)
    const rejected = await api('POST', 'avaliacoes', { necessidade_id: legacy.id, solucao_id: saved.solution.id })
    assert.equal(rejected.status(), 400)
    await go('Problemas cadastrados')
    const card = await problemCard(legacy)
    await card.getByRole('button', { name: 'Ver soluções', exact: true }).click()
    const solutionCard = page.locator('.solution-card').filter({ has: page.getByRole('heading', { name: saved.solution.nome, exact: true }) })
    await solutionCard.getByRole('button', { name: 'Avaliar solução', exact: true }).click()
    const dialog = page.getByRole('dialog')
    assert.equal(await dialog.getByRole('button', { name: 'Solicitar avaliação', exact: true }).count(), 0)
    await dialog.getByRole('button', { name: 'Completar briefing do problema', exact: true }).click()
    await page.getByRole('heading', { name: 'Editar problema', exact: true }).waitFor()
  }, ['need', 'solution', 'shop'])

  await check('UI: demanda sem solução oferece desafio sem IDs e sem duplicar loja', async () => {
    const available = await all('solucoes')
    const category = ['atendimento', 'marketing', 'logistica', 'vendas'].find(value => !available.some(item => item.categoria === value))
    assert.ok(category, 'All categories populated: cannot exercise real empty matching without modifying other data')
    await page.getByRole('button', { name: 'Cadastrar outro problema', exact: true }).click()
    await fillProblem(`${prefix} Problema sem solução`, category)
    saved.emptyNeed = await act('POST', 'necessidades', () => page.getByRole('button', { name: 'Buscar soluções', exact: true }).click(), 201)
    await track('emptyNeed', saved.emptyNeed)
    await page.getByRole('heading', { name: 'Nenhuma solução compatível por enquanto.', exact: true }).waitFor()
    assert.equal(saved.emptyNeed.lojista_id, saved.shop.id)
    const result = await get(`necessidades/${saved.emptyNeed.id}/matches`)
    assert.deepEqual(result.matches, [])
    assert.equal(result.sugerir_desafio, true)
    await page.getByRole('button', { name: 'Gerar desafio de inovação', exact: true }).click()
    assert.equal(await page.locator('main input[type="number"]').count(), 0)
    await page.getByRole('textbox', { name: 'Título do desafio', exact: true }).fill(`${prefix} Desafio`)
    await page.getByRole('combobox', { name: 'Prazo para propostas', exact: true }).selectOption('14')
    const challenge = await act('POST', 'desafios', () => page.getByRole('button', { name: 'Publicar desafio', exact: true }).click(), 201)
    await track('challenge', challenge)
    assert.equal(challenge.necessidade_id, saved.emptyNeed.id)
    assert.equal(challenge.criterios_sucesso.length, 2)
    assert.equal(challenge.prazo_dias, 14)
    await page.getByRole('heading', { name: 'Desafio de inovação publicado', exact: true }).waitFor()
  }, ['legacyNeed', 'shop'])

  await check('UI: edição de desafio mantém vínculo, critérios, prazo e aba publicada sincronizados', async () => {
    await page.getByRole('button', { name: 'Consultar e gerenciar', exact: true }).click()
    await page.getByRole('heading', { name: saved.challenge.titulo, exact: true }).waitFor()
    await page.getByRole('button', { name: `Editar ${saved.challenge.titulo}`, exact: true }).click()
    const dialog = page.getByRole('dialog', { name: 'Editar desafio', exact: true })
    await dialog.getByRole('region', { name: 'Problema vinculado', exact: true }).getByText(saved.emptyNeed.titulo, { exact: true }).waitFor()
    assert.equal(await dialog.getByLabel(/ID/).count(), 0)
    await dialog.getByRole('textbox', { name: /Título do desafio/ }).fill(`${prefix} Desafio Editado`)
    await dialog.getByRole('spinbutton', { name: /Prazo para propostas/ }).fill('21')
    saved.challenge = await act('PATCH', `desafios/${saved.challenge.id}`, () => dialog.getByRole('button', { name: 'Salvar alterações', exact: true }).click())
    assert.equal(saved.challenge.necessidade_id, saved.emptyNeed.id)
    assert.equal(saved.challenge.prazo_dias, 21)
    await page.getByRole('navigation', { name: 'Desafios de inovação', exact: true }).getByRole('button', { name: 'Desafio publicado', exact: true }).click()
    assert.equal(await page.getByRole('textbox', { name: 'Título do desafio', exact: true }).inputValue(), saved.challenge.titulo)
    assert.equal(await page.getByRole('combobox', { name: 'Prazo para propostas', exact: true }).inputValue(), '21')
    saved.challengeEdited = true
  }, ['challenge'])

  for (const field of ['resultado_esperado', 'criterios_sucesso']) {
    await check(`UI: ${field} somente com espaços não deve apagar o briefing`, async () => {
      const before = await get(`necessidades/${saved.need.id}`)
      try {
        await go('Problemas cadastrados')
        const card = await problemCard(saved.need)
        await card.getByRole('button', { name: `Editar ${saved.need.titulo}`, exact: true }).click()
        await page.getByRole('heading', { name: 'Editar problema', exact: true }).waitFor()
        await page.getByRole('textbox', { name: field === 'resultado_esperado' ? /Resultado esperado/ : /Critérios de sucesso/ }).fill('   ')
        const writesBeforeSubmit = report.requests.filter(item => item.method === 'PATCH' || item.method === 'POST').length
        await page.getByRole('button', { name: 'Buscar soluções', exact: true }).click()
        await page.waitForFunction(() =>
          !!document.querySelector('.notice.error') || Array.from(document.querySelectorAll('main h1')).some(el => el.textContent.includes('Soluções encontradas')))
        const after = await get(`necessidades/${saved.need.id}`)
        assert.deepEqual(after[field], before[field], 'Formulário aceitou espaços como um campo obrigatório válido e apagou dados persistidos do briefing')
        assert.equal(report.requests.filter(item => item.method === 'PATCH' || item.method === 'POST').length, writesBeforeSubmit, 'Frontend sent a write request for invalid mandatory fields')
      } finally {
        await payload(await api('PATCH', `necessidades/${saved.need.id}`, {
          resultado_esperado: before.resultado_esperado, criterios_sucesso: before.criterios_sucesso,
        }))
      }
    }, ['need'])
  }

  // Rejected updates must not change even updated_at on any test record.
  if (saved.startup) {
    await rejectWithoutChanging('API: nome em branco não altera startup', 'startups', saved.startup.id, { nome: '   ' })
    await rejectWithoutChanging('API: nome com 151 caracteres Unicode é rejeitado', 'startups', saved.startup.id, { nome: 'á'.repeat(151) })
  }
  await check('API: limite de 150 caracteres Unicode é aceito sem confundir caracteres e bytes', async () => {
    const before = await entity('startups', saved.startup.id)
    try {
      const result = await payload(await api('PATCH', `startups/${saved.startup.id}`, { nome: 'á'.repeat(150) }))
      assert.equal(result.nome, 'á'.repeat(150))
    } finally { await payload(await api('PATCH', `startups/${saved.startup.id}`, { nome: before.nome })) }
  }, ['startup'])
  if (saved.solution) {
    for (const [name, data] of [
      ['categoria inexistente', { categoria: 'invalida' }],
      ['startup inexistente', { startup_id: 999999999 }],
      ['31 palavras-chave', { palavras_chave: Array(31).fill('estoque') }],
      ['palavra-chave vazia', { palavras_chave: ['   '] }],
      ['tipo incorreto de startup', { startup_id: 'abc' }],
    ]) await rejectWithoutChanging(`API: solução rejeita ${name} sem alteração`, 'solucoes', saved.solution.id, data)
  }
  if (saved.need) {
    for (const [name, data] of [
      ['prioridade inválida', { prioridade: 'urgente' }],
      ['canal de venda inválido', { canal_venda: 'inexistente' }],
      ['21 critérios', { criterios_sucesso: Array(21).fill('Critério') }],
      ['critério de 301 caracteres', { criterios_sucesso: ['a'.repeat(301)] }],
      ['impacto com 2001 caracteres', { impacto: 'a'.repeat(2001) }],
      ['lojista zero', { lojista_id: 0 }],
      ['resultado esperado em branco', { resultado_esperado: '   ' }],
      ['resultado esperado vazio', { resultado_esperado: '' }],
      ['critérios de sucesso vazios', { criterios_sucesso: [] }],
      ['critério de sucesso só com espaços', { criterios_sucesso: ['   '] }],
      ['apagamento completo do briefing', { resultado_esperado: '', criterios_sucesso: [] }],
    ]) await rejectWithoutChanging(`API: necessidade rejeita ${name} sem alteração`, 'necessidades', saved.need.id, data)
  }
  if (saved.challenge) {
    for (const [name, data] of [
      ['prazo zero', { prazo_dias: 0 }],
      ['prazo de 366 dias', { prazo_dias: 366 }],
      ['status inválido', { status: 'inventado' }],
      ['título em branco', { titulo: '  ' }],
    ]) await rejectWithoutChanging(`API: desafio rejeita ${name} sem alteração`, 'desafios', saved.challenge.id, data)
  }
  if (saved.evaluation) {
    for (const [name, data] of [
      ['status inválido', { status: 'inventado' }],
      ['resposta de 4001 caracteres', { resposta: 'a'.repeat(4001) }],
      ['conclusão sem parecer', { status: 'concluida', resposta: '' }],
    ]) await rejectWithoutChanging(`API: avaliação rejeita ${name} sem alteração`, 'avaliacoes', saved.evaluation.id, data)
  }

  await check('API: JSON malformado é rejeitado sem gravar', async () => {
    const response = await api('POST', 'startups', '{invalid')
    assert.equal(response.status(), 400)
    assert.ok((await response.json()).error)
  })
  await check('API: PATCH null deve ser rejeitado sem alterar registro', async () => {
    const before = await entity('startups', saved.startup.id)
    const response = await api('PATCH', `startups/${saved.startup.id}`, 'null')
    const after = await entity('startups', saved.startup.id)
    assert.equal(response.status(), 400, `Received ${response.status()}; persisted updated_at: ${before.updated_at} -> ${after.updated_at}`)
    assert.deepEqual(after, before)
  }, ['startup'])

  for (const route of ['necessidades/0', 'necessidades/-1/matches', 'lojistas/abc', 'necessidades/1.5/desafios']) {
    await check(`API: ID inválido em ${route} retorna 400`, async () => {
      const response = await api('GET', route)
      assert.equal(response.status(), 400)
      assert.ok((await response.json()).error)
    })
  }
  for (const route of ['necessidades/999999999', 'necessidades/999999999/matches', 'necessidades/999999999/desafios']) {
    await check(`API: recurso ausente ${route} retorna 404`, async () => assert.equal((await api('GET', route)).status(), 404))
  }
  await check('API: avaliação com referências inexistentes é rejeitada sem criar', async () => {
    const before = (await all('avaliacoes')).length
    assert.equal((await api('POST', 'avaliacoes', { necessidade_id: 999999999, solucao_id: 999999999 })).status(), 404)
    assert.equal((await all('avaliacoes')).length, before)
  })

  await check('API: paginação normal separa páginas e não repete registros', async () => {
    const p1 = await get('solucoes?page=1&size=1')
    const p2 = await get('solucoes?page=2&size=1')
    assert.equal(p1.length, 1); assert.equal(p2.length, 1)
    assert.notEqual(p1[0].id, p2[0].id)
    assert.deepEqual(await get('solucoes?page=99999&size=12'), [])
  })
  for (const resource of ['startups', 'solucoes', 'desafios', 'necessidades', 'avaliacoes']) {
    await check(`API: paginação extrema em ${resource} não deve retornar à primeira página`, async () => {
      const result = await get(`${resource}?page=9223372036854775809&size=2`)
      assert.deepEqual(result, [], `Overflow leaked first-page IDs: ${result.map(row => row.id).join(', ')}`)
    })
  }
  for (const query of ['page=-1', 'page=abc', 'size=-1', 'startup_id=abc']) {
    await check(`API: query inválida ${query} é rejeitada`, async () => {
      const response = await api('GET', `avaliacoes?${query}`)
      assert.equal(response.status(), 400)
    })
  }

  await check('Segurança de renderização: HTML armazenado é exibido como texto, sem executar script', async () => {
    const before = await entity('startups', saved.otherStartup.id)
    const content = `${prefix} <img src=x onerror="window.__e2eInjected=true"><script>window.__e2eInjected=true</script>`
    await payload(await api('PATCH', `startups/${saved.otherStartup.id}`, { descricao: content }))
    try {
      await go('Startups')
      const card = page.locator('.directory-card').filter({ has: page.getByRole('heading', { name: before.nome, exact: true }) })
      await card.getByText(content, { exact: true }).waitFor()
      assert.equal(await card.locator('img,script').count(), 0)
      assert.equal(await page.evaluate(() => window.__e2eInjected), undefined)
    } finally { await payload(await api('PATCH', `startups/${saved.otherStartup.id}`, { descricao: before.descricao })) }
  }, ['otherStartup'])

  await check('Rede real offline mostra erro e recupera a lista ao voltar online', async () => {
    await context.setOffline(true)
    await go('Soluções')
    await page.getByRole('heading', { name: 'Não foi possível carregar as soluções.', exact: true }).waitFor()
    await context.setOffline(false)
    await page.getByRole('button', { name: 'Tentar novamente', exact: true }).click()
    await page.getByRole('heading', { name: saved.solution.nome, exact: true }).waitFor()
    assert.equal(await page.getByRole('alert').count(), 0)
  }, ['solution'])

  await check('Mobile: navegação, listas e modal reais sem overflow horizontal', async () => {
    await page.setViewportSize({ width: 390, height: 844 })
    for (const label of ['Início', 'Problemas cadastrados', 'Soluções', 'Startups', 'Desafios de inovação', 'Avaliações']) {
      await go(label)
      await page.getByRole('heading', { level: 1, name: label === 'Início' ? 'Inovação que começa no seu negócio' : label === 'Soluções' ? 'Catálogo de soluções' : label, exact: true }).waitFor()
      assert.equal(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth), true, `Horizontal overflow: ${label}`)
    }
    await go('Startups')
    await page.getByRole('button', { name: 'Cadastrar startup', exact: true }).click()
    await page.getByRole('dialog').waitFor({ state: 'visible' })
    assert.equal(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth), true)
    await page.getByRole('button', { name: 'Cancelar', exact: true }).click()
    await page.screenshot({ path: path.join(output, 'mobile.png'), fullPage: true })
  })

  await check('Recarregamento completo confirma dados reais e não depende do estado da sessão', async () => {
    await page.setViewportSize({ width: 1440, height: 1080 })
    await page.reload({ waitUntil: 'domcontentloaded' })
    await go('Problemas cadastrados')
    await page.getByRole('heading', { name: saved.need.titulo, exact: true }).waitFor()
    const card = await problemCard(saved.need)
    await card.getByRole('button', { name: 'Ver soluções', exact: true }).click()
    await page.getByRole('heading', { name: saved.solution.nome, exact: true }).waitFor()
    const need = await entity('necessidades', saved.need.id)
    assert.equal(need.criterios_sucesso.length, 2)
    assert.equal(need.lojista_id, saved.shop.id)
    assert.equal((await entity('avaliacoes', saved.evaluation.id)).status, 'piloto_proposto')
    assert.equal((await entity('desafios', saved.challenge.id)).prazo_dias, saved.challengeEdited ? 21 : saved.challenge.prazo_dias)
    assert.equal(await page.evaluate(() => localStorage.length), 0)
    await page.screenshot({ path: path.join(output, 'matching-after-reload.png'), fullPage: true })
  }, ['need', 'shop', 'solution', 'evaluation', 'challenge'])

  await check('Matching: entrada direta após recarregar permite selecionar um problema real', async () => {
    await page.reload({ waitUntil: 'domcontentloaded' })
    await go('Matching (soluções)')
    await page.getByRole('heading', { name: 'Matching de soluções', exact: true }).waitFor()
    const card = await problemCard(saved.need)
    await card.getByRole('button', { name: 'Analisar este problema', exact: true }).click()
    await page.getByRole('heading', { name: saved.solution.nome, exact: true }).waitFor()
    await page.screenshot({ path: path.join(output, 'matching-direct-selection.png'), fullPage: true })
  }, ['need', 'solution'])

  await check('Matching: trocar problema atualiza resultados sem reutilizar soluções anteriores', async () => {
    await page.getByRole('button', { name: 'Trocar problema', exact: true }).click()
    const card = await problemCard(saved.emptyNeed)
    await card.getByRole('button', { name: 'Analisar este problema', exact: true }).click()
    await page.getByRole('heading', { name: 'Nenhuma solução compatível por enquanto.', exact: true }).waitFor()
    assert.equal(await page.locator('.solution-card').count(), 0)
    assert.ok((await page.locator('.problem-summary').innerText()).includes(saved.emptyNeed.titulo))
    await page.setViewportSize({ width: 390, height: 844 })
    await page.getByRole('button', { name: 'Trocar problema', exact: true }).click()
    await page.getByRole('heading', { name: saved.need.titulo, exact: true }).waitFor()
    assert.equal(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth), true)
    await page.getByRole('button', { name: 'Voltar aos resultados', exact: true }).click()
    await page.getByRole('heading', { name: 'Nenhuma solução compatível por enquanto.', exact: true }).waitFor()
  }, ['emptyNeed', 'need'])

  await check('JavaScript: nenhum erro de execução durante os cenários', async () => assert.deepEqual(report.runtimeErrors, []))
} finally {
  await context.setOffline(false)
  await context.tracing.stop({ path: path.join(output, 'trace.zip') }).catch(() => {})
  report.finishedAt = new Date().toISOString()
  report.summary = {
    passed: report.cases.filter(item => item.status === 'passed').length,
    failed: report.cases.filter(item => item.status === 'failed').length,
    skipped: report.cases.filter(item => item.status === 'skipped').length,
  }
  await persisted()
  await browser.close()
  console.log(JSON.stringify({ ...report.summary, report: path.join(output, 'report.json'), records: report.records }, null, 2))
  if (report.summary.failed || report.summary.skipped) process.exitCode = 1
}
