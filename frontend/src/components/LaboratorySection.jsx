import { useEffect, useMemo, useState } from 'react'
import { getLaboratories } from '../lib/api'
import DataTable from './DataTable'
import Icon from './Icons'

function formatDateTime(value) {
  if (!value) return '-'
  const date = new Date(value)
  if (Number.isNaN(date.getTime()) || date.getFullYear() <= 1) return '-'
  return date.toLocaleString('id-ID', {
    dateStyle: 'medium',
    timeStyle: 'short',
  })
}

const columns = [
  { key: 'no', label: 'No', sortable: false },
  { key: 'name', label: 'Nama' },
  { key: 'capasity', label: 'Kapasitas' },
  { key: 'checkedAt', label: 'Terakhir Dicek' },
  { key: 'updatedAt', label: 'Diperbarui' },
]

function LaboratorySection() {
  const [laboratories, setLaboratories] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [reloadKey, setReloadKey] = useState(0)

  useEffect(() => {
    let active = true

    async function load() {
      setLoading(true)
      setError('')
      try {
        const data = await getLaboratories()
        if (!active) return
        setLaboratories(Array.isArray(data) ? data : [])
      } catch (err) {
        if (active) setError(err.message || 'Gagal memuat data laboratorium.')
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
      laboratories.map((lab) => ({
        id: lab.id,
        name: lab.name,
        locationId: lab.location_id,
        capasity: lab.capasity,
        checkedAt: formatDateTime(lab.checked_at),
        updatedAt: formatDateTime(lab.updated_at),
      })),
    [laboratories],
  )

  const reload = () => setReloadKey((value) => value + 1)

  if (loading) {
    return (
      <div className="flex items-center gap-3 rounded-xl border border-blue-100 bg-blue-50/40 p-6 text-sm text-slate-500">
        <span className="h-4 w-4 animate-spin rounded-full border-2 border-blue-200 border-t-blue-600" />
        Memuat data laboratorium…
      </div>
    )
  }

  if (error) {
    return (
      <div className="rounded-xl border border-red-100 bg-red-50/60 p-6">
        <h1 className="text-lg font-semibold text-red-600">
          Gagal memuat Laboratorium
        </h1>
        <p className="mt-1 text-sm text-red-500">{error}</p>
        <p className="mt-2 text-xs text-slate-500">
          Pastikan backend berjalan, kamu login dengan role <b>Admin</b>, dan
          token masih valid.
        </p>
        <button
          type="button"
          onClick={reload}
          className="mt-4 inline-flex items-center gap-2 rounded-lg bg-red-500 px-3 py-1.5 text-sm font-medium text-white transition hover:bg-red-600"
        >
          <Icon name="history" className="h-4 w-4" />
          Coba lagi
        </button>
      </div>
    )
  }

  if (rows.length === 0) {
    return (
      <div className="rounded-xl border border-blue-100 bg-blue-50/40 p-6 text-sm text-slate-500">
        Belum ada data laboratorium di database.
      </div>
    )
  }

  return (
    <DataTable
      title="Laboratorium Computer"
      subtitle={`${rows.length} laboratorium dari database — klik judul kolom untuk mengurutkan.`}
      columns={columns}
      rows={rows}
    />
  )
}

export default LaboratorySection
