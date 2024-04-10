package interceptor

import (
	"context"
	"log"
	"strings"

	"google.golang.org/grpc"
)

func UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func (ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log.Println("===unary interceptor===")

		// リクエスト前処理
		log.Printf("===req value: %v===", req)

		res, err := handler(ctx, req)

		// 後処理
		log.Printf("===res value: %v===", res)

		return res, err
	}
}

func StreamInterceptor() grpc.StreamServerInterceptor {
	return func (srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// リクエスト前処理
		log.Println("===stream 始まったンゴ===")
		// rpc名表示
		rpcName := strings.Split(info.FullMethod, "/")[2]
		log.Printf("===rpc name: %v===", rpcName)

		// recv, sendの前処理追加には、streamWrapperを使う
		err := handler(srv, &streamWrapper{stream})

		// 後処理
		log.Println("===stream 終わったンゴ===")

		return err
	}
}

type streamWrapper struct {
	grpc.ServerStream
}

func (s *streamWrapper) RecvMsg(m interface{}) error {
	err := s.ServerStream.RecvMsg(m)
	// recvの内容をログに出力
	log.Println("===recvしたンゴ===")
	log.Printf("===recv value: %v===", m)

	return err
}

func (s *streamWrapper) SendMsg(m interface{}) error {
	log.Println("===sendしたンゴ===")
	log.Printf("===send value: %v===", m)
	return s.ServerStream.SendMsg(m)
}
