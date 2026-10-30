import { useState } from 'react'
import { searchProblems, type Problem } from './api'
import { difficultyColor } from './difficultyColor'

export function SearchSection() {
  const [query, setQuery] = useState('')
  const [difficultyMin, setDifficultyMin] = useState('')
  const [difficultyMax, setDifficultyMax] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [results, setResults] = useState<Problem[] | null>(null)

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault()
    if (!query.trim()) {
      return
    }
    setLoading(true)
    setError(null)
    searchProblems(query, {
      difficultyMin: difficultyMin === '' ? undefined : Number(difficultyMin),
      difficultyMax: difficultyMax === '' ? undefined : Number(difficultyMax),
    })
      .then(setResults)
      .catch((e: Error) => setError(e.message))
      .finally(() => setLoading(false))
  }

  return (
    <section>
      <h2>問題を検索</h2>
      <form className="search-form" onSubmit={handleSearch}>
        <input
          type="text"
          placeholder="例: グラフの最短経路"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
        />
        <input
          type="number"
          placeholder="難易度 下限"
          value={difficultyMin}
          onChange={(e) => setDifficultyMin(e.target.value)}
        />
        <input
          type="number"
          placeholder="難易度 上限"
          value={difficultyMax}
          onChange={(e) => setDifficultyMax(e.target.value)}
        />
        <button type="submit" disabled={loading}>
          {loading ? '検索中...' : '検索'}
        </button>
      </form>

      {error && <p className="error">エラー: {error}</p>}
      {results && results.length === 0 && <p>見つかりませんでした。</p>}
      {results && results.length > 0 && (
        <ul className="problem-list">
          {results.map((p) => (
            <li key={p.problemId} className="problem-card">
              <a
                href={`https://atcoder.jp/contests/${p.contestId}/tasks/${p.problemId}`}
                target="_blank"
                rel="noreferrer"
              >
                {p.title}
              </a>
              <span
                className="difficulty-badge"
                style={{ backgroundColor: difficultyColor(p.difficulty) }}
              >
                {p.difficulty ?? '不明'}
              </span>
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
    </section>
  )
}
