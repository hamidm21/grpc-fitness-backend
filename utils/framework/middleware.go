package framework

// const APIType = "api_type"

// func recovery() gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		defer func() {
// 			if err := recover(); err != nil {
// 				log.Logger().WithFields(logrus.Fields{
// 					"err":   err,
// 					"stack": string(debug.Stack()),
// 				}).Error()
// 			}
// 		}()
// 		ctx.Next()
// 	}
// }

// func api() gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		ctx.Set(APIType, V1)
// 	}
// }
