package main

import (
	// "net"

	"encoding/json"
	"fmt"
	"time"

	// "string"
	"log"

	pb "gitlab.com/mefit/mefit-api/proto"
	"gitlab.com/mefit/mefit-server/services/storage/mongodb"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func main() {
	// defer initializer.Initialize()()

	opts := []grpc.DialOption{

		// oauth.NewOauthAccess requires the configuration of transport credentials.
		// grpc.WithTransportCredentials(
		// 	//TODO skip this for now
		// 	credentials.NewTLS(&tls.Config{InsecureSkipVerify: true}),
		// ),
		grpc.WithInsecure(),
	}
	conn, err := grpc.Dial("127.0.0.1:8443", opts...)
	// conn, err := grpc.Dial("grpc.fitex.app.yottab.io:443", opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := pb.NewMefitClient(conn)

	//Finish some workout here
	// FinishWorkout(client, &pb.FeedbackReq{Id: 3, Rate: 2, Difficulty: 2})
	// GetArticles(client, &pb.ListReq{
	// 	Page: 1,
	// })
	// ArticleInfo(client, &pb.FindByMongoIdReq{
	// 	Id: "5d46a3ddacfebc455fbcb7c6",
	// })
	// GetClass(client, &pb.Empty{})
	// GetClassMovements(client, &pb.FindByIdAndPage{
	// 	Id: 1,
	// 	Page:1,
	// })
	// GetProgram(client, &pb.Empty{

	// })
	AnonySignUp(client, &pb.AnonyReq{})
	// GetProducts(client, &pb.Empty{

	// })
	// PaymentRequest(client, &pb.PayReq{
	// 	ProductID: 1,
	// 	IsBazaar: true,
	// })
	// GetPlans(client, &pb.Empty{})
	// req := pb.ListReq{
	// 	Page: 0,
	// }
	// GetArticles(client, &req)
	// AnonySignUp(client, &pb.AnonyReq{})
}

func GetArticlesInternal() {
	articles, err := mongodb.GetArticle("5d46a2f3acfebc0f40bcb7c4")
	if err != nil {
		log.Fatal(err)
	}
	chamar, err := json.Marshal(articles)
	fmt.Print(string(chamar))
}

func GetArticles(client pb.MefitClient, lr *pb.ListReq) {
	md := metadata.New(map[string]string{"x-token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOjE5OCwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.NKYmBMpH8xdH-xgcJQz9DUeN051QWkFuKD5JiGuAnpY"})
	ctx1 := metadata.NewOutgoingContext(context.Background(), md)
	articles, err := client.GetArticles(ctx1, lr)
	if err != nil {
		log.Fatal(err)
	}
	chamar, err := json.Marshal(articles)
	fmt.Print(string(chamar))
}

func GetCurrentPlan(client pb.MefitClient, lr *pb.Empty) {
	md := metadata.New(map[string]string{"x-token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTI4OTU5MjAsImlhdCI6MTU2Njk3NTkyMCwic3ViIjoxOTF9.F5zi7ALXpTzZQ_fS2sunaqAq6aI7oMEL7xBCF_bR9x4"})
	ctx1 := metadata.NewOutgoingContext(context.Background(), md)
	articles, err := client.GetCurrentPlan(ctx1, lr)
	if err != nil {
		log.Fatal(err)
	}
	chamar, err := json.Marshal(articles)
	fmt.Print(string(chamar))
}

func ArticleInfo(client pb.MefitClient, fbi *pb.FindByMongoIdReq) {
	md := metadata.New(map[string]string{"x-token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTIzMTIyNzQsImlhdCI6MTU2NjM5MjI3NCwic3ViIjoxNjJ9.ujLLimeuaDIe02zyZbrCQtXtUBN4te9MHvuvPdPEkOU"})
	ctx1 := metadata.NewOutgoingContext(context.Background(), md)
	article, err := client.ArticleInfo(ctx1, fbi)
	if err != nil {
		log.Fatal(err)
	}
	chamar, err := json.Marshal(article)
	fmt.Println(string(chamar))
}

func GetClass(client pb.MefitClient, em *pb.Empty) {
	md := metadata.New(map[string]string{"x-token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOjM1LCJpYXQiOjE1NjYzNzYwMjcsImV4cCI6MTU2NjM3OTY0NCwianRpIjoiNWIyYjQ4ZjUtOGM2NC00YmQ1LWFlNDctZTE4NGE5YzgwZDE4In0.qMAnkUYTHD_rZINU2n-Qv9J6v_5xYvciHLw1BepKYTQ"})
	ctx1 := metadata.NewOutgoingContext(context.Background(), md)
	ctx, error := context.WithTimeout(ctx1, 30*time.Second)
	defer error()
	class, err := client.GetClasses(ctx, em)
	if err != nil {
		log.Fatal(err)
	}
	shit, err := json.Marshal(class)
	fmt.Println(string(shit))
}

func GetClassMovements(client pb.MefitClient, id *pb.FindByIdAndPage) {
	ctx, err := context.WithTimeout(context.Background(), 10*time.Second)
	defer err()
	movements, error := client.GetClassMovements(ctx, id)
	if error != nil {
		log.Fatal(err)
	}
	shit, error := json.Marshal(movements)
	fmt.Println(string(shit))
}

func GetProgram(client pb.MefitClient, em *pb.Empty) {
	// ctx, err := context.WithTimeout(context.Background(), 10*time.Second)
	// defer err()
	// levels, error := client.GetProgram(ctx, em)
	// if error != nil {
	// 	log.Fatal(error)
	// }
	// shit, error := json.Marshal(levels)
	// fmt.Println(string(shit))
}

func GetProducts(client pb.MefitClient, em *pb.Empty) {
	ctx, err := context.WithTimeout(context.Background(), 10*time.Second)
	defer err()
	levels, error := client.GetProducts(ctx, em)
	if error != nil {
		log.Fatal(error)
	}
	shit, error := json.Marshal(levels)
	fmt.Println(string(shit))
}

func PaymentRequest(client pb.MefitClient, em *pb.PayReq) {
	md := metadata.New(map[string]string{"x-token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTI4MDI1MjYsImlhdCI6MTU2Njg4MjUyNiwic3ViIjoxODF9.Rbp5Xe0fz0iLLlysuRDtgtpXFbkbsJT9nLlH7EzCZaE"})
	ctx1 := metadata.NewOutgoingContext(context.Background(), md)
	ctx, error := context.WithTimeout(ctx1, 30*time.Second)
	defer error()
	class, err := client.PaymentRequest(ctx, em)
	if err != nil {
		log.Fatal(err)
	}
	shit, err := json.Marshal(class)
	fmt.Println(string(shit))
}

func GetPlans(client pb.MefitClient, em *pb.Empty) {
	md := metadata.New(map[string]string{"x-token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTI4MDI1MjYsImlhdCI6MTU2Njg4MjUyNiwic3ViIjoxODF9.Rbp5Xe0fz0iLLlysuRDtgtpXFbkbsJT9nLlH7EzCZaE"})
	ctx1 := metadata.NewOutgoingContext(context.Background(), md)
	ctx, error := context.WithTimeout(ctx1, 30*time.Second)
	defer error()
	class, err := client.GetPlans(ctx, em)
	if err != nil {
		log.Fatal(err)
	}
	shit, err := json.Marshal(class)
	fmt.Println(string(shit))
}

func AnonySignUp(client pb.MefitClient, em *pb.AnonyReq) {
	ctx, err := context.WithTimeout(context.Background(), 10*time.Second)
	defer err()
	levels, error := client.AnonySignUp(ctx, em)
	if error != nil {
		log.Fatal(error)
	}
	shit, error := json.Marshal(levels)
	fmt.Println(string(shit))
}

func FinishWorkout(client pb.MefitClient, em *pb.FeedbackReq) {
	md := metadata.New(map[string]string{"x-token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOjE5NiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.w4furV0NF3kXKZydKer3M4HX4QfkKCAZViHQch-6uVo"})
	ctx1 := metadata.NewOutgoingContext(context.Background(), md)
	_, err := client.FinishWorkout(ctx1, em)
	if err != nil {
		log.Fatal(err)
	}
}
