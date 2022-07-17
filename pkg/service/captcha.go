package service

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/grin-ch/grin-api/api/captcha"
	"github.com/grin-ch/grin-captcha-service/pkg/model"
	"github.com/grin-ch/grin-captcha-service/pkg/util"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type captchaService struct {
	pv *model.Provider
}

func NewCaptchaServer(pv *model.Provider) captcha.CaptchaServiceServer {
	return &captchaService{
		pv: pv,
	}
}

// 异步验证码 用于短信/邮箱等
func (s *captchaService) AsyncCode(ctx context.Context, req *captcha.AsyncCodeReq) (*captcha.AsyncCodeRsp, error) {
	return nil, nil
}

// 图形验证码
func (s *captchaService) GraphCaptcha(ctx context.Context, req *captcha.GraphCaptchaReq) (*captcha.GraphCaptchaRsp, error) {
	text := util.GenFormSet(4, util.LowerSet())
	// 设置key,value
	key := md5Str(text)
	val := string(text)
	c := model.Captcha{
		Key:     key,
		Content: val,
		Purpose: int32(req.Purpose),
	}
	if err := s.pv.SetCaptcha(c); err != nil {
		log.Errorf("set captcha err:%s", err.Error())
		return nil, status.Errorf(codes.Internal, "set captcha err")
	}

	imgByte, err := util.NewImg(val)
	if err != nil {
		log.Errorf("new captcha image err:%s", err.Error())
		return nil, status.Errorf(codes.Internal, "new captcha image err")
	}
	base64str := base64.StdEncoding.EncodeToString(imgByte)
	return &captcha.GraphCaptchaRsp{
		Captcha: &captcha.Captcha{
			Key:     key,
			Purpose: req.Purpose,
			Content: base64str,
		},
	}, nil
}

// 验证码校验
func (s *captchaService) Verify(ctx context.Context, req *captcha.VerifyReq) (*captcha.VerifyRsp, error) {
	ok, err := s.pv.Verify(model.Captcha{
		Key:     req.Key,
		Content: req.Value,
		Purpose: int32(req.Purpose),
	})
	if err != nil {
		log.Errorf("verify captcha err:%s", err.Error())
		return nil, status.Errorf(codes.Internal, "verify captcha err")
	}

	return &captcha.VerifyRsp{
		Success: ok,
	}, nil
}

func md5Str(src []byte) string {
	w := md5.New()
	io.WriteString(w, string(src))
	md5str := fmt.Sprintf("%x", w.Sum(nil))
	return md5str
}
