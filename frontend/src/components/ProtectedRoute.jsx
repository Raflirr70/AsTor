import { Navigate, useLocation } from 'react-router-dom'
import { isAuthenticated } from '../lib/auth'

/**
 * Membungkus halaman yang hanya boleh diakses setelah login.
 * Kalau belum login → dialihkan ke /login.
 */
function ProtectedRoute({ children }) {
  const location = useLocation()

  if (!isAuthenticated()) {
    return <Navigate to="/login" replace state={{ from: location }} />
  }

  return children
}

export default ProtectedRoute
