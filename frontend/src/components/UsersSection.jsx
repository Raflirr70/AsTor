import { useEffect, useMemo, useState } from 'react'
import { getUsers } from '../lib/api'
import DataTable from './DataTable'
import Icon from './Icons'

function UsersSection() {
  const [users, setUsers] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [reloadKey, setReloadKey] = useState(0)

  useEffect(() => {
    let active = true

    async function load() {
      setLoading(true)
      setError('')
      try {
        const data = await getUsers()
        if (!active) return
        setUsers(Array.isArray(data) ? data : [])
      } catch (err) {
        if (active) setError(err.message || 'Gagal memuat data users.')
      } finally {
        if (active) setLoading(false)
      }
    }

    load()
    return () => {
      active = false
    }
  }, [reloadKey])

  const rows = useMemo(
    () =>
      users.map((u) => {
        const addr = u.address || '-'
        const shortAddr = addr.length > 30 ? addr.substring(0, 30) + '...' : addr
        return {
          id: u.id,
          npm: u.npm || '-',
          firstname: u.firstname || '-',
          lastname: u.lastname || '-',
          address: (
            <span title={addr} className="cursor-help">
              {shortAddr}
            </span>
          ),
          birthday: u.birthday ? u.birthday.split('T')[0] : '-',
          username: u.username || '-',
        }
      }),
    [users],
  )

  const reload = () => setReloadKey((value) => value + 1)

  const columns = [
    { key: 'no', label: 'No', sortable: false },
    { key: 'npm', label: 'NPM' },
    { key: 'firstname', label: 'First Name' },
    { key: 'lastname', label: 'Last Name' },
    { key: 'username', label: 'Username' },
    { key: 'birthday', label: 'Birthday' },
    { key: 'address', label: 'Address' },
  ]

  return (
    <div>
      <div className="mb-6 flex flex-wrap items-start justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-slate-800">Users</h1>
          <p className="mt-1 text-sm text-slate-500">
            Data pengguna dari database
            {!loading && !error ? ` — ${rows.length} pengguna` : ''}. Klik judul
            kolom untuk mengurutkan.
          </p>
        </div>
      </div>

      {loading && (
        <div className="flex items-center gap-3 rounded-xl border border-blue-100 bg-blue-50/40 p-6 text-sm text-slate-500">
          <span className="h-4 w-4 animate-spin rounded-full border-2 border-blue-200 border-t-blue-600" />
          Memuat data users…
        </div>
      )}

      {!loading && error && (
        <div className="rounded-xl border border-red-100 bg-red-50/60 p-6">
          <h2 className="text-lg font-semibold text-red-600">
            Gagal memuat Users
          </h2>
          <p className="mt-1 text-sm text-red-500">{error}</p>
          <button
            type="button"
            onClick={reload}
            className="mt-4 inline-flex items-center gap-2 rounded-lg bg-red-500 px-3 py-1.5 text-sm font-medium text-white transition hover:bg-red-600"
          >
            <Icon name="history" className="h-4 w-4" />
            Coba lagi
          </button>
        </div>
      )}

      {!loading && !error && rows.length === 0 && (
        <div className="rounded-xl border border-blue-100 bg-blue-50/40 p-6 text-sm text-slate-500">
          Belum ada data users.
        </div>
      )}

      {!loading && !error && rows.length > 0 && (
        <DataTable
          title="Daftar Users"
          subtitle={`${rows.length} pengguna dari database — klik judul kolom untuk mengurutkan.`}
          columns={columns}
          rows={rows}
          showHeader={false}
        />
      )}
    </div>
  )
}

export default UsersSection
