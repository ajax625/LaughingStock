import { useState, useEffect } from 'react'
import { api } from './api'
import { Header } from './components/Header'
import { LoginPage } from './pages/LoginPage'
import { RegisterPage } from './pages/RegisterPage'
import { DashboardPage } from './pages/DashboardPage'
import { SettingsPage } from './pages/SettingsPage'
import type { User } from './types'

type AuthPage = 'login' | 'register'
type AppPage = 'dashboard' | 'settings'

export default function App() {
  const [user, setUser] = useState<User | null>(null)
  const [authPage, setAuthPage] = useState<AuthPage>('login')
  const [appPage, setAppPage] = useState<AppPage>('dashboard')
  const [checking, setChecking] = useState(true)

  useEffect(() => {
    api.getMe()
      .then(setUser)
      .catch(() => setUser(null))
      .finally(() => setChecking(false))
  }, [])

  function handleLogin() {
    api.getMe().then(setUser)
  }

  function handleUpdated() {
    api.getMe().then(setUser)
  }

  if (checking) {
    return (
      <div className="min-h-screen bg-gray-950 flex items-center justify-center text-gray-500">
        Loading...
      </div>
    )
  }

  if (!user) {
    return authPage === 'login'
      ? <LoginPage onLogin={handleLogin} onGoRegister={() => setAuthPage('register')} />
      : <RegisterPage onRegistered={handleLogin} onGoLogin={() => setAuthPage('login')} />
  }

  return (
    <div className="min-h-screen bg-gray-950 flex flex-col">
      <Header email={user.email} page={appPage} onNavigate={p => setAppPage(p as AppPage)} />
      {appPage === 'dashboard'
        ? <DashboardPage />
        : <SettingsPage user={user} onUpdated={handleUpdated} />
      }
    </div>
  )
}
