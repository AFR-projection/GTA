package catalog

// VehicleListing is a purchasable vehicle (MVP: in-code catalog, original fiction designs).
type VehicleListing struct {
	Key         string  `json:"listing_key"`
	Label       string  `json:"label"`
	Type        string  `json:"vehicle_type"` // motorcycle | car
	Price       int64   `json:"price"`
	FuelMax     float64 `json:"fuel_max"`
	Description string  `json:"description"`
}

var Vehicles = []VehicleListing{
	{
		Key:         "motor_bebek_kota",
		Label:       "Motor Bebek Kota",
		Type:        "motorcycle",
		Price:       180,
		FuelMax:     100,
		Description: "Bebek harian buat nyelip macet (desain fiksi)",
	},
	{
		Key:         "motor_matic_ringan",
		Label:       "Motor Matic Ringan",
		Type:        "motorcycle",
		Price:       220,
		FuelMax:     100,
		Description: "Matic ringan antar gang (desain fiksi)",
	},
	{
		Key:         "mobil_city_hatch",
		Label:       "Mobil City Hatch",
		Type:        "car",
		Price:       500,
		FuelMax:     100,
		Description: "Hatchback kota hemat (desain fiksi)",
	},
	{
		Key:         "mobil_family_mpv",
		Label:       "Mobil Family MPV",
		Type:        "car",
		Price:       750,
		FuelMax:     120,
		Description: "MPV keluarga jalan tol (desain fiksi)",
	},
}

func VehicleByKey(key string) (VehicleListing, bool) {
	for _, v := range Vehicles {
		if v.Key == key {
			return v, true
		}
	}
	return VehicleListing{}, false
}

// FuelPricePerUnit — harga isi bensin di SPBU fiksi (per 1 fuel unit).
const FuelPricePerUnit int64 = 2

const SPBUShopID = "spbu"
