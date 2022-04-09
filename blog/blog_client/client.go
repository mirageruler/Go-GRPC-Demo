package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/mirageruler/grpc-go-course/blog/blogpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	fmt.Println("Blog Client")
	fmt.Println("----------------------------------------------------------------------------------------------------------------------------------------")

	opts := grpc.WithTransportCredentials(insecure.NewCredentials())

	conn, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	c := blogpb.NewBlogServiceClient(conn)

	// create a blog
	fmt.Println("Creating the blog...")

	blog := &blogpb.Blog{
		AuthorId: "Khoi",
		Title:    "My first blog",
		Content:  "Content of the first blog",
	}

	createBlogRes, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}

	fmt.Printf("Blog has been created: %v\n", createBlogRes)
	blogId := createBlogRes.GetBlog().GetId()

	// ----------------------------------------------------------------------------------------------------
	// read a blog
	fmt.Println("Reading the blog...")

	// not found
	if _, err = c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: "5bdc29e661b75adcac496cf4"}); err != nil {
		fmt.Printf("Error happened while reading: %v\n", err)
	}

	// found
	readBlogReq := &blogpb.ReadBlogRequest{BlogId: blogId}
	readBlogRes, err := c.ReadBlog(context.Background(), readBlogReq)
	if err != nil {
		fmt.Printf("Error happened while reading: %v\n", err)
	}

	fmt.Printf("Blog was read: %v\n", readBlogRes)

	// ----------------------------------------------------------------------------------------------------
	// update a blog
	fmt.Println("Updating the blog...")

	newBlog := &blogpb.Blog{
		Id:       blogId,
		AuthorId: "Changed Author",
		Title:    "My First Blog (edited)",
		Content:  "Content of my first blog, with new additions",
	}
	updatedBlogRes, err := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{Blog: newBlog})
	if err != nil {
		fmt.Printf("Error happened while updating: %v\n", err)
	}

	fmt.Printf("Updated blog : %v\n", updatedBlogRes)

	// ----------------------------------------------------------------------------------------------------
	// delete a blog
	deletedRes, err := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{BlogId: blogId})
	if err != nil {
		fmt.Printf("Error happened while deleting: %v\n", err)
	}

	fmt.Printf("Blog was deleted: %v \n", deletedRes)

	// ----------------------------------------------------------------------------------------------------
	// listing blogs
	stream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
	if err != nil {
		log.Fatalf("error while calling ListBlog RPC: %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Something happened: %v", err)
		}
		fmt.Println(res.GetBlog())
	}
}
