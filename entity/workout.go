package entity

// import "gitlab.com/mefit/mefit-server/entity"

// import (
// 	"github.com/jinzhu/gorm"
// 	"gitlab.com/mefit/mefit-server/utils"
// 	"gitlab.com/mefit/mefit-server/utils/log"
// )

//TODO: complete this using user payment records
// func WorkoutVipLock(pro *entity.Profile, item *entity.Workout) bool {
// }

func (instance Workout) NamedEntity(name string) Entity {
	instance.Name = name
	return instance
}
func (instance Workout) IDEntity(Id uint) Entity {
	instance.ID = Id
	return instance
}
func (instance Workout) UserEntity(usrID uint) Entity {
	return instance
}

//WorkoutHistory implement Entity
func (instance WorkoutHistory) NamedEntity(name string) Entity {
	//NOOP
	return instance
}
func (instance WorkoutHistory) IDEntity(Id uint) Entity {
	//NOOP
	return instance
}
func (instance WorkoutHistory) UserEntity(usrID uint) Entity {
	//TODO: preload profile here to make this working
	//instance.Profile.UserID = usrID
	return instance
}

func (instance WorkoutHistory) Valid() (valid bool) {
	if instance.Rating < 5 && instance.Rating > 0 {
		valid = true
	}
	return
}
