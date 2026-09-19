# 進捗・TODO

最終更新: 2026-09-19

## 今どこにいるか

**Phase0（基盤）**: コードはほぼ完成、データはGeminiでの分類が進行中に停滞中（1345/9532件）。
**方針転換中**: タグ分類のバックエンドをGemini → Jevに乗り換える予定（コスト・レート制限の問題を解消できる見込み。詳細は[DESIGN.md](DESIGN.md)参照）。
**Phase1（レコメンド）**: 未着手。

## 次にやること（優先順）

1. **`internal/client/statement.go`のサンプル除去ロジックの仕上げ**
   `#task-statement`から`入力例`/`出力例`セクションを`.Remove()`する実装は完了・動作確認済み。空白圧縮（`strings.Fields`で潰す）はまだ。
2. **既存スクレイピング済みデータの再クレンジング**
   再スクレイピング不要、SQLで直接:
   ```sql
   UPDATE problems
   SET statement = trim(regexp_replace(regexp_replace(statement, '(入力例|出力例).*', ''), '\s+', ' ', 'g')),
       editorial_text = trim(regexp_replace(editorial_text, '\s+', ' ', 'g'))
   WHERE tags IS NULL AND statement IS NOT NULL;
   ```
3. **`classifier`パッケージをJevベースに書き換え**
   - `state` = `{title, statement, editorial_text(あれば)}`
   - `questions.tags` = `type: choice`, `criteria`は107タグ（現状タグ名をそのまま説明文にしてる簡易版。精度次第でdescriptionを足す）
   - `probabilities`から閾値 or 上位N件でタグを決定するロジックをGo側に実装（まだ未実装、要設計）
   - Go公式SDKが無いので生HTTPで叩く（`POST https://api.typesafe.ai/v1/systemone`）
   - 実験で分かった注意点: 難問×解説無しの組み合わせは精度が落ちる（Gemini/Jev共通の限界、要受容）
4. **既存1345件も含め全件Jevで統一して再分類するか判断**（コスト差は170円程度なので、一貫性重視で全件やり直すのが良さそう）
5. **全件タグ付けを最後まで走らせる**
6. Phase0完了確認後、Phase1（レコメンドロジック）の設計に着手
   - 弱点タグの特定（`submissions`×`problems`のJOIN、タグ別正答率）
   - 難易度バンド（今のレート＋α）
   - 忘却曲線（kenkoooo Issue #532の間隔計算式を参考に）
   - レコメンド結果表示のフロントエンド（React/Vite、まだ`frontend/`は空）

## 余裕があればやりたいこと（someday/maybe）

- **Jev × Gemini/より深いモデルのハイブリッド化**: 難易度が高くeditorialが無い問題だけ、より高精度なモデル（Gemini等）にフォールバックする条件分岐。今はまだ早い、Jev運用が安定してから検討
- タグにdescription（自分の言葉でパラフレーズしたもの）を追加し、Jevの`criteria`の精度を上げる
- ルールベースの完全一致マッチング（一度検討して見送った。優先度低い）
- 分野別正答率の可視化（Phase1のレコメンドロジックに統合予定なので単体では作らない）
- 非機能要件（[DESIGN.md](DESIGN.md)の非機能要件セクション参照）: 単体テスト、CI/CD、観測性

## 決定済みで蒸し返さないこと

- kikutak1182氏の分類結果データセット（`problems_data.json`）は不採用のまま（カバレッジ・精度・ライセンスの理由）
- 複合confidenceスコア（embedding類似度等を混ぜた合成スコア）は作らない。独立したシグナルのまま持つ
- Phase0の分野別正答率は独立機能として作らず、Phase1に統合する
