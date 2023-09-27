package entity

//Profile implement Entity
func (instance Profile) NamedEntity(name string) Entity {
	instance.Name = name
	return instance
}
func (instance Profile) IDEntity(id uint) Entity {
	instance.ID = id
	return instance
}
func (instance Profile) UserEntity(usrID uint) Entity {
	instance.UserID = usrID
	return instance
}

func (instance Profile) Valid() error {
	return nil
}
