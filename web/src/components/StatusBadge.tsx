import type { Status } from '../api/types'

const statusColors: Record<Status, { bg: string; text: string }> = {
  'backlog':     { bg: '#6b7280', text: '#fff' },
  'todo':        { bg: '#3b82f6', text: '#fff' },
  'in-progress': { bg: '#f59e0b', text: '#000' },
  'review':      { bg: '#8b5cf6', text: '#fff' },
  'done':        { bg: '#10b981', text: '#fff' },
  'cancelled':   { bg: '#ef4444', text: '#fff' },
  'archived':    { bg: '#9ca3af', text: '#fff' },
}

export function StatusBadge({ status }: { status: Status }) {
  const colors = statusColors[status] || { bg: '#6b7280', text: '#fff' }
  return (
    <span
      style={{
        display: 'inline-block',
        padding: '2px 8px',
        borderRadius: '4px',
        fontSize: '12px',
        fontWeight: 600,
        backgroundColor: colors.bg,
        color: colors.text,
      }}
    >
      {status}
    </span>
  )
}
