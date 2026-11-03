import { useState } from 'react'
import type { SolvedProblem } from './api'

type Props = {
  problems: SolvedProblem[]
}

export function ReviewList({ problems }: Props) {
  const [maxOverdueDays, setMaxOverdueDays] = useState('')

  const threshold = maxOverdueDays === '' ? null : Number(maxOverdueDays)
  const filtered =
    threshold === null ? problems : problems.filter((p) => p.daysOverdue <= threshold)

  // 並び順(優先度の上限を加味したもの)はバックエンドが既に決めてるので、ここでは並べ替えない
  const sorted = filtered.slice(0, 5)

  return (
    <div>
      <label className="overdue-filter">
        <input
          type="number"
          placeholder="例: 180"
          value={maxOverdueDays}
          onChange={(e) => setMaxOverdueDays(e.target.value)}
        />
        日以上超過のものは非表示
      </label>

      {sorted.length === 0 ? (
        <p>今日復習すべき問題はありません。</p>
      ) : (
        <ul className="problem-list">
          {sorted.map((p) => (
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
              <span>{Math.floor(p.daysOverdue)}日超過</span>
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}
