package domain

type (
	ConsumerRepository interface {
		Save(consumer *Consumer) error
		Update(id uint, updates map[string]interface{}) error
		FindByID(id uint) (*Consumer, error)
		FindByNIK(nik string) (*Consumer, error)
		FindAll() ([]*Consumer, error)
		Delete(id uint) error
	}
)
