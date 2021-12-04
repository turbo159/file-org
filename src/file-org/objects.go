package main

const constTaskfile string = "./tasks.json"

type task struct {
	IsEnabled       bool     `json:"enabled"`
	Recursive       bool     `json:"recursive"`
	FileType        []string `json:"filetype"`
	FilePrefix      string   `json:"fileprefix"`
	ScriptPrefix    string   `json:"scriptprefix"`
	SourcePath      string   `json:"sourcepath"`
	DestinationPath string   `json:"destinationpath"`
	ScriptPath      string   `json:"scriptpath"`
}

type Tasklist struct {
	Tasks []task `json:"tasks"`
}

type fileObj struct {
	sha1hash   string
	duplicate  bool
	sourcepath string
	sourcename string
}
