package models

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	DisplayName  string    `json:"display_name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Character struct {
	ID         uuid.UUID `json:"id"`
	AccountID  uuid.UUID `json:"account_id"`
	Name       string    `json:"name"`
	Gender     string    `json:"gender"`
	SkinTone   int16     `json:"skin_tone"`
	HairStyle  int16     `json:"hair_style"`
	FacePreset int16     `json:"face_preset"`
	OutfitID   string    `json:"outfit_id"`
	Cash       int64     `json:"cash"`
	Bank       int64     `json:"bank"`
	PosX       float64   `json:"pos_x"`
	PosY       float64   `json:"pos_y"`
	PosZ       float64   `json:"pos_z"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type InventoryItem struct {
	ID          uuid.UUID `json:"id"`
	CharacterID uuid.UUID `json:"character_id"`
	ItemKey     string    `json:"item_key"`
	Quantity    int       `json:"quantity"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type House struct {
	ID            uuid.UUID `json:"id"`
	CharacterID   uuid.UUID `json:"character_id"`
	ListingKey    string    `json:"listing_key"`
	Label         string    `json:"label"`
	PosX          float64   `json:"pos_x"`
	PosY          float64   `json:"pos_y"`
	PosZ          float64   `json:"pos_z"`
	PurchasePrice int64     `json:"purchase_price"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Vehicle struct {
	ID            uuid.UUID `json:"id"`
	CharacterID   uuid.UUID `json:"character_id"`
	ListingKey    string    `json:"listing_key"`
	Label         string    `json:"label"`
	VehicleType   string    `json:"vehicle_type"`
	Fuel          float64   `json:"fuel"`
	FuelMax       float64   `json:"fuel_max"`
	PosX          float64   `json:"pos_x"`
	PosY          float64   `json:"pos_y"`
	PosZ          float64   `json:"pos_z"`
	PurchasePrice int64     `json:"purchase_price"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
