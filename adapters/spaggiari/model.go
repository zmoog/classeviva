package spaggiari

type Grade struct {
	Subject       string  `json:"subjectDesc"`
	Date          string  `json:"evtDate"`
	DecimalValue  float32 `json:"decimalValue"`
	DisplaylValue string  `json:"displayValue"`
	Color         string  `json:"color"`
	Description   string  `json:"skillValueDesc"`
	Notes         string  `json:"notesForFamily"`
}
