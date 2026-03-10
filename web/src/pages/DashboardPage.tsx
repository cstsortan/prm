import { useEffect, useState } from 'react'
import { getDashboard } from '../api/client'
import type { DashboardResponse, EntityType, Status, Priority } from '../api/types'

const typeOrder: EntityType[] = ['epic', 'story', 'sub-task', 'task', 'bug']
const statusOrder: Status[] = ['in-progress', 'review', 'todo', 'backlog', 'done', 'cancelled', 'archived']
const priorityOrder: Priority[] = ['critical', 'high', 'medium', 'low']

const statusColors: Record<string, string> = {
  'in-progress': '#f59e0b',
  'review': '#8b5cf6',
  'todo': '#3b82f6',
  'backlog': '#6b7280',
  'done': '#10b981',
  'cancelled': '#ef4444',
  'archived': '#9ca3af',
}

const priorityColors: Record<string, string> = {
  critical: '#991b1b',
  high: '#dc2626',
  medium: '#f59e0b',
  low: '#3b82f6',
}

function StatCard({ label, value, color }: { label: string; value: number; color?: string }) {
  return (
    <div
      style={{
        backgroundColor: '#fff',
        borderRadius: 8,
        padding: '16px 20px',
        boxShadow: '0 1px 3px rgba(0,0,0,0.1)',
        borderLeft: color ? `4px solid ${color}` : undefined,
      }}
    >
      <div style={{ fontSize: 28, fontWeight: 700, color: color || '#1e293b' }}>{value}</div>
      <div style={{ fontSize: 13, color: '#64748b', marginTop: 2 }}>{label}</div>
    </div>
  )
}

function BarChart({ data, colors }: { data: [string, number][]; colors: Record<string, string> }) {
  const max = Math.max(...data.map(([, v]) => v), 1)
  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 6 }}>
      {data.map(([label, value]) => (
        <div key={label} style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
          <div style={{ width: 90, fontSize: 13, color: '#64748b', textAlign: 'right' }}>{label}</div>
          <div style={{ flex: 1, height: 22, backgroundColor: '#f1f5f9', borderRadius: 4, overflow: 'hidden' }}>
            <div
              style={{
                width: `${(value / max) * 100}%`,
                height: '100%',
                backgroundColor: colors[label] || '#94a3b8',
                borderRadius: 4,
                transition: 'width 0.3s ease',
                minWidth: value > 0 ? 2 : 0,
              }}
            />
          </div>
          <div style={{ width: 30, fontSize: 13, fontWeight: 600 }}>{value}</div>
        </div>
      ))}
    </div>
  )
}

export function DashboardPage() {
  const [data, setData] = useState<DashboardResponse | null>(null)
  const [error, setError] = useState('')

  useEffect(() => {
    getDashboard().then(setData).catch((e) => setError(e.message))
  }, [])

  if (error) return <div style={{ color: '#ef4444', padding: 20 }}>Error: {error}</div>
  if (!data) return <div style={{ padding: 20, color: '#64748b' }}>Loading...</div>

  const { stats } = data

  const statusData: [string, number][] = statusOrder
    .filter((s) => stats.ByStatus[s])
    .map((s) => [s, stats.ByStatus[s]])

  const priorityData: [string, number][] = priorityOrder
    .filter((p) => stats.ByPriority[p])
    .map((p) => [p, stats.ByPriority[p]])

  return (
    <div>
      <h1 style={{ fontSize: 24, fontWeight: 700, marginBottom: 20 }}>{data.project}</h1>

      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(140px, 1fr))', gap: 12, marginBottom: 28 }}>
        <StatCard label="Total Items" value={stats.Total} />
        {typeOrder.map((t) =>
          stats.ByType[t] ? <StatCard key={t} label={t} value={stats.ByType[t]} /> : null
        )}
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 24 }}>
        <div style={{ backgroundColor: '#fff', borderRadius: 8, padding: 20, boxShadow: '0 1px 3px rgba(0,0,0,0.1)' }}>
          <h2 style={{ fontSize: 16, fontWeight: 600, marginBottom: 16 }}>By Status</h2>
          <BarChart data={statusData} colors={statusColors} />
        </div>
        <div style={{ backgroundColor: '#fff', borderRadius: 8, padding: 20, boxShadow: '0 1px 3px rgba(0,0,0,0.1)' }}>
          <h2 style={{ fontSize: 16, fontWeight: 600, marginBottom: 16 }}>By Priority</h2>
          <BarChart data={priorityData} colors={priorityColors} />
        </div>
      </div>
    </div>
  )
}
