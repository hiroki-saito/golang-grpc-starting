package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	pb "app/pkg/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewMyServer() *myServer {
	return &myServer{}
}

func main() {
	// 1. 8080番portのLisnterを作成
	port := 8080
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	// 2. gRPCサーバーを作成
	s := grpc.NewServer()

	// 3. gRPCサーバーにGreetingServiceを登録
	myServer := NewMyServer()
	pb.RegisterGreetingServiceServer(s, myServer)
	pb.RegisterCalculateServiceServer(s, myServer)

	// 4. サーバーリフレクションの設定
	reflection.Register(s)

	// 5. 作成したgRPCサーバーを、8080番ポートで稼働させる
	go func() {
		log.Printf("start gRPC server port: %v", port)
		s.Serve(listener)
	}()

	// 6.Ctrl+Cが入力されたらGraceful shutdownされるようにする
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("stopping gRPC server...")
	s.GracefulStop()
}

type myServer struct {
	pb.UnimplementedGreetingServiceServer
	pb.UnimplementedCalculateServiceServer
}

func (s *myServer) Hello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	// リクエストからnameフィールドを取り出して
	// "Hello, [名前]!"というレスポンスを返す
	return &pb.HelloResponse{
		Message: fmt.Sprintf("Hello, %s!", req.GetName()),
	}, nil
}

func (s *myServer) HelloServerStream(req *pb.HelloRequest, stream pb.GreetingService_HelloServerStreamServer) error {
	resCount := 5
	for i := 0; i < resCount; i++ {
		if err := stream.Send(&pb.HelloResponse{
			Message: fmt.Sprintf("[%d] Hello, %s!", i, req.GetName()),
		}); err != nil {
			return err
		}
		time.Sleep(time.Second * 1)
	}
	// return文でメソッドを終了させる=ストリームの終わり
	return nil
}

func (s *myServer) HelloClientStream(stream pb.GreetingService_HelloClientStreamServer) error {
	nameList := make([]string, 0)
	for {
		// streamのRecvメソッドを呼び出してリクエスト内容を取得する
		req, err := stream.Recv()
		// エラーがEOFなら正常終了
		if errors.Is(err, io.EOF) {
			message := fmt.Sprintf("Hello, %v!", nameList)
			return stream.SendAndClose(&pb.HelloResponse{
				Message: message,
			})
		}
		if err != nil {
			return err
		}
		nameList = append(nameList, req.GetName())
	}
}

func (s *myServer) HelloBiStreams(stream pb.GreetingService_HelloBiStreamsServer) error {
	for {
		req, err := stream.Recv()
		fmt.Println("HelloBiStreams")
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
		message := fmt.Sprintf("Hello, %v!", req.GetName())
		if err := stream.Send(&pb.HelloResponse{
			Message: message,
		}); err != nil {
			return err
		}
	}
}

func (s *myServer) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	return &pb.AddResponse{
		Value: req.GetA() + req.GetB(),
	}, nil
}

func (s *myServer) Subtraction(ctx context.Context, req *pb.SubtractionRequest) (*pb.SubtractionResponse, error) {
	return &pb.SubtractionResponse{
		Value: req.GetA() - req.GetB(),
	}, nil
}
