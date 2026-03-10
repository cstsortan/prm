import type { EntityType } from '../api/types'

const typeLabels: Record<EntityType, { icon: string; color: string }> = {
  epic:       { icon: '\u26A1', color: '#7c3aed' },
  story:      { icon: '\uD83D\uDCD6', color: '#2563eb' },
  'sub-task': { icon: '\u2611\uFE0F', color: '#059669' },
  task:       { icon: '\u2705', color: '#059669' },
  bug:        { icon: '\uD83D\uDC1B', color: '#dc2626' },
}

export function TypeIcon({ type }: { type: EntityType }) {
  const label = typeLabels[type] || { icon: '\u2753', color: '#6b7280' }
  return (
    <span title={type} style={{ color: label.color, marginRight: 4 }}>
      {label.icon}
    </span>
  )
}
