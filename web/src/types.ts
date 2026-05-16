export interface User {
  id: string
  email: string
  telegram_chat_id: string
  has_alpaca: boolean
  webhook_token: string
}

export interface UserSymbol {
  id: string
  user_id: string
  symbol: string
  trader_type: string
  min_roi: number
  min_rr: number
  min_confidence: number
  tradex_subscription_id: string
  active: boolean
}

export interface Signal {
  id: string
  symbol: string
  direction: 'LONG' | 'SHORT'
  confidence: number
  entry: number
  target: number
  stop: number
  timeframe: string
  signals_fired: string[]
  mode: string
  generated_at: string
}

export interface IntradaySignal {
  symbol: string
  direction: 'LONG' | 'SHORT'
  confidence: number
  long: TradeLevels
  short: TradeLevels
}

export interface TradeLevels {
  entry: number
  target: number
  stop: number
  roi: number
  rr: number
}

export type LiveSignal = Signal | IntradaySignal

export interface SymbolSignalEntry {
  symbol: string
  signal: Signal | null
}
