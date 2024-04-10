package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"strings"

	pb "chatdemo/chatdemo"
	interceptor "chatdemo/interceptor"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type server struct {
	pb.UnimplementedChatDemoServiceServer
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.UnaryInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			interceptor.StreamInterceptor(),
		),
	)
	pb.RegisterChatDemoServiceServer(s, &server{})
	reflection.Register(s)

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{
		Message: "Hello " + in.Name,
	}, nil
}


// しりとり用RPC
func (s *server) WordChainChat(stream pb.ChatDemoService_WordChainChatServer) error {
	wordMap, err := genWordMap()
	if err != nil {
		log.Printf("err: %v\n", err)
		return err
	}

	usedWords := make(map[string]bool)
	befPrefix := ""

	for {
		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return nil
		}

		if err != nil {
			log.Printf("err: %v\n", err)
			return err
		}

		// しりとりのルールに従っているか確認
		if befPrefix != "" {
			if string([]rune(req.GetWord())[0]) != befPrefix {
				if err := stream.Send(&pb.WordChain{
					Result: pb.Result_RESULT_LOSE,
					Message: "しりとりになってません、遊戯ボーイ、ユーの負けでーす！",
				}); err != nil {
					log.Printf("err: %v\n", err)
					return err
				}

				return nil
			}
		}

		// 単語の末尾に「ん」が含まれていたらResult_LOSEを返す
		if strings.HasSuffix(req.GetWord(), "ん") {
			if err := stream.Send(&pb.WordChain{
				Result: pb.Result_RESULT_LOSE,
				Message: "「ん」がつきました。ユーの負けでーす！",
			}); err != nil {
				log.Printf("err: %v\n", err)
				return err
			}

			return nil
		}

		// usedMapに単語があるか確認
		if ok := usedWords[req.GetWord()]; ok {
			if err := stream.Send(&pb.WordChain{
				Result: pb.Result_RESULT_LOSE,
				Message: "その単語はもう使われています、ユーの負けでーす！",
			}); err != nil {
				log.Printf("err: %v\n", err)
				return err
			}

			return nil
		}

		// usedMapに単語を追加
		usedWords[req.GetWord()] = true

		// リクエストの末尾の単語を取得
		lastWord := string([]rune(req.GetWord())[len([]rune(req.GetWord()))-1:])

		// ち、く、びの場合無条件で降参する
		if lastWord == "ち" || lastWord == "く" || lastWord == "び" {
			if err := stream.Send(&pb.WordChain{
				Result: pb.Result_RESULT_WIN,
				Message: fmt.Sprintf("「%v」には抗えません、ユーの勝ちでーす！", lastWord),
			}); err != nil {
				log.Printf("err: %v\n", err)
				return err
			}

			return nil
		}

		// 単語のマップから末尾の単語をキーにしてランダムに値を取得
		words, ok := wordMap[lastWord]
		if !ok {
			if err := stream.Send(&pb.WordChain{
				Result: pb.Result_RESULT_LOSE,
				Message: "五十音表に含まれない単語です！ユーの負けでーす！",
			}); err != nil {
				log.Printf("err: %v\n", err)
				return err
			}

			return nil
		}

		resWord := words[rand.Intn(len(words))]

		// 使用済み単語のマップに含まれている or 単語の末尾に「ん」が含まれていたらResult_WINを返す
		if ok := usedWords[resWord]; ok || strings.HasSuffix(resWord, "ん") {
			if err := stream.Send(&pb.WordChain{
				Result: pb.Result_RESULT_WIN,
				Word: resWord,
				Message: "返せる単語がありません！ユーの勝ちでーす",
			}); err != nil {
				log.Printf("err: %v\n", err)
				return err
			}

			return nil
		}

		// usedMapに単語を追加
		usedWords[resWord] = true
		befPrefix = string([]rune(resWord)[len([]rune(resWord))-1:])

		// レスポンスを返す
		if err := stream.Send(&pb.WordChain{
			Word: resWord,
		}); err != nil {
			log.Printf("err: %v\n", err)
			return err
		}
	}
}

func genWordMap() (map[string][]string, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	return viper.GetStringMapStringSlice("words"), nil
}
