package database

import (
	"time"

	"api.mijkomp.com/exception"
	"api.mijkomp.com/helpers"
	"api.mijkomp.com/models/entity"
	"gorm.io/gorm"
)

func SeedData(db *gorm.DB) {
	seedUserAdmin(db)
	seedComponentTypes(db)
}

func seedComponentTypes(db *gorm.DB) error {
	componentTypes := []entity.ComponentType{
		{Code: "Motherboard", Name: "Motherboard", Description: "Main circuit board that connects all components"},
		{Code: "CPU", Name: "Processor", Description: "Central Processing Unit"},
		{Code: "RAM", Name: "Memory (RAM)", Description: "Random Access Memory"},
		{Code: "GPU", Name: "Graphics Card (GPU)", Description: "Video output processor"},
		{Code: "PowerSupply", Name: "Power Supply (PSU)", Description: "Supplies power to all components"},
		{Code: "Storage", Name: "Storage (SSD/HDD)", Description: "Stores operating system and files"},
		{Code: "CPUCooler", Name: "CPU Cooler", Description: "Cools down the processor"},
		{Code: "Case", Name: "PC Case", Description: "Enclosure that holds all components"},
		{Code: "CaseFan", Name: "Case Fan", Description: "Cooling fans for airflow"},
		{Code: "ThermalPaste", Name: "Thermal Paste", Description: "Thermal interface for CPU heat transfer"},
		{Code: "OpticalDrive", Name: "Optical Drive", Description: "CD/DVD/Blu-ray drive"},
		{Code: "SoundCard", Name: "Sound Card", Description: "Dedicated audio processing card"},
		{Code: "NetworkCard", Name: "Network Card", Description: "LAN/WiFi adapter"},
		{Code: "Monitor", Name: "Monitor", Description: "Output display screen"},
		{Code: "Keyboard", Name: "Keyboard", Description: "Input device for typing"},
		{Code: "Mouse", Name: "Mouse", Description: "Input pointing device"},
		{Code: "Speaker", Name: "Speaker", Description: "Audio output device"},
		{Code: "Headset", Name: "Headset", Description: "Headphones with microphone"},
		{Code: "ExpansionCard", Name: "Expansion Card", Description: "Additional internal card (e.g., capture card)"},
		{Code: "M2Heatsink", Name: "M.2 Heatsink", Description: "Cooling solution for M.2 SSDs"},
		{Code: "RGBController", Name: "RGB Controller", Description: "Manages RGB lighting"},
		{Code: "Bracket", Name: "Mounting Bracket", Description: "Bracket for support/mounting components"},
	}

	for _, ct := range componentTypes {
		newComponentType := ct
		newComponentType.CreatedById = 0
		newComponentType.CreatedAt = time.Now().UTC()
		newComponentType.ModifiedById = 0
		newComponentType.ModifiedAt = time.Now().UTC()

		err := db.FirstOrCreate(&newComponentType, entity.ComponentType{Code: ct.Code}).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func seedUserAdmin(db *gorm.DB) {
	// seed admin
	var userCount int64 = 0
	db.Find(&entity.User{}).Count(&userCount)

	if userCount == 0 {

		pass := ""
		db.Save(&entity.User{
			UserName: "system",
			FullName: "System",
			Email:    "system@mail.com",
			Password: &pass,
		})

		pass, err := helpers.PasswordHash("Admin123!@#")
		exception.PanicIfNeeded(err)
		db.Save(&entity.User{
			UserName: "SuperAdmin",
			FullName: "Admin",
			Email:    "admin@mail.com",
			Password: &pass,
		})
	}
}
