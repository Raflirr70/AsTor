import { useEffect, useMemo, useState } from 'react'
import {
  createCategory,
  deleteCategory,
  getCategories,
  getSpecifications,
  updateCategory,
} from '../lib/api'
import DataTable from './DataTable'
import Icon from './Icons'
import Modal from './Modal'

function CategorySection() {
  const [categories, setCategories] = useState([])
  const [specifications, setSpecifications] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [reloadKey, setReloadKey] = useState(0)

  const [modal, setModal] = useState(null) // { mode: 'create' | 'edit', category }
  const [feedback, setFeedback] = useState(null) // { type, message }

  const [formData, setFormData] = useState({
    name: '',
    code: '',
    specification_ids: [],
  })

  const [submitting, setSubmitting] = useState(false)
  const [formError, setFormError] = useState('')

  useEffect(() => {
    let active = true

    async function load() {
      setLoading(true)
      setError('')
      try {
        const [catData, specData] = await Promise.all([
          getCategories(),
          getSpecifications().catch(() => []),
        ])
        if (!active) return
        setCategories(Array.isArray(catData) ? catData : [])
        setSpecifications(Array.isArray(specData) ? specData : [])
      } catch (err) {
        if (active) setError(err.message || 'Gagal memuat data category.')
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
      categories.map((cat) => ({
        id: cat.id,
        category: cat,
        code: cat.code,
        name: cat.name,
        specifications:
          Array.isArray(cat.specifications) && cat.specifications.length > 0
            ? cat.specifications
                .map((s) => `${s.name}${s.unit ? ` (${s.unit})` : ''}`)
                .join(', ')
            : '-',
      })),
    [categories],
  )

  const reload = () => setReloadKey((value) => value + 1)

  const handleOpenCreate = () => {
    setFormData({ name: '', code: '', specification_ids: [] })
    setFormError('')
    setModal({ mode: 'create', category: null })
  }

  const handleOpenEdit = (cat) => {
    setFormData({
      name: cat.name || '',
      code: cat.code || '',
      specification_ids: Array.isArray(cat.specifications)
        ? cat.specifications.map((s) => s.id)
        : [],
    })
    setFormError('')
    setModal({ mode: 'edit', category: cat })
  }

  const handleDelete = async (cat) => {
    const confirmed = window.confirm(
      `Hapus category "${cat.name}" (${cat.code})?`,
    )
    if (!confirmed) return

    setFeedback(null)
    try {
      await deleteCategory(cat.id)
      setFeedback({ type: 'success', message: 'Category berhasil dihapus.' })
      reload()
    } catch (err) {
      setFeedback({
        type: 'error',
        message: err.message || 'Gagal menghapus category.',
      })
    }
  }

  const handleCheckboxToggle = (specId) => {
    setFormData((prev) => {
      const exists = prev.specification_ids.includes(specId)
      return {
        ...prev,
        specification_ids: exists
          ? prev.specification_ids.filter((id) => id !== specId)
          : [...prev.specification_ids, specId],
      }
    })
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    setFormError('')
    setSubmitting(true)

    try {
      if (modal.mode === 'create') {
        await createCategory(formData)
        setFeedback({ type: 'success', message: 'Category berhasil ditambahkan.' })
      } else {
        await updateCategory(modal.category.id, formData)
        setFeedback({ type: 'success', message: 'Category berhasil diperbarui.' })
      }
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
    { key: 'code', label: 'Kode' },
    { key: 'name', label: 'Nama' },
    { key: 'specifications', label: 'Spesifikasi' },
    {
      key: 'actions',
      label: 'Aksi',
      sortable: false,
      render: (row) => (
        <div className="flex items-center gap-1.5">
          <button
            type="button"
            onClick={() => handleOpenEdit(row.category)}
            title="Update category"
            className="rounded-lg border border-blue-200 p-1.5 text-blue-600 transition hover:bg-blue-50"
          >
            <Icon name="edit" className="h-4 w-4" />
          </button>
          <button
            type="button"
            onClick={() => handleDelete(row.category)}
            title="Hapus category"
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
          <h1 className="text-2xl font-bold text-slate-800">Category</h1>
          <p className="mt-1 text-sm text-slate-500">
            Data kategori dari database
            {!loading && !error ? ` — ${rows.length} kategori` : ''}. Klik judul
            kolom untuk mengurutkan.
          </p>
        </div>
        <button
          type="button"
          onClick={handleOpenCreate}
          className="inline-flex items-center gap-2 rounded-lg bg-blue-600 px-4 py-2.5 text-sm font-medium text-white transition hover:bg-blue-700"
        >
          <Icon name="plus" className="h-4 w-4" />
          Tambah Category
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
          Memuat data category…
        </div>
      )}

      {!loading && error && (
        <div className="rounded-xl border border-red-100 bg-red-50/60 p-6">
          <h2 className="text-lg font-semibold text-red-600">
            Gagal memuat Category
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
          Belum ada data category. Klik <b>Tambah Category</b> untuk menambahkan.
        </div>
      )}

      {!loading && !error && rows.length > 0 && (
        <DataTable
          title="Daftar Category"
          subtitle={`${rows.length} kategori dari database — klik judul kolom untuk mengurutkan.`}
          columns={columns}
          rows={rows}
          showHeader={false}
        />
      )}

      {modal && (
        <Modal
          title={modal.mode === 'create' ? 'Tambah Category Baru' : 'Edit Category'}
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
                Nama Kategori
              </label>
              <input
                type="text"
                required
                value={formData.name}
                onChange={(e) =>
                  setFormData((prev) => ({ ...prev, name: e.target.value }))
                }
                placeholder="cth. Cable HDMI"
                className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm text-slate-800 focus:border-blue-500 focus:outline-none"
              />
            </div>

            <div>
              <label className="mb-1 block text-xs font-semibold uppercase tracking-wider text-slate-600">
                Kode Kategori
              </label>
              <input
                type="text"
                required
                value={formData.code}
                onChange={(e) =>
                  setFormData((prev) => ({ ...prev, code: e.target.value }))
                }
                placeholder="cth. CBH"
                className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm text-slate-800 focus:border-blue-500 focus:outline-none"
              />
            </div>

            <div>
              <label className="mb-2 block text-xs font-semibold uppercase tracking-wider text-slate-600">
                Spesifikasi (Dapat memilih lebih dari 1)
              </label>
              <div className="max-h-48 space-y-2 overflow-y-auto rounded-lg border border-slate-200 p-3">
                {specifications.length === 0 ? (
                  <p className="text-xs text-slate-400">
                    Tidak ada spesifikasi tersedia.
                  </p>
                ) : (
                  specifications.map((spec) => {
                    const checked = formData.specification_ids.includes(spec.id)
                    return (
                      <label
                        key={spec.id}
                        className="flex items-center gap-2.5 text-sm text-slate-700 cursor-pointer select-none"
                      >
                        <input
                          type="checkbox"
                          checked={checked}
                          onChange={() => handleCheckboxToggle(spec.id)}
                          className="h-4 w-4 rounded border-slate-300 text-blue-600 focus:ring-blue-500"
                        />
                        <span>
                          {spec.name} {spec.unit ? `(${spec.unit})` : ''}
                        </span>
                      </label>
                    )
                  })
                )}
              </div>
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

export default CategorySection