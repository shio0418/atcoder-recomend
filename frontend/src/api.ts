export type SolvedProblem = {
  problemId: string
  acCount: number
  lastAcAt: string
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

const API_BASE_URL = 'http://localhost:8080'

export async function fetchRecommendations(): Promise<RecommendResponse> {
  const res = await fetch(`${API_BASE_URL}/api/recommend`)
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
