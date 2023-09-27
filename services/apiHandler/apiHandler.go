package apiHandler

import (
	"fmt"
	"gitlab.com/mefit/mefit-server/utils/config"
	"github.com/kataras/iris"
	// "gitlab.com/mefit/mefit-server/entity"
	"gitlab.com/mefit/mefit-server/services/payment"
	// "github.com/kataras/iris/context"
	"gitlab.com/mefit/mefit-server/utils"
	// "gitlab.com/mefit/mefit-server/controller"
)

func StartApiServer() {
	app := iris.Default();

	type Body struct {
		productID  	uint `json:"productId"`
		userID		uint `json:"userId"`
	}

	app.Post("/payment", func (ctx iris.Context) {
		body := Body{}
		ctx.ReadJSON(&body)
		fmt.Println(body)
		// product := entity.Product{}
		// product.ID = uint(body.productID)
		callbackURL := fmt.Sprintf(":%s", config.Config().GetString(utils.CallbackURL))
		payment.NewPaymentRequest(int(10000), callbackURL)
		// user := entity.User{}
		// user.ID = uint(body.userID)
		// if err := entity.SimpleCrud(user).Get(&user); err != nil {
		// 	ctx.Writef("payment request failed")
		// }
		ctx.Writef("successful %s" , body.productID)
	})

	// app.Post("/paymentResult{Authority:string}", controller)


	Port := fmt.Sprintf(":%s", config.Config().GetString(utils.APIPort))
	fmt.Println("api port --------> %s", Port)
	app.Run(iris.Addr(Port))
}