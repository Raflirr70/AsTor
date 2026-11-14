import { useEffect, useMemo, useState } from 'react'
import { deleteAsset, getAssets, getCategories } from '../lib/api'
import AssetFormModal from './AssetFormModal'
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

function formatSpecs(specs) {
  if (!Array.isArray(specs) || specs.length === 0) return '-'
  return specs
    .map((spec) =>
      [spec.name, spec.value, spec.unit].filter(Boolean).join(' ').trim(),
    )
    .filter(Boolean)
    .join(', ')
}

function AssetsSection() {
  const [assets, setAssets] = useState([])
  const [categories, setCategories] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [reloadKey, setReloadKey] = useState(0)

  const [modal, setModal] = useState(null) // { mode: 'create' | 'edit', asset }
  const [feedback, setFeedback] = useState(null) // { type, message }

  useEffect(() => {
    let active = true

    async function load() {
      setLoading(true)
      setError('')
      try {
        const [assetData, categoryData] = await Promise.all([
          getAssets(),
          // Kategori bersifat pelengkap; kalau gagal, tetap tampilkan assets
          getCategories().catch(() => []),
        ])
        if (!active) return
        setAssets(Array.isArray(assetData) ? assetData : [])
        setCategories(Array.isArray(categoryData) ? categoryData : [])
      } catch (err) {
        if (active) setError(err.message || 'Gagal memuat data assets.')
      } finally {
        if (active) setLoading(false)
      }
    }

    load()
    return () => {
      active = false
    }
  }, [reloadKey])

  const categoryMap = useMemo(() => {
    const map = {}
    for (const category of categories) {
      map[category.id] = category.name
    }
    return map
  }, [categories])

  const rows = useMemo(
    () =>
      assets.map((asset) => ({
        id: asset.id,
        asset,
        code: asset.code,
        label: asset.label,
        categoryName: categoryMap[asset.category_id] ?? `#${asset.category_id}`,
        status: asset.status,
        specs: formatSpecs(asset.asset_specifications),
        checkedAt: formatDateTime(asset.checked_at),
        createdAt: formatDateTime(asset.created_at),
      })),
    [assets, categoryMap],
  )

  const reload = () => setReloadKey((value) => value + 1)

  const handleDelete = async (asset) => {
    const confirmed = window.confirm(
      `Hapus asset "${asset.label}" (${asset.code})?`,
    )
    if (!confirmed) return

    setFeedback(null)
    try {
      await deleteAsset(asset.id)
      setFeedback({ type: 'success', message: 'Asset berhasil dihapus.' })
      reload()
    } catch (err) {
      setFeedback({
        type: 'error',
        message: err.message || 'Gagal menghapus asset.',
      })
    }
  }

  const handleSuccess = (message) => {
    setModal(null)
    setFeedback({ type: 'success', message })
    reload()
  }

  const columns = [
    { key: 'no', label: 'No', sortable: false },
    { key: 'code', label: 'Kode' },
    { key: 'label', label: 'Label' },
    { key: 'categoryName', label: 'Kategori' },
    { key: 'status', label: 'Status' },
    { key: 'specs', label: 'Spesifikasi' },
    { key: 'checkedAt', label: 'Terakhir Dicek' },
    { key: 'createdAt', label: 'Barang Masuk' },
    {
      key: 'actions',
      label: 'Aksi',
      sortable: false,
      render: (row) => (
        <div className="flex items-center gap-1.5">
          <button
            type="button"
            onClick={() => setModal({ mode: 'edit', asset: row.asset })}
            title="Update asset"
            className="rounded-lg border border-blue-200 p-1.5 text-blue-600 transition hover:bg-blue-50"
          >
            <Icon name="edit" className="h-4 w-4" />
          </button>
          <button
            type="button"
            onClick={() => handleDelete(row.asset)}
            title="Hapus asset"
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
          <h1 className="text-2xl font-bold text-slate-800">Assets</h1>
          <p className="mt-1 text-sm text-slate-500">
            Data dari database
            {!loading && !error ? ` — ${rows.length} asset` : ''}. Klik judul
            kolom untuk mengurutkan.
          </p>
        </div>
        <button
          type="button"
          onClick={() => setModal({ mode: 'create', asset: null })}
          className="inline-flex items-center gap-2 rounded-lg bg-blue-600 px-4 py-2.5 text-sm font-medium text-white transition hover:bg-blue-700"
        >
          <Icon name="plus" className="h-4 w-4" />
          Tambah Asset
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
          Memuat data assets…
        </div>
      )}

      {!loading && error && (
        <div className="rounded-xl border border-red-100 bg-red-50/60 p-6">
          <h2 className="text-lg font-semibold text-red-600">
            Gagal memuat Assets
          </h2>
          <p className="mt-1 text-sm text-red-500">{error}</p>
          <p className="mt-2 text-xs text-slate-500">
            Pastikan backend berjalan, kamu login dengan role <b>PJ Gudang</b>,
            dan token masih valid.
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
      )}

      {!loading && !error && rows.length === 0 && (
        <div className="rounded-xl border border-blue-100 bg-blue-50/40 p-6 text-sm text-slate-500">
          Belum ada data assets. Klik <b>Tambah Asset</b> untuk menambahkan.
        </div>
      )}

      {!loading && !error && rows.length > 0 && (
        <DataTable
          title="Daftar Asset"
          subtitle={`${rows.length} asset dari database — klik judul kolom untuk mengurutkan.`}
          columns={columns}
          rows={rows}
          showHeader={false}
        />
      )}

      {modal && (
        <AssetFormModal
          mode={modal.mode}
          asset={modal.asset}
          categories={categories}
          onClose={() => setModal(null)}
          onSuccess={handleSuccess}
        />
      )}
    </div>
  )
}

export default AssetsSection
