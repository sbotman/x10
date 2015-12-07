package api

type X10Device struct {
	Id          int32  `json:"id"`
	Created     int32  `json:"created"`
	Title       string `json:"title"`
	State       string `json:"state"`
	Room        string `json:"room"`
	Code		string `json:"code"`
}

const (
	OnStatus  string = "on"
	OffStatus string = "off"
)
