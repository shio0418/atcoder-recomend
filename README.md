# atcoder_recomend

AtCoderの提出履歴を蓄積し、弱点タグ・忘却曲線に基づいて次に解くべき問題を推薦する個人開発プロジェクト（進行中）。

- 設計・技術選定の理由: [docs/DESIGN.md](docs/DESIGN.md)
- 今の進捗・次のTODO: [docs/PROGRESS.md](docs/PROGRESS.md)

## 構成

```
atcoder_recomend/
├── backend/          # Go
│   ├── cmd/
│   │   ├── server/   # AtCoder提出履歴・問題データの同期
│   │   ├── classify/ # LLMによる分野分類バッチ
│   │   └── migrate/  # DBマイグレーション（goose）
│   ├── internal/
│   ├── migration/    # gooseマイグレーションファイル
│   └── go.mod
├── frontend/         # React/Vite（未着手）
└── docker-compose.yml
```

## セットアップ

```bash
docker compose up -d
cd backend
go run ./cmd/migrate   # マイグレーション適用
go run ./cmd/server    # AtCoderからデータ同期
go run ./cmd/classify  # 分野分類バッチ
```

`.env`に`DATABASE_URL`・`GEMINI_API_KEY`（または今後Jevのキー）が必要。
