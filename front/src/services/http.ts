const baseURL = (import.meta.env.VITE_API_URL || '/api/v1').replace(/\/$/, '')

interface ErrorMessage {
  Property: string
  Message: string
}

export async function request<T>(path: string, options?: RequestInit): Promise<T> {
  let response: Response
  try {
    response = await fetch(`${baseURL}${path}`, {
      ...options,
      headers: { 'Content-Type': 'application/json', ...options?.headers },
    })
  } catch {
    throw new Error('Não foi possível conectar à API. Confira se o backend está em execução.')
  }
  const payload: { data?: T; error?: ErrorMessage | ErrorMessage[] } = await response.json()
    .catch(() => { throw new Error('A API retornou uma resposta inválida.') })
  if (!response.ok) {
    const messages = Array.isArray(payload.error) ? payload.error : payload.error ? [payload.error] : []
    throw new Error(messages.map(error => `${error.Property}: ${error.Message}`).join(' • ') || `Erro HTTP ${response.status}`)
  }
  if (payload.data === undefined) {
    throw new Error('A resposta da API não contém os dados esperados.')
  }
  return payload.data
}

export function errorMessage(error: unknown): string {
  return error instanceof Error ? error.message : 'Não foi possível concluir a operação.'
}
