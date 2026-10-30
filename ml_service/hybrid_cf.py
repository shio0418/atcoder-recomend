import pandas as pd
from train_cf import apply_popularity_weighting, compute_similarities, select_similar_users, recommend_problems

def find_weak_tags(conn, threshold=0.5):
    df = pd.read_sql("""
        SELECT tag, count(*) AS attempts, count(*) FILTER (WHERE result = 'AC') AS ac_count
        FROM submissions
        JOIN problems ON submissions.problem_id = problems.problem_id
        CROSS JOIN unnest(problems.tags) AS tag
        GROUP BY tag
    """, conn)
    df['ac_rate'] = df['ac_count'] / df['attempts']
    return set(df[df['ac_rate'] < threshold]['tag'])

def load_problem_tags(conn):
    df = pd.read_sql("SELECT tags, problem_id FROM problems", conn)
    result = {}

    for _, row in df.iterrows():
        tags = row['tags']
        if not isinstance(tags, list):
            tags = []
        result[row['problem_id']] = tags
    return result

def matches_weak_tag(problem_index, problem_cat, problem_tags, weak_tags):
    problem_id = problem_cat.categories[problem_index]
    tags = problem_tags.get(problem_id, [])
    return any(tag in weak_tags for tag in tags)

def rerank_by_weak_tag(top, problem_cat, problem_tags, weak_tags, k):
    ranked = sorted(
        top,
        key=lambda idx: not matches_weak_tag(idx, problem_cat, problem_tags, weak_tags)
    )
    return ranked[:k]

def make_hybrid_score_fn(weak_tags, problem_cat, problem_tags):
    def score(hidden_matrix, user_index, held_out, k=10):
        weighted_matrix = apply_popularity_weighting(hidden_matrix)
        similarities = compute_similarities(weighted_matrix, user_index)

        similar_indices = select_similar_users(hidden_matrix, similarities, user_index)
        top, scores = recommend_problems(weighted_matrix, similarities, user_index, similar_indices, top_k=k * 3)

        ranked = rerank_by_weak_tag(top, problem_cat, problem_tags, weak_tags, k)
        success = [
            pid for pid in ranked
            if pid in held_out
        ]
        return len(success) / k
    return score