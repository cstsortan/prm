import { useEffect, useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import ReactMarkdown from 'react-markdown'
import { getEntity, updateStatus, addComment } from '../api/client'
import type { EntityDetailResponse, Status } from '../api/types'
import { PriorityBadge } from '../components/PriorityBadge'
import { TypeIcon } from '../components/TypeIcon'

const allStatuses: Status[] = ['backlog', 'todo', 'in-progress', 'review', 'done', 'cancelled']

export function EntityDetailPage() {
  const { id } = useParams<{ id: string }>()
  const [data, setData] = useState<EntityDetailResponse | null>(null)
  const [error, setError] = useState('')
  const [commentText, setCommentText] = useState('')

  const load = () => {
    if (!id) return
    getEntity(id).then(setData).catch((e) => setError(e.message))
  }

  useEffect(load, [id])

  const handleStatusChange = async (newStatus: string) => {
    if (!id) return
    try {
      const updated = await updateStatus(id, newStatus)
      setData((prev) => prev ? { ...prev, entity: updated } : prev)
    } catch (e: any) {
      setError(e.message)
    }
  }

  const handleAddComment = async () => {
    if (!id || !commentText.trim()) return
    try {
      const updated = await addComment(id, commentText.trim())
      setData((prev) => prev ? { ...prev, entity: updated } : prev)
      setCommentText('')
    } catch (e: any) {
      setError(e.message)
    }
  }

  if (error) return <div style={{ color: '#ef4444', padding: 20 }}>Error: {error}</div>
  if (!data) return <div style={{ padding: 20, color: '#64748b' }}>Loading...</div>

  const { entity, readme, dependencies, children_entities } = data

  return (
    <div style={{ maxWidth: 800 }}>
      <div style={{ marginBottom: 20 }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: 10, marginBottom: 8 }}>
          <TypeIcon type={entity.type} />
          <span style={{ fontSize: 12, color: '#94a3b8', textTransform: 'uppercase', fontWeight: 600 }}>{entity.type}</span>
          <span style={{ fontSize: 12, color: '#cbd5e1' }}>{entity.id.slice(0, 8)}</span>
        </div>
        <h1 style={{ fontSize: 24, fontWeight: 700, margin: 0 }}>{entity.title}</h1>
      </div>

      {/* Metadata grid */}
      <div style={{
        display: 'grid', gridTemplateColumns: 'auto 1fr', gap: '8px 16px',
        backgroundColor: '#fff', padding: 16, borderRadius: 8,
        boxShadow: '0 1px 3px rgba(0,0,0,0.1)', marginBottom: 20,
        fontSize: 14,
      }}>
        <span style={labelStyle}>Status</span>
        <span>
          <select
            value={entity.status}
            onChange={(e) => handleStatusChange(e.target.value)}
            style={{ padding: '2px 6px', borderRadius: 4, border: '1px solid #e2e8f0', fontSize: 13 }}
          >
            {allStatuses.map((s) => <option key={s} value={s}>{s}</option>)}
          </select>
        </span>
        <span style={labelStyle}>Priority</span>
        <span><PriorityBadge priority={entity.priority} /></span>
        <span style={labelStyle}>Slug</span>
        <span style={{ fontFamily: 'monospace', fontSize: 13 }}>{entity.slug}</span>
        <span style={labelStyle}>Created</span>
        <span>{new Date(entity.created_at).toLocaleString()}</span>
        <span style={labelStyle}>Updated</span>
        <span>{new Date(entity.updated_at).toLocaleString()}</span>
        {entity.tags && entity.tags.length > 0 && <>
          <span style={labelStyle}>Tags</span>
          <span style={{ display: 'flex', gap: 4 }}>
            {entity.tags.map((t) => (
              <span key={t} style={{ padding: '1px 6px', backgroundColor: '#e2e8f0', borderRadius: 3, fontSize: 12 }}>{t}</span>
            ))}
          </span>
        </>}
        {children_entities && children_entities.length > 0 && <>
          <span style={labelStyle}>Children</span>
          <span>{children_entities.length} items</span>
        </>}
        {entity.parent_id && <>
          <span style={labelStyle}>Parent</span>
          <Link to={`/entity/${entity.parent_id}`} style={{ color: '#2563eb' }}>
            {entity.parent_id.slice(0, 8)}
          </Link>
        </>}
        {dependencies && Object.keys(dependencies).length > 0 && <>
          <span style={labelStyle}>Depends on</span>
          <span style={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
            {Object.entries(dependencies).map(([depId, title]) => (
              <Link key={depId} to={`/entity/${depId}`} style={{ color: '#2563eb', fontSize: 13 }}>
                {title}
              </Link>
            ))}
          </span>
        </>}
        {entity.severity && <>
          <span style={labelStyle}>Severity</span>
          <span>{entity.severity}</span>
        </>}
      </div>

      {/* Children */}
      {children_entities && children_entities.length > 0 && (
        <div style={{
          backgroundColor: '#fff', padding: 20, borderRadius: 8,
          boxShadow: '0 1px 3px rgba(0,0,0,0.1)', marginBottom: 20,
        }}>
          <h2 style={{ fontSize: 16, fontWeight: 600, marginBottom: 12 }}>
            Children ({children_entities.length})
          </h2>
          <div style={{ display: 'flex', flexDirection: 'column', gap: 6 }}>
            {children_entities.map((child) => (
              <Link
                key={child.id}
                to={`/entity/${child.id}`}
                style={{
                  display: 'flex', alignItems: 'center', gap: 10,
                  padding: '8px 12px', borderRadius: 6,
                  border: '1px solid #f1f5f9', textDecoration: 'none', color: 'inherit',
                  transition: 'background-color 0.15s',
                }}
                onMouseEnter={(e) => (e.currentTarget.style.backgroundColor = '#f8fafc')}
                onMouseLeave={(e) => (e.currentTarget.style.backgroundColor = 'transparent')}
              >
                <TypeIcon type={child.type} />
                <span style={{ flex: 1, fontWeight: 500, fontSize: 14 }}>{child.title}</span>
                <span style={{
                  fontSize: 11, padding: '2px 8px', borderRadius: 4,
                  backgroundColor: child.status === 'done' ? '#dcfce7' : child.status === 'in-progress' ? '#dbeafe' : '#f1f5f9',
                  color: child.status === 'done' ? '#166534' : child.status === 'in-progress' ? '#1e40af' : '#64748b',
                }}>{child.status}</span>
                <PriorityBadge priority={child.priority} />
              </Link>
            ))}
          </div>
        </div>
      )}

      {/* README */}
      {readme && (
        <div style={{
          backgroundColor: '#fff', padding: 20, borderRadius: 8,
          boxShadow: '0 1px 3px rgba(0,0,0,0.1)', marginBottom: 20,
        }}>
          <h2 style={{ fontSize: 16, fontWeight: 600, marginBottom: 12 }}>Description</h2>
          <div className="markdown-body" style={{ fontSize: 14, lineHeight: 1.6 }}>
            <ReactMarkdown>{readme}</ReactMarkdown>
          </div>
        </div>
      )}

      {/* Comments */}
      <div style={{
        backgroundColor: '#fff', padding: 20, borderRadius: 8,
        boxShadow: '0 1px 3px rgba(0,0,0,0.1)',
      }}>
        <h2 style={{ fontSize: 16, fontWeight: 600, marginBottom: 12 }}>
          Comments ({entity.comments?.length || 0})
        </h2>
        {entity.comments?.map((c, i) => (
          <div key={i} style={{ borderBottom: '1px solid #f1f5f9', padding: '10px 0' }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 4 }}>
              <span style={{ fontWeight: 600, fontSize: 13 }}>{c.author}</span>
              <span style={{ fontSize: 12, color: '#94a3b8' }}>{new Date(c.created_at).toLocaleString()}</span>
            </div>
            <div style={{ fontSize: 14, color: '#334155' }}>{c.text}</div>
          </div>
        ))}
        <div style={{ display: 'flex', gap: 8, marginTop: 12 }}>
          <input
            value={commentText}
            onChange={(e) => setCommentText(e.target.value)}
            placeholder="Add a comment..."
            onKeyDown={(e) => e.key === 'Enter' && handleAddComment()}
            style={{
              flex: 1, padding: '8px 12px', borderRadius: 6,
              border: '1px solid #e2e8f0', fontSize: 14,
            }}
          />
          <button
            onClick={handleAddComment}
            style={{
              padding: '8px 16px', borderRadius: 6, border: 'none',
              backgroundColor: '#334155', color: '#fff', cursor: 'pointer',
              fontWeight: 500,
            }}
          >
            Post
          </button>
        </div>
      </div>
    </div>
  )
}

const labelStyle: React.CSSProperties = {
  fontWeight: 600,
  color: '#64748b',
  fontSize: 13,
}
