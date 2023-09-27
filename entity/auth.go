package entity

// import (
	// "github.com/jinzhu/gorm"
// 	"gitlab.com/mefit/mefit-server/utils"
// 	"gitlab.com/mefit/mefit-server/utils/log"
// )

// func UserActivation(tx *gorm.DB, phoneNo, activationToken string) (*User, error) {
// 	usr := &User{Phone: phoneNo}

// 	if err := SimpleCrud(*usr).WithTransaction(tx).Get(usr); err != nil {
// 		log.Logger().Errorf("While query for %s: %v", phoneNo, err)
// 		return nil, utils.ErrNotFound
// 	}
// 	if !utils.CheckPasswordHash(activationToken, usr.ActivationToken) {
// 		log.Logger().Errorf("While check activation token for %v", phoneNo)
// 		return nil, utils.ErrInvalidActivationToken
// 	}
// 	//Clear user token (in transaction)
// 	usr.ActivationToken = ""
// 	// usr.PhoneConfirmed = true
// 	if err := SimpleCrud(usr).WithTransaction(tx).Save(); err != nil {
// 		log.Logger().Errorf("While expiring activation token for %s: %v", phoneNo, err)
// 		return nil, utils.ErrInternal
// 	}

// 	//Create profile if there is none
// 	pro := &Profile{UserID: usr.ID, DaysToWorkout: One2Three}
// 	if err := SimpleCrud(pro).WithTransaction(tx).FirstOrCreate(); err != nil {
// 		return nil, err
// 	}
// 	return usr, nil
// }

// func UserSignIn(tx *gorm.DB, email, password string) (*User, error) {
// 	usr := &User{}
// }


//User implement Entity
func (instance User) NamedEntity(name string) Entity {
	// instance.Phone = name
	return instance
}
func (instance User) IDEntity(id uint) Entity {
	instance.ID = id
	return instance
}
func (instance User) UserEntity(usrID uint) Entity {
	instance.ID = usrID
	return instance
}
