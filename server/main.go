package main

import (
	"bytes"
	"context"
	"fmt"
	"grpc-practice/pb"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedFileServiceServer
}

func (*server) UploadAndNotifyProgress(stream pb.FileService_UploadAndNotifyProgressServer) error {
	fmt.Println("UploadAndNotifyProgress was invoked")
	size := 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		data := req.GetData()
		log.Printf("received %v", data)
		size += len(data)

		res := &pb.UploadAndNotifyProgressResponse{
			Msg: fmt.Sprintf("received %v bytes", size),
		}

		err = stream.Send(res)

		if err != nil {
			return err
		}
	}

}

func (*server) ListFiles(ctx context.Context, req *pb.ListFilesRequest) (*pb.ListFilesResponse, error) {
	fmt.Println("ListFiles was invoked")

	dir := "/Users/ippei/grpc-practice/storage"

	paths, err := ioutil.ReadDir(dir)

	if err != nil {
		return nil, err
	}

	var filenames []string

	for _, path := range paths {
		if !path.IsDir() {
			filenames = append(filenames, path.Name())
		}
	}

	return &pb.ListFilesResponse{FileNames: filenames}, nil
}

func (*server) Download(req *pb.DownloadRequest, stream pb.FileService_DownloadServer) error {
	fmt.Println("Download was invoked")

	filename := req.GetFileName()

	path := "/Users/ippei/grpc-practice/storage/" + filename

	file, err := os.Open(path)

	if err != nil {
		return err
	}

	defer file.Close()

	buf := make([]byte, 5)

	for {
		n, err := file.Read(buf)

		if n == 0 || err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		res := &pb.DownloadResponse{
			Data: buf[:n],
		}

		sendErr := stream.Send(res)

		if sendErr != nil {
			return sendErr
		}

		time.Sleep(5 * time.Second)
	}

	return nil
}

func (*server) Upload(stream pb.FileService_UploadServer) error {
	fmt.Println("Upload was invoked")

	var buffer bytes.Buffer
	for {
		req, err := stream.Recv()

		if err == io.EOF {
			res := &pb.UploadResponse{
				Size: int32(buffer.Len()),
			}
			return stream.SendAndClose(res)
		}

		if err != nil {
			return err
		}

		data := req.GetData()
		log.Printf("Received data : %v", data)
		log.Printf("Received data : %v", string(data))

		buffer.Write(data)
	}
}

func main() {
	lis, err := net.Listen("tcp", "localhost:50051")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterFileServiceServer(s, &server{})

	fmt.Println("Server is running on port :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
