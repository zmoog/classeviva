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
	Description   string  `json:"skillValueDesc,omitempty"`
	Notes         string  `json:"notesForFamily,omitempty"`
}

type AgendaEntry struct {
	ID            int    `json:"evtId"`            // "evtId": 502508,
	Code          string `json:"evtCode"`          // "evtCode": "AGHW",
	DatetimeBegin string `json:"evtDatetimeBegin"` // "evtDatetimeBegin": "2022-04-04T09:00:00+02:00",
	DatetimeEnd   string `json:"evtDatetimeEnd"`   //	"evtDatetimeEnd": "2022-04-04T10:00:00+02:00",
	Notes         string `json:"notes"`            // "notes": "Page 201 tutti gli esercizi",
	AuthorName    string `json:"authorName"`       // "authorName": "PESANDO MARGHERITA",
	Subject       string `json:"subjectDesc"`      // "INGLESE",
	// "classDesc": "2E MUSICALE",
	// "subjectId": 396137,
	// "homeworkId": null
}
