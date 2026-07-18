package catalog

// JobShift is a simple MVP earn loop. Server decides payout (client cannot set amount).
type JobShift struct {
	Key         string `json:"job_key"`
	Label       string `json:"label"`
	Payout      int64  `json:"payout"`
	Description string `json:"description"`
}

var Jobs = []JobShift{
	{Key: "ojol_shift", Label: "Shift Ojol", Payout: 120, Description: "Antar penumpang keliling kota (simulasi)"},
	{Key: "kurir_shift", Label: "Shift Kurir", Payout: 100, Description: "Antar paket warung ke rumah"},
	{Key: "kasir_shift", Label: "Shift Kasir Minimarket", Payout: 80, Description: "Jaga kasir 1 shift"},
}

// JobCooldownSeconds — anti-spam payout (server-side).
const JobCooldownSeconds = 60

func JobByKey(key string) (JobShift, bool) {
	for _, j := range Jobs {
		if j.Key == key {
			return j, true
		}
	}
	return JobShift{}, false
}
