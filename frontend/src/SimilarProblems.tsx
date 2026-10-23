import { useState } from 'react'
import { fetchSimilarProblems, type Problem } from './api'

type Props = {
  problemId: string
}

export function SimilarProblems({ problemId }: Props) {
  const [expanded, setExpanded] = useState(false)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [similar, setSimilar] = useState<Problem[] | null>(null)

  const handleToggle = () => {
    if (expanded) {
      setExpanded(false)
      return
    }
    setExpanded(true)
    if (similar !== null) {
      return
    }
    setLoading(true)
    setError(null)
    fetchSimilarProblems(problemId)
      .then(setSimilar)
      .catch((e: Error) => setError(e.message))
      .finally(() => setLoading(false))
  }

  return (
    <div className="similar-problems">
      <button type="button" className="similar-toggle" onClick={handleToggle}>
        {expanded ? '似た問題を隠す' : '似た問題を見る'}
      </button>
      {expanded && (
        <div className="similar-list">
          {loading && <p>読み込み中...</p>}
          {error && <p className="error">エラー: {error}</p>}
          {similar && similar.length === 0 && <p>似た問題が見つかりませんでした。</p>}
          {similar && similar.length > 0 && (
            <ul>
              {similar.map((p) => (
                <li key={p.problemId}>
                  <a
                    href={`https://atcoder.jp/contests/${p.contestId}/tasks/${p.problemId}`}
                    target="_blank"
                    rel="noreferrer"
                  >
                    {p.title}
                  </a>
                  <span className="tags">
                    {p.tags.map((tag) => (
                      <span key={tag} className="tag-badge">
                        {tag}
                      </span>
                    ))}
                  </span>
                </li>
              ))}
            </ul>
          )}
        </div>
      )}
    </div>
  )
}
