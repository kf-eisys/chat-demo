package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	pb "chatdemo/chatdemo"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	port = flag.Int("port", 50051, "The server port")

	reply = []string{
		"それはズッコケましたね！笑",
		"本当に？！それはちょっとユニークですね。",
		"それはクレイジーな話ですね、信じられない！",
		"うわー、それはまさに天才的な発想です！",
	}
	ngWord = []string{"ちんちん", "ちんこ", "ちんぽこ", "ぽこちん"}
)

type server struct {
	pb.UnimplementedChatServiceServer
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			unaryInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			streamInterceptor(),
		),
	)
	pb.RegisterChatServiceServer(s, &server{}) // 何してるの？
	reflection.Register(s)

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func unaryInterceptor() grpc.UnaryServerInterceptor {
	return func (ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log.Println("===unary interceptor===")

		// rpc名表示
		rpcName := strings.Split(info.FullMethod, "/")[2]
		log.Printf("===req name: %v===", rpcName)

		// リクエスト表示
		log.Printf("===req value: %v===", req)

		return handler(ctx, req)
	}
}

func streamInterceptor() grpc.StreamServerInterceptor {
	return func (srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		log.Println("===stream interceptor===")

		// rpc名表示
		rpcName := strings.Split(info.FullMethod, "/")[2]
		log.Printf("===req name: %v===", rpcName)

		// リクエスト表示
		// log.Printf("req value: %v\n", stream)

		return handler(srv, stream)
	}
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{
		Message: "Hello " + in.Name,
	}, nil
}

func (s *server) SendMessage(stream pb.ChatService_SendMessageServer) error {
	for {
		req, err := stream.Recv()
		// チャットだとこれで良いか要確認
		if errors.Is(err, io.EOF) {
			return nil
		}

		if err != nil {
			log.Printf("err: %v\n", err)
			return err
		}

		log.Printf("req: %v\n", req)

		if req.GetMessage() == "さようなら" {
			// 終端のメッセージを送信して終了
			if err := stream.Send(&pb.SendMessageResponse{Message: "さようなら、今までありがとう"}); err != nil {
				log.Printf("err: %v\n", err)
				return err
			}

			return nil
		}

		resp := fmt.Sprintf("Received: %v", req.GetMessage())
		if err := stream.Send(&pb.SendMessageResponse{Message: resp}); err != nil {
			log.Printf("err: %v\n", err)
			return err
		}

		// var resp string

		// // レスポンスを返す
		// defer func() {
		// 	if err := stream.Send(&pb.SendMessageResponse{Message: resp}); err != nil {
		// 		log.Printf("err: %v\n", err)
		// 	}
		// }()

		// if initFlg {
		// 	resp = "こんにちは、私はチャットボットです。終了する場合は「さようなら」と言ってください"
		// 	initFlg = false
		// 	continue
		// }

		// userWord := req.GetMessage()
		// if lo.Contains(ngWord, userWord) {
		// 	resp = "下ネタを言う人とは話しません"
		// 	return nil
		// }

		// if userWord == "さようなら" {
		// 	resp = "さようなら、今までありがとう"
		// 	return nil
		// }

		// resp = reply[rand.Intn(len(reply))]
	}
}
