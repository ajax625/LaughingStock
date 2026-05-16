import { useState } from 'react'
import { api } from '../api'
import { ResearchModal } from './ResearchModal'
import type { UserSymbol, Signal, IntradaySignal } from '../types'

interface Props {
  userSymbol: UserSymbol
  currentPrice: number | null
  signal: Signal | null
  liveSignal: IntradaySignal | null
  alerted: boolean
  onDelete: (id: string) => void
}

function isIntraday(s: Signal | IntradaySignal): s is IntradaySignal {
  return 'long' in s && 'short' in s
}

export function SymbolCard({ userSymbol, currentPrice, signal, liveSignal, alerted, onDelete }: Props) {
  const [qty, setQty] = useState('1')
  const [trading, setTrading] = useState(false)
  const [tradeResult, setTradeResult] = useState<string | null>(null)
  const [showResearch, setShowResearch] = useState(false)

  const displaySignal = liveSignal ?? signal

  const direction = displaySignal?.direction
  const confidence = displaySignal ? Math.round(displaySignal.confidence * 100) : null

  let entry: number | null = null
  let target: number | null = null
  let stop: number | null = null
  let roi: number | null = null
  let rr: number | null = null

  if (displaySignal) {
    if (isIntraday(displaySignal)) {
      const levels = displaySignal.direction === 'LONG' ? displaySignal.long : displaySignal.short
      entry = levels.entry
      target = levels.target
      stop = levels.stop
      roi = levels.roi
      rr = levels.rr
    } else {
      entry = displaySignal.entry
      target = displaySignal.target
      stop = displaySignal.stop
    }
  }

  async function handleExecute(side: 'buy' | 'sell') {
    const q = parseFloat(qty)
    if (isNaN(q) || q <= 0) return
    setTrading(true)
    setTradeResult(null)
    try {
      const result = await api.executeTrade(userSymbol.symbol, side, q)
      setTradeResult(`Order placed: ${result.order_id} (${result.status})`)
    } catch (e: unknown) {
      setTradeResult(`Failed: ${e instanceof Error ? e.message : 'Unknown error'}`)
    } finally {
      setTrading(false)
    }
  }

  const borderColor = alerted
    ? direction === 'LONG' ? 'border-green-500 shadow-green-500/20 shadow-lg' : 'border-red-500 shadow-red-500/20 shadow-lg'
    : 'border-gray-800'

  const directionBadge = direction === 'LONG'
    ? 'bg-green-900/50 text-green-400 border border-green-700'
    : direction === 'SHORT'
      ? 'bg-red-900/50 text-red-400 border border-red-700'
      : 'bg-gray-800 text-gray-400 border border-gray-700'

  return (
    <>
      <div className={`bg-gray-900 rounded-xl border ${borderColor} p-5 flex flex-col gap-4 transition-all duration-300`}>
        {/* Header */}
        <div className="flex items-start justify-between">
          <div>
            <div className="flex items-baseline gap-2">
              <h3 className="text-xl font-bold text-white">{userSymbol.symbol}</h3>
              {currentPrice !== null && (
                <span className="text-lg font-mono font-semibold text-gray-100">
                  ${currentPrice.toFixed(2)}
                </span>
              )}
            </div>
            <span className="text-xs text-gray-500">{userSymbol.trader_type} trader</span>
          </div>
          <div className="flex items-center gap-2">
            {direction && (
              <span className={`text-xs font-semibold px-2 py-1 rounded-full ${directionBadge}`}>
                {direction}
              </span>
            )}
            {confidence !== null && (
              <span className="text-xs text-gray-400">{confidence}%</span>
            )}
          </div>
        </div>

        {/* Signal levels */}
        {entry !== null ? (
          <div className="grid grid-cols-3 gap-2 text-sm">
            <div className="bg-gray-800 rounded-lg p-2 text-center">
              <div className="text-gray-500 text-xs mb-1">Entry</div>
              <div className="text-white font-mono font-semibold">${entry.toFixed(2)}</div>
            </div>
            <div className="bg-gray-800 rounded-lg p-2 text-center">
              <div className="text-gray-500 text-xs mb-1">Target</div>
              <div className="text-green-400 font-mono font-semibold">${target!.toFixed(2)}</div>
            </div>
            <div className="bg-gray-800 rounded-lg p-2 text-center">
              <div className="text-gray-500 text-xs mb-1">Stop</div>
              <div className="text-red-400 font-mono font-semibold">${stop!.toFixed(2)}</div>
            </div>
            {roi !== null && rr !== null && (
              <>
                <div className="bg-gray-800 rounded-lg p-2 text-center">
                  <div className="text-gray-500 text-xs mb-1">ROI</div>
                  <div className="text-blue-400 font-mono font-semibold">{roi.toFixed(1)}%</div>
                </div>
                <div className="bg-gray-800 rounded-lg p-2 text-center col-span-2">
                  <div className="text-gray-500 text-xs mb-1">R/R</div>
                  <div className="text-blue-400 font-mono font-semibold">{rr.toFixed(1)}</div>
                </div>
              </>
            )}
          </div>
        ) : (
          <div className="text-center text-gray-600 text-sm py-2">Awaiting signal...</div>
        )}

        {/* Thresholds */}
        <div className="grid grid-cols-3 gap-2">
          <div className="bg-gray-800/60 rounded-lg p-2 text-center">
            <div className="text-gray-600 text-xs mb-0.5">Min Conf</div>
            <div className="text-gray-300 text-sm font-semibold">{Math.round(userSymbol.min_confidence * 100)}%</div>
          </div>
          <div className="bg-gray-800/60 rounded-lg p-2 text-center">
            <div className="text-gray-600 text-xs mb-0.5">Min ROI</div>
            <div className="text-gray-300 text-sm font-semibold">{userSymbol.min_roi}%</div>
          </div>
          <div className="bg-gray-800/60 rounded-lg p-2 text-center">
            <div className="text-gray-600 text-xs mb-0.5">Min R/R</div>
            <div className="text-gray-300 text-sm font-semibold">{userSymbol.min_rr}</div>
          </div>
        </div>

        {/* Alert: execute controls */}
        {alerted && (
          <div className="bg-gray-800 rounded-lg p-3 border border-yellow-600/40">
            <div className="text-yellow-400 text-xs font-semibold mb-2">Signal Alert — Ready to Execute</div>
            <div className="flex gap-2">
              <input
                type="number"
                value={qty}
                onChange={e => setQty(e.target.value)}
                min="1"
                step="1"
                className="w-20 bg-gray-700 border border-gray-600 rounded px-2 py-1 text-sm text-white text-center"
                placeholder="Qty"
              />
              <button
                onClick={() => handleExecute('buy')}
                disabled={trading}
                className="flex-1 bg-green-700 hover:bg-green-600 disabled:opacity-50 text-white text-sm font-semibold rounded px-3 py-1 transition-colors"
              >
                Buy
              </button>
              <button
                onClick={() => handleExecute('sell')}
                disabled={trading}
                className="flex-1 bg-red-700 hover:bg-red-600 disabled:opacity-50 text-white text-sm font-semibold rounded px-3 py-1 transition-colors"
              >
                Sell
              </button>
            </div>
          </div>
        )}

        {/* Manual execute (always available when signal present) */}
        {!alerted && entry !== null && (
          <div className="flex gap-2">
            <input
              type="number"
              value={qty}
              onChange={e => setQty(e.target.value)}
              min="1"
              className="w-20 bg-gray-800 border border-gray-700 rounded px-2 py-1 text-sm text-white text-center"
              placeholder="Qty"
            />
            <button
              onClick={() => handleExecute('buy')}
              disabled={trading}
              className="flex-1 bg-gray-700 hover:bg-green-700 disabled:opacity-50 text-white text-sm rounded px-3 py-1 transition-colors"
            >
              Buy
            </button>
            <button
              onClick={() => handleExecute('sell')}
              disabled={trading}
              className="flex-1 bg-gray-700 hover:bg-red-700 disabled:opacity-50 text-white text-sm rounded px-3 py-1 transition-colors"
            >
              Sell
            </button>
          </div>
        )}

        {tradeResult && (
          <div className="text-xs text-gray-400 text-center">{tradeResult}</div>
        )}

        {/* Footer actions */}
        <div className="flex items-center justify-between pt-1 border-t border-gray-800">
          <button
            onClick={() => setShowResearch(true)}
            className="text-xs text-blue-400 hover:text-blue-300 transition-colors"
          >
            Fundamentals →
          </button>
          <button
            onClick={() => onDelete(userSymbol.id)}
            className="text-xs text-gray-600 hover:text-red-400 transition-colors"
          >
            Remove
          </button>
        </div>
      </div>

      {showResearch && (
        <ResearchModal
          symbolId={userSymbol.id}
          symbol={userSymbol.symbol}
          onClose={() => setShowResearch(false)}
        />
      )}
    </>
  )
}
