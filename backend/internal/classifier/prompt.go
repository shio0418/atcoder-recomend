package classifier

import (
	"atcoder-recomend/internal/model"
	"fmt"
	"strings"
)

var TagCandidates = []string{
	"全探索", "貪欲法", "二分探索", "累積和", "動的計画法",
	"グラフ", "木", "最短経路", "Union-Find", "数論",
	"組合せ論", "文字列", "幾何", "構築", "ゲーム理論",
}

func BuildPrompt(problems []model.Problem) string {
	var sb strings.Builder

	sb.WriteString("あなたは競技プログラミングの問題分類器です。\n")
	sb.WriteString("以下の各問題タイトルについて、タグ候補の中から最も当てはまるものを最大2個選んでください。\n")
	sb.WriteString("タグ候補にないタグを新しく作らないでください。当てはまるものがなければ空配列にしてください。\n\n")

	sb.WriteString("タグ候補: " + strings.Join(TagCandidates, ", ") + "\n\n")

	sb.WriteString("問題一覧:\n")
	for _, p := range problems {
		sb.WriteString(fmt.Sprintf("- problem_id: %s, title: %s\n", p.ProblemID, p.Title))
	}

	sb.WriteString("\n出力形式（この形式のJSON配列のみを出力し、説明文やコードブロック記号は含めないこと）:\n")
	sb.WriteString(`[{"problem_id": "abc123_a", "tags": ["全探索"]}, ...]`)

	return sb.String()
}
