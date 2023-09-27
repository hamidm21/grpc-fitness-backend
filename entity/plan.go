package entity

// import (
// 	"github.com/jinzhu/gorm"
// 	"gitlab.com/mefit/mefit-server/utils"
// 	"gitlab.com/mefit/mefit-server/utils/log"
// )

func (instance Plan) NamedEntity(name string) Entity {
	instance.Name = name
	return instance
}
func (instance Plan) IDEntity(Id uint) Entity {
	instance.ID = Id
	return instance
}
func (instance Plan) UserEntity(usrID uint) Entity {
	return instance
}
