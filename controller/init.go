package controller

import (
	pb "gitlab.com/mefit/mefit-api/proto"
	"gitlab.com/mefit/mefit-server/utils"
	"gitlab.com/mefit/mefit-server/utils/config"
)

//Controller is used to implement helloworld.GreeterServer.
type Controller struct {
	pb.MefitServer
	// pb_admin.AdminServer
}

var secretKey string

func init() {
	secretKey = config.Config().GetString(utils.KeySecretKey)
}
