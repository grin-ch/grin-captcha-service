package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/grin-ch/grin-api/api/captcha"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestClient(t *testing.T) {
	conn, err := grpc.Dial("192.168.1.102:8083", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	//创建endpoint，并指明grpc调用的接口名和方法名
	client := captcha.NewCaptchaServiceClient(conn)
	rsp, err := client.GraphCaptcha(context.Background(), &captcha.GraphCaptchaReq{
		Purpose: captcha.Purpose_SIGN_IN,
	})

	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%s\n", rsp.Captcha.Key)
}
