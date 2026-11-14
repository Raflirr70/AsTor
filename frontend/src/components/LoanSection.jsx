import { useEffect, useMemo, useState } from 'react'
import {
  createLoan,
  getAssets,
  getLoans,
  returnLoan,
} from '../lib/api'
import DataTable from './DataTable'
import Icon from './Icons'
import Modal from './Modal'

function formatDateTime(value) {
  if (!value) return '-'
  const date = new Date(value)
  if (Number.isNaN(date.getTime()) || date.getFullYear() <= 1) return '-'
  return date.toLocaleDateString('id-ID', {
    day: 'numeric',
    month: 'short',
    year: '2-digit',
  })
}

function LoanSection() {
  const [loans, setLoans] = useState([])
  const [assets, setAssets] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [reloadKey, setReloadKey] = useState(0)

  const [modal, setModal] = useState(null) // { mode: 'create' | 'return', loan }
  const [feedback, setFeedback] = useState(null) // { type, message }

  const [formData, setFormData] = useState({
    loans_name: '',
    loans_npm: '',
    items: [{ asset_code: '', quantity: 1 }],
  })

  const [returnStatuses, setReturnStatuses] = useState({}) // asset_id -> status
  const [submitting, setSubmitting] = useState(false)
  const [formError, setFormError] = useState('')

  useEffect(() => {
    let active = true

    async function load() {
      setLoading(true)
      setError('')
      try {
        const [loanData, assetData] = await Promise.all([
          getLoans(),
          getAssets().catch(() => []),
        ])
        if (!active) return
        setLoans(Array.isArray(loanData) ? loanData : [])
        setAssets(Array.isArray(assetData) ? assetData : [])
      } catch (err) {
        if (active) setError(err.message || 'Gagal memuat data loan.')
      } finally {
        if (active) setLoading(false)
      }
    }

    load()
    return () => {
      active = false
    }
  }, [reloadKey])

  // Filter aset yang boleh dipinjam:
  // - Status bukan 'Rusak', 'Dipinjam', atau 'Hilang'
  // - Belum dipilih di baris item lain dalam form
  const availableAssetsForIndex = (currentIndex) => {
    const selectedCodes = new Set(
      formData.items
        .map((item, idx) => (idx !== currentIndex ? item.asset_code : null))
        .filter(Boolean),
    )

    return assets.filter((a) => {
      const st = (a.status || '').toLowerCase()
      if (st.includes('rusak') || st.includes('dipinjam') || st.includes('hilang') || st.includes('broken')) {
        return false
      }
      if (selectedCodes.has(a.code)) {
        return false
      }
      return true
    })
  }

  const rows = useMemo(
    () =>
      loans.map((loan, idx) => {
        const details = Array.isArray(loan.loan_details) ? loan.loan_details : []
        const combinedAsset = details.length > 0
          ? details.map(d => {
              const code = d.asset_code || d.asset?.code || '-'
              const name = d.asset_name || d.asset?.label || '-'
              return code !== '-' && name !== '-' ? `${code} (${name})` : (code !== '-' ? code : name)
            }).join(', ')
          : '-'

        return {
          id: loan.id,
          no: idx + 1,
          loan,
          loans_name: loan.loans_name,
          loans_npm: loan.loans_npm,
          combinedAsset,
          userLoanName: loan.user_loan_name || '-',
          userReturnName: loan.user_return_name || '-',
          status: loan.status,
          createdAt: formatDateTime(loan.created_at),
          deadline: formatDateTime(loan.deadline),
          returnedAt: formatDateTime(loan.returned_at),
        }
      }),
    [loans],
  )

  const reload = () => setReloadKey((value) => value + 1)

  const handleOpenCreate = () => {
    setFormData({
      loans_name: '',
      loans_npm: '',
      items: [{ asset_code: '', quantity: 1 }],
    })
    setFormError('')
    setModal({ mode: 'create', loan: null })
  }

  const handleOpenReturn = (loan) => {
    const initialStatuses = {}
    if (Array.isArray(loan.loan_details)) {
      loan.loan_details.forEach(d => {
        initialStatuses[d.asset_id] = 'Baik'
      })
    }
    setReturnStatuses(initialStatuses)
    setFormError('')
    setModal({ mode: 'return', loan })
  }

  const handleAddItem = () => {
    setFormData((prev) => ({
      ...prev,
      items: [...prev.items, { asset_code: '', quantity: 1 }],
    }))
  }

  const handleRemoveItem = (index) => {
    setFormData((prev) => ({
      ...prev,
      items: prev.items.filter((_, i) => i !== index),
    }))
  }

  const handleItemChange = (index, field, value) => {
    setFormData((prev) => {
      const newItems = [...prev.items]
      newItems[index] = { ...newItems[index], [field]: value }
      return { ...prev, items: newItems }
    })
  }

  const assetHasQuantitySpec = (assetCode) => {
    const asset = assets.find((a) => a.code === assetCode)
    if (!asset || !Array.isArray(asset.asset_specifications)) return false
    return asset.asset_specifications.some(
      (spec) =>
        spec.name?.toLowerCase().includes('quantity') ||
        spec.name?.toLowerCase().includes('jumlah') ||
        spec.unit?.toLowerCase() === 'pcs',
    )
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    setFormError('')
    setSubmitting(true)

    try {
      if (modal.mode === 'return') {
        await returnLoan(modal.loan.id, { statuses: returnStatuses, status: 'Baik' })
        setFeedback({ type: 'success', message: 'Asset berhasil dikembalikan.' })
      } else {
        const payload = {
          loans_name: formData.loans_name,
          loans_npm: formData.loans_npm,
          items: formData.items.map((item) => ({
            asset_code: item.asset_code,
            quantity: assetHasQuantitySpec(item.asset_code) ? Number(item.quantity) || 1 : 1,
          })),
        }

        await createLoan(payload)
        setFeedback({ type: 'success', message: 'Loan berhasil ditambahkan.' })
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
    { key: 'loans_name', label: 'Peminjam' },
    { key: 'loans_npm', label: 'NPM' },
    { key: 'combinedAsset', label: 'Barang & Kode' },
    { key: 'userLoanName', label: 'Saksi' },
    { key: 'userReturnName', label: 'Pengembalian' },
    { key: 'createdAt', label: 'Tgl pinjam' },
    { key: 'deadline', label: 'Deadline' },
    { key: 'returnedAt', label: 'Dikembalikan' },
    {
      key: 'status',
      label: 'Status',
      render: (row) => {
        const isDone = row.loan.status === 'Done'
        return isDone ? (
          <span className="inline-flex rounded-full bg-emerald-50 px-2.5 py-1 text-xs font-semibold text-emerald-600">
            Done
          </span>
        ) : (
          <button
            type="button"
            onClick={() => handleOpenReturn(row.loan)}
            title="Pengembalian (Return)"
            className="rounded-lg border border-amber-200 bg-amber-50 px-3 py-1 text-xs font-semibold text-amber-600 transition hover:bg-amber-100"
          >
            Return
          </button>
        )
      }
    },
  ]

  return (
    <div>
      <div className="mb-6 flex flex-wrap items-start justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-slate-800">Loans (Peminjaman)</h1>
          <p className="mt-1 text-sm text-slate-500">
            Data peminjaman dari database
            {!loading && !error ? ` — ${rows.length} peminjaman` : ''}. Klik judul
            kolom untuk mengurutkan.
          </p>
        </div>
        <button
          type="button"
          onClick={handleOpenCreate}
          className="inline-flex items-center gap-2 rounded-lg bg-blue-600 px-4 py-2.5 text-sm font-medium text-white transition hover:bg-blue-700"
        >
          <Icon name="plus" className="h-4 w-4" />
          Tambah Loan
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
          Memuat data loans…
        </div>
      )}

      {!loading && error && (
        <div className="rounded-xl border border-red-100 bg-red-50/60 p-6">
          <h2 className="text-lg font-semibold text-red-600">
            Gagal memuat Loans
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
          Belum ada data loan. Klik <b>Tambah Loan</b> untuk menambahkan.
        </div>
      )}

      {!loading && !error && rows.length > 0 && (
        <DataTable
          title="Daftar Loans"
          subtitle={`${rows.length} peminjaman dari database — klik judul kolom untuk mengurutkan.`}
          columns={columns}
          rows={rows}
          showHeader={false}
        />
      )}

      {modal && (
        <Modal
          title={modal.mode === 'create' ? 'Tambah Loan Baru' : 'Pengembalian Aset (Return)'}
          onClose={() => setModal(null)}
        >
          <form onSubmit={handleSubmit} className="space-y-4">
            {formError && (
              <div className="rounded-lg bg-red-50 p-3 text-sm text-red-600">
                {formError}
              </div>
            )}

            {modal.mode === 'return' ? (
              <div className="space-y-4">
                <p className="text-xs font-medium text-slate-500">
                  Pilih kondisi pengembalian untuk masing-masing barang yang dipinjam:
                </p>
                <div className="space-y-3">
                  {Array.isArray(modal.loan.loan_details) && modal.loan.loan_details.map((d) => {
                    const code = d.asset_code || d.asset?.code || '-'
                    const name = d.asset_name || d.asset?.label || '-'
                    return (
                      <div key={d.id || d.asset_id} className="flex items-center justify-between gap-4 rounded-lg border border-slate-200 p-3 bg-slate-50/50">
                        <div>
                          <p className="text-sm font-semibold text-slate-800">{code} ({name})</p>
                          <p className="text-xs text-slate-500">Qty: {d.quantity || 1}</p>
                        </div>
                        <div className="w-36">
                          <select
                            value={returnStatuses[d.asset_id] || 'Baik'}
                            onChange={(e) => setReturnStatuses(prev => ({ ...prev, [d.asset_id]: e.target.value }))}
                            className="w-full rounded-lg border border-slate-200 px-3 py-1.5 text-sm text-slate-800 focus:border-blue-500 focus:outline-none bg-white"
                          >
                            <option value="Baik">Baik</option>
                            <option value="Rusak">Rusak</option>
                            <option value="Hilang">Hilang</option>
                          </select>
                        </div>
                      </div>
                    )
                  })}
                </div>
              </div>
            ) : (
              <>
                <div>
                  <label htmlFor="loan-name" className="mb-1 block text-xs font-semibold uppercase tracking-wider text-slate-600">
                    Nama Peminjam
                  </label>
                  <input
                    id="loan-name"
                    name="loans_name"
                    type="text"
                    required
                    value={formData.loans_name}
                    onChange={(e) =>
                      setFormData((prev) => ({ ...prev, loans_name: e.target.value }))
                    }
                    placeholder="cth. Rafli"
                    className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm text-slate-800 focus:border-blue-500 focus:outline-none"
                  />
                </div>

                <div>
                  <label htmlFor="loan-npm" className="mb-1 block text-xs font-semibold uppercase tracking-wider text-slate-600">
                    NPM
                  </label>
                  <input
                    id="loan-npm"
                    name="loans_npm"
                    type="text"
                    required
                    value={formData.loans_npm}
                    onChange={(e) =>
                      setFormData((prev) => ({ ...prev, loans_npm: e.target.value }))
                    }
                    placeholder="cth. 552012"
                    className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm text-slate-800 focus:border-blue-500 focus:outline-none"
                  />
                </div>

                <div>
                  <div className="mb-2 flex items-center justify-between">
                    <label className="text-xs font-semibold uppercase tracking-wider text-slate-600">
                      Daftar Barang (Items)
                    </label>
                    <button
                      type="button"
                      onClick={handleAddItem}
                      className="text-xs font-medium text-blue-600 hover:underline"
                    >
                      + Tambah Barang
                    </button>
                  </div>

                  <div className="space-y-3">
                    {formData.items.map((item, index) => {
                      const availableList = availableAssetsForIndex(index)
                      const hasQuantity = assetHasQuantitySpec(item.asset_code)
                      return (
                        <div key={index} className="flex items-end gap-2 rounded-lg border border-slate-200 p-3 bg-slate-50/50">
                          <div className="flex-1">
                            <label className="mb-1 block text-[11px] font-medium text-slate-500">
                              Pilih Aset #{index + 1}
                            </label>
                            <select
                              required
                              value={item.asset_code}
                              onChange={(e) => handleItemChange(index, 'asset_code', e.target.value)}
                              className="w-full rounded-lg border border-slate-200 px-3 py-1.5 text-sm text-slate-800 focus:border-blue-500 focus:outline-none bg-white"
                            >
                              <option value="">-- Pilih Aset --</option>
                              {availableList.map((a) => (
                                <option key={a.id} value={a.code}>
                                  {a.code} - {a.label} ({a.status})
                                </option>
                              ))}
                            </select>
                          </div>

                          {hasQuantity && (
                            <div className="w-24">
                              <label className="mb-1 block text-[11px] font-medium text-slate-500">
                                Quantity
                              </label>
                              <input
                                type="number"
                                min="1"
                                required
                                value={item.quantity}
                                onChange={(e) => handleItemChange(index, 'quantity', e.target.value)}
                                className="w-full rounded-lg border border-slate-200 px-3 py-1.5 text-sm text-slate-800 focus:border-blue-500 focus:outline-none bg-white"
                              />
                            </div>
                          )}

                          {formData.items.length > 1 && (
                            <button
                              type="button"
                              onClick={() => handleRemoveItem(index)}
                              className="rounded-lg border border-red-200 p-2 text-red-500 transition hover:bg-red-50"
                              title="Hapus item"
                            >
                              <Icon name="trash" className="h-4 w-4" />
                            </button>
                          )}
                        </div>
                      )
                    })}
                  </div>
                </div>
              </>
            )}

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

export default LoanSection