import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { getTree } from '../api/client'
import type { TreeNode } from '../api/types'
import { StatusBadge } from '../components/StatusBadge'
import { PriorityBadge } from '../components/PriorityBadge'
import { TypeIcon } from '../components/TypeIcon'

function TreeNodeComponent({ node, depth = 0 }: { node: TreeNode; depth?: number }) {
  const [expanded, setExpanded] = useState(depth < 2)
  const navigate = useNavigate()
  const hasChildren = node.children && node.children.length > 0

  return (
    <div style={{ marginLeft: depth * 20 }}>
      <div
        style={{
          display: 'flex',
          alignItems: 'center',
          gap: 8,
          padding: '6px 10px',
          borderRadius: 6,
          cursor: 'pointer',
          transition: 'background-color 0.15s',
        }}
        onMouseEnter={(e) => (e.currentTarget.style.backgroundColor = '#f1f5f9')}
        onMouseLeave={(e) => (e.currentTarget.style.backgroundColor = '')}
      >
        {hasChildren ? (
          <span
            onClick={(e) => { e.stopPropagation(); setExpanded(!expanded) }}
            style={{ width: 16, fontSize: 12, color: '#94a3b8', cursor: 'pointer', userSelect: 'none' }}
          >
            {expanded ? '\u25BC' : '\u25B6'}
          </span>
        ) : (
          <span style={{ width: 16 }} />
        )}
        <span onClick={() => navigate(`/entity/${node.entity.id}`)}>
          <TypeIcon type={node.entity.type} />
        </span>
        <span
          onClick={() => navigate(`/entity/${node.entity.id}`)}
          style={{ fontWeight: 500, fontSize: 14, flex: 1 }}
        >
          {node.entity.title}
        </span>
        <StatusBadge status={node.entity.status} />
        <PriorityBadge priority={node.entity.priority} />
      </div>
      {expanded && hasChildren && (
        <div>
          {node.children!.map((child) => (
            <TreeNodeComponent key={child.entity.id} node={child} depth={depth + 1} />
          ))}
        </div>
      )}
    </div>
  )
}

export function TreePage() {
  const [trees, setTrees] = useState<TreeNode[]>([])
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    getTree()
      .then(setTrees)
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false))
  }, [])

  if (error) return <div style={{ color: '#ef4444', padding: 20 }}>Error: {error}</div>
  if (loading) return <div style={{ padding: 20, color: '#64748b' }}>Loading...</div>

  return (
    <div>
      <h1 style={{ fontSize: 24, fontWeight: 700, marginBottom: 16 }}>Epic Tree</h1>
      <div style={{
        backgroundColor: '#fff', borderRadius: 8, padding: 16,
        boxShadow: '0 1px 3px rgba(0,0,0,0.1)',
      }}>
        {trees.length === 0 ? (
          <div style={{ color: '#94a3b8', textAlign: 'center', padding: 20 }}>No epics found</div>
        ) : (
          trees.map((tree) => <TreeNodeComponent key={tree.entity.id} node={tree} />)
        )}
      </div>
    </div>
  )
}
