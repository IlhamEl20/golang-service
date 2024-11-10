package utils

import (
	"encoding/csv"
	"fmt"
	"os"
	"provinsi/models"
	"strconv"

	"gorm.io/gorm"
)

// Fungsi untuk mengimpor data Provinsi dari CSV
func ImportProvinsi(db *gorm.DB, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open provinsi.csv: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read provinsi.csv: %w", err)
	}

	for _, record := range records[1:] { // Lewati header
		id, _ := strconv.Atoi(record[0])
		provinsi := models.Provinsi{
			ID:   uint(id),
			Name: record[1],
		}
		db.Create(&provinsi)
	}

	return nil
}

// Fungsi untuk mengimpor data Kota/Kabupaten dari CSV
func ImportKota(db *gorm.DB, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open kabupaten_kota.csv: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read kabupaten_kota.csv: %w", err)
	}

	for _, record := range records[1:] { // Lewati header
		id, _ := strconv.Atoi(record[0])
		provinsiID, _ := strconv.Atoi(record[2]) // Assuming provinsiID is in the third column
		kota := models.Kota{
			ID:         uint(id),
			Name:       record[1],
			ProvinsiID: uint(provinsiID),
		}
		db.Create(&kota)
	}

	return nil
}

// Fungsi untuk mengimpor data Kecamatan dari CSV
func ImportKecamatan(db *gorm.DB, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open kecamatan.csv: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read kecamatan.csv: %w", err)
	}

	for _, record := range records[1:] { // Lewati header
		id, _ := strconv.Atoi(record[0])
		kotaID, _ := strconv.Atoi(record[2]) // Assuming kotaID is in the third column
		kecamatan := models.Kecamatan{
			ID:     uint(id),
			Name:   record[1],
			KotaID: uint(kotaID),
		}
		db.Create(&kecamatan)
	}

	return nil
}

// Fungsi untuk mengimpor data Kelurahan dari CSV
func ImportKelurahan(db *gorm.DB, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open kelurahan.csv: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read kelurahan.csv: %w", err)
	}

	for _, record := range records[1:] { // Lewati header
		id, _ := strconv.Atoi(record[0])
		kecamatanID, _ := strconv.Atoi(record[2]) // Assuming kecamatanID is in the third column
		kelurahan := models.Kelurahan{
			ID:          uint(id),
			Name:        record[1],
			KecamatanID: uint(kecamatanID),
		}
		db.Create(&kelurahan)
	}

	return nil
}
