import type { Priority } from '../api/types'

const priorityColors: Record<Priority, { bg: string; text: string }> = {
  low:      { bg: '#dbeafe', text: '#1d4ed8' },
  medium:   { bg: '#fef3c7', text: '#92400e' },
  high:     { bg: '#fee2e2', text: '#991b1b' },
  critical: { bg: '#991b1b', text: '#fff' },
}

export function PriorityBadge({ priority }: { priority: Priority }) {
  const colors = priorityColors[priority] || { bg: '#e5e7eb', text: '#374151' }
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
      {priority}
    </span>
  )
}
