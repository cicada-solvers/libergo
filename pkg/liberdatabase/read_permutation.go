package liberdatabase

type ReadPermutation struct {
	ID                   string `json:"id"`
	StartArray           []byte `json:"start_array"`
	EndArray             []byte `json:"end_array"`
	PackageName          string `json:"package_name"`
	PermName             string `json:"perm_name"`
	ReportedToAPI        bool   `json:"reported_to_api"`
	Processed            bool   `json:"processed"`
	ArrayLength          int    `json:"array_length"`
	NumberOfPermutations int64  `json:"number_of_permutations"`
}
