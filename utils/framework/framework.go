package framework

// var (
// 	routes []Router
// )

// type Type string

// const (
// 	Template Type = "template"
// 	V1            = "v1"
// )

// func Run() {
// 	engine := gin.New()
// 	engine.Use(recovery())
// 	apiGrpc := engine.Group(V1, api())

// 	for i := range routes {
// 		// for template
// 		routes[i].Route(engine.Group(""))

// 		// for api v1
// 		routes[i].Route(apiGrpc)
// 	}

// 	go func() {
// 		err := engine.Run(":3412")
// 		assert.Nil(err)
// 	}()

// }
