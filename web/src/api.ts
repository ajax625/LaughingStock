const base = ''

async function request<T>(method: string, path: string, body?: unknown): Promise<T> {
  const res = await fetch(base + path, {
    method,
    credentials: 'include',
    headers: body ? { 'Content-Type': 'application/json' } : undefined,
    body: body ? JSON.stringify(body) : undefined,
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }))
    throw new Error(err.error || res.statusText)
  }
  return res.json()
}

export const api = {
  // Auth
  register: (email: string, password: string) =>
    request('POST', '/auth/register', { email, password }),
  login: (email: string, password: string) =>
    request<{ id: string; email: string }>('POST', '/auth/login', { email, password }),
  logout: () => request('POST', '/auth/logout'),

  // Me
  getMe: () => request<{ id: string; email: string; telegram_chat_id: string; has_alpaca: boolean; webhook_token: string }>('GET', '/me'),
  updateMe: (data: { telegram_chat_id?: string; alpaca_key?: string; alpaca_secret?: string }) =>
    request('PATCH', '/me', data),

  // Symbols
  getSymbols: () => request<{ symbols: import('./types').UserSymbol[] }>('GET', '/me/symbols'),
  addSymbol: (data: { symbol: string; trader_type: string; min_roi: number; min_rr: number; min_confidence: number }) =>
    request('POST', '/me/symbols', data),
  updateSymbol: (id: string, data: { min_roi: number; min_rr: number; min_confidence: number }) =>
    request('PATCH', `/me/symbols/${id}`, data),
  deleteSymbol: (id: string) => request('DELETE', `/me/symbols/${id}`),

  // Signals
  getSignals: () => request<{ signals: import('./types').SymbolSignalEntry[] }>('GET', '/me/signals'),

  // Quotes
  getQuotes: () => request<Record<string, number>>('GET', '/me/quotes'),

  // Trade
  executeTrade: (symbol: string, side: 'buy' | 'sell', qty: number) =>
    request<{ order_id: string; symbol: string; side: string; qty: number; status: string }>(
      'POST', '/me/trade', { symbol, side, qty }
    ),
}
