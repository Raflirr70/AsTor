import {
  assetHistories,
  assetMovements,
  damageByLabkom,
  notifications,
  overviewStats,
  userActivities,
} from '../data/dummyData'
import Icon from './Icons'

function StatCard({ label, value, icon }) {
  return (
    <div className="rounded-xl border border-blue-100 bg-blue-50/40 p-5">
      <div className="flex items-center justify-between">
        <p className="text-sm text-slate-500">{label}</p>
        <span className="rounded-lg bg-white p-2 text-blue-600 shadow-sm">
          <Icon name={icon} className="h-5 w-5" />
        </span>
      </div>
      <p className="mt-3 text-3xl font-bold text-blue-700">{value}</p>
    </div>
  )
}

function DamageCard({ label, value }) {
  return (
    <div className="rounded-xl border border-red-100 bg-red-50/40 p-5">
      <p className="text-sm text-slate-500">Laporan Kerusakan {label}</p>
      <p className="mt-2 text-3xl font-bold text-red-500">{value}</p>
    </div>
  )
}

function HistoryList({ title, icon, items }) {
  return (
    <section className="flex flex-col rounded-xl border border-blue-100 bg-white">
      <header className="flex items-center gap-2 border-b border-blue-100 px-4 py-3">
        <span className="text-blue-600">
          <Icon name={icon} className="h-5 w-5" />
        </span>
        <h2 className="text-sm font-semibold text-slate-700">{title}</h2>
      </header>
      <ul className="divide-y divide-blue-50">
        {items.map((item) => (
          <li key={item.id} className="flex items-start gap-3 px-4 py-3">
            <span
              className={`mt-1.5 h-2 w-2 shrink-0 rounded-full ${
                item.unread ? 'bg-blue-500' : 'bg-blue-100'
              }`}
            />
            <div className="min-w-0 flex-1">
              <p className="truncate text-sm font-medium text-slate-700">
                {item.primary}
              </p>
              <p className="truncate text-xs text-slate-500">{item.secondary}</p>
            </div>
            <span className="shrink-0 text-xs text-slate-400">{item.time}</span>
          </li>
        ))}
      </ul>
    </section>
  )
}

function DashboardOverview() {
  const notificationItems = notifications.map((n) => ({
    id: n.id,
    primary: n.title,
    secondary: n.message,
    time: n.time,
    unread: !n.read,
  }))

  const activityItems = userActivities.map((a) => ({
    id: a.id,
    primary: a.user,
    secondary: a.action,
    time: a.time,
  }))

  const assetHistoryItems = assetHistories.map((h) => ({
    id: h.id,
    primary: h.asset,
    secondary: h.action,
    time: h.time,
  }))

  const movementItems = assetMovements.map((m) => ({
    id: m.id,
    primary: m.asset,
    secondary: `${m.from} → ${m.to}`,
    time: m.time,
  }))

  return (
    <div>
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-slate-800">Dashboard</h1>
        <p className="mt-1 text-sm text-slate-500">
          Ringkasan inventaris gudang laboratorium komputer.
        </p>
      </div>

      {/* Statistik utama */}
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {overviewStats.map((stat) => (
          <StatCard key={stat.key} {...stat} />
        ))}
      </div>

      {/* Laporan kerusakan per labkom */}
      <h2 className="mt-8 mb-3 text-sm font-semibold text-slate-600">
        Laporan Kerusakan Labkom
      </h2>
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {damageByLabkom.map((item) => (
          <DamageCard key={item.label} {...item} />
        ))}
      </div>

      {/* Riwayat */}
      <h2 className="mt-8 mb-3 text-sm font-semibold text-slate-600">
        Riwayat
      </h2>
      <div className="grid gap-4 lg:grid-cols-2">
        <HistoryList
          title="Riwayat Notifikasi"
          icon="bell"
          items={notificationItems}
        />
        <HistoryList
          title="Riwayat Aktivitas User"
          icon="activity"
          items={activityItems}
        />
        <HistoryList
          title="Riwayat Asset"
          icon="history"
          items={assetHistoryItems}
        />
        <HistoryList
          title="Riwayat Pemindahan Asset"
          icon="move"
          items={movementItems}
        />
      </div>
    </div>
  )
}

export default DashboardOverview
