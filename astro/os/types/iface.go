package types

type Iface struct {
	Interface string `json:"interface"`
	Gateway   string `json:"gateway,omitempty"`
	IP        string `json:"ip"`
	IsPublic  bool   `json:"is_public"`
	Network   string `json:"network"`
	Mask      string `json:"mask"`
}
