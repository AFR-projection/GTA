package catalog

// ShopItem is a sellable catalog entry. MVP: in-code, not DB.
// Swap to DB later without changing purchase API contract.
type ShopItem struct {
	Key         string `json:"item_key"`
	Name        string `json:"name"`
	Price       int64  `json:"price"`
	Description string `json:"description"`
}

const WarungShopID = "warung"

// Warung is the starter convenience shop (fiksi Indonesia vibe).
var Warung = []ShopItem{
	{Key: "nasi_bungkus", Name: "Nasi Bungkus", Price: 50, Description: "Bekal warung pinggir jalan"},
	{Key: "kopi_tubruk", Name: "Kopi Tubruk", Price: 20, Description: "Kopi panas warkop"},
	{Key: "air_mineral", Name: "Air Mineral", Price: 10, Description: "Botol kecil"},
	{Key: "mie_instan", Name: "Mie Instan", Price: 30, Description: "Siap seduh"},
	{Key: "roti_bakar", Name: "Roti Bakar", Price: 25, Description: "Snack sore"},
}

func WarungByKey(key string) (ShopItem, bool) {
	for _, item := range Warung {
		if item.Key == key {
			return item, true
		}
	}
	return ShopItem{}, false
}
