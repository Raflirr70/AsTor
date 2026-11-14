import { useEffect, useMemo, useState } from 'react'
import {
  createSpecification,
  deleteSpecification,
  getSpecifications,
} from '../lib/api'
import DataTable from './DataTable'
import Icon from './Icons'
import Modal from './Modal'

function SpecificationSection() {
  const [specifications, setSpecifications] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [reloadKey, setReloadKey] = useState(0)

  const [modal, setModal] = useState(null) // { mode: 'create' }
  const [feedback, setFeedback] = useState(null) // { type, message }

  const [formData, setFormData] = useState({
    name: '',
    unit: '',
  })

  const [submitting, setSubmitting] = useState(false)
  const [formError, setFormError] = useState('')

  useEffect(() => {
    let active = true

    async function load() {
      setLoading(true)
      setError('')
      try {
        const specData = await getSpecifications()
        if (!active) return
        setSpecifications(Array.isArray(specData) ? specData : [])
      } catch (err) {
        if (active) setError(err.message || 'Gagal memuat data spesifikasi.')
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
      specifications.map((spec) => ({
        id: spec.id,
        specification: spec,
        name: spec.name || '-',
        unit: spec.unit || '-',
      })),
    [specifications],
  )

  const reload = () => setReloadKey((value) => value + 1)

  const handleOpenCreate = () => {
    setFormData({
      name: '',
      unit: '',
    })
    setFormError('')
    setModal({ mode: 'create' })
  }

  const handleDelete = async (spec) => {
    const confirmed = window.confirm(
      `Hapus spesifikasi "${spec.name}"?`,
    )
    if (!confirmed) return

    setFeedback(null)
    try {
      await deleteSpecification(spec.id)
      setFeedback({ type: 'success', message: 'Spesifikasi berhasil dihapus.' })
      reload()
    } catch (err) {
      setFeedback({
        type: 'error',
        message: err.message || 'Gagal menghapus spesifikasi.',
      })
    }
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    setFormError('')
    setSubmitting(true)

    try {
      await createSpecification(formData)
      setFeedback({ type: 'success', message: 'Spesifikasi berhasil ditambahkan.' })
      setModal(null)
      reload()
    } catch (err) {
      setFormError(err.message || 'Gagal menyimpan data.')
    } finally {
      setSubmitting(false)
    }
  }

  const columns = [
    { key: 'no', label: 'No', sortable: false },
    { key: 'name', label: 'Nama' },
    { key: 'unit', label: 'Unit' },
    {
      key: 'actions',
      label: 'Aksi',
      sortable: false,
      render: (row) => (
        <div className="flex items-center gap-1.5">
          <button
            type="button"
            onClick={() => handleDelete(row.specification)}
            title="Hapus spesifikasi"
            className="rounded-lg border border-red-200 p-1.5 text-red-500 transition hover:bg-red-50"
          >
            <Icon name="trash" className="h-4 w-4" />
          </button>
        </div>
      ),
    },
  ]

  return (
    <div>
      <div className="mb-6 flex flex-wrap items-start justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-slate-800">Specification</h1>
          <p className="mt-1 text-sm text-slate-500">
            Data spesifikasi dari database
            {!loading && !error ? ` — ${rows.length} spesifikasi` : ''}. Klik judul
            kolom untuk mengurutkan.
          </p>
        </div>
        <button
          type="button"
          onClick={handleOpenCreate}
          className="inline-flex items-center gap-2 rounded-lg bg-blue-600 px-4 py-2.5 text-sm font-medium text-white transition hover:bg-blue-700"
        >
          <Icon name="plus" className="h-4 w-4" />
          Tambah Specification
        </button>
      </div>

      {feedback && (
        <div
          className={`mb-4 rounded-lg px-4 py-2.5 text-sm ${
            feedback.type === 'success'
              ? 'bg-emerald-50 text-emerald-600'
              : 'bg-red-50 text-red-600'
          }`}
        >
          {feedback.message}
        </div>
      )}

      {loading && (
        <div className="flex items-center gap-3 rounded-xl border border-blue-100 bg-blue-50/40 p-6 text-sm text-slate-500">
          <span className="h-4 w-4 animate-spin rounded-full border-2 border-blue-200 border-t-blue-600" />
          Memuat data spesifikasi…
        </div>
      )}

      {!loading && error && (
        <div className="rounded-xl border border-red-100 bg-red-50/60 p-6">
          <h2 className="text-lg font-semibold text-red-600">
            Gagal memuat Specification
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
          Belum ada data spesifikasi. Klik <b>Tambah Specification</b> untuk menambahkan.
        </div>
      )}

      {!loading && !error && rows.length > 0 && (
        <DataTable
          title="Daftar Specification"
          subtitle={`${rows.length} spesifikasi dari database — klik judul kolom untuk mengurutkan.`}
          columns={columns}
          rows={rows}
          showHeader={false}
        />
      )}

      {modal && (
        <Modal
          title="Tambah Specification Baru"
          onClose={() => setModal(null)}
        >
          <form onSubmit={handleSubmit} className="space-y-4">
            {formError && (
              <div className="rounded-lg bg-red-50 p-3 text-sm text-red-600">
                {formError}
              </div>
            )}

            <div>
              <label className="mb-1 block text-xs font-semibold uppercase tracking-wider text-slate-600">
                Nama
              </label>
              <input
                type="text"
                required
                value={formData.name}
                onChange={(e) =>
                  setFormData((prev) => ({ ...prev, name: e.target.value }))
                }
                placeholder="cth. RAM"
                className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm text-slate-800 focus:border-blue-500 focus:outline-none"
              />
            </div>

            <div>
              <label className="mb-1 block text-xs font-semibold uppercase tracking-wider text-slate-600">
                Unit
              </label>
              <input
                type="text"
                required
                value={formData.unit}
                onChange={(e) =>
                  setFormData((prev) => ({ ...prev, unit: e.target.value }))
                }
                placeholder="cth. GB"
                className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm text-slate-800 focus:border-blue-500 focus:outline-none"
              />
            </div>

            <div className="flex justify-end gap-3 pt-2">
              <button
                type="button"
                onClick={() => setModal(null)}
                className="rounded-lg border border-slate-200 px-4 py-2 text-sm font-medium text-slate-600 transition hover:bg-slate-50"
              >
                Batal
              </button>
              <button
                type="submit"
                disabled={submitting}
                className="rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white transition hover:bg-blue-700 disabled:opacity-50"
              >
                {submitting ? 'Menyimpan...' : 'Simpan'}
              </button>
            </div>
          </form>
        </Modal>
      )}
    </div>
  )
}

export default SpecificationSection
