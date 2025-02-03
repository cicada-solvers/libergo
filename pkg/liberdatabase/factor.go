package liberdatabase

type Factor struct {
	ID        string `json:"id"`
	Factor    string `json:"factor"`
	MainId    string `json:"mainid"`
	SeqNumber int64  `json:"seqnumber"`
}
