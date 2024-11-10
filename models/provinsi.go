package models

type Provinsi struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"type:varchar(100);not null"`
	Kotas []Kota
}

type Kota struct {
	ID         uint   `gorm:"primaryKey"`
	Name       string `gorm:"type:varchar(100);not null"`
	ProvinsiID uint
	Provinsi   Provinsi `gorm:"foreignKey:ProvinsiID"`
	Kecamatans []Kecamatan
}

type Kecamatan struct {
	ID         uint   `gorm:"primaryKey"`
	Name       string `gorm:"type:varchar(100);not null"`
	KotaID     uint
	Kota       Kota `gorm:"foreignKey:KotaID"`
	Kelurahans []Kelurahan
}

type Kelurahan struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"type:varchar(100);not null"`
	KecamatanID uint
	Kecamatan   Kecamatan `gorm:"foreignKey:KecamatanID"`
}
