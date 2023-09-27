package entity

// import (
// 	"github.com/jinzhu/gorm"
// 	"gitlab.com/mefit/mefit-server/utils"
// 	"gitlab.com/mefit/mefit-server/utils/log"
// )

func (instance Payment) NamedEntity(name string) Entity {
	// instance.Name = name
	return instance
}
func (instance Payment) IDEntity(Id uint) Entity {
	instance.ID = Id
	return instance
}
func (instance Payment) UserEntity(usrID uint) Entity {
	return instance
}


func (instance Product) NamedEntity(name string) Entity {
	instance.Name = name
	return instance
}
func (instance Product) IDEntity(Id uint) Entity {
	instance.ID = Id
	return instance
}
func (instance Product) UserEntity(usrID uint) Entity {
	return instance
}


func (instance BazaarPayment) NamedEntity(name string) Entity {
	return instance
}
func (instance BazaarPayment) IDEntity(Id uint) Entity {
	instance.ID = Id
	return instance
}
func (instance BazaarPayment) UserEntity(usrID uint) Entity {
	return instance
}


func (instance PurchasedProduct) NamedEntity(name string) Entity {
	return instance
}
func (instance PurchasedProduct) IDEntity(Id uint) Entity {
	instance.ID = Id
	return instance
}
func (instance PurchasedProduct) UserEntity(usrID uint) Entity {
	return instance
}