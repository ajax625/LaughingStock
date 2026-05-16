import { api } from '../api'

interface Props {
  email: string
  page: string
  onNavigate: (page: string) => void
}

export function Header({ email, page, onNavigate }: Props) {
  async function handleLogout() {
    await api.logout()
    window.location.reload()
  }

  return (
    <header className="bg-gray-900 border-b border-gray-800 px-6 py-4 flex items-center justify-between">
      <div className="flex items-center gap-2">
        <span className="text-xl font-bold text-white tracking-tight">LaughingStock</span>
      </div>
      <nav className="flex items-center gap-6">
        <button
          onClick={() => onNavigate('dashboard')}
          className={`text-sm font-medium transition-colors ${page === 'dashboard' ? 'text-white' : 'text-gray-400 hover:text-white'}`}
        >
          Dashboard
        </button>
        <button
          onClick={() => onNavigate('settings')}
          className={`text-sm font-medium transition-colors ${page === 'settings' ? 'text-white' : 'text-gray-400 hover:text-white'}`}
        >
          Settings
        </button>
        <span className="text-gray-600 text-sm">{email}</span>
        <button
          onClick={handleLogout}
          className="text-sm text-gray-400 hover:text-red-400 transition-colors"
        >
          Logout
        </button>
      </nav>
    </header>
  )
}
