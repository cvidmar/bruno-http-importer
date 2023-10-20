package importer

type BrunoConfig struct {
	Meta    *Meta    `json:"meta"`
	Call    *Call    `json:"call"`
	Headers []string `json:"headers"`
	Body    *Body    `json:"body"`
}

type Meta struct {
	Name string `json:"name"`
	Verb string `json:"verb"`
	Seq  int    `json:"seq"`
}

type Call struct {
	Url    string `json:"url"`
	Method string `json:"method"`
	Body   string `json:"body"`
}

type Body struct {
	Mode string   `json:"mode"`
	Raw  []string `json:"raw"`
}
