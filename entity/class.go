package entity

import "gitlab.com/mefit/mefit-server/services/asset"

// import (
// 	"github.com/jinzhu/gorm"
// 	"gitlab.com/mefit/mefit-server/utils"
// 	"gitlab.com/mefit/mefit-server/utils/log"
// )

func (instance Class) NamedEntity(name string) Entity {
	instance.Name = name
	return instance
}
func (instance Class) IDEntity(Id uint) Entity {
	instance.ID = Id
	return instance
}
func (instance Class) UserEntity(usrID uint) Entity {
	return instance
}

func (instance Class) GetCoverUrl() string {
	return asset.GetAssetManager().MediaUrl(instance.CoverUrl)
}
