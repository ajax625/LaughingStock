import { useState, useEffect, useCallback } from 'react'
import { api } from '../api'
import { useWebSocket } from '../hooks/useWebSocket'
import { SymbolCard } from '../components/SymbolCard'
import { AddSymbolModal } from '../components/AddSymbolModal'
import type { UserSymbol, Signal, IntradaySignal, SymbolSignalEntry } from '../types'

function isIntraday(s: unknown): s is IntradaySignal {
  return typeof s === 'object' && s !== null && 'long' in s && 'short' in s
}

export function DashboardPage() {
  const [symbols, setSymbols] = useState<UserSymbol[]>([])
  const [signalMap, setSignalMap] = useState<Record<string, Signal | null>>({})
  const [liveSignalMap, setLiveSignalMap] = useState<Record<string, IntradaySignal>>({})
  const [alertedSymbols, setAlertedSymbols] = useState<Set<string>>(new Set())
  const [showAddModal, setShowAddModal] = useState(false)
  const [loading, setLoading] = useState(true)

  const loadData = useCallback(async () => {
    const [symRes, sigRes] = await Promise.all([
      api.getSymbols(),
      api.getSignals(),
    ])
    setSymbols(symRes.symbols ?? [])
    const map: Record<string, Signal | null> = {}
    ;(sigRes.signals ?? []).forEach((entry: SymbolSignalEntry) => {
      map[entry.symbol] = entry.signal
    })
    setSignalMap(map)
    setLoading(false)
  }, [])

  useEffect(() => { loadData() }, [loadData])

  // WebSocket — receives live signals from TradeX via LaughingStock
  useWebSocket((data) => {
    const payload = data as Record<string, unknown>
    const sym = payload.symbol as string | undefined
    if (!sym) return

    if (isIntraday(payload)) {
      setLiveSignalMap(prev => ({ ...prev, [sym]: payload as IntradaySignal }))
    } else {
      setSignalMap(prev => ({ ...prev, [sym]: payload as unknown as Signal }))
    }

    // Trigger alert state — auto-dismiss after 60s
    setAlertedSymbols(prev => new Set([...prev, sym]))
    setTimeout(() => {
      setAlertedSymbols(prev => {
        const next = new Set(prev)
        next.delete(sym)
        return next
      })
    }, 60_000)
  })

  async function handleDelete(id: string) {
    await api.deleteSymbol(id)
    loadData()
  }

  if (loading) {
    return (
      <div className="flex-1 flex items-center justify-center text-gray-500">
        Loading...
      </div>
    )
  }

  return (
    <div className="flex-1 p-6">
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-xl font-bold text-white">Watchlist</h2>
        <button
          onClick={() => setShowAddModal(true)}
          className="bg-blue-600 hover:bg-blue-500 text-white text-sm font-semibold px-4 py-2 rounded-lg transition-colors"
        >
          + Add Symbol
        </button>
      </div>

      {symbols.length === 0 ? (
        <div className="text-center py-20 text-gray-500">
          <p className="text-lg mb-2">No symbols yet</p>
          <p className="text-sm">Add a symbol to start receiving signals</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
          {symbols.map(us => (
            <SymbolCard
              key={us.id}
              userSymbol={us}
              signal={signalMap[us.symbol] ?? null}
              liveSignal={liveSignalMap[us.symbol] ?? null}
              alerted={alertedSymbols.has(us.symbol)}
              onDelete={handleDelete}
            />
          ))}
        </div>
      )}

      {showAddModal && (
        <AddSymbolModal
          onClose={() => setShowAddModal(false)}
          onAdded={loadData}
        />
      )}
    </div>
  )
}
