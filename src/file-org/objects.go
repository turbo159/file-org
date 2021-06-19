package main

type config struct {
	Simulation bool   `json:"simulation"`
	Loglevel   int    `json:"loglevel"`
	SafeMove   bool   `json:"safemove"`
	Taskfile   string `json:"taskfile"`
}

type taskfile struct {
	IsEnabled  bool     `json:"enabled"`
	Sourcepath string   `json:"sourcepath"`
	Recursive  bool     `json:"recursive"`
	Filetype   []string `json:"filetype"`
	Destpath   string   `json:"destpath"`
}

type Tasklist struct {
	Tasks []taskfile `json:"tasks"`
}

type fileObj struct {
	sourcepath string
	sourcename string
	sha1hash   string
	targetpath string
	targetname string
}
