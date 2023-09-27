package entity

import (
	"gitlab.com/mefit/mefit-server/services/asset"
)

// import (
// 	"github.com/jinzhu/gorm"
// 	"gitlab.com/mefit/mefit-server/utils"
// 	"gitlab.com/mefit/mefit-server/utils/log"
// )

func (instance Movement) NamedEntity(name string) Entity {
	instance.Name = name
	return instance
}
func (instance Movement) IDEntity(Id uint) Entity {
	instance.ID = Id
	return instance
}
func (instance Movement) UserEntity(usrID uint) Entity {
	return instance
}

func (instance Movement) GetVideoUrl() string {
	return asset.GetAssetManager().MediaUrl(instance.VideoUrl)
}

func (instance Movement) GetThumbnailUrl() string {
	return asset.GetAssetManager().MediaUrl(instance.ThumbnailUrl)
}
