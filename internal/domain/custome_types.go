package domain

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// JSONDate adalah tipe kustom untuk menangani format tanggal YYYY-MM-DD dalam JSON.
type JSONDate time.Time

// MarshalJSON mengimplementasikan interface json.Marshaler.
// Ini akan dipanggil saat Gin mengubah struct menjadi JSON.
func (j JSONDate) MarshalJSON() ([]byte, error) {
	// Format waktu ke "YYYY-MM-DD" dan tambahkan tanda kutip ganda agar valid sebagai string JSON.
	stamp := fmt.Sprintf(`"%s"`, time.Time(j).Format("2006-01-02"))
	return []byte(stamp), nil
}

// Value mengimplementasikan interface driver.Valuer.
// Ini memberitahu GORM cara menyimpan tipe ini ke database.
func (j JSONDate) Value() (driver.Value, error) {
	return time.Time(j), nil
}

// Scan mengimplementasikan interface sql.Scanner.
// Ini memberitahu GORM cara membaca tipe ini dari database.
func (j *JSONDate) Scan(value interface{}) error {
	if t, ok := value.(time.Time); ok {
		*j = JSONDate(t)
		return nil
	}
	return fmt.Errorf("could not scan type %T into JSONDate", value)
}
