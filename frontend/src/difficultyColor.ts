// AtCoderのレーティング色分けに準拠(灰→茶→緑→水→青→黄→橙→赤)
export function difficultyColor(difficulty: number | null): string {
  if (difficulty === null) return '#808080'
  if (difficulty < 400) return '#808080' // 灰
  if (difficulty < 800) return '#804000' // 茶
  if (difficulty < 1200) return '#008000' // 緑
  if (difficulty < 1600) return '#00C0C0' // 水
  if (difficulty < 2000) return '#0000FF' // 青
  if (difficulty < 2400) return '#C0C000' // 黄
  if (difficulty < 2800) return '#FF8000' // 橙
  return '#FF0000' // 赤
}
