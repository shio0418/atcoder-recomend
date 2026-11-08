# atcoder_recomend

AtCoderの提出履歴を蓄積し、問題ごとのタグや忘却曲線に基づいて次に解くべき問題を推薦するwebアプリ
協調フィルタリング・embedding検索などを組み合わせてみた

デモリンク：https://atcoder-recommend.shio0418.dev/

## 構成

```
atcoder_recomend/
├── backend/                    # Go
│   ├── cmd/
│   │   ├── api/                 # HTTP APIサーバー（フロントから叩かれる本体）
│   │   ├── server/              # AtCoder提出履歴・問題データの同期
│   │   ├── classify/            # LLMによる分野分類バッチ
│   │   ├── embed/               # 問題文のembedding生成バッチ
│   │   ├── recommend/           # 推薦ロジック動作確認用CLI
│   │   ├── collection_cf_data/  # 協調フィルタリング用の他ユーザーデータ収集
│   │   └── migrate/             # DBマイグレーション（goose）
│   ├── internal/
│   ├── migration/               # gooseマイグレーションファイル
│   └── go.mod
├── ml_service/                  # Python（gRPCサービス）
│   ├── server.py                # embedding検索サービス
│   ├── cf_server.py             # 協調フィルタリングサービス
│   ├── train_cf.py / hybrid_cf.py
│   ├── evaluate_cf.py / evaluate_search.py  # 評価スクリプト（MLflow記録）
│   └── requirements.txt
├── frontend/                    # React / TypeScript / Vite
├── infra/terraform/              # VPS（Vultr）+ DNS/TLS（Cloudflare）のIaC
├── docker-compose.yml            # ローカル開発用
└── docker-compose.prod.yml       # 本番（VPS）用
```

## 技術スタック

- **バックエンド**: Go / Echo / pgx / gRPC
- **ML・データ処理**: Python / gRPC / sentence-transformers（multilingual-e5-small）/ scikit-learn / pandas / MLflow / Gemini API（google-genai）
- **フロントエンド**: React / TypeScript / Vite
- **データベース**: PostgreSQL + pgvector
- **インフラ**: Docker / Docker Compose / Terraform / Vultr（VPS）/ Cloudflare（DNS・TLS）

## 技術的判断など
- 忘却曲線は、[kenkooooさん](https://github.com/kenkoooo/AtCoderProblems/issues/532)の数式を参考に自分で調整できるようにしました
- 協調フィルタリングについて、精度が安定していませんが実験的機能として追加しています
    - [こちらの記事](https://sucrose.hatenablog.com/entry/2014/04/26/000334)を参考に、多くのユーザーに説かれている問題ほどスコアが低くなるように調整しました
- 文章による問題検索については、geminiによるLLM-as-Judgeで精度を調べてみたところ、タグベースのものより精度が低くなるという結果でした

## セットアップ

```bash
# 1. DBを起動
docker compose up -d db

# 2. バックエンド（Go）
cd backend
go run ./cmd/migrate   # マイグレーション適用
go run ./cmd/server    # AtCoderからデータ同期
go run ./cmd/classify  # 分野分類バッチ
go run ./cmd/embed     # embedding生成バッチ

# 3. MLサービス（Python、別タブで2つ起動）
cd ml_service
source venv/bin/activate
python server.py       # embedding検索サービス（:5051）
python cf_server.py    # 協調フィルタリングサービス（:50052）

# 4. APIサーバー（別タブ）
cd backend
go run ./cmd/api       # :8080

# 5. フロントエンド（別タブ）
cd frontend
npm install
npm run dev            # :5173
```

`backend/.env`に`DATABASE_URL`・`GEMINI_API_KEY`（分野分類用）が必要。
`ml_service/.env`に`DATABASE_URL`・`GEMINI_API_KEY`（検索評価用）が必要。
