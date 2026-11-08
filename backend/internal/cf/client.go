package cf

import (
	cfpb "atcoder-recomend/internal/pb/cf/v1"
	"context"
	"fmt"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Recommend(ctx context.Context, userHandle string, k int) ([]string, error) {
	addr := os.Getenv("CF_SERVICE_ADDR")
	if addr == "" {
		addr = "localhost:50052" // ローカル開発時のデフォルト
	}
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("接続失敗: %w", err)
	}
	defer conn.Close()

	client := cfpb.NewCFServiceClient(conn)
	response, err := client.Recommend(ctx, &cfpb.RecommendRequest{
		UserHandle: userHandle,
		K:          int32(k),
	})
	if err != nil {
		return nil, fmt.Errorf("推薦結果の取得失敗: %w", err)
	}

	return response.ProblemId, nil
}
