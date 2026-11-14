import { useMemo, useState } from 'react'
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

// Bandingkan dua nilai: angka dibanding angka, sisanya sebagai teks
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
  // sort.key = kolom yang aktif diurutkan, dir = 'asc' | 'desc'
  const [sort, setSort] = useState({ key: null, dir: 'asc' })

  const handleSort = (key) => {
    setSort((prev) => {
      if (prev.key !== key) return { key, dir: 'asc' } // kolom baru → naik
      if (prev.dir === 'asc') return { key, dir: 'desc' } // klik lagi → turun
      return { key: null, dir: 'asc' } // klik ketiga → kembali urutan asli
    })
  }

  const sortedRows = useMemo(() => {
    if (!sort.key) return rows
    const result = [...rows].sort((a, b) =>
      compareValues(a[sort.key], b[sort.key]),
    )
    return sort.dir === 'desc' ? result.reverse() : result
  }, [rows, sort])

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
              `Data dummy — ${rows.length} baris. Klik judul kolom untuk mengurutkan.`}
          </p>
        </div>
      )}

      <div className="overflow-x-auto rounded-xl border border-blue-100">
        <table className="w-full min-w-[640px] border-collapse text-left text-sm">
          <thead>
            <tr className="bg-blue-50/60 text-xs uppercase tracking-wide text-blue-700">
              {columns.map((col) => {
                const active = sort.key === col.key
                const sortable = col.sortable !== false
                if (!sortable) {
                  return (
                    <th key={col.key} className="px-4 py-3 font-semibold">
                      {col.label}
                    </th>
                  )
                }
                return (
                  <th key={col.key} className="px-4 py-3 font-semibold">
                    <button
                      type="button"
                      onClick={() => handleSort(col.key)}
                      title={`Urutkan berdasarkan ${col.label}`}
                      className={`inline-flex items-center gap-1.5 rounded-md px-1.5 py-1 transition hover:bg-blue-100/70 ${
                        active ? 'text-blue-800' : 'text-blue-700'
                      }`}
                    >
                      <span>{col.label}</span>
                      <Icon
                        name={sortIconName(col.key)}
                        className={`h-3.5 w-3.5 ${
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
            {sortedRows.map((row, index) => (
              <tr key={index} className="transition hover:bg-blue-50/40">
                {columns.map((col) => (
                  <td key={col.key} className="px-4 py-3 text-slate-700">
                    {col.key === 'no' ? index + 1 : (col.render ? col.render(row) : <Cell value={row[col.key]} />)}
                  </td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}

export default DataTable
