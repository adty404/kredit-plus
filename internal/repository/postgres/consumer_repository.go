package postgres

import (
	"github.com/adty404/kredit-plus/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// consumerRepository adalah implementasi konkret dari domain.ConsumerRepository untuk database PostgreSQL.
type consumerRepository struct {
	db *gorm.DB
}

// NewConsumerRepository adalah factory function untuk membuat instance baru dari consumerRepository.
func NewConsumerRepository(db *gorm.DB) domain.ConsumerRepository {
	return &consumerRepository{db: db}
}

func (r *consumerRepository) WithTx(tx *gorm.DB) domain.ConsumerRepository {
	return &consumerRepository{db: tx}
}

// FindByIDForUpdate mencari konsumen berdasarkan ID dan mengunci barisnya menggunakan 'SELECT ... FOR UPDATE'.
func (r *consumerRepository) FindByIDForUpdate(id uint) (*domain.Consumer, error) {
	var consumer domain.Consumer
	err := r.db.Clauses(clause.Locking{Strength: "UPDATE"}).First(&consumer, id).Error
	if err != nil {
		return nil, err
	}
	return &consumer, nil
}

// Save menyimpan data konsumen baru ke database.
func (r *consumerRepository) Save(consumer *domain.Consumer) error {
	return r.db.Create(consumer).Error
}

// Update memperbarui data konsumen yang sudah ada di database.
func (r *consumerRepository) Update(id uint, updates map[string]interface{}) error {
	return r.db.Model(&domain.Consumer{}).Where("id = ?", id).Updates(updates).Error
}

// FindByID mencari satu konsumen berdasarkan ID mereka.
func (r *consumerRepository) FindByID(id uint) (*domain.Consumer, error) {
	var consumer domain.Consumer
	err := r.db.Preload("CreditLimits").Preload("Transactions").First(&consumer, id).Error
	if err != nil {
		return nil, err
	}
	return &consumer, nil
}

// FindByNIK mencari satu konsumen berdasarkan NIK mereka.
func (r *consumerRepository) FindByNIK(nik string) (*domain.Consumer, error) {
	var consumer domain.Consumer
	err := r.db.Where("nik = ?", nik).First(&consumer).Error
	if err != nil {
		return nil, err
	}
	return &consumer, nil
}

// FindAll mengambil semua data konsumen dari database.
func (r *consumerRepository) FindAll() ([]*domain.Consumer, error) {
	var consumers []*domain.Consumer
	err := r.db.Preload("CreditLimits").Preload("Transactions").Find(&consumers).Error
	if err != nil {
		return nil, err
	}
	return consumers, nil
}

// Delete menghapus data konsumen dari database berdasarkan ID.
func (r *consumerRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Consumer{}, id).Error
}
