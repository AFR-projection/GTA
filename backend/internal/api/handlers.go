package api

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/AFR-projection/GTA/backend/internal/auth"
	"github.com/AFR-projection/GTA/backend/internal/catalog"
	"github.com/AFR-projection/GTA/backend/internal/httpx"
	"github.com/AFR-projection/GTA/backend/internal/middleware"
	"github.com/AFR-projection/GTA/backend/internal/models"
	"github.com/AFR-projection/GTA/backend/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	Store  *store.Store
	Tokens *auth.TokenService
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		httpx.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	r.Route("/v1", func(r chi.Router) {
		r.Post("/auth/register", h.Register)
		r.Post("/auth/login", h.Login)

		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(h.Tokens))
			r.Get("/me", h.Me)
			r.Get("/characters", h.ListCharacters)
			r.Post("/characters", h.CreateCharacter)
			r.Get("/characters/{id}", h.GetCharacter)
			r.Post("/characters/{id}/bank/deposit", h.DepositBank)
			r.Post("/characters/{id}/bank/withdraw", h.WithdrawBank)
			r.Get("/characters/{id}/inventory", h.ListInventory)
			r.Post("/characters/{id}/shops/{shopID}/purchase", h.PurchaseFromShop)
			r.Get("/shops/warung", h.ListWarungCatalog)
			r.Get("/housing/listings", h.ListHousing)
			r.Get("/characters/{id}/houses", h.ListHouses)
			r.Post("/characters/{id}/houses/buy", h.BuyHouse)
		})
	})

	return r
}

type registerRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	Account   any       `json:"account"`
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid json body")
		return
	}
	req.Email = strings.TrimSpace(req.Email)
	req.DisplayName = strings.TrimSpace(req.DisplayName)
	if req.Email == "" || len(req.Password) < 8 || req.DisplayName == "" {
		httpx.Error(w, http.StatusBadRequest, "email, display_name, and password (min 8) required")
		return
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "could not hash password")
		return
	}

	account, err := h.Store.CreateAccount(r.Context(), req.Email, hash, req.DisplayName)
	if errors.Is(err, store.ErrConflict) {
		httpx.Error(w, http.StatusConflict, "email already registered")
		return
	}
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "could not create account")
		return
	}

	token, exp, err := h.Tokens.Issue(account.ID, account.Email)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "could not issue token")
		return
	}

	httpx.JSON(w, http.StatusCreated, authResponse{
		Token:     token,
		ExpiresAt: exp,
		Account:   account,
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid json body")
		return
	}

	account, err := h.Store.GetAccountByEmail(r.Context(), req.Email)
	if errors.Is(err, store.ErrNotFound) || !auth.CheckPassword(account.PasswordHash, req.Password) {
		httpx.Error(w, http.StatusUnauthorized, "invalid email or password")
		return
	}
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "login failed")
		return
	}

	token, exp, err := h.Tokens.Issue(account.ID, account.Email)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "could not issue token")
		return
	}

	httpx.JSON(w, http.StatusOK, authResponse{
		Token:     token,
		ExpiresAt: exp,
		Account:   account,
	})
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	accountID, ok := middleware.AccountIDFromContext(r.Context())
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	account, err := h.Store.GetAccountByID(r.Context(), accountID)
	if errors.Is(err, store.ErrNotFound) {
		httpx.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "could not load account")
		return
	}
	httpx.JSON(w, http.StatusOK, account)
}

type createCharacterRequest struct {
	Name       string `json:"name"`
	Gender     string `json:"gender"`
	SkinTone   int16  `json:"skin_tone"`
	HairStyle  int16  `json:"hair_style"`
	FacePreset int16  `json:"face_preset"`
	OutfitID   string `json:"outfit_id"`
}

func (h *Handler) CreateCharacter(w http.ResponseWriter, r *http.Request) {
	accountID, ok := middleware.AccountIDFromContext(r.Context())
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req createCharacterRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid json body")
		return
	}

	c, err := h.Store.CreateCharacter(r.Context(), accountID, store.CreateCharacterInput{
		Name:       req.Name,
		Gender:     req.Gender,
		SkinTone:   req.SkinTone,
		HairStyle:  req.HairStyle,
		FacePreset: req.FacePreset,
		OutfitID:   req.OutfitID,
	})
	switch {
	case errors.Is(err, store.ErrInvalidInput):
		httpx.Error(w, http.StatusBadRequest, "invalid character payload")
	case errors.Is(err, store.ErrLimitReached):
		httpx.Error(w, http.StatusConflict, "mvp allows only 1 character per account")
	case errors.Is(err, store.ErrConflict):
		httpx.Error(w, http.StatusConflict, "character name already used")
	case err != nil:
		httpx.Error(w, http.StatusInternalServerError, "could not create character")
	default:
		httpx.JSON(w, http.StatusCreated, c)
	}
}

func (h *Handler) ListCharacters(w http.ResponseWriter, r *http.Request) {
	accountID, ok := middleware.AccountIDFromContext(r.Context())
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	list, err := h.Store.ListCharacters(r.Context(), accountID)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "could not list characters")
		return
	}
	httpx.JSON(w, http.StatusOK, map[string]any{"characters": list})
}

func (h *Handler) GetCharacter(w http.ResponseWriter, r *http.Request) {
	accountID, ok := middleware.AccountIDFromContext(r.Context())
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid character id")
		return
	}
	c, err := h.Store.GetCharacterForAccount(r.Context(), accountID, id)
	if errors.Is(err, store.ErrNotFound) {
		httpx.Error(w, http.StatusNotFound, "character not found")
		return
	}
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "could not load character")
		return
	}
	httpx.JSON(w, http.StatusOK, c)
}

type moneyRequest struct {
	Amount int64 `json:"amount"`
}

func (h *Handler) DepositBank(w http.ResponseWriter, r *http.Request) {
	h.mutateBank(w, r, true)
}

func (h *Handler) WithdrawBank(w http.ResponseWriter, r *http.Request) {
	h.mutateBank(w, r, false)
}

func (h *Handler) mutateBank(w http.ResponseWriter, r *http.Request, deposit bool) {
	accountID, ok := middleware.AccountIDFromContext(r.Context())
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid character id")
		return
	}
	var req moneyRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid json body")
		return
	}

	var c models.Character
	if deposit {
		c, err = h.Store.DepositBank(r.Context(), accountID, id, req.Amount)
	} else {
		c, err = h.Store.WithdrawBank(r.Context(), accountID, id, req.Amount)
	}

	switch {
	case errors.Is(err, store.ErrInvalidInput):
		httpx.Error(w, http.StatusBadRequest, "amount must be > 0")
	case errors.Is(err, store.ErrNotFound):
		httpx.Error(w, http.StatusNotFound, "character not found")
	case errors.Is(err, store.ErrInsufficientFunds):
		httpx.Error(w, http.StatusConflict, "insufficient funds")
	case err != nil:
		httpx.Error(w, http.StatusInternalServerError, "bank operation failed")
	default:
		httpx.JSON(w, http.StatusOK, c)
	}
}

func (h *Handler) ListWarungCatalog(w http.ResponseWriter, _ *http.Request) {
	httpx.JSON(w, http.StatusOK, map[string]any{
		"shop_id": catalog.WarungShopID,
		"items":   catalog.Warung,
	})
}

func (h *Handler) ListInventory(w http.ResponseWriter, r *http.Request) {
	accountID, ok := middleware.AccountIDFromContext(r.Context())
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid character id")
		return
	}
	items, err := h.Store.ListInventory(r.Context(), accountID, id)
	if errors.Is(err, store.ErrNotFound) {
		httpx.Error(w, http.StatusNotFound, "character not found")
		return
	}
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "could not list inventory")
		return
	}
	httpx.JSON(w, http.StatusOK, map[string]any{"items": items})
}

type purchaseRequest struct {
	ItemKey  string `json:"item_key"`
	Quantity int    `json:"quantity"`
}

func (h *Handler) PurchaseFromShop(w http.ResponseWriter, r *http.Request) {
	accountID, ok := middleware.AccountIDFromContext(r.Context())
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid character id")
		return
	}
	shopID := chi.URLParam(r, "shopID")
	if shopID != catalog.WarungShopID {
		httpx.Error(w, http.StatusNotFound, "shop not found")
		return
	}

	var req purchaseRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid json body")
		return
	}
	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	item, ok := catalog.WarungByKey(req.ItemKey)
	if !ok {
		httpx.Error(w, http.StatusBadRequest, "unknown item_key for this shop")
		return
	}

	result, err := h.Store.PurchaseItem(r.Context(), accountID, id, shopID, item.Key, item.Price, req.Quantity)
	switch {
	case errors.Is(err, store.ErrInvalidInput):
		httpx.Error(w, http.StatusBadRequest, "invalid purchase")
	case errors.Is(err, store.ErrNotFound):
		httpx.Error(w, http.StatusNotFound, "character not found")
	case errors.Is(err, store.ErrInsufficientFunds):
		httpx.Error(w, http.StatusConflict, "insufficient cash")
	case err != nil:
		httpx.Error(w, http.StatusInternalServerError, "purchase failed")
	default:
		httpx.JSON(w, http.StatusOK, result)
	}
}

func (h *Handler) ListHousing(w http.ResponseWriter, _ *http.Request) {
	httpx.JSON(w, http.StatusOK, map[string]any{
		"listings": catalog.Housing,
		"note":     "MVP: 1 house per character (buy ownership)",
	})
}

func (h *Handler) ListHouses(w http.ResponseWriter, r *http.Request) {
	accountID, ok := middleware.AccountIDFromContext(r.Context())
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid character id")
		return
	}
	houses, err := h.Store.ListHouses(r.Context(), accountID, id)
	if errors.Is(err, store.ErrNotFound) {
		httpx.Error(w, http.StatusNotFound, "character not found")
		return
	}
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "could not list houses")
		return
	}
	httpx.JSON(w, http.StatusOK, map[string]any{"houses": houses})
}

type buyHouseRequest struct {
	ListingKey string `json:"listing_key"`
}

func (h *Handler) BuyHouse(w http.ResponseWriter, r *http.Request) {
	accountID, ok := middleware.AccountIDFromContext(r.Context())
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid character id")
		return
	}
	var req buyHouseRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid json body")
		return
	}
	listing, found := catalog.HousingByKey(req.ListingKey)
	if !found {
		httpx.Error(w, http.StatusBadRequest, "unknown listing_key")
		return
	}

	result, err := h.Store.BuyHouse(
		r.Context(), accountID, id,
		listing.Key, listing.Label, listing.Price,
		listing.PosX, listing.PosY, listing.PosZ,
	)
	switch {
	case errors.Is(err, store.ErrInvalidInput):
		httpx.Error(w, http.StatusBadRequest, "invalid purchase")
	case errors.Is(err, store.ErrNotFound):
		httpx.Error(w, http.StatusNotFound, "character not found")
	case errors.Is(err, store.ErrInsufficientFunds):
		httpx.Error(w, http.StatusConflict, "insufficient cash")
	case errors.Is(err, store.ErrAlreadyOwnsHouse):
		httpx.Error(w, http.StatusConflict, "mvp allows only 1 house per character")
	case errors.Is(err, store.ErrConflict):
		httpx.Error(w, http.StatusConflict, "house already owned")
	case err != nil:
		httpx.Error(w, http.StatusInternalServerError, "house purchase failed")
	default:
		httpx.JSON(w, http.StatusCreated, result)
	}
}
