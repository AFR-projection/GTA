package catalog

// HouseListing is a purchasable home definition (MVP: in-code catalog).
type HouseListing struct {
	Key         string  `json:"listing_key"`
	Label       string  `json:"label"`
	Price       int64   `json:"price"`
	Description string  `json:"description"`
	PosX        float64 `json:"pos_x"`
	PosY        float64 `json:"pos_y"`
	PosZ        float64 `json:"pos_z"`
}

// Housing listings — nuansa Indonesia, desain fiksi.
var Housing = []HouseListing{
	{
		Key:         "kontrakan_gang",
		Label:       "Kontrakan Gang",
		Price:       200,
		Description: "Kontrakan 1 pintu di gang sempit",
		PosX:        120, PosY: 0, PosZ: -40,
	},
	{
		Key:         "subsidi_blok_a",
		Label:       "Rumah Subsidi Blok A",
		Price:       350,
		Description: "Rumah tapak sederhana perumahan subsidi",
		PosX:        250, PosY: 0, PosZ: 80,
	},
	{
		Key:         "ruko_pinggir",
		Label:       "Ruko Pinggir Jalan",
		Price:       450,
		Description: "Ruko 2 lantai buat tinggal + usaha kecil",
		PosX:        40, PosY: 0, PosZ: 160,
	},
}

func HousingByKey(key string) (HouseListing, bool) {
	for _, h := range Housing {
		if h.Key == key {
			return h, true
		}
	}
	return HouseListing{}, false
}
