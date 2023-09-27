package main

import (
	"fmt"

	"net/http"

	"github.com/qor/admin"
	"github.com/qor/qor"
	"gitlab.com/mefit/mefit-server/entity"
	"gitlab.com/mefit/mefit-server/utils"
	"gitlab.com/mefit/mefit-server/utils/config"
	"gitlab.com/mefit/mefit-server/utils/initializer"
	"gitlab.com/mefit/mefit-server/utils/log"
)

func main() {
	// Initalize
	defer initializer.Initialize()()
	Admin := admin.New(&admin.AdminConfig{DB: entity.GetDB()})

	// Allow to use Admin to manage User, Product
	user := Admin.AddResource(&entity.User{}, &admin.Config{Menu: []string{"User Management"}})
	user.Meta(&admin.Meta{Name: "Password",
		Type:   "password",
		Valuer: func(interface{}, *qor.Context) interface{} { return "" },
		//TODO: add with auth service
		// Setter: func(record interface{}, metaValue *resource.MetaValue, context *qor.Context) {
		// 	if newPassword := utils.ToString(metaValue.Value); newPassword != "" {
		// 		bcryptPassword, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		// 		record.(*models.User).EncryptedPassword = string(bcryptPassword)
		// 	}
		// },
	})
	proRes := Admin.AddResource(&entity.Profile{}, &admin.Config{Menu: []string{"User Management"}})
	proRes.Meta(&admin.Meta{Name: "Gender", Type: "select_one",
		Config: &admin.SelectOneConfig{
			Collection: [][]string{[]string{"1", "Unknown"}, []string{"2", "Male"}, []string{"3", "Female"}},
		},
		Valuer: func(in interface{}, ctx *qor.Context) interface{} {
			return in.(*entity.Profile).GetGender()
		},
	})
	// Admin.AddResource(&entity.Program{}, &admin.Config{Menu: []string{"Level Management"}})
	wtRes := Admin.AddResource(&entity.WorkoutType{}, &admin.Config{Menu: []string{"Plan Management"}})
	wtRes.Meta(&admin.Meta{Name: "Description", Type: "rich_editor"})
	// Admin.AddResource(&entity.Level{}, &admin.Config{Menu: []string{"Level Management"}})
	// Admin.AddResource(&entity.SubLevel{}, &admin.Config{Menu: []string{"Level Management"}})
	Admin.AddResource(&entity.Product{}, &admin.Config{Menu: []string{"Finance Management"}})
	Admin.AddResource(&entity.Payment{}, &admin.Config{Menu: []string{"Finance Management"}})
	Admin.AddResource(&entity.BazaarPayment{}, &admin.Config{Menu: []string{"Finance Mnagement"}})
	Admin.AddResource(&entity.PurchasedProduct{}, &admin.Config{Menu: []string{"Finance Management"}})
	planRes := Admin.AddResource(&entity.Plan{}, &admin.Config{Menu: []string{"Plan Management"}})
	planRes.Meta(&admin.Meta{Name: "Description", Type: "rich_editor"})
	workoutRes := Admin.AddResource(&entity.Workout{}, &admin.Config{Menu: []string{"Plan Management"}})
	workoutRes.Meta(&admin.Meta{Name: "Description", Type: "rich_editor"})
	workoutRes.Meta(&admin.Meta{Name: "Instruction", Type: "rich_editor"})
	workoutRes.Meta(&admin.Meta{Name: "Duration", Label: "Duration in second"})
	Admin.AddResource(&entity.ExerciseSection{}, &admin.Config{Menu: []string{"Plan Management"}})
	exerciseResource := Admin.AddResource(&entity.Exercise{}, &admin.Config{Menu: []string{"Plan Management"}})
	exerciseResource.Meta(&admin.Meta{Name: "ExerciseType", Type: "select_one",
		Config: &admin.SelectOneConfig{
			Collection: [][]string{[]string{"1", "Duration"}, []string{"2", "Repitation"}},
		}},
	)
	Admin.AddResource(&entity.Class{}, &admin.Config{Menu: []string{"Movement Management"}})
	movementResource := Admin.AddResource(&entity.Movement{}, &admin.Config{Menu: []string{"Movement Management"}})
	movementResource.SearchAttrs("Name")
	movementResource.Meta(&admin.Meta{Name: "Instruction", Type: "rich_editor"})
	movementResource.Meta(&admin.Meta{Name: "Description", Type: "rich_editor"})
	movementResource.Meta(&admin.Meta{Name: "Tips", Type: "rich_editor"})

	// articleRes := Admin.AddResource(&entity.Article{}, &admin.Config{Menu: []string{"Article Management"}})
	// articleRes.Meta(&admin.Meta{Name: "Body", Type: "rich_editor"})
	Admin.AddResource(&entity.Keyword{})
	Admin.AddResource(&entity.MuscleGroup{}, &admin.Config{Menu: []string{"Movement Management"}})

	// initalize an HTTP request multiplexer
	mux := http.NewServeMux()

	// Mount admin interface to mux
	Admin.MountTo("/admin", mux)
	adminPort := config.Config().GetString(utils.KeyPort)
	log.Logger().Infof("Listening on: %s", adminPort)
	panic(http.ListenAndServe(fmt.Sprintf(":%s", adminPort), mux))
}
