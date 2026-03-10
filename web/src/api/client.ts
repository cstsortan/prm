import type {
  DashboardResponse,
  Entity,
  EntityDetailResponse,
  TreeNode,
  SearchResult,
} from './types'

const BASE = '/api'

async function fetchJSON<T>(url: string, init?: RequestInit): Promise<T> {
  const res = await fetch(url, init)
  if (!res.ok) {
    const body = await res.json().catch(() => ({ error: res.statusText }))
    throw new Error(body.error || res.statusText)
  }
  return res.json()
}

export async function getDashboard(): Promise<DashboardResponse> {
  return fetchJSON(`${BASE}/dashboard`)
}

export interface ListParams {
  type?: string
  status?: string
  priority?: string
  tags?: string
  sort?: string
  desc?: boolean
  archived?: boolean
}

export async function listEntities(params: ListParams = {}): Promise<Entity[]> {
  const q = new URLSearchParams()
  if (params.type) q.set('type', params.type)
  if (params.status) q.set('status', params.status)
  if (params.priority) q.set('priority', params.priority)
  if (params.tags) q.set('tags', params.tags)
  if (params.sort) q.set('sort', params.sort)
  if (params.desc) q.set('desc', 'true')
  if (params.archived) q.set('archived', 'true')
  return fetchJSON(`${BASE}/entities?${q}`)
}

export async function getEntity(id: string): Promise<EntityDetailResponse> {
  return fetchJSON(`${BASE}/entities/${id}`)
}

export async function getTree(id?: string): Promise<TreeNode[]> {
  const path = id ? `${BASE}/tree/${id}` : `${BASE}/tree`
  return fetchJSON(path)
}

export async function search(query: string): Promise<SearchResult[]> {
  return fetchJSON(`${BASE}/search?q=${encodeURIComponent(query)}`)
}

export async function updateStatus(id: string, status: string): Promise<Entity> {
  return fetchJSON(`${BASE}/entities/${id}/status`, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ status }),
  })
}

export async function addComment(id: string, text: string, author = 'web'): Promise<Entity> {
  return fetchJSON(`${BASE}/entities/${id}/comments`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ author, text }),
  })
}
