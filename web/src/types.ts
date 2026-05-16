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

export interface FundamentalSummary {
  pe_ratio: number
  eps_growth_yoy: number
  revenue_growth_yoy: number
  analyst_rating: string
  signal: string
  strength: number
}

export interface InstitutionalSummary {
  current_filings: number
  prior_filings: number
  signal: string
  strength: number
}

export interface InsiderSummary {
  recent_filings: number
  prior_filings: number
  signal: string
  strength: number
}

export interface SmartMoneySummary {
  institutional_13f: InstitutionalSummary
  insider_form4: InsiderSummary
}

export interface DeepResearch {
  id: string
  symbol: string
  status: string
  result?: {
    direction: string
    confidence: number
    fundamental?: FundamentalSummary
    smart_money?: SmartMoneySummary
    signals_fired: string[]
    generated_at: string
  }
}
