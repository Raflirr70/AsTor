import { useMemo, useState } from 'react'
import { createAsset, updateAsset } from '../lib/api'
import Modal from './Modal'

const inputClass =
  'w-full rounded-lg border border-blue-100 bg-blue-50/40 px-3.5 py-2.5 text-sm text-slate-800 placeholder:text-slate-400 transition focus:border-blue-500 focus:bg-white focus:outline-none focus:ring-2 focus:ring-blue-200'

const STATUS_OPTIONS = ['Baik', 'Rusak']

// Cocokkan status lama (bisa beda kapitalisasi) ke salah satu pilihan
function normalizeStatus(value) {
  const match = STATUS_OPTIONS.find(
    (option) => option.toLowerCase() === String(value ?? '').toLowerCase(),
  )
  return match ?? STATUS_OPTIONS[0]
}

function Field({ label, htmlFor, error, children, hint }) {
  return (
    <div>
      <label htmlFor={htmlFor} className="mb-1.5 block text-sm font-medium text-slate-700">
        {label}
      </label>
      {children}
      {hint && !error && <p className="mt-1 text-xs text-slate-400">{hint}</p>}
      {error && <p className="mt-1 text-xs text-red-500">{error}</p>}
    </div>
  )
}

function AssetFormModal({ mode, asset, categories, onClose, onSuccess }) {
  const isCreate = mode === 'create'

  const [categoryId, setCategoryId] = useState(
    isCreate ? '' : String(asset?.category_id ?? ''),
  )
  const [label, setLabel] = useState(asset?.label ?? '')
  const [status, setStatus] = useState(normalizeStatus(asset?.status))
  const [specValues, setSpecValues] = useState({})
  const [fieldErrors, setFieldErrors] = useState({})
  const [formError, setFormError] = useState('')
  const [submitting, setSubmitting] = useState(false)

  const selectedCategory = useMemo(
    () => categories.find((c) => c.id === Number(categoryId)),
    [categories, categoryId],
  )

  // Spesifikasi yang wajib diisi mengikuti kategori yang dipilih
  const specs = useMemo(
    () => (isCreate ? selectedCategory?.specifications ?? [] : []),
    [isCreate, selectedCategory],
  )

  // Reset nilai spesifikasi tiap kali kategori berganti
  const handleCategoryChange = (value) => {
    setCategoryId(value)
    setSpecValues({})
    setFieldErrors((prev) => ({ ...prev, category: undefined }))
  }

  const handleSubmit = async (event) => {
    event.preventDefault()
    setFormError('')

    const errors = {}
    if (isCreate && !categoryId) errors.category = 'Kategori wajib dipilih.'
    if (!label.trim()) errors.label = 'Label wajib diisi.'
    if (!status.trim()) errors.status = 'Status wajib diisi.'
    if (isCreate) {
      for (const spec of specs) {
        if (!String(specValues[spec.id] ?? '').trim()) {
          errors[`spec-${spec.id}`] = `${spec.name} wajib diisi.`
        }
      }
    }

    if (Object.keys(errors).length > 0) {
      setFieldErrors(errors)
      return
    }

    setSubmitting(true)
    try {
      if (isCreate) {
        await createAsset({
          category_id: Number(categoryId),
          label: label.trim(),
          status: status.trim(),
          asset_specifications: specs.map((spec) => ({
            specification_category_id: spec.id,
            value: String(specValues[spec.id]).trim(),
          })),
        })
        onSuccess('Asset berhasil ditambahkan.')
      } else {
        await updateAsset(asset.id, {
          label: label.trim(),
          status: status.trim(),
        })
        onSuccess('Asset berhasil diperbarui.')
      }
    } catch (err) {
      setFormError(err.message || 'Gagal menyimpan asset.')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <Modal
      title={isCreate ? 'Tambah Asset' : `Update Asset — ${asset?.code ?? ''}`}
      onClose={onClose}
    >
      <form onSubmit={handleSubmit} className="space-y-4">
        <Field label="Kategori" htmlFor="category_id" error={fieldErrors.category}>
          {isCreate ? (
            <select
              id="category_id"
              value={categoryId}
              onChange={(event) => handleCategoryChange(event.target.value)}
              className={inputClass}
            >
              <option value="">— Pilih kategori —</option>
              {categories.map((category) => (
                <option key={category.id} value={category.id}>
                  {category.name}
                </option>
              ))}
            </select>
          ) : (
            <input
              id="category_id"
              type="text"
              value={selectedCategory?.name ?? `#${asset?.category_id}`}
              readOnly
              className={`${inputClass} cursor-not-allowed text-slate-500`}
            />
          )}
        </Field>

        <Field label="Label" htmlFor="label" error={fieldErrors.label}>
          <input
            id="label"
            type="text"
            value={label}
            onChange={(event) => setLabel(event.target.value)}
            placeholder="mis. Monitor tak 2"
            className={inputClass}
          />
        </Field>

        <Field label="Status" htmlFor="status" error={fieldErrors.status}>
          <select
            id="status"
            value={status}
            onChange={(event) => setStatus(event.target.value)}
            className={inputClass}
          >
            {STATUS_OPTIONS.map((option) => (
              <option key={option} value={option}>
                {option}
              </option>
            ))}
          </select>
        </Field>

        {/* Spesifikasi dinamis (hanya saat create) */}
        {isCreate && categoryId && (
          <div className="rounded-xl border border-blue-100 bg-blue-50/30 p-4">
            <p className="mb-3 text-sm font-semibold text-slate-700">
              Spesifikasi
            </p>
            {specs.length === 0 ? (
              <p className="text-xs text-slate-500">
                Kategori ini tidak memiliki spesifikasi.
              </p>
            ) : (
              <div className="space-y-3">
                {specs.map((spec) => (
                  <Field
                    key={spec.id}
                    label={`${spec.name}${spec.unit ? ` (${spec.unit})` : ''}`}
                    htmlFor={`spec-${spec.id}`}
                    error={fieldErrors[`spec-${spec.id}`]}
                  >
                    <input
                      id={`spec-${spec.id}`}
                      type="text"
                      value={specValues[spec.id] ?? ''}
                      onChange={(event) =>
                        setSpecValues((prev) => ({
                          ...prev,
                          [spec.id]: event.target.value,
                        }))
                      }
                      placeholder={`Masukkan ${spec.name}`}
                      className={inputClass}
                    />
                  </Field>
                ))}
              </div>
            )}
          </div>
        )}

        {/* Spesifikasi yang sudah ada (hanya saat update, read-only) */}
        {!isCreate && asset?.asset_specifications?.length > 0 && (
          <div className="rounded-xl border border-blue-100 bg-blue-50/30 p-4">
            <p className="mb-2 text-sm font-semibold text-slate-700">
              Spesifikasi saat ini
            </p>
            <ul className="space-y-1 text-xs text-slate-500">
              {asset.asset_specifications.map((spec) => (
                <li key={spec.id}>
                  {[spec.name, spec.value, spec.unit].filter(Boolean).join(' ')}
                </li>
              ))}
            </ul>
          </div>
        )}

        {formError && (
          <p className="rounded-lg bg-red-50 px-4 py-2.5 text-sm text-red-600">
            {formError}
          </p>
        )}

        <div className="flex justify-end gap-2 pt-2">
          <button
            type="button"
            onClick={onClose}
            className="rounded-lg border border-blue-200 px-4 py-2 text-sm font-medium text-slate-600 transition hover:bg-blue-50"
          >
            Batal
          </button>
          <button
            type="submit"
            disabled={submitting}
            className="rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white transition hover:bg-blue-700 disabled:cursor-not-allowed disabled:opacity-60"
          >
            {submitting ? 'Menyimpan…' : 'Simpan'}
          </button>
        </div>
      </form>
    </Modal>
  )
}

export default AssetFormModal
