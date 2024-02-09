package records

type Record struct {
	ID       int    `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Site     string `json:"site,omitempty"`
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
	Info     string `json:"info,omitempty"`
}
