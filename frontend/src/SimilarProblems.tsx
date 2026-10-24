import { useState } from 'react'
import {
  fetchSimilarProblems,
  fetchSimilarProblemsByTag,
  type Problem,
} from './api'

type Props = {
  problemId: string
}

type FetchState = {
  loading: boolean
  error: string | null
  problems: Problem[] | null
}

const initialFetchState: FetchState = { loading: false, error: null, problems: null }

function ProblemSubList({ title, state }: { title: string; state: FetchState }) {
  return (
    <div className="similar-section">
      <h3>{title}</h3>
      {state.loading && <p>読み込み中...</p>}
      {state.error && <p className="error">エラー: {state.error}</p>}
      {state.problems && state.problems.length === 0 && (
        <p>似た問題が見つかりませんでした。</p>
      )}
      {state.problems && state.problems.length > 0 && (
        <ul>
          {state.problems.map((p) => (
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
  )
}

export function SimilarProblems({ problemId }: Props) {
  const [expanded, setExpanded] = useState(false)
  const [byEmbedding, setByEmbedding] = useState<FetchState>(initialFetchState)
  const [byTag, setByTag] = useState<FetchState>(initialFetchState)

  const handleToggle = () => {
    if (expanded) {
      setExpanded(false)
      return
    }
    setExpanded(true)
    if (byEmbedding.problems !== null) {
      return
    }

    setByEmbedding({ loading: true, error: null, problems: null })
    fetchSimilarProblems(problemId)
      .then((problems) => setByEmbedding({ loading: false, error: null, problems }))
      .catch((e: Error) =>
        setByEmbedding({ loading: false, error: e.message, problems: null }),
      )

    setByTag({ loading: true, error: null, problems: null })
    fetchSimilarProblemsByTag(problemId)
      .then((problems) => setByTag({ loading: false, error: null, problems }))
      .catch((e: Error) => setByTag({ loading: false, error: e.message, problems: null }))
  }

  return (
    <div className="similar-problems">
      <button type="button" className="similar-toggle" onClick={handleToggle}>
        {expanded ? '似た問題を隠す' : '似た問題を見る'}
      </button>
      {expanded && (
        <div className="similar-list">
          <ProblemSubList title="意味の近さ（embedding）" state={byEmbedding} />
          <ProblemSubList title="タグの一致" state={byTag} />
        </div>
      )}
    </div>
  )
}
