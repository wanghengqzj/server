package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	blobpb "server/blob/api/gen/v1"
)

//func main() {
//	// 将 examplebucket-1250000000 和 COS_REGION 修改为用户真实的信息
//	// 存储桶名称，由 bucketname-appid 组成，appid 必须填入，可以在 COS 控制台查看存储桶名称。https://console.cloud.tencent.com/cos5/bucket
//	// COS_REGION 可以在控制台查看，https://console.cloud.tencent.com/cos5/bucket, 关于地域的详情见 https://cloud.tencent.com/document/product/436/6224
//	u, err := url.Parse("https://coolcarlearn-1317438012.cos.ap-beijing.myqcloud.com")
//	if err != nil {
//		panic(err)
//	}
//	// 用于 Get Service 查询，默认全地域 service.cos.myqcloud.com
//	su, _ := url.Parse("https://cos.COS_REGION.myqcloud.com")
//	b := &cos.BaseURL{BucketURL: u, ServiceURL: su}
//	// 1.永久密钥
//	secID := "AKIDHVDmDQPDtSyVRXju12tdhcMvJXYMyQoK"
//	secKEY := "XvAKSi087xfGkMsIZlprofO5e2JRYYN6"
//	client := cos.NewClient(b, &http.Client{
//		Transport: &cos.AuthorizationTransport{
//			SecretID:  os.Getenv(secID),  // 用户的 SecretId，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参考 https://cloud.tencent.com/document/product/598/37140
//			SecretKey: os.Getenv(secKEY), // 用户的 SecretKey，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参考 https://cloud.tencent.com/document/product/598/37140
//		},
//	})
//
//	name := "96f8b6d852e3b71b286d711c55f10b6.jpg"
//
//	//获取预签名 URL
//	presignedURL, err := client.Object.GetPresignedURL(context.Background(),
//		http.MethodGet, name, secID, secKEY, 40*time.Second, nil)
//	if err != nil {
//		panic(err)
//	}
//
//	//name := "abj.jpg"
//	//presignedURL, err := client.Object.GetPresignedURL(context.Background(),
//	//	http.MethodPut, name, secID, secKEY, 1000*time.Second, nil)
//	//if err != nil {
//	//	panic(err)
//	//}
//	fmt.Println(presignedURL)
//}

func main() {
	conn, err := grpc.Dial("localhost:8083", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	fmt.Println("success")
	c := blobpb.NewBlobServiceClient(conn)
	ctx := context.Background()
	//res, err := c.CreateBlob(ctx, &blobpb.CreateBlobRequest{
	//	AccountId:           "account2",
	//	UploadUrlTimeoutSec: 1000,
	//})

	//res, err := c.GetBlob(ctx, &blobpb.GetBlobRequest{Id: "64632fa3ceab9c7571ff409e"})

	res, err := c.GetBlobURL(ctx, &blobpb.GetBlobURLRequest{
		Id:         "64632fa3ceab9c7571ff409e",
		TimeoutSec: 100,
	})

	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", res)
}
