import Icon from './Icons'
import Logo from './Logo'

/**
 * Sidebar kiri dengan tombol panah sebagai toggle buka/tutup.
 * Saat tertutup hanya menampilkan ikon (rail).
 */
function Sidebar({ items, activeKey, onSelect, open, onToggle }) {
  return (
    <aside
      className={`sticky top-0 flex h-screen shrink-0 flex-col border-r border-blue-100 bg-white transition-[width] duration-300 ${
        open ? 'w-64' : 'w-[76px]'
      }`}
    >
      {/* Header + tombol toggle */}
      <div className="flex h-16 items-center border-b border-blue-100 px-3">
        {open && (
          <div className="flex min-w-0 flex-1 items-center gap-2">
            <Logo size={30} />
            <span className="truncate text-base font-semibold text-blue-700">
              AsmenStore
            </span>
          </div>
        )}
        <button
          type="button"
          onClick={onToggle}
          aria-label={open ? 'Tutup sidebar' : 'Buka sidebar'}
          title={open ? 'Tutup sidebar' : 'Buka sidebar'}
          className={`rounded-lg p-2 text-slate-500 transition hover:bg-blue-50 hover:text-blue-700 ${
            open ? 'ml-auto' : 'mx-auto'
          }`}
        >
          <Icon name={open ? 'chevronLeft' : 'chevronRight'} className="h-5 w-5" />
        </button>
      </div>

      {/* Menu */}
      <nav className="flex-1 space-y-1 overflow-y-auto px-3 py-4">
        {items.map((item) => {
          const active = item.key === activeKey
          return (
            <button
              key={item.key}
              type="button"
              onClick={() => onSelect(item.key)}
              title={item.label}
              className={`flex w-full items-center gap-3 rounded-lg py-2.5 text-sm font-medium transition ${
                open ? 'px-3' : 'justify-center px-0'
              } ${
                active
                  ? 'bg-blue-600 text-white shadow-sm shadow-blue-200'
                  : 'text-slate-600 hover:bg-blue-50 hover:text-blue-700'
              }`}
            >
              <Icon name={item.icon} className="h-5 w-5 shrink-0" />
              {open && <span className="truncate">{item.label}</span>}
            </button>
          )
        })}
      </nav>
    </aside>
  )
}

export default Sidebar
