import { useState } from 'react'
import { api } from '../api'
import type { User } from '../types'

interface Props {
  user: User
  onUpdated: () => void
}

export function SettingsPage({ user, onUpdated }: Props) {
  const [telegramChatID, setTelegramChatID] = useState(user.telegram_chat_id || '')
  const [alpacaKey, setAlpacaKey] = useState('')
  const [alpacaSecret, setAlpacaSecret] = useState('')
  const [saved, setSaved] = useState(false)
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setError('')
    setSaved(false)
    setLoading(true)
    try {
      await api.updateMe({
        telegram_chat_id: telegramChatID,
        ...(alpacaKey ? { alpaca_key: alpacaKey, alpaca_secret: alpacaSecret } : {}),
      })
      setSaved(true)
      onUpdated()
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : 'Update failed')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="max-w-lg mx-auto py-10 px-4">
      <h2 className="text-2xl font-bold text-white mb-6">Settings</h2>

      <div className="bg-gray-900 border border-gray-800 rounded-xl p-6 mb-6">
        <h3 className="text-sm font-semibold text-gray-400 uppercase tracking-wide mb-1">Webhook Token</h3>
        <p className="text-xs text-gray-500 mb-2">This is your unique webhook URL registered with TradeX.</p>
        <code className="block bg-gray-800 rounded px-3 py-2 text-xs text-blue-300 break-all">
          https://laughingstock.informatrixlabs.com/webhook/{user.webhook_token}
        </code>
      </div>

      <form onSubmit={handleSubmit} className="bg-gray-900 border border-gray-800 rounded-xl p-6 flex flex-col gap-5">
        <div>
          <h3 className="text-sm font-semibold text-gray-400 uppercase tracking-wide mb-3">Telegram</h3>
          <label className="text-xs text-gray-400 mb-1 block">
            Chat ID{' '}
            <a
              href="https://t.me/userinfobot"
              target="_blank"
              rel="noreferrer"
              className="text-blue-400 hover:underline"
            >
              (get yours here)
            </a>
          </label>
          <input
            value={telegramChatID}
            onChange={e => setTelegramChatID(e.target.value)}
            placeholder="123456789"
            className="w-full bg-gray-800 border border-gray-700 rounded-lg px-3 py-2 text-white placeholder-gray-600 focus:outline-none focus:border-blue-500"
          />
        </div>

        <div>
          <h3 className="text-sm font-semibold text-gray-400 uppercase tracking-wide mb-3">
            Alpaca Paper Trading {user.has_alpaca && <span className="text-green-400 normal-case font-normal">(configured)</span>}
          </h3>
          <div className="flex flex-col gap-3">
            <div>
              <label className="text-xs text-gray-400 mb-1 block">API Key</label>
              <input
                value={alpacaKey}
                onChange={e => setAlpacaKey(e.target.value)}
                placeholder={user.has_alpaca ? '••••••••••••' : 'Enter API key'}
                className="w-full bg-gray-800 border border-gray-700 rounded-lg px-3 py-2 text-white placeholder-gray-600 focus:outline-none focus:border-blue-500"
              />
            </div>
            <div>
              <label className="text-xs text-gray-400 mb-1 block">API Secret</label>
              <input
                type="password"
                value={alpacaSecret}
                onChange={e => setAlpacaSecret(e.target.value)}
                placeholder={user.has_alpaca ? '••••••••••••' : 'Enter API secret'}
                className="w-full bg-gray-800 border border-gray-700 rounded-lg px-3 py-2 text-white placeholder-gray-600 focus:outline-none focus:border-blue-500"
              />
            </div>
          </div>
        </div>

        {error && <p className="text-red-400 text-sm">{error}</p>}
        {saved && <p className="text-green-400 text-sm">Settings saved.</p>}

        <button
          type="submit"
          disabled={loading}
          className="bg-blue-600 hover:bg-blue-500 disabled:opacity-50 text-white rounded-lg py-2 font-semibold transition-colors"
        >
          {loading ? 'Saving...' : 'Save Settings'}
        </button>
      </form>
    </div>
  )
}
