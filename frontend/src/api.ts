export type SolvedProblem = {
  problemId: string
  acCount: number
  lastAcAt: string
  daysOverdue: number
}

export type Problem = {
  id: number
  problemId: string
  title: string
  difficulty: number | null
  tags: string[]
  contestId: string
  hasEditorial: boolean
}

export type RecommendResponse = {
  review: SolvedProblem[]
  new: Problem[]
}

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080'

export type ReviewParams = {
  baseInterval?: number
  growthRate?: number
  difficultyWeight?: number
  overdueCap?: number
}

export async function fetchRecommendations(
  params: ReviewParams = {},
): Promise<RecommendResponse> {
  const query = new URLSearchParams()
  if (params.baseInterval !== undefined) {
    query.set('base_interval', String(params.baseInterval))
  }
  if (params.growthRate !== undefined) {
    query.set('growth_rate', String(params.growthRate))
  }
  if (params.difficultyWeight !== undefined) {
    query.set('difficulty_weight', String(params.difficultyWeight))
  }
  if (params.overdueCap !== undefined) {
    query.set('overdue_cap', String(params.overdueCap))
  }
  const qs = query.toString()
  const res = await fetch(`${API_BASE_URL}/api/recommend${qs ? `?${qs}` : ''}`)
  if (!res.ok) {
    throw new Error(`レコメンド取得に失敗しました (status: ${res.status})`)
  }
  return res.json()
}

export async function fetchSimilarProblems(problemId: string): Promise<Problem[]> {
  const res = await fetch(
    `${API_BASE_URL}/api/similar?problem_id=${encodeURIComponent(problemId)}`,
  )
  if (!res.ok) {
    throw new Error(`類似問題取得に失敗しました (status: ${res.status})`)
  }
  return res.json()
}

export async function fetchSimilarProblemsByTag(problemId: string): Promise<Problem[]> {
  const res = await fetch(
    `${API_BASE_URL}/api/similarByTag?problem_id=${encodeURIComponent(problemId)}`,
  )
  if (!res.ok) {
    throw new Error(`タグによる類似問題取得に失敗しました (status: ${res.status})`)
  }
  return res.json()
}

export type SearchFilters = {
  difficultyMin?: number
  difficultyMax?: number
}

export async function searchProblems(
  query: string,
  filters: SearchFilters = {},
): Promise<Problem[]> {
  const params = new URLSearchParams({ q: query })
  if (filters.difficultyMin !== undefined) {
    params.set('difficulty_min', String(filters.difficultyMin))
  }
  if (filters.difficultyMax !== undefined) {
    params.set('difficulty_max', String(filters.difficultyMax))
  }

  const res = await fetch(`${API_BASE_URL}/api/search?${params.toString()}`)
  if (!res.ok) {
    throw new Error(`検索に失敗しました (status: ${res.status})`)
  }
  return res.json()
}

export async function fetchCFRecommendations(): Promise<string[]> {
  const res = await fetch(`${API_BASE_URL}/api/cfRecommend`)
  if (!res.ok) {
    throw new Error(`協調フィルタリング推薦の取得に失敗しました (status: ${res.status})`)
  }
  return res.json()
}
