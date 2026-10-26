import psycopg2
import os
import pandas as pd
from dotenv import load_dotenv
from scipy.sparse import coo_matrix
import numpy as np
from sklearn.metrics.pairwise import cosine_similarity

USER_NAME = "uwaaaaaaa"

def load_external_solves(conn):
    return pd.read_sql("SELECT user_handle, problem_id FROM external_solves", conn)

def build_matrix(df):
    user_cat = pd.Categorical(df["user_handle"])
    problem_cat = pd.Categorical(df["problem_id"])
    matrix = coo_matrix(
        ([1] * len(df), (user_cat.codes, problem_cat.codes)),
        shape=(len(user_cat.categories), len(problem_cat.categories)),
    ).tocsr()
    return matrix, user_cat, problem_cat

def apply_popularity_weighting(matrix):
    solve_counts = np.asarray(matrix.sum(axis=0)).flatten()
    weights = 1 / np.log1p(solve_counts)
    return matrix.multiply(weights).tocsr()

def compute_similarities(weighted_matrix, user_index):
    return cosine_similarity(weighted_matrix[user_index], weighted_matrix).flatten()

def is_subset_user(matrix, target_index, other_index):
    target_row = matrix[target_index]
    other_row = matrix[other_index]
    overlap = other_row.multiply(target_row).sum()
    return overlap == other_row.sum()

def select_similar_users(matrix, similarities, target_index, top_n=50):
    candidates = [
        i for i in np.argsort(similarities)[::-1]
        if i != target_index and not is_subset_user(matrix, target_index, i)    
    ]
    return candidates[:top_n]

def recommend_problems(matrix, similarities, target_index, similar_indices, top_k=10):
    already_solved = matrix[target_index].toarray().flatten()
    scores = similarities[similar_indices] @ matrix[similar_indices].toarray()
    scores[already_solved > 0] = -1
    top = np.argsort(scores)[::-1][:top_k]
    return top, scores[top] 

if __name__ == "__main__":
    # 今までのメイン処理をここに
    load_dotenv()
    conn = psycopg2.connect(os.environ["DATABASE_URL"])
    cur = conn.cursor()

    query = """
    INSERT INTO external_solves (user_handle, problem_id)
    SELECT 'uwaaaaaaa', problem_id FROM submissions WHERE result = 'AC'
    ON CONFLICT (user_handle, problem_id) DO NOTHING;
    """
    cur.execute(query)
    conn.commit()
    
    df = load_external_solves(conn)
    matrix, user_cat, problem_cat = build_matrix(df)

    weighted_matrix = apply_popularity_weighting(matrix)
    my_index = user_cat.categories.get_loc(USER_NAME)
    similarities = compute_similarities(weighted_matrix, my_index)

    similar_indicies = select_similar_users(matrix, similarities, my_index)
    top, scores = recommend_problems(matrix, similarities, my_index, similar_indicies)

    for idx, score in zip(top, scores):
        print(problem_cat.categories[idx], score)

    conn.close()