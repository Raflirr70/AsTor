// ============================================================
// Data dummy — dipakai sementara sampai tersambung ke API backend
// ============================================================

export const menuItems = [
  { key: 'dashboard', label: 'Dashboard', icon: 'dashboard' },
  { key: 'assets', label: 'Assets', icon: 'assets' },
  { key: 'labkom', label: 'Laboratorium Computer', icon: 'labkom' },
  { key: 'users', label: 'Users', icon: 'users' },
  { key: 'category', label: 'Category', icon: 'category' },
  { key: 'computer', label: 'Computer', icon: 'computer' },
  { key: 'loan', label: 'Loan', icon: 'loans' },
  { key: 'report', label: 'Report', icon: 'report' },
  { key: 'specification', label: 'Specification', icon: 'specification' },
]

export const overviewStats = [
  { key: 'assets', label: 'Jumlah Asset', value: 128, icon: 'assets' },
  { key: 'reports', label: 'Jumlah Report', value: 24, icon: 'report' },
  { key: 'loans', label: 'Jumlah Peminjaman', value: 57, icon: 'loans' },
]

export const damageByLabkom = [
  { label: 'Labkom 1', value: 3 },
  { label: 'Labkom 2', value: 5 },
  { label: 'Labkom 3', value: 2 },
]

export const notifications = [
  { id: 1, title: 'Peminjaman disetujui', message: 'Proyektor Epson disetujui oleh PJ Gudang', time: '5 menit lalu', read: false },
  { id: 2, title: 'Laporan kerusakan baru', message: 'Monitor LG 24" dilaporkan rusak di Labkom 2', time: '1 jam lalu', read: false },
  { id: 3, title: 'Batas pengembalian', message: 'Peminjaman Labkom 1 jatuh tempo besok', time: '3 jam lalu', read: true },
  { id: 4, title: 'Aset dipindahkan', message: 'Router Mikrotik dipindah ke Labkom 3', time: 'Kemarin', read: true },
  { id: 5, title: 'Stok menipis', message: 'Kabel HDMI tersisa 4 unit di Gudang Utama', time: '2 hari lalu', read: true },
]

export const userActivities = [
  { id: 1, user: 'Admin', action: 'Menambahkan asset AST-006 (Kabel HDMI 5m)', time: '10 menit lalu' },
  { id: 2, user: 'PJ Gudang', action: 'Menyetujui peminjaman #LN-014', time: '40 menit lalu' },
  { id: 3, user: 'Asisten', action: 'Membuat laporan kerusakan #RP-021', time: '2 jam lalu' },
  { id: 4, user: 'PJ Labkom', action: 'Memperbarui spesifikasi RAM DDR4', time: '5 jam lalu' },
  { id: 5, user: 'Admin', action: 'Menambahkan user baru (Budi Santoso)', time: 'Kemarin' },
]

export const assetHistories = [
  { id: 1, asset: 'AST-003 — Monitor LG 24"', action: 'Dipasang di Labkom 1', time: '1 jam lalu' },
  { id: 2, asset: 'AST-002 — Proyektor Epson', action: 'Dipinjam oleh Rina', time: '4 jam lalu' },
  { id: 3, asset: 'AST-004 — Printer HP LaserJet', action: 'Ditandai rusak', time: 'Kemarin' },
  { id: 4, asset: 'AST-001 — Laptop Dell Latitude', action: 'Dikembalikan ke Gudang Utama', time: '2 hari lalu' },
  { id: 5, asset: 'AST-005 — Router Mikrotik', action: 'Ditambahkan ke sistem', time: '3 hari lalu' },
]

export const assetMovements = [
  { id: 1, asset: 'AST-005 — Router Mikrotik', from: 'Labkom 2', to: 'Labkom 3', time: '30 menit lalu' },
  { id: 2, asset: 'AST-003 — Monitor LG 24"', from: 'Gudang Utama', to: 'Labkom 1', time: '3 jam lalu' },
  { id: 3, asset: 'AST-006 — Kabel HDMI 5m', from: 'Labkom 1', to: 'Gudang Utama', time: 'Kemarin' },
  { id: 4, asset: 'AST-001 — Laptop Dell Latitude', from: 'Labkom 3', to: 'Gudang Utama', time: '2 hari lalu' },
]

// Tabel untuk tiap menu di sidebar (Assets diambil dari API, lihat AssetsSection.jsx)
export const tables = {
  users: {
    title: 'Users',
    columns: [
      { key: 'name', label: 'Nama' },
      { key: 'username', label: 'Username' },
      { key: 'npm', label: 'NPM' },
      { key: 'role', label: 'Role' },
    ],
    rows: [
      { name: 'Andi Pratama', username: 'admin', npm: '2211102001', role: 'Admin' },
      { name: 'Budi Santoso', username: 'budi', npm: '2211102002', role: 'PJ Gudang' },
      { name: 'Citra Dewi', username: 'citra', npm: '2211102003', role: 'PJ Labkom' },
      { name: 'Rina Melati', username: 'rina', npm: '2211102004', role: 'Asisten' },
      { name: 'Dodi Hermawan', username: 'dodi', npm: '2211102005', role: 'Asisten' },
    ],
  },
  category: {
    title: 'Category',
    columns: [],
    rows: [],
  },
  computer: {
    title: 'Computer',
    columns: [],
    rows: [],
  },
  report: {
    title: 'Report',
    columns: [
      { key: 'code', label: 'Kode' },
      { key: 'asset', label: 'Asset' },
      { key: 'reporter', label: 'Pelapor' },
      { key: 'description', label: 'Deskripsi' },
      { key: 'status', label: 'Status' },
    ],
    rows: [
      { code: 'RP-021', asset: 'Monitor LG 24"', reporter: 'Rina Melati', description: 'Layar berkedip', status: 'Menunggu' },
      { code: 'RP-020', asset: 'Printer HP LaserJet', reporter: 'Dodi Hermawan', description: 'Tidak bisa mencetak', status: 'Diproses' },
      { code: 'RP-019', asset: 'PC-L2-01', reporter: 'Citra Dewi', description: 'Tidak menyala', status: 'Selesai' },
      { code: 'RP-018', asset: 'Keyboard Logitech', reporter: 'Budi Santoso', description: 'Beberapa tombol rusak', status: 'Selesai' },
    ],
  },
  specification: {
    title: 'Specification',
    columns: [
      { key: 'code', label: 'Kode' },
      { key: 'name', label: 'Nama' },
      { key: 'category', label: 'Kategori' },
      { key: 'unit', label: 'Satuan' },
      { key: 'dataType', label: 'Tipe Data' },
    ],
    rows: [
      { code: 'SP-01', name: 'Ukuran Layar', category: 'Monitor', unit: 'inch', dataType: 'Angka' },
      { code: 'SP-02', name: 'Kapasitas RAM', category: 'RAM', unit: 'GB', dataType: 'Angka' },
      { code: 'SP-03', name: 'Kapasitas SSD', category: 'SSD', unit: 'GB', dataType: 'Angka' },
      { code: 'SP-04', name: 'Tipe Keyboard', category: 'Keyboard', unit: '-', dataType: 'Teks' },
      { code: 'SP-05', name: 'Jumlah Tombol', category: 'Mouse', unit: 'tombol', dataType: 'Angka' },
    ],
  },
}
