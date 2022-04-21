package spaggiari

type Identity struct {
	Ident     string `json:"ident"`
	ID        string
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Token     string `json:"token"`
	Release   string `json:"release"`
	Expire    string `json:"expire"`
}

type Grade struct {
	Subject       string  `json:"subjectDesc"`
	Date          string  `json:"evtDate"`
	DecimalValue  float32 `json:"decimalValue"`
	DisplaylValue string  `json:"displayValue"`
	Color         string  `json:"color"`
	Description   string  `json:"skillValueDesc"`
	Notes         string  `json:"notesForFamily"`
}
