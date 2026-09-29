import type { Problem } from './api'

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
          <span>難易度: {p.difficulty ?? '不明'}</span>
          <span className="tags">{p.tags.join(', ')}</span>
        </li>
      ))}
    </ul>
  )
}
