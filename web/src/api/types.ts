export type EntityType = 'epic' | 'story' | 'sub-task' | 'task' | 'bug'
export type Status = 'backlog' | 'todo' | 'in-progress' | 'review' | 'done' | 'cancelled' | 'archived'
export type Priority = 'low' | 'medium' | 'high' | 'critical'
export type Severity = 'cosmetic' | 'minor' | 'major' | 'blocker'

export interface Comment {
  author: string
  text: string
  created_at: string
}

export interface Entity {
  id: string
  type: EntityType
  slug: string
  title: string
  description?: string
  status: Status
  priority: Priority
  tags?: string[]
  created_at: string
  updated_at: string
  started_at?: string
  completed_at?: string
  due_date?: string
  dependencies?: string[]
  comments?: Comment[]
  parent_id?: string
  children?: string[]
  severity?: Severity
  steps_to_reproduce?: string
}

export interface DashboardStats {
  Total: number
  ByType: Record<EntityType, number>
  ByStatus: Record<Status, number>
  ByPriority: Record<Priority, number>
  BySeverity: Record<Severity, number>
}

export interface DashboardResponse {
  project: string
  stats: DashboardStats
}

export interface EntityDetailResponse {
  entity: Entity
  readme: string
  dependencies: Record<string, string>
  children_entities?: Entity[]
}

export interface TreeNode {
  entity: Entity
  children?: TreeNode[]
}

export interface SearchResult {
  entity: Entity
  score: number
}
