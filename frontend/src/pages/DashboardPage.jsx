import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import AssetsSection from '../components/AssetsSection'
import CategorySection from '../components/CategorySection'
import ComputerSection from '../components/ComputerSection'
import DashboardOverview from '../components/DashboardOverview'
import LoanSection from '../components/LoanSection'
import SpecificationSection from '../components/SpecificationSection'
import UsersSection from '../components/UsersSection'
import DataTable from '../components/DataTable'
import Icon from '../components/Icons'
import LaboratorySection from '../components/LaboratorySection'
import Logo from '../components/Logo'
import Sidebar from '../components/Sidebar'
import { menuItems, tables } from '../data/dummyData'
import { clearSession, getUser } from '../lib/auth'

// Kunci penyimpanan supaya pilihan menu & status sidebar tidak tereset saat refresh
const ACTIVE_MENU_KEY = 'asmenstore:activeMenu'
const SIDEBAR_KEY = 'asmenstore:sidebarOpen'

function DashboardPage() {
  const navigate = useNavigate()
  const user = getUser()

  // Ambil nilai awal dari localStorage (bertahan saat halaman di-refresh)
  const [activeMenu, setActiveMenu] = useState(
    () => localStorage.getItem(ACTIVE_MENU_KEY) || 'dashboard',
  )
  const [sidebarOpen, setSidebarOpen] = useState(
    () => localStorage.getItem(SIDEBAR_KEY) !== 'false',
  )

  useEffect(() => {
    localStorage.setItem(ACTIVE_MENU_KEY, activeMenu)
  }, [activeMenu])

  useEffect(() => {
    localStorage.setItem(SIDEBAR_KEY, String(sidebarOpen))
  }, [sidebarOpen])

  const fullName = [user?.firstname, user?.lastname].filter(Boolean).join(' ')

  const handleLogout = () => {
    clearSession()
    navigate('/login', { replace: true })
  }

  const renderContent = () => {
    if (activeMenu === 'dashboard') {
      return <DashboardOverview />
    }
    if (activeMenu === 'assets') {
      return <AssetsSection />
    }
    if (activeMenu === 'labkom') {
      return <LaboratorySection />
    }
    if (activeMenu === 'category') {
      return <CategorySection />
    }
    if (activeMenu === 'computer') {
      return <ComputerSection />
    }
    if (activeMenu === 'loan') {
      return <LoanSection />
    }
    if (activeMenu === 'specification') {
      return <SpecificationSection />
    }
    if (activeMenu === 'users') {
      return <UsersSection />
    }
    const table = tables[activeMenu]
    if (!table) return null
    return <DataTable title={table.title} columns={table.columns} rows={table.rows} />
  }

  return (
    <div className="flex min-h-screen bg-white">
      <Sidebar
        items={menuItems}
        activeKey={activeMenu}
        onSelect={setActiveMenu}
        open={sidebarOpen}
        onToggle={() => setSidebarOpen((value) => !value)}
      />

      <div className="flex min-w-0 flex-1 flex-col">
        <header className="sticky top-0 z-10 border-b border-blue-100 bg-white/90 backdrop-blur">
          <div className="flex h-16 items-center justify-between gap-4 px-6">
            <div className="flex items-center gap-2 text-sm text-slate-500">
              <Logo size={24} className="lg:hidden" />
              <span className="font-medium text-slate-700">AsmenStore</span>
              <span className="text-blue-300">/</span>
              <span className="capitalize">{activeMenu}</span>
            </div>

            <div className="flex items-center gap-4">
              <div className="text-right">
                <p className="text-sm font-medium text-slate-800">
                  {fullName || user?.username || 'Pengguna'}
                </p>
                <p className="text-xs text-slate-500">
                  {user?.role?.name ?? 'Tanpa Role'}
                </p>
              </div>
              <button
                type="button"
                onClick={handleLogout}
                className="flex items-center gap-2 rounded-lg border border-blue-200 px-3 py-1.5 text-sm font-medium text-blue-700 transition hover:bg-blue-50"
              >
                <Icon name="logout" className="h-4 w-4" />
                Keluar
              </button>
            </div>
          </div>
        </header>

        <main className="flex-1 px-6 py-8">{renderContent()}</main>
      </div>
    </div>
  )
}

export default DashboardPage
