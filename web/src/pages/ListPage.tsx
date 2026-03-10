import { useEffect, useState } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { listEntities, updateStatus } from '../api/client'
import type { Entity, EntityType, Status, Priority } from '../api/types'
import { PriorityBadge } from '../components/PriorityBadge'
import { TypeIcon } from '../components/TypeIcon'

const allTypes: EntityType[] = ['epic', 'story', 'sub-task', 'task', 'bug']
const allStatuses: Status[] = ['backlog', 'todo', 'in-progress', 'review', 'done', 'cancelled']
const allPriorities: Priority[] = ['critical', 'high', 'medium', 'low']

function FilterBar({
  types, statuses, priorities,
  onTypesChange, onStatusesChange, onPrioritiesChange,
}: {
  types: string[]; statuses: string[]; priorities: string[]
  onTypesChange: (v: string[]) => void
  onStatusesChange: (v: string[]) => void
  onPrioritiesChange: (v: string[]) => void
}) {
  const toggle = (list: string[], val: string, setter: (v: string[]) => void) => {
    setter(list.includes(val) ? list.filter((v) => v !== val) : [...list, val])
  }

  const chipStyle = (active: boolean): React.CSSProperties => ({
    padding: '3px 10px',
    borderRadius: 4,
    border: '1px solid #e2e8f0',
    backgroundColor: active ? '#334155' : '#fff',
    color: active ? '#fff' : '#475569',
    cursor: 'pointer',
    fontSize: 12,
    fontWeight: 500,
    transition: 'all 0.15s',
  })

  return (
    <div style={{ display: 'flex', gap: 16, flexWrap: 'wrap', marginBottom: 16 }}>
      <div style={{ display: 'flex', gap: 4, alignItems: 'center' }}>
        <span style={{ fontSize: 12, color: '#94a3b8', marginRight: 4 }}>Type:</span>
        {allTypes.map((t) => (
          <button key={t} style={chipStyle(types.includes(t))} onClick={() => toggle(types, t, onTypesChange)}>{t}</button>
        ))}
      </div>
      <div style={{ display: 'flex', gap: 4, alignItems: 'center' }}>
        <span style={{ fontSize: 12, color: '#94a3b8', marginRight: 4 }}>Status:</span>
        {allStatuses.map((s) => (
          <button key={s} style={chipStyle(statuses.includes(s))} onClick={() => toggle(statuses, s, onStatusesChange)}>{s}</button>
        ))}
      </div>
      <div style={{ display: 'flex', gap: 4, alignItems: 'center' }}>
        <span style={{ fontSize: 12, color: '#94a3b8', marginRight: 4 }}>Priority:</span>
        {allPriorities.map((p) => (
          <button key={p} style={chipStyle(priorities.includes(p))} onClick={() => toggle(priorities, p, onPrioritiesChange)}>{p}</button>
        ))}
      </div>
    </div>
  )
}

export function ListPage() {
  const [entities, setEntities] = useState<Entity[]>([])
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(true)
  const [searchParams, setSearchParams] = useSearchParams()
  const navigate = useNavigate()

  const types = searchParams.get('type')?.split(',').filter(Boolean) || []
  const statuses = searchParams.get('status')?.split(',').filter(Boolean) || []
  const priorities = searchParams.get('priority')?.split(',').filter(Boolean) || []

  const updateParams = (key: string, values: string[]) => {
    const next = new URLSearchParams(searchParams)
    if (values.length > 0) {
      next.set(key, values.join(','))
    } else {
      next.delete(key)
    }
    setSearchParams(next)
  }

  useEffect(() => {
    setLoading(true)
    listEntities({
      type: types.join(',') || undefined,
      status: statuses.join(',') || undefined,
      priority: priorities.join(',') || undefined,
    })
      .then(setEntities)
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false))
  }, [searchParams.toString()])

  const handleStatusChange = async (id: string, newStatus: string) => {
    try {
      const updated = await updateStatus(id, newStatus)
      setEntities((prev) => prev.map((e) => (e.id === id ? updated : e)))
    } catch (e: any) {
      setError(e.message)
    }
  }

  if (error) return <div style={{ color: '#ef4444', padding: 20 }}>Error: {error}</div>

  return (
    <div>
      <h1 style={{ fontSize: 24, fontWeight: 700, marginBottom: 16 }}>All Items</h1>
      <FilterBar
        types={types} statuses={statuses} priorities={priorities}
        onTypesChange={(v) => updateParams('type', v)}
        onStatusesChange={(v) => updateParams('status', v)}
        onPrioritiesChange={(v) => updateParams('priority', v)}
      />
      {loading ? (
        <div style={{ color: '#64748b', padding: 20 }}>Loading...</div>
      ) : (
        <div style={{ backgroundColor: '#fff', borderRadius: 8, boxShadow: '0 1px 3px rgba(0,0,0,0.1)', overflow: 'hidden' }}>
          <table style={{ width: '100%', borderCollapse: 'collapse', fontSize: 14 }}>
            <thead>
              <tr style={{ backgroundColor: '#f8fafc', borderBottom: '1px solid #e2e8f0' }}>
                <th style={thStyle}>Type</th>
                <th style={{ ...thStyle, textAlign: 'left' }}>Title</th>
                <th style={thStyle}>Status</th>
                <th style={thStyle}>Priority</th>
                <th style={thStyle}>Updated</th>
              </tr>
            </thead>
            <tbody>
              {entities.map((e) => (
                <tr
                  key={e.id}
                  style={{ borderBottom: '1px solid #f1f5f9', cursor: 'pointer' }}
                  onClick={() => navigate(`/entity/${e.id}`)}
                  onMouseEnter={(ev) => (ev.currentTarget.style.backgroundColor = '#f8fafc')}
                  onMouseLeave={(ev) => (ev.currentTarget.style.backgroundColor = '')}
                >
                  <td style={tdStyle}><TypeIcon type={e.type} /> {e.type}</td>
                  <td style={{ ...tdStyle, fontWeight: 500 }}>{e.title}</td>
                  <td style={tdStyle}>
                    <StatusSelect
                      status={e.status}
                      onChange={(s) => { handleStatusChange(e.id, s) }}
                    />
                  </td>
                  <td style={tdStyle}><PriorityBadge priority={e.priority} /></td>
                  <td style={{ ...tdStyle, color: '#64748b', fontSize: 12 }}>
                    {new Date(e.updated_at).toLocaleDateString()}
                  </td>
                </tr>
              ))}
              {entities.length === 0 && (
                <tr><td colSpan={5} style={{ padding: 20, textAlign: 'center', color: '#94a3b8' }}>No items found</td></tr>
              )}
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}

function StatusSelect({ status, onChange }: { status: Status; onChange: (s: string) => void }) {
  return (
    <span onClick={(e) => e.stopPropagation()}>
      <select
        value={status}
        onChange={(e) => onChange(e.target.value)}
        style={{
          padding: '2px 6px',
          borderRadius: 4,
          border: '1px solid #e2e8f0',
          fontSize: 12,
          backgroundColor: '#fff',
          cursor: 'pointer',
        }}
      >
        {allStatuses.map((s) => (
          <option key={s} value={s}>{s}</option>
        ))}
      </select>
    </span>
  )
}

const thStyle: React.CSSProperties = {
  padding: '10px 12px',
  fontSize: 12,
  fontWeight: 600,
  color: '#64748b',
  textAlign: 'center',
}

const tdStyle: React.CSSProperties = {
  padding: '10px 12px',
  textAlign: 'center',
}
