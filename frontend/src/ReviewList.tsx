import type { SolvedProblem } from './api'

type Props = {
  problems: SolvedProblem[]
}

export function ReviewList({ problems }: Props) {
  if (problems.length === 0) {
    return <p>今日復習すべき問題はありません。</p>
  }

  return (
    <ul className="problem-list">
      {problems.map((p) => (
        <li key={p.problemId} className="problem-card">
          <a
            href={`https://atcoder.jp/contests/${p.problemId.split('_')[0]}/tasks/${p.problemId}`}
            target="_blank"
            rel="noreferrer"
          >
            {p.problemId}
          </a>
          <span>AC回数: {p.acCount}</span>
          <span>最終AC: {new Date(p.lastAcAt).toLocaleDateString()}</span>
        </li>
      ))}
    </ul>
  )
}
