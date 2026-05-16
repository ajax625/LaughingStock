import { useState } from 'react'
import { api } from '../api'

interface Props {
  onClose: () => void
  onAdded: () => void
}

export function AddSymbolModal({ onClose, onAdded }: Props) {
  const [symbol, setSymbol] = useState('')
  const [minConf, setMinConf] = useState('60')
  const [minRoi, setMinRoi] = useState('0.5')
  const [minRr, setMinRr] = useState('2.0')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      await api.addSymbol({
        symbol: symbol.toUpperCase().trim(),
        trader_type: 'day',
        min_roi: parseFloat(minRoi),
        min_rr: parseFloat(minRr),
        min_confidence: parseFloat(minConf) / 100,
      })
      onAdded()
      onClose()
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : 'Failed to add symbol')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="fixed inset-0 bg-black/60 flex items-center justify-center z-50" onClick={onClose}>
      <div className="bg-gray-900 border border-gray-700 rounded-xl p-6 w-full max-w-sm" onClick={e => e.stopPropagation()}>
        <h2 className="text-lg font-bold text-white mb-4">Add Symbol</h2>
        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
          <div>
            <label className="text-xs text-gray-400 mb-1 block">Symbol</label>
            <input
              value={symbol}
              onChange={e => setSymbol(e.target.value)}
              required
              placeholder="TSLA"
              className="w-full bg-gray-800 border border-gray-700 rounded-lg px-3 py-2 text-white uppercase placeholder-gray-600 focus:outline-none focus:border-blue-500"
            />
          </div>
          <div className="grid grid-cols-3 gap-3">
            <div>
              <label className="text-xs text-gray-400 mb-1 block">Min Conf %</label>
              <input
                type="number"
                value={minConf}
                onChange={e => setMinConf(e.target.value)}
                min="1" max="100"
                className="w-full bg-gray-800 border border-gray-700 rounded-lg px-3 py-2 text-white focus:outline-none focus:border-blue-500"
              />
            </div>
            <div>
              <label className="text-xs text-gray-400 mb-1 block">Min ROI %</label>
              <input
                type="number"
                value={minRoi}
                onChange={e => setMinRoi(e.target.value)}
                step="0.1" min="0"
                className="w-full bg-gray-800 border border-gray-700 rounded-lg px-3 py-2 text-white focus:outline-none focus:border-blue-500"
              />
            </div>
            <div>
              <label className="text-xs text-gray-400 mb-1 block">Min R/R</label>
              <input
                type="number"
                value={minRr}
                onChange={e => setMinRr(e.target.value)}
                step="0.1" min="0"
                className="w-full bg-gray-800 border border-gray-700 rounded-lg px-3 py-2 text-white focus:outline-none focus:border-blue-500"
              />
            </div>
          </div>
          {error && <p className="text-red-400 text-sm">{error}</p>}
          <div className="flex gap-3 pt-2">
            <button
              type="button"
              onClick={onClose}
              className="flex-1 bg-gray-800 hover:bg-gray-700 text-gray-300 rounded-lg py-2 text-sm transition-colors"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={loading}
              className="flex-1 bg-blue-600 hover:bg-blue-500 disabled:opacity-50 text-white rounded-lg py-2 text-sm font-semibold transition-colors"
            >
              {loading ? 'Adding...' : 'Add Symbol'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}
