// Base URL backend. Bisa dioverride lewat frontend/.env → VITE_API_BASE_URL
import { getToken } from './auth'

const BASE_URL =
  import.meta.env.VITE_API_BASE_URL ?? 'https://radio-unpopular-saline.ngrok-free.dev'

/**
 * Wrapper fetch ke backend.
 * - Otomatis set Content-Type JSON
 * - Header `ngrok-skip-browser-warning` supaya request tidak kena halaman
 *   peringatan ngrok (khusus tunnel gratis)
 * - Melempar Error berisi pesan dari backend kalau response tidak sukses
 */
async function apiFetch(path, { method = 'GET', body, token } = {}) {
  const headers = {
    'ngrok-skip-browser-warning': 'true',
  }
  if (body !== undefined) {
    headers['Content-Type'] = 'application/json'
  }
  if (token) {
    headers.Authorization = `Bearer ${token}`
  }

  const res = await fetch(`${BASE_URL}${path}`, {
    method,
    headers,
    body: body !== undefined ? JSON.stringify(body) : undefined,
  })

  let payload = null
  try {
    payload = await res.json()
  } catch {
    // response bukan JSON (mis. halaman HTML ngrok) → biarkan payload null
  }

  if (!res.ok || payload?.success === false) {
    const base = payload?.message || `Terjadi kesalahan (HTTP ${res.status})`
    const details =
      Array.isArray(payload?.errors) && payload.errors.length > 0
        ? `: ${payload.errors.join(', ')}`
        : ''
    throw new Error(base + details)
  }

  // Struktur response backend: { success, code, message, data }
  return payload?.data ?? payload
}

export function getUsers() {
  return apiFetch('/admin/users', { token: getToken() })
}
export function login({ username, password }) {
  return apiFetch('/login', {
    method: 'POST',
    body: { username, password },
  })
}

// ===== Assets (butuh login + role "PJ Gudang") =====
export function getAssets() {
  return apiFetch('/gudang/assets', { token: getToken() })
}

// Dipakai untuk menerjemahkan category_id → nama kategori
export function getCategories() {
  return apiFetch('/gudang/categories', { token: getToken() })
}

export function createCategory(payload) {
  return apiFetch('/gudang/categories', {
    method: 'POST',
    body: payload,
    token: getToken(),
  })
}

export function updateCategory(id, payload) {
  return apiFetch(`/gudang/category/${id}`, {
    method: 'PUT',
    body: payload,
    token: getToken(),
  })
}

export function deleteCategory(id) {
  return apiFetch(`/gudang/category/${id}`, {
    method: 'DELETE',
    token: getToken(),
  })
}

export function getSpecifications() {
  return apiFetch('/gudang/specifications', { token: getToken() })
}

export function createSpecification(payload) {
  return apiFetch('/gudang/specifications', {
    method: 'POST',
    body: payload,
    token: getToken(),
  })
}

export function deleteSpecification(id) {
  return apiFetch(`/gudang/specification/${id}`, {
    method: 'DELETE',
    token: getToken(),
  })
}

export function createAsset(payload) {
  return apiFetch('/gudang/assets', {
    method: 'POST',
    body: payload,
    token: getToken(),
  })
}

export function updateAsset(id, payload) {
  return apiFetch(`/gudang/asset/${id}`, {
    method: 'PUT',
    body: payload,
    token: getToken(),
  })
}

export function deleteAsset(id) {
  return apiFetch(`/gudang/asset/${id}`, {
    method: 'DELETE',
    token: getToken(),
  })
}

// ===== Laboratories (butuh login + role "PJ Gudang" / "Admin") =====
export function getLaboratories() {
  return apiFetch('/gudang/labkoms', { token: getToken() })
}
export function getComputers() {
  return apiFetch('/gudang/computers', { token: getToken() })
}

export function createComputer(payload) {
  return apiFetch('/gudang/computers', {
    method: 'POST',
    body: payload,
    token: getToken(),
  })
}

export function updateComputer(id, payload) {
  return apiFetch(`/gudang/computer/${id}`, {
    method: 'PUT',
    body: payload,
    token: getToken(),
  })
}

export function deleteComputer(id) {
  return apiFetch(`/gudang/computer/${id}`, {
    method: 'DELETE',
    token: getToken(),
  })
}

// ===== Loans =====
export function getLoans() {
  return apiFetch('/asisten/loans', { token: getToken() })
}

export function createLoan(payload) {
  return apiFetch('/asisten/loans', {
    method: 'POST',
    body: payload,
    token: getToken(),
  })
}

export function updateLoan(id, payload) {
  return apiFetch(`/asisten/loan/${id}`, {
    method: 'PUT',
    body: payload,
    token: getToken(),
  })
}

export function deleteLoan(id) {
  return apiFetch(`/asisten/loan/${id}`, {
    method: 'DELETE',
    token: getToken(),
  })
}

export function returnLoan(id, payload) {
  return apiFetch(`/asisten/loan/${id}/returned`, {
    method: 'PATCH',
    body: payload,
    token: getToken(),
  })
}
