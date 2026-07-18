package store

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/AFR-projection/GTA/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrConflict      = errors.New("conflict")
	ErrLimitReached  = errors.New("character limit reached")
	ErrInvalidInput  = errors.New("invalid input")
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

var ErrInsufficientFunds = errors.New("insufficient funds")

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

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
