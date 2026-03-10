import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { Nav } from './components/Nav'
import { DashboardPage } from './pages/DashboardPage'
import { ListPage } from './pages/ListPage'
import { EntityDetailPage } from './pages/EntityDetailPage'
import { TreePage } from './pages/TreePage'
import { SearchPage } from './pages/SearchPage'

export default function App() {
  return (
    <BrowserRouter>
      <div style={{ display: 'flex', minHeight: '100vh' }}>
        <Nav />
        <main style={{ flex: 1, padding: '24px 32px', backgroundColor: '#f8fafc' }}>
          <Routes>
            <Route path="/" element={<DashboardPage />} />
            <Route path="/list" element={<ListPage />} />
            <Route path="/entity/:id" element={<EntityDetailPage />} />
            <Route path="/tree" element={<TreePage />} />
            <Route path="/search" element={<SearchPage />} />
          </Routes>
        </main>
      </div>
    </BrowserRouter>
  )
}
