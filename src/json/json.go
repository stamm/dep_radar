package json

type Result struct {
	App  string
	Libs []Lib
}

type Lib struct {
	Version string
	Hash    string
}
