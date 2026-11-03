from google import genai
from google.genai import types
import requests
import os
import numpy as np
import psycopg2
from dotenv import load_dotenv
import mlflow

TEST_QUERIES = [
    "二分探索を使う問題",
    "動的計画法で解ける問題",
    "グラフの最短経路を求める問題",
    "累積和を使うと楽になる問題",
    "貪欲法で解ける問題",
    "木構造に関する問題",
    "文字列のパターンマッチングに関する問題",
    "尺取り法を使う問題",
    "区間に対するクエリを処理する問題",
    "素数判定や素因数分解に関する問題",
    "UnionFindを使う問題",
    "二部グラフの判定に関する問題",
    "座標圧縮が必要な問題",
    "ビット演算を使う問題",
    "シミュレーションをそのまま実装する問題",
]

load_dotenv()
client = genai.Client(api_key=os.environ["GEMINI_API_KEY"])

def judge_relevance(query: str, title: str, tags: list[str]) -> bool:
    prompt = f"クエリ: {query}\n問題タイトル: {title}\nタグ: {', '.join(tags)}\nこの問題はクエリに関連していますか？"
    response = client.models.generate_content(
        model="gemini-3.8-flash",
        contents=prompt,
        config=types.GenerateContentConfig(
            response_mime_type="application/json",
            response_schema={
                "type": "OBJECT",
                "properties": {"relevant": {"type": "BOOLEAN"}},
                "required": ["relevant"],
            },
        ),
    )
    import json
    return json.loads(response.text)["relevant"]

def load_all_problems(conn):
    cur = conn.cursor()
    cur.execute("""
        SELECT problem_id, title, tags
        FROM problems
        WHERE statement IS NOT NULL AND trim(statement) != ''
    """)
    rows = cur.fetchall()
    result = [{"problem_id": r[0], "title": r[1], "tags": r[2] or []} for r in rows]
    return result

def search_problems(query, k):
    res = requests.get("http://localhost:8080/api/search", params={"q": query})
    return res.json()[:k]

def keyword_baseline(query, all_problems, k):
    matches = [
        p for p in all_problems 
        if query in p["title"] or any(tag in query for tag in p["tags"])
    ]
    return matches[:k]

def precision_at_k(query, results, k):
    num_relevance = 0
    for result in results:
        if judge_relevance(query, result["title"], result["tags"]):
            num_relevance += 1
    return num_relevance / k

def evaluate_all(queries, search_fn, k):
    precisions = []
    for query in queries:
        print(f"評価中: {query}")
        results = search_fn(query, k)
        precision = precision_at_k(query, results, k)
        precisions.append(precision)
    return np.mean(precisions), precisions

if __name__ == "__main__":
    conn = psycopg2.connect(os.environ["DATABASE_URL"])
    all_problems = load_all_problems(conn)
    conn.close()

    with mlflow.start_run(run_name="search_eval"):
        embedding_precision, _ = evaluate_all(TEST_QUERIES, search_problems, k=5)
        baseline_precision, _ = evaluate_all(
            TEST_QUERIES,
            lambda q, k: keyword_baseline(q, all_problems, k),
            k=5,
        )
        mlflow.log_param("n_queries", len(TEST_QUERIES))
        mlflow.log_param("k", 5)
        mlflow.log_metric("embedding_precision_at_5", embedding_precision)
        mlflow.log_metric("baseline_precision_at_5", baseline_precision)

        print(f"Embedding precision@5: {embedding_precision:.3f}")
        print(f"Baseline precision@5: {baseline_precision:.3f}")
    