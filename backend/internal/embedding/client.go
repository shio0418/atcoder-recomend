package embedding

import (
	embedpb "atcoder-recomend/internal/pb/embedding/v1"
	"context"
	"fmt"
	"io"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func EmbedBatch(ctx context.Context, texts map[string]string) (map[string]*embedpb.EmbedResponse, error) {
	addr := os.Getenv("EMBEDDING_SERVICE_ADDR")
	if addr == "" {
		addr = "localhost:5051" // ローカル開発時のデフォルト
	}
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, fmt.Errorf("接続失敗: %w", err)
	}
	defer conn.Close()
	client := embedpb.NewEmbeddingServiceClient(conn)
	stream, err := client.Embed(ctx)
	if err != nil {
		return nil, fmt.Errorf("接続失敗: %w", err)
	}

	go func() {
		for problemID, text := range texts {
			if err := stream.Send(&embedpb.EmbedRequest{ProblemId: problemID, Text: text}); err != nil {
				return
			}
		}
		stream.CloseSend()
	}()

	results := make(map[string]*embedpb.EmbedResponse)
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		results[resp.ProblemId] = resp
	}
	return results, nil
}
