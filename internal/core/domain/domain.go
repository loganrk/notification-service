package domain

type UserActivation struct {
	Type    string            `json:"type"`
	To      string            `json:"to"`
	Subject string            `json:"subject"`
	Macros  map[string]string `json:"macros"`
}

type UserPasswordReset struct {
	Type    string            `json:"type"`
	To      string            `json:"to"`
	Subject string            `json:"subject"`
	Macros  map[string]string `json:"macros"`
}
