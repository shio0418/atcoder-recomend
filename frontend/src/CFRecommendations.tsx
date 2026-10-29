import { useEffect, useState } from 'react'
import { fetchCFRecommendations } from './api'

export function CFRecommendations() {
  const [problemIds, setProblemIds] = useState<string[] | null>(null)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    fetchCFRecommendations()
      .then(setProblemIds)
      .catch((e: Error) => setError(e.message))
  }, [])

  return (
    <section>
      <h2>協調フィルタリングによる推薦</h2>
      <p className="section-note">
        似た解答傾向のユーザーが解いている問題から推薦（実験的機能）
      </p>
      {error && <p className="error">エラー: {error}</p>}
      {!error && !problemIds && <p>読み込み中...</p>}
      {problemIds && problemIds.length === 0 && <p>推薦できる問題が見つかりませんでした。</p>}
      {problemIds && problemIds.length > 0 && (
        <ul className="problem-list">
          {problemIds.map((problemId) => (
            <li key={problemId} className="problem-card">
              <a
                href={`https://atcoder.jp/contests/${problemId.split('_')[0]}/tasks/${problemId}`}
                target="_blank"
                rel="noreferrer"
              >
                {problemId}
              </a>
            </li>
          ))}
        </ul>
      )}
    </section>
  )
}
