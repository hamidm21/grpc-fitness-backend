package entity

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/mefit/mefit-server/utils"
	"gitlab.com/mefit/mefit-server/utils/log"
)

func (in *User) AfterCreate(tx *gorm.DB) (err error) {
	return
}

func (in *User) AfterUpdate(tx *gorm.DB) (err error) {
	return
}

func (in *WorkoutHistory) AfterCreate(tx *gorm.DB) (err error) {
	//Profile update to next workout
	pro := Profile{}
	pro.ID = in.ProfileID

	if err := SimpleCrud(pro).Get(&pro, "Plan"); err != nil {
		return err
	}
	if pro.Plan == nil {
		return utils.ErrAppNotInstalled
	}

	CurrentWorkoutNo := 0
	if len(pro.Plan.WorkoutArray) > CurrentWorkoutNo {
		CurrentWorkoutNo = pro.CurrentWorkoutNo + 1
	}
	query := &Profile{}
	query.ID = pro.ID
	if err := SimpleCrud(query).WithTransaction(tx).Updates(Profile{CurrentWorkoutNo: CurrentWorkoutNo}); err != nil {
		return utils.ErrInternal
	}

	return

}

func (in *Profile) AfterSave(tx *gorm.DB) (err error) {
	if err = in.Valid(); err != nil {
		return err
	}
	return
}

// func (in *BazaarPayment) AfterUpdate(tx *gorm.DB) (err error) {
// 	prof := Profile{}
// 	prof.UserID = in.UserID
// 	if err := SimpleCrud(prof).Get(&prof); err != nil {
// 		return err
// 	}
// 	purchased := PurchasedProduct{
// 		PaymentType: "bazaar",
// 		PaymentID:   in.ID,
// 		ProfileID:   prof.ID,
// 		ProductID:   in.ProductID,
// 	}
// 	if err := SimpleCrud(PurchasedProduct{ProfileID: prof.ID, ProductID: in.ProductID}).WithTransaction(tx).FirstOrCreate(&purchased); err != nil {
// 		return err
// 	}

// 	return nil
// }

func (in *Payment) AfterUpdate(tx *gorm.DB) (err error) {
	product := Product{}
	product.ID = in.ProductID
	if err := SimpleCrud(product).Get(&product); err != nil {
		return err
	}
	prof := Profile{}
	prof.UserID = in.UserID
	if err := SimpleCrud(prof).Get(&prof); err != nil {
		return err
	}
	purchased := PurchasedProduct{
		PaymentType: "zarin",
		PaymentID:   in.ID,
		// Profile:   prof,
		ProfileID: prof.ID,
		// Product:   product,
		ProductID: product.ID,
	}
	if err := SimpleCrud(&purchased).WithTransaction(tx).Create(); err != nil {
		log.Logger().Error(err)
		return err
	}
	return nil
}
