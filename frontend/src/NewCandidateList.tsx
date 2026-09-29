import type { Problem } from './api'
import { difficultyColor } from './difficultyColor'

type Props = {
  problems: Problem[]
}

export function NewCandidateList({ problems }: Props) {
  if (problems.length === 0) {
    return <p>今の条件に合う新規問題はありません。</p>
  }

  return (
    <ul className="problem-list">
      {problems.map((p) => (
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
  )
}
