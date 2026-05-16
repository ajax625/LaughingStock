import { useEffect, useState } from 'react'
import { api } from '../api'
import type { DeepResearch } from '../types'

interface Props {
  symbolId: string
  symbol: string
  onClose: () => void
}

function SignalBadge({ signal }: { signal: string }) {
  const cls = signal === 'BULLISH'
    ? 'text-green-400 bg-green-900/40 border border-green-700'
    : signal === 'BEARISH'
      ? 'text-red-400 bg-red-900/40 border border-red-700'
      : 'text-gray-400 bg-gray-800 border border-gray-700'
  return <span className={`text-xs font-semibold px-2 py-0.5 rounded-full ${cls}`}>{signal || 'N/A'}</span>
}

export function ResearchModal({ symbolId, symbol, onClose }: Props) {
  const [data, setData] = useState<DeepResearch | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    api.getSymbolResearch(symbolId).then(setData).finally(() => setLoading(false))
  }, [symbolId])

  const result = data?.result

  return (
    <div className="fixed inset-0 bg-black/70 flex items-center justify-center z-50 p-4" onClick={onClose}>
      <div
        className="bg-gray-900 border border-gray-700 rounded-2xl w-full max-w-lg max-h-[90vh] overflow-y-auto"
        onClick={e => e.stopPropagation()}
      >
        {/* Header */}
        <div className="flex items-center justify-between p-5 border-b border-gray-800">
          <div>
            <h2 className="text-lg font-bold text-white">{symbol} — Fundamentals</h2>
            {result?.generated_at && (
              <p className="text-xs text-gray-500 mt-0.5">
                {new Date(result.generated_at).toLocaleDateString(undefined, { dateStyle: 'medium' })}
              </p>
            )}
          </div>
          <button onClick={onClose} className="text-gray-500 hover:text-white transition-colors text-xl leading-none">×</button>
        </div>

        <div className="p-5 space-y-5">
          {loading && <p className="text-center text-gray-500 py-8">Loading...</p>}

          {!loading && (!data || data.status === 'none' || data.status === 'unavailable') && (
            <p className="text-center text-gray-500 py-8">No research available yet. It runs automatically when you add a symbol.</p>
          )}

          {!loading && data?.status === 'pending' && (
            <p className="text-center text-gray-500 py-8">Research is running — check back in a minute.</p>
          )}

          {result && (
            <>
              {/* Overall */}
              <div className="bg-gray-800 rounded-xl p-4 flex items-center justify-between">
                <span className="text-sm text-gray-400 font-medium">Overall</span>
                <div className="flex items-center gap-3">
                  <SignalBadge signal={result.direction} />
                  <span className="text-sm text-gray-300">{Math.round(result.confidence * 100)}% confidence</span>
                </div>
              </div>

              {/* Fundamentals */}
              {result.fundamental && (
                <div>
                  <h3 className="text-xs text-gray-500 uppercase tracking-wider mb-2">Company Fundamentals</h3>
                  <div className="bg-gray-800 rounded-xl p-4 space-y-3">
                    <div className="flex justify-between items-center">
                      <span className="text-sm text-gray-400">Signal</span>
                      <SignalBadge signal={result.fundamental.signal} />
                    </div>
                    <div className="grid grid-cols-2 gap-3 text-sm">
                      <div>
                        <div className="text-gray-500 text-xs mb-0.5">P/E Ratio</div>
                        <div className="text-white font-mono">{result.fundamental.pe_ratio?.toFixed(1) || '—'}</div>
                      </div>
                      <div>
                        <div className="text-gray-500 text-xs mb-0.5">Analyst Rating</div>
                        <div className="text-white">{result.fundamental.analyst_rating || '—'}</div>
                      </div>
                      <div>
                        <div className="text-gray-500 text-xs mb-0.5">EPS Growth YoY</div>
                        <div className={result.fundamental.eps_growth_yoy > 0 ? 'text-green-400 font-mono' : 'text-red-400 font-mono'}>
                          {result.fundamental.eps_growth_yoy > 0 ? '+' : ''}{(result.fundamental.eps_growth_yoy * 100).toFixed(1)}%
                        </div>
                      </div>
                      <div>
                        <div className="text-gray-500 text-xs mb-0.5">Revenue Growth YoY</div>
                        <div className={result.fundamental.revenue_growth_yoy > 0 ? 'text-green-400 font-mono' : 'text-red-400 font-mono'}>
                          {result.fundamental.revenue_growth_yoy > 0 ? '+' : ''}{(result.fundamental.revenue_growth_yoy * 100).toFixed(1)}%
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              )}

              {/* Smart Money */}
              {result.smart_money && (
                <div>
                  <h3 className="text-xs text-gray-500 uppercase tracking-wider mb-2">Smart Money</h3>
                  <div className="bg-gray-800 rounded-xl p-4 space-y-3">
                    <div className="flex justify-between items-center">
                      <span className="text-sm text-gray-400">Institutional (13F)</span>
                      <div className="flex items-center gap-2">
                        <span className="text-xs text-gray-500">
                          {result.smart_money.institutional_13f.prior_filings} → {result.smart_money.institutional_13f.current_filings} filings
                        </span>
                        <SignalBadge signal={result.smart_money.institutional_13f.signal} />
                      </div>
                    </div>
                    <div className="flex justify-between items-center">
                      <span className="text-sm text-gray-400">Insider (Form 4)</span>
                      <div className="flex items-center gap-2">
                        <span className="text-xs text-gray-500">
                          {result.smart_money.insider_form4.recent_filings} recent filings
                        </span>
                        <SignalBadge signal={result.smart_money.insider_form4.signal} />
                      </div>
                    </div>
                  </div>
                </div>
              )}
            </>
          )}
        </div>
      </div>
    </div>
  )
}
