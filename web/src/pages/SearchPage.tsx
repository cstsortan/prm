import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { search } from '../api/client'
import type { SearchResult } from '../api/types'
import { StatusBadge } from '../components/StatusBadge'
import { PriorityBadge } from '../components/PriorityBadge'
import { TypeIcon } from '../components/TypeIcon'

export function SearchPage() {
  const [query, setQuery] = useState('')
  const [results, setResults] = useState<SearchResult[]>([])
  const [loading, setLoading] = useState(false)
  const [searched, setSearched] = useState(false)
  const [error, setError] = useState('')
  const navigate = useNavigate()

  const doSearch = async () => {
    if (!query.trim()) return
    setLoading(true)
    setError('')
    try {
      const r = await search(query.trim())
      setResults(r)
      setSearched(true)
    } catch (e: any) {
      setError(e.message)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div>
      <h1 style={{ fontSize: 24, fontWeight: 700, marginBottom: 16 }}>Search</h1>
      <div style={{ display: 'flex', gap: 8, marginBottom: 20 }}>
        <input
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          onKeyDown={(e) => e.key === 'Enter' && doSearch()}
          placeholder="Search entities..."
          style={{
            flex: 1, padding: '10px 14px', borderRadius: 6,
            border: '1px solid #e2e8f0', fontSize: 14,
          }}
          autoFocus
        />
        <button
          onClick={doSearch}
          disabled={loading}
          style={{
            padding: '10px 20px', borderRadius: 6, border: 'none',
            backgroundColor: '#334155', color: '#fff', cursor: 'pointer',
            fontWeight: 500,
          }}
        >
          {loading ? 'Searching...' : 'Search'}
        </button>
      </div>

      {error && <div style={{ color: '#ef4444', marginBottom: 12 }}>Error: {error}</div>}

      {searched && (
        <div style={{
          backgroundColor: '#fff', borderRadius: 8,
          boxShadow: '0 1px 3px rgba(0,0,0,0.1)', overflow: 'hidden',
        }}>
          {results.length === 0 ? (
            <div style={{ padding: 20, textAlign: 'center', color: '#94a3b8' }}>
              No results for "{query}"
            </div>
          ) : (
            results.map((r) => (
              <div
                key={r.entity.id}
                onClick={() => navigate(`/entity/${r.entity.id}`)}
                style={{
                  display: 'flex', alignItems: 'center', gap: 10,
                  padding: '12px 16px', borderBottom: '1px solid #f1f5f9',
                  cursor: 'pointer', transition: 'background-color 0.15s',
                }}
                onMouseEnter={(e) => (e.currentTarget.style.backgroundColor = '#f8fafc')}
                onMouseLeave={(e) => (e.currentTarget.style.backgroundColor = '')}
              >
                <TypeIcon type={r.entity.type} />
                <span style={{ flex: 1, fontWeight: 500 }}>{r.entity.title}</span>
                <StatusBadge status={r.entity.status} />
                <PriorityBadge priority={r.entity.priority} />
                <span style={{ fontSize: 11, color: '#94a3b8' }}>score: {r.score}</span>
              </div>
            ))
          )}
        </div>
      )}
    </div>
  )
}
