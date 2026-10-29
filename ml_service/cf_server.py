import grpc
from cf.v1 import cf_pb2, cf_pb2_grpc
from hybrid_cf import rerank_by_weak_tag, find_weak_tags, load_problem_tags
from train_cf import apply_popularity_weighting, compute_similarities, select_similar_users, recommend_problems, build_matrix, load_external_solves
import psycopg2
import os
import pandas as pd
from dotenv import load_dotenv
from concurrent import futures

MODEL_VERSION = "hybrid_cf"

class CFServicer(cf_pb2_grpc.CFServiceServicer):
    def __init__(self, weighted_matrix, matrix, user_cat, problem_cat, weak_tags, problem_tags):
        self.weighted_matrix = weighted_matrix
        self.matrix = matrix
        self.user_cat = user_cat
        self.problem_cat = problem_cat
        self.weak_tags = weak_tags
        self.problem_tags =problem_tags

    def Recommend(self, request, context):
        my_index = self.user_cat.categories.get_loc(request.user_handle)
        similarities = compute_similarities(self.weighted_matrix, my_index)
        similar_indicies = select_similar_users(self.matrix, similarities, my_index)
        top, scores = recommend_problems(self.weighted_matrix, similarities, my_index, similar_indicies, request.k)
        reranked = rerank_by_weak_tag(top, self.problem_cat, self.problem_tags, self.weak_tags)

        problem_ids = [self.problem_cat.categories[idx] for idx in reranked]
        return cf_pb2.RecommendResponse(
            problem_id = problem_ids,
            model_version = MODEL_VERSION
        )

def serve():
    load_dotenv()
    conn = psycopg2.connect(os.environ["DATABASE_URL"])
    df = load_external_solves(conn)
    matrix, user_cat, problem_cat = build_matrix(df)
    weighted_matrix = apply_popularity_weighting(matrix)
    weak_tags = find_weak_tags(conn)
    problem_tags = load_problem_tags(conn)

    conn.close()
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=4))
    cf_pb2_grpc.add_CFServiceServicer_to_server(
        CFServicer(weighted_matrix, matrix, user_cat, problem_cat, weak_tags, problem_tags), server
    )
    server.add_insecure_port('[::]:50052')
    server.start()
    server.wait_for_termination()

if __name__ == "__main__":
    serve()