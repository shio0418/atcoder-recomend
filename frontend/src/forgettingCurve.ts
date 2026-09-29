// backendの internal/recommend/forgetting_curve.go (ForgettingCurveInterval) と同じ式。
// 表示の並び替え専用で、レコメンド対象の絞り込み自体はbackend側で完結している。
function forgettingCurveIntervalDays(acCount: number): number {
  return 5 * Math.pow(2, Math.max(acCount - 1, 0))
}

// 期限をどれだけ過ぎているか(日数)。値が大きいほど復習の優先度が高い。
export function daysOverdue(lastAcAt: string, acCount: number, now: Date = new Date()): number {
  const last = new Date(lastAcAt)
  const daysSinceAc = (now.getTime() - last.getTime()) / (1000 * 60 * 60 * 24)
  return daysSinceAc - forgettingCurveIntervalDays(acCount)
}
