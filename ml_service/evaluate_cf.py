from train_cf import build_matrix, apply_popularity_weighting, compute_similarities, select_similar_users, recommend_problems, load_external_solves
import numpy as np
import psycopg2
import os
from dotenv import load_dotenv
import mlflow
from hybrid_cf import make_hybrid_score_fn, find_weak_tags, load_problem_tags

rng = np.random.default_rng(42)

def hold_out_problem(matrix, user_index, holdout_ratio=0.2):
    user_row = matrix[user_index]
    solved_indices = user_row.indices

    n_holdout = max(1, int(len(solved_indices) * holdout_ratio))
    held_out = rng.choice(solved_indices, n_holdout, replace = False)

    modified = matrix.tolil()
    modified[user_index, held_out] = 0
    modified = modified.tocsr()

    return modified, held_out

def select_users(matrix, n, min_solved=5):
    solve_counts = np.asarray(matrix.sum(axis=1)).flatten()
    canditates = [
        i for i in range(len(solve_counts))
        if solve_counts[i] > min_solved
    ]
    return rng.choice(canditates, n, replace=False)

def evaluate_user(hidden_matrix, user_index, held_out, k=10):
    weighted_matrix = apply_popularity_weighting(hidden_matrix)
    similarities = compute_similarities(weighted_matrix, user_index)

    similar_indicies = select_similar_users(hidden_matrix, similarities, user_index)
    top, scores = recommend_problems(weighted_matrix, similarities, user_index, similar_indicies, k)

    success = [
        pid for pid in top
        if pid in held_out
    ]
    return len(success) / k

def prepare_eval_cases(matrix, n_users=50, holdout_ratio=0.2):
    test_users = select_users(matrix, n_users)
    cases = []
    for user in test_users:
        hidden_matrix, held_out = hold_out_problem(matrix, user, holdout_ratio)
        cases.append((user, hidden_matrix, held_out))
    return cases

def evaluate_all(cases, score_fn, k=10):
    precisions = [score_fn(hidden_matrix, user, held_out, k) for user, hidden_matrix, held_out in cases]
    return np.mean(precisions), precisions

def popular_baseline_user(hidden_matrix, user_index, held_out, k=10):
    solve_counts = np.asarray(hidden_matrix.sum(axis=0)).flatten()
    already_solved = hidden_matrix[user_index].toarray().flatten()
    solve_counts[already_solved > 0] = -1
    top = np.argsort(solve_counts)[::-1][:k]
    success = [pid for pid in top if pid in held_out]
    return len(success) / k

if __name__ == "__main__":
    load_dotenv()
    conn = psycopg2.connect(os.environ["DATABASE_URL"])
    
    df = load_external_solves(conn)
    matrix, user_cat, problem_cat = build_matrix(df)

    weak_tags = find_weak_tags(conn)
    problem_tags = load_problem_tags(conn)
    conn.close()

    with mlflow.start_run(run_name="cf_comparison"):
        mlflow.log_param("n_users", 50)
        mlflow.log_param("k", 10)
        mlflow.log_param("holdout_ratio", 0.2)
        mlflow.log_param("method", "user_based_cf_with_popularity_weighting")
        cases = prepare_eval_cases(matrix, n_users=50, holdout_ratio=0.2)

        cf_precision, _ = evaluate_all(cases, evaluate_user)
        baseline_precision, _ = evaluate_all(cases, popular_baseline_user)

        hybrid_score_fn = make_hybrid_score_fn(weak_tags, problem_cat, problem_tags)
        hybrid_precision, _ = evaluate_all(cases, hybrid_score_fn)

        mlflow.log_metric("cf_precision_at_10", cf_precision)
        mlflow.log_metric("baseline_precision_at_10", baseline_precision)
        mlflow.log_metric("hybrid_precision_at_10", hybrid_precision)

        print(f"CF precision@10: {cf_precision:.3f}")
        print(f"Baseline precision@10: {baseline_precision:.3f}")
        print(f"Hybrid precision@10: {hybrid_precision:.3f}")
