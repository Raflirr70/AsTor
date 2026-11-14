import { useMemo, useState, useEffect } from 'react'
import Icon from './Icons'

const toneByStatus = {
  baik: 'emerald',
  tersedia: 'emerald',
  aktif: 'emerald',
  selesai: 'emerald',
  available: 'emerald',
  good: 'emerald',
  active: 'emerald',
  terpasang: 'blue',
  diproses: 'blue',
  in_use: 'blue',
  used: 'blue',
  dipinjam: 'amber',
  menunggu: 'amber',
  maintenance: 'amber',
  pending: 'amber',
  borrowed: 'amber',
  rusak: 'red',
  'perlu perbaikan': 'red',
  broken: 'red',
  damaged: 'red',
}

const toneClass = {
  emerald: 'bg-emerald-50 text-emerald-600',
  blue: 'bg-blue-50 text-blue-600',
  amber: 'bg-amber-50 text-amber-600',
  red: 'bg-red-50 text-red-600',
}

function Cell({ value }) {
  const tone =
    typeof value === 'string' ? toneClass[toneByStatus[value.toLowerCase()]] : null
  if (tone) {
    return (
      <span className={`inline-block rounded-full px-2.5 py-0.5 text-xs font-medium ${tone}`}>
        {value}
      </span>
    )
  }
  return value
}

function compareValues(a, b) {
  if (typeof a === 'number' && typeof b === 'number') {
    return a - b
  }
  return String(a).localeCompare(String(b), 'id', {
    numeric: true,
    sensitivity: 'base',
  })
}

function DataTable({ title, subtitle, columns, rows, showHeader = true }) {
  const [sort, setSort] = useState({ key: null, dir: 'asc' })
  const [currentPage, setCurrentPage] = useState(() => {
    const saved = localStorage.getItem(`dt_page_${title}`)
    return saved ? Number(saved) : 1
  })
  const pageSize = 10

  useEffect(() => {
    localStorage.setItem(`dt_page_${title}`, currentPage)
  }, [currentPage, title])

  const handleSort = (key) => {
    setSort((prev) => {
      if (prev.key !== key) return { key, dir: 'asc' }
      if (prev.dir === 'asc') return { key, dir: 'desc' }
      return { key: null, dir: 'asc' }
    })
  }

  const sortedRows = useMemo(() => {
    if (!sort.key) return rows
    const result = [...rows].sort((a, b) =>
      compareValues(a[sort.key], b[sort.key]),
    )
    return sort.dir === 'desc' ? result.reverse() : result
  }, [rows, sort])

  const totalPages = Math.ceil(sortedRows.length / pageSize) || 1

  useEffect(() => {
    if (currentPage > totalPages) {
      setCurrentPage(1)
    }
  }, [totalPages, currentPage])

  const paginatedRows = useMemo(() => {
    const start = (currentPage - 1) * pageSize
    return sortedRows.slice(start, start + pageSize)
  }, [sortedRows, currentPage])

  const sortIconName = (key) => {
    if (sort.key !== key) return 'sortNeutral'
    return sort.dir === 'asc' ? 'sortAsc' : 'sortDesc'
  }

  return (
    <div>
      {showHeader && (
        <div className="mb-6">
          <h1 className="text-2xl font-bold text-slate-800">{title}</h1>
          <p className="mt-1 text-sm text-slate-500">
            {subtitle ??
              `Data — ${rows.length} baris. Klik judul kolom untuk mengurutkan.`}
          </p>
        </div>
      )}

      <div className="overflow-x-auto rounded-xl border border-blue-100">
        <table className="w-full min-w-[640px] border-collapse text-left text-xs">
          <thead>
            <tr className="bg-blue-50/60 uppercase tracking-wide text-blue-700">
              {columns.map((col) => {
                const active = sort.key === col.key
                const sortable = col.sortable !== false
                if (!sortable) {
                  return (
                    <th key={col.key} className="px-3 py-2.5 font-semibold">
                      {col.label}
                    </th>
                  )
                }
                return (
                  <th key={col.key} className="px-3 py-2.5 font-semibold">
                    <button
                      type="button"
                      onClick={() => handleSort(col.key)}
                      title={`Urutkan berdasarkan ${col.label}`}
                      className={`inline-flex items-center gap-1 rounded-md px-1 py-0.5 transition hover:bg-blue-100/70 ${
                        active ? 'text-blue-800' : 'text-blue-700'
                      }`}
                    >
                      <span>{col.label}</span>
                      <Icon
                        name={sortIconName(col.key)}
                        className={`h-3 w-3 ${
                          active ? 'text-blue-700' : 'text-blue-300'
                        }`}
                      />
                    </button>
                  </th>
                )
              })}
            </tr>
          </thead>
          <tbody className="divide-y divide-blue-50">
            {paginatedRows.length === 0 ? (
              <tr>
                <td colSpan={columns.length} className="px-4 py-6 text-center text-slate-400">
                  Tidak ada data.
                </td>
              </tr>
            ) : (
              paginatedRows.map((row, index) => {
                const absoluteIndex = (currentPage - 1) * pageSize + index
                return (
                  <tr key={absoluteIndex} className="transition hover:bg-blue-50/40">
                    {columns.map((col) => (
                      <td key={col.key} className="px-3 py-2 text-slate-700 whitespace-nowrap">
                        {col.key === 'no' ? absoluteIndex + 1 : (col.render ? col.render(row) : <Cell value={row[col.key]} />)}
                      </td>
                    ))}
                  </tr>
                )
              })
            )}
          </tbody>
        </table>
      </div>

      {totalPages > 1 && (
        <div className="mt-4 flex items-center justify-between text-xs text-slate-500">
          <div>
            Menampilkan {(currentPage - 1) * pageSize + 1} - {Math.min(currentPage * pageSize, sortedRows.length)} dari {sortedRows.length} data
          </div>
          <div className="flex items-center gap-1">
            <button
              type="button"
              disabled={currentPage === 1}
              onClick={() => setCurrentPage((p) => Math.max(p - 1, 1))}
              className="rounded-lg border border-slate-200 px-3 py-1 font-medium transition hover:bg-slate-50 disabled:opacity-40"
            >
              Sebelumnya
            </button>
            <span className="px-2 py-1 font-medium text-slate-700">
              {currentPage} / {totalPages}
            </span>
            <button
              type="button"
              disabled={currentPage === totalPages}
              onClick={() => setCurrentPage((p) => Math.min(p + 1, totalPages))}
              className="rounded-lg border border-slate-200 px-3 py-1 font-medium transition hover:bg-slate-50 disabled:opacity-40"
            >
              Selanjutnya
            </button>
          </div>
        </div>
      )}
    </div>
  )
}

export default DataTable