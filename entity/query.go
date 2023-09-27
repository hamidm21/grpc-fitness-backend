package entity

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/mefit/mefit-server/utils"
	"gitlab.com/mefit/mefit-server/utils/log"
)

//Entity implemented by all entitiess
type Entity interface {
	NamedEntity(name string) Entity
	IDEntity(id uint) Entity
	UserEntity(usrID uint) Entity
}

type Binding interface {
	Entity
	ForApplication(id uint) Binding
	BindingField() string
	GetBound() Entity
}

//DBHelper will embedd an entity and a db instance
type DBHelper struct {
	*gorm.DB
	in   Entity
	page uint32
}

func (instance *DBHelper) List(out interface{}, customQuery *string, preloads ...string) error {
	db := instance.DB
	for _, preload := range preloads {
		db = db.Preload(preload)
	}
	db = db.Where(instance.in)

	if customQuery != nil {
		db = db.Where(*customQuery)
	}
	err := db.Order("created_at desc").Find(out).Error
	if err != nil {

		log.Logger().Errorf("cant fetch list of type %T duo unxpected error: %v", out, err)
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			return utils.ErrNotFound
		}
		return err
	}
	return nil
}

func (instance *DBHelper) LimitedList(out interface{}, page uint32, customQuery string, preloads ...string) error {
	db := instance.DB
	for _, preload := range preloads {
		db = db.Preload(preload)
	}

	if customQuery != "" {
		db = db.Where(customQuery)
	}

	db = db.Where(instance.in).Limit(20).Offset(page - 1)

	err := db.Order("updated_at desc").Find(out).Error
	if err != nil {

		log.Logger().Errorf("cant fetch list of type %T duo unxpected error: %v", out, err)
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			return utils.ErrNotFound
		}
		return err
	}
	return nil
}

// func (instance *DBHelper) Limit(page uint) error {
// 	db := instance.DB

// }

func (instance *DBHelper) Names() ([]string, error) {
	db := instance.DB
	rows, err := db.Model(instance.in).Where(instance.in).Select("name").Rows()

	if err != nil {
		log.Logger().Errorf("cant list names duo unxpected error: %v", err)
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			return nil, utils.ErrNotFound
		}
		return nil, err
	}

	names := []string{}
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			return nil, utils.ErrInternal
		}
		names = append(names, name)
	}
	return names, nil
}

//FIXME: This method seems to have bug because of Model.Where bug in gorm
func (instance *DBHelper) Count() (uint32, error) {
	db := instance.DB
	var (
		count uint32
		err   error
	)
	err = db.Model(instance.in).Where(instance.in).Count(&count).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Logger().Errorf("cant fetch list of type %T", instance.in)
		return 0, utils.ErrInternal
	}
	return count, nil
}

func (instance *DBHelper) Get(out interface{}, preloads ...string) error {
	db := instance.DB
	for _, preload := range preloads {
		db = db.Preload(preload)
	}
	err := db.Where(instance.in).First(out).Error
	if err != nil {

		log.Logger().Errorf("cant fetch item of type %T duo unxpected error: %v", out, err)
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			return utils.ErrNotFound
		}
		return err
	}
	return nil
}

func (instance *DBHelper) Related(out interface{}, foriegnKeys ...string) error {
	db := instance.DB
	err := db.Model(instance.in).Related(out, foriegnKeys...).Error
	if err != nil {
		log.Logger().Errorf("cant fetch list of type %T duo unxpected error: %v", out, err)
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			return utils.ErrNotFound
		}
		return err
	}
	return nil
}

func (instance *DBHelper) Create() error {
	db := instance.DB
	err := db.Create(instance.in).Error
	if err != nil {
		log.Logger().Errorf("cant create of type %T duo unxpected error: %v", instance.in, err)
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			return utils.ErrNotFound
		}
		if q := db.First(instance.in); !q.RecordNotFound() {
			return utils.ErrExists
		}
		return err
	}
	return nil
}
func (instance *DBHelper) FirstOrCreate(out interface{}) error {
	db := instance.DB
	err := db.Where(instance.in).FirstOrCreate(out).Error
	if err != nil {
		log.Logger().Errorf("cant create of type %T duo unxpected error: %v", instance.in, err)
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			return utils.ErrNotFound
		}
		if q := db.First(instance.in); !q.RecordNotFound() {
			return utils.ErrExists
		}
		return err
	}
	return nil
}

//TODO: more error handling needed
func (instance *DBHelper) UpdateOrCreate(out interface{}) error {
	db := instance.DB
	err := db.Where(instance.in).First(out).Error
	if err != nil && err.Error() == gorm.ErrRecordNotFound.Error() {
		return db.Create(out).Error
	}
	return db.Model(instance.in).Where(instance.in).Updates(out).Error
}

func (instance *DBHelper) Save() error {
	db := instance.DB
	err := db.Save(instance.in).Error
	if err != nil {
		log.Logger().Errorf("cant save of type %T duo unxpected error: %v", instance.in, err)
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			return utils.ErrNotFound
		}
		if q := db.First(instance.in); !q.RecordNotFound() {
			return utils.ErrExists
		}
		return err
	}
	return nil
}

func (instance *DBHelper) Updates(updates interface{}) error {
	db := instance.DB
	err := db.Model(instance.in).Updates(updates).Error
	if err != nil {
		log.Logger().Errorf("cant update of type %T duo unxpected error: %v", instance.in, err)
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			return utils.ErrNotFound
		}

		return err
	}
	return nil
}

func (instance *DBHelper) Delete(unscopped bool) error {
	db := instance.DB
	//FIXME: For now delete imediatly
	// if unscopped {
	db = db.Unscoped()
	// }

	//Double check for ID because of https://github.com/jinzhu/gorm/issues/334
	err := db.Where(instance.in).First(instance.in).Error
	if err != nil && err.Error() == gorm.ErrRecordNotFound.Error() {
		log.Logger().Errorf("cant fetch list of type %T", instance.in)
		return utils.ErrNotFound
	}

	err = db.Where(instance.in).Delete(instance.in).Error
	if err != nil {
		log.Logger().Errorf("cant delete of type %T duo unxpected error: %v", instance.in, err)
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			return utils.ErrNotFound
		}
		return err
	}
	return nil
}

//DBHelper is an per user authorized DBHelper
func Crud(in Entity, usrID uint) *DBHelper {
	in = in.UserEntity(usrID)
	return SimpleCrud(in)
}

//DBHelper is an per user authorized DBHelper
func SimpleCrud(in Entity) *DBHelper {
	return &DBHelper{DB: db, in: in}
}

func (instance *DBHelper) WithTransaction(tx *gorm.DB) *DBHelper {
	instance.DB = tx
	return instance
}

func (instance *DBHelper) ID(id uint32) *DBHelper {
	instance.in = instance.in.IDEntity(uint(id))
	return instance
}

func (instance *DBHelper) Name(name string) *DBHelper {
	instance.in = instance.in.NamedEntity(name)
	return instance
}

func (instance *DBHelper) Limit(page uint32) *DBHelper {
	instance.page = page
	return instance
}
