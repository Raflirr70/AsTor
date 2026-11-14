import { useEffect, useMemo, useState } from 'react'
import {
  createComputer,
  deleteComputer,
  getAssets,
  getComputers,
  updateComputer,
} from '../lib/api'
import DataTable from './DataTable'
import Icon from './Icons'
import Modal from './Modal'

function formatDateTime(value) {
  if (!value) return '-'
  const date = new Date(value)
  if (Number.isNaN(date.getTime()) || date.getFullYear() <= 1) return '-'
  return date.toLocaleString('id-ID', {
    dateStyle: 'medium',
    timeStyle: 'short',
  })
}

function ComputerSection() {
  const [computers, setComputers] = useState([])
  const [assets, setAssets] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [reloadKey, setReloadKey] = useState(0)

  const [modal, setModal] = useState(null) // { mode: 'create' | 'edit', computer }
  const [feedback, setFeedback] = useState(null) // { type, message }

  const [formData, setFormData] = useState({
    label: '',
    status: 'Baik',
    ram_id: '',
    psu_id: '',
    cpu_id: '',
    vga_id: '',
  })
  const [submitting, setSubmitting] = useState(false)
  const [formError, setFormError] = useState('')

  useEffect(() => {
    let active = true

    async function load() {
      setLoading(true)
      setError('')
      try {
        const [compData, assetData] = await Promise.all([
          getComputers(),
          getAssets().catch(() => []),
        ])
        if (!active) return
        setComputers(Array.isArray(compData) ? compData : [])
        setAssets(Array.isArray(assetData) ? assetData : [])
      } catch (err) {
        if (active) setError(err.message || 'Gagal memuat data computer.')
      } finally {
        if (active) setLoading(false)
      }
    }

    load()
    return () => {
      active = false
    }
  }, [reloadKey])

  // Kumpulkan ID asset yang sudah terpasang di komputer lain
  // Kecuali jika sedang mode edit, asset milik komputer yang sedang diedit boleh dipilih kembali
  const usedAssetIds = useMemo(() => {
    const ids = new Set()
    for (const comp of computers) {
      if (modal?.mode === 'edit' && modal?.computer?.id === comp.id) continue
      if (Array.isArray(comp.component)) {
        for (const c of comp.component) {
          if (c.asset_id) ids.add(c.asset_id)
        }
      }
    }
    return ids
  }, [computers, modal])

  // Filter aset yang belum terpakai
  const availableAssets = useMemo(
    () => assets.filter((a) => !usedAssetIds.has(a.id)),
    [assets, usedAssetIds],
  )

  // Filter berdasarkan kode kategori (cocokkan awalan kode atau nama kategori)
  const filterByCategoryCode = (item, codePrefixes) => {
    const code = (item.code || '').toUpperCase()
    const cat = (item.category || '').toUpperCase()
    for (const prefix of codePrefixes) {
      const p = prefix.toUpperCase()
      if (code.startsWith(p) || code.includes(`-${p}-`) || cat.includes(p)) {
        return true
      }
    }
    return false
  }

  const ramAssets = useMemo(
    () => availableAssets.filter((a) => filterByCategoryCode(a, ['RAM'])),
    [availableAssets],
  )
  const psuAssets = useMemo(
    () => availableAssets.filter((a) => filterByCategoryCode(a, ['PSU', 'CBP', 'POWER'])),
    [availableAssets],
  )
  const cpuAssets = useMemo(
    () => availableAssets.filter((a) => filterByCategoryCode(a, ['CPU', 'PROC'])),
    [availableAssets],
  )
  const vgaAssets = useMemo(
    () => availableAssets.filter((a) => filterByCategoryCode(a, ['VGA', 'GPU'])),
    [availableAssets],
  )

  const rows = useMemo(
    () =>
      computers.map((comp) => {
        const codesText = Array.isArray(comp.component)
          ? comp.component.map((c) => c.asset_code).filter(Boolean).join(', ')
          : '-'

        return {
          id: comp.id,
          computer: comp,
          code: comp.code || `#${comp.id}`,
          label: comp.label,
          category: comp.category || 'Computer',
          status: comp.status || 'Baik',
          componentsSummary: codesText || '-',
          updatedAt: formatDateTime(comp.updated_at),
        }
      }),
    [computers],
  )

  const reload = () => setReloadKey((value) => value + 1)

  const handleOpenCreate = () => {
    setFormData({
      label: '',
      status: 'Baik',
      ram_id: '',
      psu_id: '',
      cpu_id: '',
      vga_id: '',
    })
    setFormError('')
    setModal({ mode: 'create', computer: null })
  }

  const handleOpenEdit = (comp) => {
    // Map komponen berdasarkan kategori atau urutan
    let ramId = ''
    let psuId = ''
    let cpuId = ''
    let vgaId = ''

    if (Array.isArray(comp.component)) {
      for (const c of comp.component) {
        const assetObj = assets.find((a) => a.id === c.asset_id)
        const code = (assetObj?.code || '').toUpperCase()
        const cat = (assetObj?.category || '').toUpperCase()

        if (code.includes('RAM') || cat.includes('RAM')) {
          ramId = c.asset_id
        } else if (code.includes('PSU') || code.includes('CBP') || cat.includes('PSU') || cat.includes('POWER')) {
          psuId = c.asset_id
        } else if (code.includes('CPU') || code.includes('PROC') || cat.includes('CPU') || cat.includes('PROCESSOR')) {
          cpuId = c.asset_id
        } else if (code.includes('VGA') || code.includes('GPU') || cat.includes('VGA') || cat.includes('GRAPHICS')) {
          vgaId = c.asset_id
        } else {
          // Fallback urutan
          if (!ramId) ramId = c.asset_id
          else if (!psuId) psuId = c.asset_id
          else if (!cpuId) cpuId = c.asset_id
          else if (!vgaId) vgaId = c.asset_id
        }
      }
    }

    setFormData({
      label: comp.label || '',
      status: comp.status || 'Aktif',
      ram_id: ramId,
      psu_id: psuId,
      cpu_id: cpuId,
      vga_id: vgaId,
    })
    setFormError('')
    setModal({ mode: 'edit', computer: comp })
  }

  const handleDelete = async (comp) => {
    const compCodes = Array.isArray(comp.component) && comp.component.length > 0
      ? comp.component.map((c) => c.asset_code).filter(Boolean).join(', ')
      : 'tidak ada'

    const message = `Hapus computer "${comp.label}"?\n\nKomponen aset berikut yang terpasang juga akan ikut terhapus:\n- ${compCodes}`
    const confirmed = window.confirm(message)
    if (!confirmed) return

    setFeedback(null)
    try {
      await deleteComputer(comp.id)
      setFeedback({ type: 'success', message: 'Computer beserta komponennya berhasil dihapus.' })
      reload()
    } catch (err) {
      setFeedback({
        type: 'error',
        message: err.message || 'Gagal menghapus computer.',
      })
    }
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    setFormError('')
    setSubmitting(true)

    try {
      const components = [
        formData.ram_id ? Number(formData.ram_id) : null,
        formData.psu_id ? Number(formData.psu_id) : null,
        formData.cpu_id ? Number(formData.cpu_id) : null,
        formData.vga_id ? Number(formData.vga_id) : null,
      ].filter(Boolean)

      const payload = {
        label: formData.label,
        status: formData.status,
        components,
      }

      if (modal.mode === 'create') {
        await createComputer(payload)
        setFeedback({ type: 'success', message: 'Computer berhasil ditambahkan.' })
      } else {
        await updateComputer(modal.computer.id, payload)
        setFeedback({ type: 'success', message: 'Computer berhasil diperbarui.' })
      }
      setModal(null)
      reload()
    } catch (err) {
      setFormError(err.message || 'Gagal menyimpan computer.')
    } finally {
      setSubmitting(false)
    }
  }

  const columns = [
    { key: 'no', label: 'No', sortable: false },
    { key: 'code', label: 'Kode' },
    { key: 'label', label: 'Label' },
    { key: 'status', label: 'Status' },
    { key: 'componentsSummary', label: 'Komponen / Spesifikasi' },
    { key: 'updatedAt', label: 'Diperbarui' },
    {
      key: 'actions',
      label: 'Aksi',
      sortable: false,
      render: (row) => (
        <div className="flex items-center gap-1.5">
          <button
            type="button"
            onClick={() => handleOpenEdit(row.computer)}
            title="Update computer"
            className="rounded-lg border border-blue-200 p-1.5 text-blue-600 transition hover:bg-blue-50"
          >
            <Icon name="edit" className="h-4 w-4" />
          </button>
          <button
            type="button"
            onClick={() => handleDelete(row.computer)}
            title="Hapus computer"
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
          <h1 className="text-2xl font-bold text-slate-800">Computer</h1>
          <p className="mt-1 text-sm text-slate-500">
            Data komputer dari database
            {!loading && !error ? ` — ${rows.length} komputer` : ''}. Klik judul
            kolom untuk mengurutkan.
          </p>
        </div>
        <button
          type="button"
          onClick={handleOpenCreate}
          className="inline-flex items-center gap-2 rounded-lg bg-blue-600 px-4 py-2.5 text-sm font-medium text-white transition hover:bg-blue-700"
        >
          <Icon name="plus" className="h-4 w-4" />
          Tambah Computer
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
          Memuat data computer…
        </div>
      )}

      {!loading && error && (
        <div className="rounded-xl border border-red-100 bg-red-50/60 p-6">
          <h2 className="text-lg font-semibold text-red-600">
            Gagal memuat Computer
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
          Belum ada data computer. Klik <b>Tambah Computer</b> untuk menambahkan.
        </div>
      )}

      {!loading && !error && rows.length > 0 && (
        <DataTable
          title="Daftar Computer"
          subtitle={`${rows.length} komputer dari database — klik judul kolom untuk mengurutkan.`}
          columns={columns}
          rows={rows}
          showHeader={false}
        />
      )}

      {modal && (
        <Modal
          title={modal.mode === 'create' ? 'Tambah Computer Baru' : 'Edit Computer'}
          onClose={() => setModal(null)}
        >
          <form onSubmit={handleSubmit} className="space-y-4">
            {formError && (
              <div className="rounded-lg bg-red-50 p-3 text-sm text-red-600">
                {formError}
              </div>
            )}

            <div>
              <label htmlFor="comp-label" className="mb-1 block text-xs font-semibold uppercase tracking-wider text-slate-600">
                Label Komputer
              </label>
              <input
                id="comp-label"
                name="label"
                type="text"
                required
                value={formData.label}
                onChange={(e) =>
                  setFormData((prev) => ({ ...prev, label: e.target.value }))
                }
                placeholder="cth. PC-L1-01"
                className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm text-slate-800 focus:border-blue-500 focus:outline-none"
              />
            </div>

            <div>
              <label htmlFor="comp-status" className="mb-1 block text-xs font-semibold uppercase tracking-wider text-slate-600">
                Status
              </label>
              <select
                id="comp-status"
                name="status"
                value={formData.status}
                onChange={(e) =>
                  setFormData((prev) => ({ ...prev, status: e.target.value }))
                }
                className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm text-slate-800 focus:border-blue-500 focus:outline-none"
              >
                <option value="Baik">Baik</option>
                <option value="Maintenance">Maintenance</option>
                <option value="Rusak">Rusak</option>
              </select>
            </div>

            <div>
              <label htmlFor="comp-ram" className="mb-1 block text-xs font-semibold uppercase tracking-wider text-slate-600">
                Pilihan RAM
              </label>
              <select
                id="comp-ram"
                name="ram_id"
                value={formData.ram_id}
                onChange={(e) =>
                  setFormData((prev) => ({ ...prev, ram_id: e.target.value }))
                }
                className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm text-slate-800 focus:border-blue-500 focus:outline-none"
              >
                <option value="">-- Pilih RAM --</option>
                {ramAssets.map((a) => (
                  <option key={a.id} value={a.id}>
                    {a.label} ({a.code})
                  </option>
                ))}
              </select>
            </div>

            <div>
              <label htmlFor="comp-psu" className="mb-1 block text-xs font-semibold uppercase tracking-wider text-slate-600">
                Pilihan PSU
              </label>
              <select
                id="comp-psu"
                name="psu_id"
                value={formData.psu_id}
                onChange={(e) =>
                  setFormData((prev) => ({ ...prev, psu_id: e.target.value }))
                }
                className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm text-slate-800 focus:border-blue-500 focus:outline-none"
              >
                <option value="">-- Pilih PSU --</option>
                {psuAssets.map((a) => (
                  <option key={a.id} value={a.id}>
                    {a.label} ({a.code})
                  </option>
                ))}
              </select>
            </div>

            <div>
              <label htmlFor="comp-cpu" className="mb-1 block text-xs font-semibold uppercase tracking-wider text-slate-600">
                Pilihan CPU
              </label>
              <select
                id="comp-cpu"
                name="cpu_id"
                value={formData.cpu_id}
                onChange={(e) =>
                  setFormData((prev) => ({ ...prev, cpu_id: e.target.value }))
                }
                className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm text-slate-800 focus:border-blue-500 focus:outline-none"
              >
                <option value="">-- Pilih CPU --</option>
                {cpuAssets.map((a) => (
                  <option key={a.id} value={a.id}>
                    {a.label} ({a.code})
                  </option>
                ))}
              </select>
            </div>

            <div>
              <label htmlFor="comp-vga" className="mb-1 block text-xs font-semibold uppercase tracking-wider text-slate-600">
                Pilihan VGA
              </label>
              <select
                id="comp-vga"
                name="vga_id"
                value={formData.vga_id}
                onChange={(e) =>
                  setFormData((prev) => ({ ...prev, vga_id: e.target.value }))
                }
                className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm text-slate-800 focus:border-blue-500 focus:outline-none"
              >
                <option value="">-- Pilih VGA --</option>
                {vgaAssets.map((a) => (
                  <option key={a.id} value={a.id}>
                    {a.label} ({a.code})
                  </option>
                ))}
              </select>
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

export default ComputerSection