package cmd

import (
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func recoveryFunc(p interface{}) (err error) {
	log.Errorf("panic triggered: %v", p)
	return status.Errorf(codes.Unknown, "unknow error")
}
