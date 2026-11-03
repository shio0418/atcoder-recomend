import { useState } from 'react'
import type { ReviewParams } from './api'

type Props = {
  onApply: (params: ReviewParams) => void
  loading: boolean
}

export function ReviewSettings({ onApply, loading }: Props) {
  const [baseInterval, setBaseInterval] = useState('5')
  const [growthRate, setGrowthRate] = useState('2')
  const [difficultyWeight, setDifficultyWeight] = useState('1000')

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    onApply({
      baseInterval: Number(baseInterval),
      growthRate: Number(growthRate),
      difficultyWeight: Number(difficultyWeight),
    })
  }

  return (
    <form className="review-settings" onSubmit={handleSubmit}>
      <p className="formula">
        復習間隔(日) = {baseInterval || 'a'} × {growthRate || 'b'}
        <sup>(AC回数-1)</sup> ÷ (1+WA回数) ÷ (1 + 難易度 ÷ {difficultyWeight || 'c'})
      </p>
      <div className="review-settings-inputs">
        <label>
          基準日数 (a)
          <input
            type="number"
            value={baseInterval}
            onChange={(e) => setBaseInterval(e.target.value)}
          />
        </label>
        <label>
          成長率 (b)
          <input
            type="number"
            step="0.1"
            value={growthRate}
            onChange={(e) => setGrowthRate(e.target.value)}
          />
        </label>
        <label>
          難易度の重み (c)
          <input
            type="number"
            value={difficultyWeight}
            onChange={(e) => setDifficultyWeight(e.target.value)}
          />
        </label>
        <button type="submit" disabled={loading}>
          {loading ? '更新中...' : '復習間隔を更新'}
        </button>
      </div>
    </form>
  )
}
