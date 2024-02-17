package main

import (
	"context"
	"fmt"
	"grpc-practice/pb"
	"io"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		panic(err)
	}

	defer conn.Close()

	client := pb.NewFileServiceClient(conn)

	// callUploadAndNotifyProgress(client)
	// callUpload(client)
	callDownload(client)
	// callListFiles(client)

}

func callUploadAndNotifyProgress(client pb.FileServiceClient) {
	filename := "sports.txt"
	path := "/Users/ippei/grpc-practice/storage/" + filename

	file, err := os.Open(path)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	stream, err := client.UploadAndNotifyProgress(context.Background())

	if err != nil {
		panic(err)
	}

	buf := make([]byte, 5)
	go func() {
		for {
			n, err := file.Read(buf)

			if n == 0 || err == io.EOF {
				break
			}

			if err != nil {
				panic(err)
			}

			req := &pb.UploadAndNotifyProgressRequest{
				Data: buf[:n],
			}

			err = stream.Send(req)
			if err != nil {
				log.Fatalln(err)
			}
			time.Sleep(2 * time.Second)
		}
		err := stream.CloseSend()

		if err != nil {
			panic(err)
		}
	}()

	ch := make(chan struct{})
	go func() {
		for {
			res, err := stream.Recv()

			if err == io.EOF {
				break
			}

			if err != nil {
				panic(err)
			}

			fmt.Println(res.GetMsg())
		}
		close(ch)
	}()
	<-ch

}

func callUpload(client pb.FileServiceClient) {
	filename := "sports.txt"
	path := "/Users/ippei/grpc-practice/storage/" + filename

	file, err := os.Open(path)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	stream, err := client.Upload(context.Background())

	if err != nil {
		panic(err)
	}

	buf := make([]byte, 5)

	for {
		n, err := file.Read(buf)

		if n == 0 || err == io.EOF {
			break
		}

		if err != nil {
			panic(err)
		}

		err = stream.Send(&pb.UploadRequest{Data: buf})
		if err != nil {
			log.Fatalln(err)
		}
		time.Sleep(5 * time.Second)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		panic(err)
	}

	fmt.Println(res.GetSize())
}

func callDownload(client pb.FileServiceClient) {
	stream, err := client.Download(context.Background(), &pb.DownloadRequest{FileName: "sportss.txt"})

	if err != nil {
		log.Fatalln(err)
	}

	for {
		res, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			resErr, ok := status.FromError(err)
			if ok {
				if resErr.Code() == codes.NotFound {
					log.Fatalf("Error Code: %v, Error Message: %v", resErr.Code(), resErr.Message())
				} else {
					log.Fatalln("unknown error")
				}
			} else {
				log.Fatalln(err)
			}
		}

		fmt.Println(res.GetData())
		fmt.Println(string(res.GetData()))
	}
}

func callListFiles(client pb.FileServiceClient) {
	md := metadata.New(map[string]string{"authorization": "Bearer ttest-token"})

	ctx := metadata.NewOutgoingContext(context.Background(), md)

	res, err := client.ListFiles(ctx, &pb.ListFilesRequest{})

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(res.GetFileNames())
}
