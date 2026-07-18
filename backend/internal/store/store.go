package store

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/AFR-projection/GTA/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNotFound          = errors.New("not found")
	ErrConflict          = errors.New("conflict")
	ErrLimitReached      = errors.New("character limit reached")
	ErrInvalidInput      = errors.New("invalid input")
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrAlreadyOwnsHouse  = errors.New("already owns a house")
	ErrOutOfFuel         = errors.New("out of fuel")
	ErrJobOnCooldown     = errors.New("job on cooldown")
	ErrNotEnoughItems    = errors.New("not enough items")
)

type Store struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool}
}

func (s *Store) CreateAccount(ctx context.Context, email, passwordHash, displayName string) (models.Account, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	displayName = strings.TrimSpace(displayName)
	if email == "" || passwordHash == "" || displayName == "" {
		return models.Account{}, ErrInvalidInput
	}

	var a models.Account
	err := s.pool.QueryRow(ctx, `
		INSERT INTO accounts (email, password_hash, display_name)
		VALUES ($1, $2, $3)
		RETURNING id, email, password_hash, display_name, created_at, updated_at
	`, email, passwordHash, displayName).Scan(
		&a.ID, &a.Email, &a.PasswordHash, &a.DisplayName, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return models.Account{}, ErrConflict
		}
		return models.Account{}, fmt.Errorf("create account: %w", err)
	}
	return a, nil
}

func (s *Store) GetAccountByEmail(ctx context.Context, email string) (models.Account, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	var a models.Account
	err := s.pool.QueryRow(ctx, `
		SELECT id, email, password_hash, display_name, created_at, updated_at
		FROM accounts WHERE email = $1
	`, email).Scan(
		&a.ID, &a.Email, &a.PasswordHash, &a.DisplayName, &a.CreatedAt, &a.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Account{}, ErrNotFound
	}
	if err != nil {
		return models.Account{}, fmt.Errorf("get account by email: %w", err)
	}
	return a, nil
}

func (s *Store) GetAccountByID(ctx context.Context, id uuid.UUID) (models.Account, error) {
	var a models.Account
	err := s.pool.QueryRow(ctx, `
		SELECT id, email, password_hash, display_name, created_at, updated_at
		FROM accounts WHERE id = $1
	`, id).Scan(
		&a.ID, &a.Email, &a.PasswordHash, &a.DisplayName, &a.CreatedAt, &a.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Account{}, ErrNotFound
	}
	if err != nil {
		return models.Account{}, fmt.Errorf("get account by id: %w", err)
	}
	return a, nil
}

type CreateCharacterInput struct {
	Name       string
	Gender     string
	SkinTone   int16
	HairStyle  int16
	FacePreset int16
	OutfitID   string
}

func (s *Store) CreateCharacter(ctx context.Context, accountID uuid.UUID, in CreateCharacterInput) (models.Character, error) {
	in.Name = strings.TrimSpace(in.Name)
	in.Gender = strings.ToLower(strings.TrimSpace(in.Gender))
	if in.OutfitID == "" {
		in.OutfitID = "starter_01"
	}
	if in.Name == "" || (in.Gender != "male" && in.Gender != "female") {
		return models.Character{}, ErrInvalidInput
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return models.Character{}, err
	}
	defer tx.Rollback(ctx)

	var count int
	if err := tx.QueryRow(ctx, `SELECT COUNT(*) FROM characters WHERE account_id = $1`, accountID).Scan(&count); err != nil {
		return models.Character{}, err
	}
	// MVP: 1 character per account
	if count >= 1 {
		return models.Character{}, ErrLimitReached
	}

	var c models.Character
	err = tx.QueryRow(ctx, `
		INSERT INTO characters (
			account_id, name, gender, skin_tone, hair_style, face_preset, outfit_id
		) VALUES ($1,$2,$3,$4,$5,$6,$7)
		RETURNING id, account_id, name, gender, skin_tone, hair_style, face_preset, outfit_id,
		          cash, bank, pos_x, pos_y, pos_z, created_at, updated_at
	`, accountID, in.Name, in.Gender, in.SkinTone, in.HairStyle, in.FacePreset, in.OutfitID).Scan(
		&c.ID, &c.AccountID, &c.Name, &c.Gender, &c.SkinTone, &c.HairStyle, &c.FacePreset, &c.OutfitID,
		&c.Cash, &c.Bank, &c.PosX, &c.PosY, &c.PosZ, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return models.Character{}, ErrConflict
		}
		return models.Character{}, fmt.Errorf("create character: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO transactions (character_id, kind, amount, balance_cash, balance_bank, meta)
		VALUES ($1, 'starting_cash', $2, $2, 0, '{"source":"character_create"}'::jsonb)
	`, c.ID, c.Cash); err != nil {
		return models.Character{}, fmt.Errorf("audit starting cash: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return models.Character{}, err
	}
	return c, nil
}

func (s *Store) ListCharacters(ctx context.Context, accountID uuid.UUID) ([]models.Character, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, account_id, name, gender, skin_tone, hair_style, face_preset, outfit_id,
		       cash, bank, pos_x, pos_y, pos_z, created_at, updated_at
		FROM characters WHERE account_id = $1 ORDER BY created_at ASC
	`, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]models.Character, 0)
	for rows.Next() {
		var c models.Character
		if err := rows.Scan(
			&c.ID, &c.AccountID, &c.Name, &c.Gender, &c.SkinTone, &c.HairStyle, &c.FacePreset, &c.OutfitID,
			&c.Cash, &c.Bank, &c.PosX, &c.PosY, &c.PosZ, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (s *Store) GetCharacterForAccount(ctx context.Context, accountID, characterID uuid.UUID) (models.Character, error) {
	var c models.Character
	err := s.pool.QueryRow(ctx, `
		SELECT id, account_id, name, gender, skin_tone, hair_style, face_preset, outfit_id,
		       cash, bank, pos_x, pos_y, pos_z, created_at, updated_at
		FROM characters WHERE id = $1 AND account_id = $2
	`, characterID, accountID).Scan(
		&c.ID, &c.AccountID, &c.Name, &c.Gender, &c.SkinTone, &c.HairStyle, &c.FacePreset, &c.OutfitID,
		&c.Cash, &c.Bank, &c.PosX, &c.PosY, &c.PosZ, &c.CreatedAt, &c.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Character{}, ErrNotFound
	}
	if err != nil {
		return models.Character{}, err
	}
	return c, nil
}

func (s *Store) UpdatePosition(ctx context.Context, accountID, characterID uuid.UUID, x, y, z float64) (models.Character, error) {
	var c models.Character
	err := s.pool.QueryRow(ctx, `
		UPDATE characters
		SET pos_x = $1, pos_y = $2, pos_z = $3, updated_at = NOW()
		WHERE id = $4 AND account_id = $5
		RETURNING id, account_id, name, gender, skin_tone, hair_style, face_preset, outfit_id,
		          cash, bank, pos_x, pos_y, pos_z, created_at, updated_at
	`, x, y, z, characterID, accountID).Scan(
		&c.ID, &c.AccountID, &c.Name, &c.Gender, &c.SkinTone, &c.HairStyle, &c.FacePreset, &c.OutfitID,
		&c.Cash, &c.Bank, &c.PosX, &c.PosY, &c.PosZ, &c.CreatedAt, &c.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Character{}, ErrNotFound
	}
	if err != nil {
		return models.Character{}, err
	}
	return c, nil
}

// DepositBank moves cash → bank. Amount must be > 0.
func (s *Store) DepositBank(ctx context.Context, accountID, characterID uuid.UUID, amount int64) (models.Character, error) {
	if amount <= 0 {
		return models.Character{}, ErrInvalidInput
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return models.Character{}, err
	}
	defer tx.Rollback(ctx)

	var c models.Character
	err = tx.QueryRow(ctx, `
		SELECT id, account_id, name, gender, skin_tone, hair_style, face_preset, outfit_id,
		       cash, bank, pos_x, pos_y, pos_z, created_at, updated_at
		FROM characters WHERE id = $1 AND account_id = $2
		FOR UPDATE
	`, characterID, accountID).Scan(
		&c.ID, &c.AccountID, &c.Name, &c.Gender, &c.SkinTone, &c.HairStyle, &c.FacePreset, &c.OutfitID,
		&c.Cash, &c.Bank, &c.PosX, &c.PosY, &c.PosZ, &c.CreatedAt, &c.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Character{}, ErrNotFound
	}
	if err != nil {
		return models.Character{}, err
	}
	if c.Cash < amount {
		return models.Character{}, ErrInsufficientFunds
	}

	c.Cash -= amount
	c.Bank += amount
	err = tx.QueryRow(ctx, `
		UPDATE characters SET cash = $1, bank = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING updated_at
	`, c.Cash, c.Bank, c.ID).Scan(&c.UpdatedAt)
	if err != nil {
		return models.Character{}, err
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO transactions (character_id, kind, amount, balance_cash, balance_bank, meta)
		VALUES ($1, 'bank_deposit', $2, $3, $4, '{}'::jsonb)
	`, c.ID, amount, c.Cash, c.Bank); err != nil {
		return models.Character{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return models.Character{}, err
	}
	return c, nil
}

// WithdrawBank moves bank → cash. Amount must be > 0.
func (s *Store) WithdrawBank(ctx context.Context, accountID, characterID uuid.UUID, amount int64) (models.Character, error) {
	if amount <= 0 {
		return models.Character{}, ErrInvalidInput
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return models.Character{}, err
	}
	defer tx.Rollback(ctx)

	var c models.Character
	err = tx.QueryRow(ctx, `
		SELECT id, account_id, name, gender, skin_tone, hair_style, face_preset, outfit_id,
		       cash, bank, pos_x, pos_y, pos_z, created_at, updated_at
		FROM characters WHERE id = $1 AND account_id = $2
		FOR UPDATE
	`, characterID, accountID).Scan(
		&c.ID, &c.AccountID, &c.Name, &c.Gender, &c.SkinTone, &c.HairStyle, &c.FacePreset, &c.OutfitID,
		&c.Cash, &c.Bank, &c.PosX, &c.PosY, &c.PosZ, &c.CreatedAt, &c.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Character{}, ErrNotFound
	}
	if err != nil {
		return models.Character{}, err
	}
	if c.Bank < amount {
		return models.Character{}, ErrInsufficientFunds
	}

	c.Bank -= amount
	c.Cash += amount
	err = tx.QueryRow(ctx, `
		UPDATE characters SET cash = $1, bank = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING updated_at
	`, c.Cash, c.Bank, c.ID).Scan(&c.UpdatedAt)
	if err != nil {
		return models.Character{}, err
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO transactions (character_id, kind, amount, balance_cash, balance_bank, meta)
		VALUES ($1, 'bank_withdraw', $2, $3, $4, '{}'::jsonb)
	`, c.ID, amount, c.Cash, c.Bank); err != nil {
		return models.Character{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return models.Character{}, err
	}
	return c, nil
}

func (s *Store) ListInventory(ctx context.Context, accountID, characterID uuid.UUID) ([]models.InventoryItem, error) {
	// ownership check
	if _, err := s.GetCharacterForAccount(ctx, accountID, characterID); err != nil {
		return nil, err
	}

	rows, err := s.pool.Query(ctx, `
		SELECT id, character_id, item_key, quantity, created_at, updated_at
		FROM inventory_items WHERE character_id = $1 ORDER BY item_key ASC
	`, characterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]models.InventoryItem, 0)
	for rows.Next() {
		var it models.InventoryItem
		if err := rows.Scan(&it.ID, &it.CharacterID, &it.ItemKey, &it.Quantity, &it.CreatedAt, &it.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, it)
	}
	return out, rows.Err()
}

type UseItemResult struct {
	Inventory []models.InventoryItem `json:"inventory"`
	ItemKey   string                 `json:"item_key"`
	Quantity  int                    `json:"quantity_used"`
}

// UseInventoryItem consumes qty of a usable item (server validates catalog).
func (s *Store) UseInventoryItem(ctx context.Context, accountID, characterID uuid.UUID, itemKey string, qty int) (UseItemResult, error) {
	if itemKey == "" || qty <= 0 {
		return UseItemResult{}, ErrInvalidInput
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return UseItemResult{}, err
	}
	defer tx.Rollback(ctx)

	var cID uuid.UUID
	err = tx.QueryRow(ctx, `
		SELECT id FROM characters WHERE id = $1 AND account_id = $2 FOR UPDATE
	`, characterID, accountID).Scan(&cID)
	if errors.Is(err, pgx.ErrNoRows) {
		return UseItemResult{}, ErrNotFound
	}
	if err != nil {
		return UseItemResult{}, err
	}

	var itemID uuid.UUID
	var have int
	err = tx.QueryRow(ctx, `
		SELECT id, quantity FROM inventory_items
		WHERE character_id = $1 AND item_key = $2
		FOR UPDATE
	`, characterID, itemKey).Scan(&itemID, &have)
	if errors.Is(err, pgx.ErrNoRows) {
		return UseItemResult{}, ErrNotEnoughItems
	}
	if err != nil {
		return UseItemResult{}, err
	}
	if have < qty {
		return UseItemResult{}, ErrNotEnoughItems
	}

	left := have - qty
	if left == 0 {
		if _, err := tx.Exec(ctx, `DELETE FROM inventory_items WHERE id = $1`, itemID); err != nil {
			return UseItemResult{}, err
		}
	} else {
		if _, err := tx.Exec(ctx, `
			UPDATE inventory_items SET quantity = $1, updated_at = NOW() WHERE id = $2
		`, left, itemID); err != nil {
			return UseItemResult{}, err
		}
	}

	meta := fmt.Sprintf(`{"item_key":%q,"qty":%d}`, itemKey, qty)
	if _, err := tx.Exec(ctx, `
		INSERT INTO transactions (character_id, kind, amount, balance_cash, balance_bank, meta)
		VALUES ($1, 'item_use', 0, NULL, NULL, $2::jsonb)
	`, characterID, meta); err != nil {
		return UseItemResult{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return UseItemResult{}, err
	}

	inv, err := s.ListInventory(ctx, accountID, characterID)
	if err != nil {
		return UseItemResult{}, err
	}
	return UseItemResult{Inventory: inv, ItemKey: itemKey, Quantity: qty}, nil
}

type PurchaseResult struct {
	Character models.Character       `json:"character"`
	Inventory []models.InventoryItem `json:"inventory"`
	ItemKey   string                 `json:"item_key"`
	Quantity  int                    `json:"quantity"`
	TotalPaid int64                  `json:"total_paid"`
}

// PurchaseItem spends cash and upserts inventory. unitPrice/itemKey come from trusted catalog (server).
func (s *Store) PurchaseItem(ctx context.Context, accountID, characterID uuid.UUID, shopID, itemKey string, unitPrice int64, qty int) (PurchaseResult, error) {
	if qty <= 0 || unitPrice < 0 || itemKey == "" || shopID == "" {
		return PurchaseResult{}, ErrInvalidInput
	}
	total := unitPrice * int64(qty)

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return PurchaseResult{}, err
	}
	defer tx.Rollback(ctx)

	var c models.Character
	err = tx.QueryRow(ctx, `
		SELECT id, account_id, name, gender, skin_tone, hair_style, face_preset, outfit_id,
		       cash, bank, pos_x, pos_y, pos_z, created_at, updated_at
		FROM characters WHERE id = $1 AND account_id = $2
		FOR UPDATE
	`, characterID, accountID).Scan(
		&c.ID, &c.AccountID, &c.Name, &c.Gender, &c.SkinTone, &c.HairStyle, &c.FacePreset, &c.OutfitID,
		&c.Cash, &c.Bank, &c.PosX, &c.PosY, &c.PosZ, &c.CreatedAt, &c.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return PurchaseResult{}, ErrNotFound
	}
	if err != nil {
		return PurchaseResult{}, err
	}
	if c.Cash < total {
		return PurchaseResult{}, ErrInsufficientFunds
	}

	c.Cash -= total
	if err := tx.QueryRow(ctx, `
		UPDATE characters SET cash = $1, updated_at = NOW() WHERE id = $2 RETURNING updated_at
	`, c.Cash, c.ID).Scan(&c.UpdatedAt); err != nil {
		return PurchaseResult{}, err
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO inventory_items (character_id, item_key, quantity)
		VALUES ($1, $2, $3)
		ON CONFLICT (character_id, item_key)
		DO UPDATE SET quantity = inventory_items.quantity + EXCLUDED.quantity, updated_at = NOW()
	`, c.ID, itemKey, qty); err != nil {
		return PurchaseResult{}, err
	}

	meta := fmt.Sprintf(`{"shop_id":%q,"item_key":%q,"qty":%d,"unit_price":%d}`, shopID, itemKey, qty, unitPrice)
	if _, err := tx.Exec(ctx, `
		INSERT INTO transactions (character_id, kind, amount, balance_cash, balance_bank, meta)
		VALUES ($1, 'shop_purchase', $2, $3, $4, $5::jsonb)
	`, c.ID, total, c.Cash, c.Bank, meta); err != nil {
		return PurchaseResult{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return PurchaseResult{}, err
	}

	inv, err := s.ListInventory(ctx, accountID, characterID)
	if err != nil {
		return PurchaseResult{}, err
	}

	return PurchaseResult{
		Character: c,
		Inventory: inv,
		ItemKey:   itemKey,
		Quantity:  qty,
		TotalPaid: total,
	}, nil
}

func (s *Store) ListHouses(ctx context.Context, accountID, characterID uuid.UUID) ([]models.House, error) {
	if _, err := s.GetCharacterForAccount(ctx, accountID, characterID); err != nil {
		return nil, err
	}
	rows, err := s.pool.Query(ctx, `
		SELECT id, character_id, listing_key, label, pos_x, pos_y, pos_z, purchase_price, created_at, updated_at
		FROM houses WHERE character_id = $1 ORDER BY created_at ASC
	`, characterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]models.House, 0)
	for rows.Next() {
		var h models.House
		if err := rows.Scan(
			&h.ID, &h.CharacterID, &h.ListingKey, &h.Label,
			&h.PosX, &h.PosY, &h.PosZ, &h.PurchasePrice, &h.CreatedAt, &h.UpdatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, h)
	}
	return out, rows.Err()
}

type BuyHouseResult struct {
	Character models.Character       `json:"character"`
	House     models.House           `json:"house"`
	Inventory []models.InventoryItem `json:"inventory"`
}

// BuyHouse purchases from trusted catalog. MVP: max 1 house per character.
func (s *Store) BuyHouse(ctx context.Context, accountID, characterID uuid.UUID, listingKey, label string, price int64, posX, posY, posZ float64) (BuyHouseResult, error) {
	if listingKey == "" || label == "" || price < 0 {
		return BuyHouseResult{}, ErrInvalidInput
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return BuyHouseResult{}, err
	}
	defer tx.Rollback(ctx)

	var c models.Character
	err = tx.QueryRow(ctx, `
		SELECT id, account_id, name, gender, skin_tone, hair_style, face_preset, outfit_id,
		       cash, bank, pos_x, pos_y, pos_z, created_at, updated_at
		FROM characters WHERE id = $1 AND account_id = $2
		FOR UPDATE
	`, characterID, accountID).Scan(
		&c.ID, &c.AccountID, &c.Name, &c.Gender, &c.SkinTone, &c.HairStyle, &c.FacePreset, &c.OutfitID,
		&c.Cash, &c.Bank, &c.PosX, &c.PosY, &c.PosZ, &c.CreatedAt, &c.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return BuyHouseResult{}, ErrNotFound
	}
	if err != nil {
		return BuyHouseResult{}, err
	}

	var owned int
	if err := tx.QueryRow(ctx, `SELECT COUNT(*) FROM houses WHERE character_id = $1`, c.ID).Scan(&owned); err != nil {
		return BuyHouseResult{}, err
	}
	if owned >= 1 {
		return BuyHouseResult{}, ErrAlreadyOwnsHouse
	}

	if c.Cash < price {
		return BuyHouseResult{}, ErrInsufficientFunds
	}

	c.Cash -= price
	if err := tx.QueryRow(ctx, `
		UPDATE characters SET cash = $1, updated_at = NOW() WHERE id = $2 RETURNING updated_at
	`, c.Cash, c.ID).Scan(&c.UpdatedAt); err != nil {
		return BuyHouseResult{}, err
	}

	var h models.House
	err = tx.QueryRow(ctx, `
		INSERT INTO houses (character_id, listing_key, label, pos_x, pos_y, pos_z, purchase_price)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		RETURNING id, character_id, listing_key, label, pos_x, pos_y, pos_z, purchase_price, created_at, updated_at
	`, c.ID, listingKey, label, posX, posY, posZ, price).Scan(
		&h.ID, &h.CharacterID, &h.ListingKey, &h.Label, &h.PosX, &h.PosY, &h.PosZ, &h.PurchasePrice, &h.CreatedAt, &h.UpdatedAt,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return BuyHouseResult{}, ErrConflict
		}
		return BuyHouseResult{}, err
	}

	keyItem := "house_key_" + listingKey
	if _, err := tx.Exec(ctx, `
		INSERT INTO inventory_items (character_id, item_key, quantity)
		VALUES ($1, $2, 1)
		ON CONFLICT (character_id, item_key)
		DO UPDATE SET quantity = inventory_items.quantity + 1, updated_at = NOW()
	`, c.ID, keyItem); err != nil {
		return BuyHouseResult{}, err
	}

	meta := fmt.Sprintf(`{"listing_key":%q,"label":%q}`, listingKey, label)
	if _, err := tx.Exec(ctx, `
		INSERT INTO transactions (character_id, kind, amount, balance_cash, balance_bank, meta)
		VALUES ($1, 'house_purchase', $2, $3, $4, $5::jsonb)
	`, c.ID, price, c.Cash, c.Bank, meta); err != nil {
		return BuyHouseResult{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return BuyHouseResult{}, err
	}

	inv, err := s.ListInventory(ctx, accountID, characterID)
	if err != nil {
		return BuyHouseResult{}, err
	}

	return BuyHouseResult{Character: c, House: h, Inventory: inv}, nil
}

func (s *Store) ListVehicles(ctx context.Context, accountID, characterID uuid.UUID) ([]models.Vehicle, error) {
	if _, err := s.GetCharacterForAccount(ctx, accountID, characterID); err != nil {
		return nil, err
	}
	rows, err := s.pool.Query(ctx, `
		SELECT id, character_id, listing_key, label, vehicle_type, fuel, fuel_max,
		       pos_x, pos_y, pos_z, purchase_price, created_at, updated_at
		FROM vehicles WHERE character_id = $1 ORDER BY created_at ASC
	`, characterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]models.Vehicle, 0)
	for rows.Next() {
		var v models.Vehicle
		if err := rows.Scan(
			&v.ID, &v.CharacterID, &v.ListingKey, &v.Label, &v.VehicleType, &v.Fuel, &v.FuelMax,
			&v.PosX, &v.PosY, &v.PosZ, &v.PurchasePrice, &v.CreatedAt, &v.UpdatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}

type BuyVehicleResult struct {
	Character models.Character `json:"character"`
	Vehicle   models.Vehicle   `json:"vehicle"`
}

func (s *Store) BuyVehicle(ctx context.Context, accountID, characterID uuid.UUID, listingKey, label, vehicleType string, price int64, fuelMax float64) (BuyVehicleResult, error) {
	if listingKey == "" || label == "" || price < 0 || fuelMax <= 0 {
		return BuyVehicleResult{}, ErrInvalidInput
	}
	if vehicleType != "motorcycle" && vehicleType != "car" {
		return BuyVehicleResult{}, ErrInvalidInput
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return BuyVehicleResult{}, err
	}
	defer tx.Rollback(ctx)

	var c models.Character
	err = tx.QueryRow(ctx, `
		SELECT id, account_id, name, gender, skin_tone, hair_style, face_preset, outfit_id,
		       cash, bank, pos_x, pos_y, pos_z, created_at, updated_at
		FROM characters WHERE id = $1 AND account_id = $2
		FOR UPDATE
	`, characterID, accountID).Scan(
		&c.ID, &c.AccountID, &c.Name, &c.Gender, &c.SkinTone, &c.HairStyle, &c.FacePreset, &c.OutfitID,
		&c.Cash, &c.Bank, &c.PosX, &c.PosY, &c.PosZ, &c.CreatedAt, &c.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return BuyVehicleResult{}, ErrNotFound
	}
	if err != nil {
		return BuyVehicleResult{}, err
	}
	if c.Cash < price {
		return BuyVehicleResult{}, ErrInsufficientFunds
	}

	c.Cash -= price
	if err := tx.QueryRow(ctx, `
		UPDATE characters SET cash = $1, updated_at = NOW() WHERE id = $2 RETURNING updated_at
	`, c.Cash, c.ID).Scan(&c.UpdatedAt); err != nil {
		return BuyVehicleResult{}, err
	}

	var v models.Vehicle
	err = tx.QueryRow(ctx, `
		INSERT INTO vehicles (
			character_id, listing_key, label, vehicle_type, fuel, fuel_max, purchase_price,
			pos_x, pos_y, pos_z
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		RETURNING id, character_id, listing_key, label, vehicle_type, fuel, fuel_max,
		          pos_x, pos_y, pos_z, purchase_price, created_at, updated_at
	`, c.ID, listingKey, label, vehicleType, fuelMax*0.4, fuelMax, price, c.PosX, c.PosY, c.PosZ).Scan(
		&v.ID, &v.CharacterID, &v.ListingKey, &v.Label, &v.VehicleType, &v.Fuel, &v.FuelMax,
		&v.PosX, &v.PosY, &v.PosZ, &v.PurchasePrice, &v.CreatedAt, &v.UpdatedAt,
	)
	if err != nil {
		return BuyVehicleResult{}, err
	}

	keyItem := "vehicle_key_" + listingKey
	if _, err := tx.Exec(ctx, `
		INSERT INTO inventory_items (character_id, item_key, quantity)
		VALUES ($1, $2, 1)
		ON CONFLICT (character_id, item_key)
		DO UPDATE SET quantity = inventory_items.quantity + 1, updated_at = NOW()
	`, c.ID, keyItem); err != nil {
		return BuyVehicleResult{}, err
	}

	meta := fmt.Sprintf(`{"listing_key":%q,"vehicle_type":%q}`, listingKey, vehicleType)
	if _, err := tx.Exec(ctx, `
		INSERT INTO transactions (character_id, kind, amount, balance_cash, balance_bank, meta)
		VALUES ($1, 'vehicle_purchase', $2, $3, $4, $5::jsonb)
	`, c.ID, price, c.Cash, c.Bank, meta); err != nil {
		return BuyVehicleResult{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return BuyVehicleResult{}, err
	}
	return BuyVehicleResult{Character: c, Vehicle: v}, nil
}

type RefuelResult struct {
	Character models.Character `json:"character"`
	Vehicle   models.Vehicle   `json:"vehicle"`
	FuelAdded float64          `json:"fuel_added"`
	TotalPaid int64            `json:"total_paid"`
}

func (s *Store) RefuelVehicle(ctx context.Context, accountID, characterID, vehicleID uuid.UUID, units float64, pricePerUnit int64) (RefuelResult, error) {
	if units <= 0 || pricePerUnit < 0 {
		return RefuelResult{}, ErrInvalidInput
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return RefuelResult{}, err
	}
	defer tx.Rollback(ctx)

	var c models.Character
	err = tx.QueryRow(ctx, `
		SELECT id, account_id, name, gender, skin_tone, hair_style, face_preset, outfit_id,
		       cash, bank, pos_x, pos_y, pos_z, created_at, updated_at
		FROM characters WHERE id = $1 AND account_id = $2
		FOR UPDATE
	`, characterID, accountID).Scan(
		&c.ID, &c.AccountID, &c.Name, &c.Gender, &c.SkinTone, &c.HairStyle, &c.FacePreset, &c.OutfitID,
		&c.Cash, &c.Bank, &c.PosX, &c.PosY, &c.PosZ, &c.CreatedAt, &c.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return RefuelResult{}, ErrNotFound
	}
	if err != nil {
		return RefuelResult{}, err
	}

	var v models.Vehicle
	err = tx.QueryRow(ctx, `
		SELECT id, character_id, listing_key, label, vehicle_type, fuel, fuel_max,
		       pos_x, pos_y, pos_z, purchase_price, created_at, updated_at
		FROM vehicles WHERE id = $1 AND character_id = $2
		FOR UPDATE
	`, vehicleID, characterID).Scan(
		&v.ID, &v.CharacterID, &v.ListingKey, &v.Label, &v.VehicleType, &v.Fuel, &v.FuelMax,
		&v.PosX, &v.PosY, &v.PosZ, &v.PurchasePrice, &v.CreatedAt, &v.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return RefuelResult{}, ErrNotFound
	}
	if err != nil {
		return RefuelResult{}, err
	}

	room := v.FuelMax - v.Fuel
	if room <= 0 {
		return RefuelResult{}, ErrInvalidInput
	}
	if units > room {
		units = room
	}
	// bill whole units rounded up for simplicity (ceil to int fuel units charged)
	chargeUnits := int64(units + 0.999999) // ceil
	if chargeUnits < 1 {
		chargeUnits = 1
	}
	total := chargeUnits * pricePerUnit
	if c.Cash < total {
		return RefuelResult{}, ErrInsufficientFunds
	}

	c.Cash -= total
	v.Fuel += units
	if v.Fuel > v.FuelMax {
		v.Fuel = v.FuelMax
	}

	if err := tx.QueryRow(ctx, `
		UPDATE characters SET cash = $1, updated_at = NOW() WHERE id = $2 RETURNING updated_at
	`, c.Cash, c.ID).Scan(&c.UpdatedAt); err != nil {
		return RefuelResult{}, err
	}
	if err := tx.QueryRow(ctx, `
		UPDATE vehicles SET fuel = $1, updated_at = NOW() WHERE id = $2 RETURNING updated_at
	`, v.Fuel, v.ID).Scan(&v.UpdatedAt); err != nil {
		return RefuelResult{}, err
	}

	meta := fmt.Sprintf(`{"vehicle_id":%q,"units":%v,"price_per_unit":%d}`, v.ID.String(), units, pricePerUnit)
	if _, err := tx.Exec(ctx, `
		INSERT INTO transactions (character_id, kind, amount, balance_cash, balance_bank, meta)
		VALUES ($1, 'spbu_refuel', $2, $3, $4, $5::jsonb)
	`, c.ID, total, c.Cash, c.Bank, meta); err != nil {
		return RefuelResult{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return RefuelResult{}, err
	}
	return RefuelResult{Character: c, Vehicle: v, FuelAdded: units, TotalPaid: total}, nil
}

func (s *Store) CharacterSummary(ctx context.Context, accountID, characterID uuid.UUID) (map[string]any, error) {
	c, err := s.GetCharacterForAccount(ctx, accountID, characterID)
	if err != nil {
		return nil, err
	}
	inv, err := s.ListInventory(ctx, accountID, characterID)
	if err != nil {
		return nil, err
	}
	houses, err := s.ListHouses(ctx, accountID, characterID)
	if err != nil {
		return nil, err
	}
	vehicles, err := s.ListVehicles(ctx, accountID, characterID)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"character": c,
		"inventory": inv,
		"houses":    houses,
		"vehicles":  vehicles,
	}, nil
}

func (s *Store) UpdateVehiclePosition(ctx context.Context, accountID, characterID, vehicleID uuid.UUID, x, y, z float64) (models.Vehicle, error) {
	if _, err := s.GetCharacterForAccount(ctx, accountID, characterID); err != nil {
		return models.Vehicle{}, err
	}
	var v models.Vehicle
	err := s.pool.QueryRow(ctx, `
		UPDATE vehicles
		SET pos_x = $1, pos_y = $2, pos_z = $3, updated_at = NOW()
		WHERE id = $4 AND character_id = $5
		RETURNING id, character_id, listing_key, label, vehicle_type, fuel, fuel_max,
		          pos_x, pos_y, pos_z, purchase_price, created_at, updated_at
	`, x, y, z, vehicleID, characterID).Scan(
		&v.ID, &v.CharacterID, &v.ListingKey, &v.Label, &v.VehicleType, &v.Fuel, &v.FuelMax,
		&v.PosX, &v.PosY, &v.PosZ, &v.PurchasePrice, &v.CreatedAt, &v.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Vehicle{}, ErrNotFound
	}
	if err != nil {
		return models.Vehicle{}, err
	}
	return v, nil
}

func (s *Store) ConsumeFuel(ctx context.Context, accountID, characterID, vehicleID uuid.UUID, amount float64) (models.Vehicle, error) {
	if amount <= 0 {
		return models.Vehicle{}, ErrInvalidInput
	}
	if _, err := s.GetCharacterForAccount(ctx, accountID, characterID); err != nil {
		return models.Vehicle{}, err
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return models.Vehicle{}, err
	}
	defer tx.Rollback(ctx)

	var v models.Vehicle
	err = tx.QueryRow(ctx, `
		SELECT id, character_id, listing_key, label, vehicle_type, fuel, fuel_max,
		       pos_x, pos_y, pos_z, purchase_price, created_at, updated_at
		FROM vehicles WHERE id = $1 AND character_id = $2
		FOR UPDATE
	`, vehicleID, characterID).Scan(
		&v.ID, &v.CharacterID, &v.ListingKey, &v.Label, &v.VehicleType, &v.Fuel, &v.FuelMax,
		&v.PosX, &v.PosY, &v.PosZ, &v.PurchasePrice, &v.CreatedAt, &v.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Vehicle{}, ErrNotFound
	}
	if err != nil {
		return models.Vehicle{}, err
	}
	if v.Fuel <= 0 {
		return models.Vehicle{}, ErrOutOfFuel
	}
	v.Fuel -= amount
	if v.Fuel < 0 {
		v.Fuel = 0
	}
	if err := tx.QueryRow(ctx, `
		UPDATE vehicles SET fuel = $1, updated_at = NOW() WHERE id = $2 RETURNING updated_at
	`, v.Fuel, v.ID).Scan(&v.UpdatedAt); err != nil {
		return models.Vehicle{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return models.Vehicle{}, err
	}
	return v, nil
}

// CompleteJobShift — MVP money loop (fiksi). Server decides payout + cooldown.
func (s *Store) CompleteJobShift(ctx context.Context, accountID, characterID uuid.UUID, jobKey string, payout int64, cooldownSec int) (models.Character, error) {
	if jobKey == "" || payout <= 0 || cooldownSec < 0 {
		return models.Character{}, ErrInvalidInput
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return models.Character{}, err
	}
	defer tx.Rollback(ctx)

	var c models.Character
	err = tx.QueryRow(ctx, `
		SELECT id, account_id, name, gender, skin_tone, hair_style, face_preset, outfit_id,
		       cash, bank, pos_x, pos_y, pos_z, created_at, updated_at
		FROM characters WHERE id = $1 AND account_id = $2
		FOR UPDATE
	`, characterID, accountID).Scan(
		&c.ID, &c.AccountID, &c.Name, &c.Gender, &c.SkinTone, &c.HairStyle, &c.FacePreset, &c.OutfitID,
		&c.Cash, &c.Bank, &c.PosX, &c.PosY, &c.PosZ, &c.CreatedAt, &c.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Character{}, ErrNotFound
	}
	if err != nil {
		return models.Character{}, err
	}

	var lastCompleted time.Time
	err = tx.QueryRow(ctx, `
		SELECT last_completed_at FROM job_cooldowns
		WHERE character_id = $1 AND job_key = $2
		FOR UPDATE
	`, c.ID, jobKey).Scan(&lastCompleted)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return models.Character{}, err
	}
	if err == nil {
		elapsed := time.Since(lastCompleted)
		if elapsed < time.Duration(cooldownSec)*time.Second {
			return models.Character{}, ErrJobOnCooldown
		}
	}

	c.Cash += payout
	if err := tx.QueryRow(ctx, `
		UPDATE characters SET cash = $1, updated_at = NOW() WHERE id = $2 RETURNING updated_at
	`, c.Cash, c.ID).Scan(&c.UpdatedAt); err != nil {
		return models.Character{}, err
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO job_cooldowns (character_id, job_key, last_completed_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (character_id, job_key)
		DO UPDATE SET last_completed_at = NOW()
	`, c.ID, jobKey); err != nil {
		return models.Character{}, err
	}

	meta := fmt.Sprintf(`{"job_key":%q}`, jobKey)
	if _, err := tx.Exec(ctx, `
		INSERT INTO transactions (character_id, kind, amount, balance_cash, balance_bank, meta)
		VALUES ($1, 'job_payout', $2, $3, $4, $5::jsonb)
	`, c.ID, payout, c.Cash, c.Bank, meta); err != nil {
		return models.Character{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return models.Character{}, err
	}
	return c, nil
}

type TransferResult struct {
	From models.Character `json:"from"`
	To   models.Character `json:"to"`
}

// TransferCash moves cash from one character to another (P2P). Amount > 0.
func (s *Store) TransferCash(ctx context.Context, fromAccountID, fromCharacterID, toCharacterID uuid.UUID, amount int64) (TransferResult, error) {
	if amount <= 0 {
		return TransferResult{}, ErrInvalidInput
	}
	if fromCharacterID == toCharacterID {
		return TransferResult{}, ErrInvalidInput
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return TransferResult{}, err
	}
	defer tx.Rollback(ctx)

	// Lock both characters in stable UUID order to avoid deadlocks.
	first, second := fromCharacterID, toCharacterID
	if second.String() < first.String() {
		first, second = second, first
	}
	if _, err := tx.Exec(ctx, `SELECT id FROM characters WHERE id IN ($1,$2) FOR UPDATE`, first, second); err != nil {
		return TransferResult{}, err
	}

	var from models.Character
	err = tx.QueryRow(ctx, `
		SELECT id, account_id, name, gender, skin_tone, hair_style, face_preset, outfit_id,
		       cash, bank, pos_x, pos_y, pos_z, created_at, updated_at
		FROM characters WHERE id = $1 AND account_id = $2
	`, fromCharacterID, fromAccountID).Scan(
		&from.ID, &from.AccountID, &from.Name, &from.Gender, &from.SkinTone, &from.HairStyle, &from.FacePreset, &from.OutfitID,
		&from.Cash, &from.Bank, &from.PosX, &from.PosY, &from.PosZ, &from.CreatedAt, &from.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return TransferResult{}, ErrNotFound
	}
	if err != nil {
		return TransferResult{}, err
	}

	var to models.Character
	err = tx.QueryRow(ctx, `
		SELECT id, account_id, name, gender, skin_tone, hair_style, face_preset, outfit_id,
		       cash, bank, pos_x, pos_y, pos_z, created_at, updated_at
		FROM characters WHERE id = $1
	`, toCharacterID).Scan(
		&to.ID, &to.AccountID, &to.Name, &to.Gender, &to.SkinTone, &to.HairStyle, &to.FacePreset, &to.OutfitID,
		&to.Cash, &to.Bank, &to.PosX, &to.PosY, &to.PosZ, &to.CreatedAt, &to.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return TransferResult{}, ErrNotFound
	}
	if err != nil {
		return TransferResult{}, err
	}

	if from.Cash < amount {
		return TransferResult{}, ErrInsufficientFunds
	}

	from.Cash -= amount
	to.Cash += amount

	if err := tx.QueryRow(ctx, `
		UPDATE characters SET cash = $1, updated_at = NOW() WHERE id = $2 RETURNING updated_at
	`, from.Cash, from.ID).Scan(&from.UpdatedAt); err != nil {
		return TransferResult{}, err
	}
	if err := tx.QueryRow(ctx, `
		UPDATE characters SET cash = $1, updated_at = NOW() WHERE id = $2 RETURNING updated_at
	`, to.Cash, to.ID).Scan(&to.UpdatedAt); err != nil {
		return TransferResult{}, err
	}

	metaFrom := fmt.Sprintf(`{"to_character_id":%q,"to_name":%q}`, to.ID.String(), to.Name)
	metaTo := fmt.Sprintf(`{"from_character_id":%q,"from_name":%q}`, from.ID.String(), from.Name)
	if _, err := tx.Exec(ctx, `
		INSERT INTO transactions (character_id, kind, amount, balance_cash, balance_bank, meta)
		VALUES ($1, 'transfer_out', $2, $3, $4, $5::jsonb)
	`, from.ID, amount, from.Cash, from.Bank, metaFrom); err != nil {
		return TransferResult{}, err
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO transactions (character_id, kind, amount, balance_cash, balance_bank, meta)
		VALUES ($1, 'transfer_in', $2, $3, $4, $5::jsonb)
	`, to.ID, amount, to.Cash, to.Bank, metaTo); err != nil {
		return TransferResult{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return TransferResult{}, err
	}
	return TransferResult{From: from, To: to}, nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
