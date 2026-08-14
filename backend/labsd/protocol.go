package labsd

type Request struct {
	Operation string `json:"operation"`
	App       string `json:"app,omitempty"`
	Lines     int    `json:"lines,omitempty"`
	Compose   string `json:"compose,omitempty"`
}

type Response struct {
	OK      bool     `json:"ok"`
	Message string   `json:"message,omitempty"`
	Status  string   `json:"status,omitempty"`
	Apps    []string `json:"apps,omitempty"`
	Logs    string   `json:"logs,omitempty"`
}
