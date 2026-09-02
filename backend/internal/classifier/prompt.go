package classifier

import (
	"atcoder-recomend/internal/model"
	"fmt"
	"strings"
)

var TagCandidates = []string{
	"動的計画法", "桁DP", "最大流", "グラフ", "木構造",
	"幾何", "Union-Find", "セグメント木", "Fenwick Tree", "素数",
	"期待値", "累積和", "尺取り法", "遅延セグ木", "平衡二分木",
	"深さ優先探索", "幅優先探索", "ダイクストラ法", "ワーシャルフロイド法", "ベルマンフォード法",
	"最小共通祖先", "木DP", "全方位木DP", "強連結成分", "最小費用流",
	"高速フーリエ変換", "形式的冪級数", "確率・期待値DP", "主客転倒", "ダブリング",
	"二分探索", "半分全列挙", "平方分割", "貪欲法", "クエリ先読み",
	"ビット全探索", "bitDP", "接尾辞配列", "Trie木", "最長共通接頭辞",
	"ローリングハッシュ", "Z-algorithm", "写像12相", "Grundy数", "ゲーム理論",
	"座標回転", "牛ゲー", "調和級数", "マンハッタン距離", "構築問題",
	"括弧列", "Bitset", "最大公約数", "Heavy-Light分解", "XOR演算",
	"DAG", "カタラン数", "シミュレーション", "座標圧縮", "トポロジカルソート",
	"ハミルトン路", "バックトラック", "ビット演算", "メモ化再帰", "ユークリッドの互除法",
	"三分探索", "中国剰余定理", "優先度付きキュー", "全列挙", "全探索",
	"再帰関数", "凸包", "分割統治法", "区間スケジューリング", "回文",
	"数学", "最小シュタイナー木", "最小公倍数", "最適化問題", "素因数分解",
	"組み合わせ論", "編集距離", "行列", "辞書順", "鳩の巣原理",
	"Manacher法", "0-1 BFS", "二部マッチング", "二重連結成分", "オイラーツアー",
	"フロー問題", "最小カット", "Monge性", "ゼータ変換", "Alien DP",
	"包除原理", "約数列挙", "Wavelet Matrix", "永続セグメント木", "Mo's algorithm",
	"いもす法", "最長増加部分列", "巡回セールスマン問題", "OEIS参照", "双方向リスト",
	"ハッシュ", "2-SAT",
}

func BuildPrompt(problems []model.Problem) string {
	var sb strings.Builder

	sb.WriteString("あなたは競技プログラミングの問題分類器です。\n")
	sb.WriteString("以下の各問題タイトルについて、タグ候補の中から最も当てはまるものを最大5個選んでください。\n")
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
