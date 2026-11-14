import { useState } from 'react'
import type { FormEvent } from 'react'
import { emptyProblem } from '../data/problem'
import { createLojista, getLojista, updateLojista } from '../services/lojistaService'
import { createNecessidade, getMatches, getNecessidade, updateNecessidade } from '../services/necessidadeService'
import { createDesafio } from '../services/desafioService'
import { errorMessage } from '../services/http'
import type { DesafioResponse, LojistaResponse, MatchingResponse, NecessidadeResponse } from '../types/domain'
import type { ChallengeDraft, ProblemForm, Screen } from '../types/workflow'

export function useInnovationFlow() {
  const [screen, setScreen] = useState<Screen>('cadastro')
  const [form, setForm] = useState<ProblemForm>({ ...emptyProblem })
  const [shop, setShop] = useState<LojistaResponse | null>(null)
  const [need, setNeed] = useState<NecessidadeResponse | null>(null)
  const [matches, setMatches] = useState<MatchingResponse | null>(null)
  const [challengeProblem, setChallengeProblem] = useState<NecessidadeResponse | null>(null)
  const [draft, setDraft] = useState<ChallengeDraft | null>(null)
  const [published, setPublished] = useState<DesafioResponse | null>(null)
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState('')
  const [message, setMessage] = useState('')

  function navigate(next: Screen) {
    if (busy) return
    if (next === 'desafio') setChallengeProblem(null)
    setError(''); setMessage(''); setScreen(next)
  }

  function newProblem() {
    setNeed(null); setMatches(null); setDraft(null); setPublished(null); setChallengeProblem(null)
    setForm(current => ({ ...emptyProblem, nomeLoja: current.nomeLoja, segmento: current.segmento }))
    navigate('cadastro')
  }

  async function submitProblem(event: FormEvent) {
    event.preventDefault()
    setBusy(true); setError(''); setMessage('')
    try {
      if (!form.nomeLoja.trim() || !form.titulo.trim() || !form.descricao.trim()) {
        throw new Error('Preencha o nome da loja, o título e a descrição do problema.')
      }
      const resultadoEsperado = form.resultadoEsperado.trim()
      const criterios = form.criteriosSucesso.split('\n').map(value => value.trim()).filter(Boolean)
      if (!resultadoEsperado) throw new Error('Preencha o resultado esperado com um texto válido.')
      if (!criterios.length) throw new Error('Defina pelo menos um critério de sucesso com um texto válido.')
      if (!form.canalVenda || !form.impacto.trim()) throw new Error('Preencha o perfil da operação e o impacto no negócio.')
      if (criterios.length > 20 || criterios.some(value => Array.from(value).length > 300)) {
        throw new Error('Informe no máximo 20 critérios de sucesso, com até 300 caracteres cada.')
      }
      let currentShop = shop
      if (!currentShop) {
        currentShop = await createLojista({ nome: form.nomeLoja, segmento: form.segmento })
      } else if (currentShop.nome !== form.nomeLoja || currentShop.segmento !== form.segmento) {
        currentShop = await updateLojista(currentShop.id, { nome: form.nomeLoja, segmento: form.segmento })
      }
      setShop(currentShop)
      const request = {
        lojista_id: currentShop.id, titulo: form.titulo, descricao: form.descricao, categoria: form.categoria, prioridade: form.prioridade,
        impacto: form.impacto, resultado_esperado: resultadoEsperado, restricoes: form.restricoes,
        canal_venda: form.canalVenda, criterios_sucesso: criterios,
      }
      const saved = need ? await updateNecessidade(need.id, request) : await createNecessidade(request)
      setNeed(saved); setMatches(null); setDraft(null); setPublished(null); setChallengeProblem(null)
      setScreen('matching')
      setMatches(await getMatches(saved.id))
    } catch (err) { setError(errorMessage(err)) }
    finally { setBusy(false) }
  }

  async function refreshMatches() {
    if (!need) return
    setBusy(true); setError(''); setMessage('')
    try { setMatches(await getMatches(need.id)) }
    catch (err) { setError(errorMessage(err)) }
    finally { setBusy(false) }
  }

  function prepareChallenge() {
    if (!need) return
    if (!draft) {
      setDraft({
        titulo: `Solução para: ${need.titulo}`.slice(0, 200),
        descricao: (need.resultado_esperado || `Desenvolver uma solução que atenda à necessidade: ${need.descricao}`).slice(0, 5000),
        criterios: (need.criterios_sucesso || []).join('\n'),
        prazo: 7,
      })
    }
    navigate('desafio-draft')
  }

  async function publishChallenge() {
    if (!need || !draft || published) return
    setBusy(true); setError(''); setMessage('')
    try {
      if (!draft.titulo.trim() || !draft.descricao.trim()) throw new Error('Preencha o título e o objetivo do desafio.')
      const criterios = draft.criterios.split('\n').map(value => value.trim()).filter(Boolean)
      if (!criterios.length) throw new Error('Defina pelo menos um critério de sucesso para o desafio.')
      const request = { necessidade_id: need.id, titulo: draft.titulo, descricao: draft.descricao, criterios_sucesso: criterios, prazo_dias: draft.prazo }
      const result = await createDesafio(request)
      setPublished(result)
      setMessage(`Desafio #${result.id} publicado com sucesso.`)
    } catch (err) { setError(errorMessage(err)) }
    finally { setBusy(false) }
  }

  async function openProblem(id: number, destination: 'cadastro' | 'matching' | 'desafio' = 'cadastro') {
    setBusy(true); setError(''); setMessage('')
    try {
      const saved = await getNecessidade(id)
      const currentShop = await getLojista(saved.lojista_id)
      setNeed(saved); setShop(currentShop)
      setForm({
        nomeLoja: currentShop.nome, segmento: currentShop.segmento, titulo: saved.titulo, descricao: saved.descricao, categoria: saved.categoria, prioridade: saved.prioridade,
        impacto: saved.impacto || '', resultadoEsperado: saved.resultado_esperado || '', restricoes: saved.restricoes || '',
        canalVenda: saved.canal_venda || '', criteriosSucesso: (saved.criterios_sucesso || []).join('\n'),
      })
      setMatches(null); setDraft(null); setPublished(null)
      setChallengeProblem(destination === 'desafio' ? saved : null)
      setScreen(destination)
      if (destination === 'matching') setMatches(await getMatches(saved.id))
      return true
    } catch (err) { setError(errorMessage(err)); return false }
    finally { setBusy(false) }
  }

  function challengeUpdated(updated: DesafioResponse) {
    if (published?.id !== updated.id) return
    setPublished(updated)
    setDraft({ titulo: updated.titulo, descricao: updated.descricao, criterios: updated.criterios_sucesso.join('\n'), prazo: updated.prazo_dias })
  }

  return {
    screen, navigate, form, setForm, shop, need, matches, challengeProblem,
    draft, setDraft, published, busy, error, message, newProblem, submitProblem,
    refreshMatches, prepareChallenge, publishChallenge, openProblem, challengeUpdated,
  }
}

export type InnovationFlow = ReturnType<typeof useInnovationFlow>
