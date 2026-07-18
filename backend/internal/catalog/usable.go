package catalog

// UsableItem can be consumed from inventory (makan/minum).
type UsableItem struct {
	Key         string `json:"item_key"`
	Label       string `json:"label"`
	Description string `json:"description"`
}

// Usable — subset of warung foods/drinks that can be "eaten".
var Usable = map[string]UsableItem{
	"nasi_bungkus": {Key: "nasi_bungkus", Label: "Nasi Bungkus", Description: "Mengisi perut"},
	"kopi_tubruk":  {Key: "kopi_tubruk", Label: "Kopi Tubruk", Description: "Minum kopi"},
	"air_mineral":  {Key: "air_mineral", Label: "Air Mineral", Description: "Minum air"},
	"mie_instan":   {Key: "mie_instan", Label: "Mie Instan", Description: "Makan mie"},
	"roti_bakar":   {Key: "roti_bakar", Label: "Roti Bakar", Description: "Camilan"},
}

func UsableByKey(key string) (UsableItem, bool) {
	u, ok := Usable[key]
	return u, ok
}
