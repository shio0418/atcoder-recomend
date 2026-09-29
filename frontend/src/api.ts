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
