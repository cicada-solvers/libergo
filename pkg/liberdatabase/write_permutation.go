package liberdatabase

type WritePermutation struct {
	ID                   string `json:"id"`
	StartArray           string `json:"start_array"`
	EndArray             string `json:"end_array"`
	PackageName          string `json:"package_name"`
	PermName             string `json:"perm_name"`
	ReportedToAPI        bool   `json:"reported_to_api"`
	Processed            bool   `json:"processed"`
	ArrayLength          int    `json:"array_length"`
	NumberOfPermutations int64  `json:"number_of_permutations"`
}
