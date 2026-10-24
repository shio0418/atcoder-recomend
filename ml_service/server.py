import grpc
from embedding.v1 import embedding_pb2
from embedding.v1 import embedding_pb2_grpc
from sentence_transformers import SentenceTransformer
from concurrent import futures

MODEL_VERSION = "intfloat/multilingual-e5-small"

class EmbeddingServicer(embedding_pb2_grpc.EmbeddingServiceServicer):
    def __init__(self, model):
        self.model = model

    def Embed(self, request_iterator, context):
        for request in request_iterator:
            vector = self.model.encode(request.text)
            yield embedding_pb2.EmbedResponse(
                problem_id = request.problem_id,
                values = vector.tolist(),
                model_version=MODEL_VERSION,
            )

def serve():
    model = SentenceTransformer(MODEL_VERSION)

    server = grpc.server(futures.ThreadPoolExecutor(max_workers=4))
    embedding_pb2_grpc.add_EmbeddingServiceServicer_to_server(
        EmbeddingServicer(model), server
    )
    server.add_insecure_port('[::]:5051')
    server.start()
    server.wait_for_termination()

if __name__ == "__main__":
    serve()