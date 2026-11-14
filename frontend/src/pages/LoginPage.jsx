import { useState } from 'react'
import { Navigate, useLocation, useNavigate } from 'react-router-dom'
import Logo from '../components/Logo'
import { login } from '../lib/api'
import { isAuthenticated, saveSession } from '../lib/auth'

function LoginPage() {
  const navigate = useNavigate()
  const location = useLocation()
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  // Kalau sudah login, tidak perlu melihat halaman login lagi
  if (isAuthenticated()) {
    return <Navigate to="/dashboard" replace />
  }

  const handleSubmit = async (event) => {
    event.preventDefault()
    setError('')
    setLoading(true)

    try {
      const data = await login({ username, password })
      // data = { token, user }
      saveSession(data)
      // Kembali ke halaman yang tadi dituju (kalau ada), atau ke dashboard
      const redirectTo = location.state?.from?.pathname || '/dashboard'
      navigate(redirectTo, { replace: true })
    } catch (err) {
      setError(err.message || 'Login gagal, silakan coba lagi.')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="flex min-h-screen bg-white">
      {/* Panel kiri — 40% */}
      <aside className="hidden w-2/5 flex-col items-center justify-center border-r border-blue-100 bg-gradient-to-br from-blue-50 to-white px-10 lg:flex">
        <Logo size={72} />
        <h1 className="mt-6 text-4xl font-bold tracking-tight text-blue-700">
          AsmenStore
        </h1>
        <p className="mt-3 max-w-xs text-center text-sm text-slate-500">
          Sistem Inventaris Gudang Laboratorium Komputer
        </p>
      </aside>

      {/* Panel kanan — 60% */}
      <main className="flex w-full items-center justify-center px-6 py-12 lg:w-3/5">
        <div className="w-full max-w-md rounded-2xl border border-blue-100 bg-white p-8 shadow-sm shadow-blue-100/60">
          <div className="mb-8">
            <h2 className="text-2xl font-semibold text-slate-800">
              Selamat datang kembali
            </h2>
            <p className="mt-2 text-sm text-slate-500">
              Silakan masuk untuk melanjutkan ke akun Anda.
            </p>
          </div>

          <form onSubmit={handleSubmit} className="space-y-5">
            <div>
              <label
                htmlFor="username"
                className="mb-1.5 block text-sm font-medium text-slate-700"
              >
                Username
              </label>
              <input id="username" name="username" type="text" autoComplete="username" required
                value={username}
                onChange={(event) => setUsername(event.target.value)}
                placeholder="Masukkan username"
                className="w-full rounded-lg border border-blue-100 bg-blue-50/40 px-4 py-2.5 text-slate-800 placeholder:text-slate-400 transition focus:border-blue-500 focus:bg-white focus:outline-none focus:ring-2 focus:ring-blue-200"
              />
            </div>

            <div>
              <label
                htmlFor="password"
                className="mb-1.5 block text-sm font-medium text-slate-700"
              >
                Password
              </label>
              <input
                id="password"
                name="password"
                type="password"
                autoComplete="current-password"
                required
                value={password}
                onChange={(event) => setPassword(event.target.value)}
                placeholder="Masukkan password"
                className="w-full rounded-lg border border-blue-100 bg-blue-50/40 px-4 py-2.5 text-slate-800 placeholder:text-slate-400 transition focus:border-blue-500 focus:bg-white focus:outline-none focus:ring-2 focus:ring-blue-200"
              />
            </div>

            {error && (
              <p className="rounded-lg bg-red-50 px-4 py-2.5 text-sm text-red-600">
                {error}
              </p>
            )}

            <button
              type="submit"
              disabled={loading}
              className="w-full rounded-lg bg-blue-600 px-4 py-2.5 font-medium text-white transition hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-300 focus:ring-offset-2 active:bg-blue-800 disabled:cursor-not-allowed disabled:opacity-60"
            >
              {loading ? 'Memproses…' : 'Login'}
            </button>
          </form>
        </div>
      </main>
    </div>
  )
}

export default LoginPage
