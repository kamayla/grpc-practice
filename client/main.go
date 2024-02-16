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
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		panic(err)
	}

	defer conn.Close()

	client := pb.NewFileServiceClient(conn)

	callUploadAndNotifyProgress(client)
	// callUpload(client)
	// callDownload(client)
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
	stream, err := client.Download(context.Background(), &pb.DownloadRequest{FileName: "sports.txt"})

	if err != nil {
		panic(err)
	}

	for {
		fmt.Println("1")
		res, err := stream.Recv()

		fmt.Println("2")

		if err == io.EOF {
			break
		}

		if err != nil {
			panic(err)
		}

		fmt.Println("3")

		fmt.Println(res.GetData())
		fmt.Println("4")
		fmt.Println(string(res.GetData()))
		fmt.Println("5")
	}
}

func callListFiles(client pb.FileServiceClient) {
	res, err := client.ListFiles(context.Background(), &pb.ListFilesRequest{})

	if err != nil {
		panic(err)
	}

	fmt.Println(res.GetFileNames())
}
