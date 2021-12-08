package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

func init() {

	log.Info("*********************************************")
	log.Info("Starting file-org...")
	log.Info()
}

func main() {
	var filecache = make(map[string]fileObj)
	var taskFile string
	var iTaskCount int
	var flgHelp bool

	flag.StringVar(&taskFile, "t", constTaskfile, "Specify task file path/name.")
	flag.BoolVar(&flgHelp, "h", false, "Command help.")
	flag.Parse()

	if flgHelp {
		flag.PrintDefaults()
		os.Exit(1)
	}

	tl, e := loadTaskFile(taskFile)
	if e != nil {
		log.Error("Error loading task file: ", taskFile)
		os.Exit(0)
	}

	if len(tl.Tasks) < 1 {
		log.Info("Tasklist empty")
	}

	for _, t := range tl.Tasks {
		if t.IsEnabled {
			iTaskCount++
		}
	}
	var runScriptCollection string
	for i, value := range tl.Tasks {
		if value.IsEnabled {
			log.Info()
			log.Info("Executing task/enabled tasks: ", i+1, "/", iTaskCount)
			value.DestinationPath = cleanPath(value.DestinationPath)
			err := buidFileCache(filecache, value.SourcePath, value.FileType)
			if err != nil {
				log.Error("Build fileCache failed: ", err)
			} else {
				// Export scripts
				var copyscript, dupscript string
				log.Info("Preparing scripts...")
				for _, file := range filecache {
					if !file.duplicate {
						copyscript += "cp " + file.sourcepath + "/" + file.sourcename + " " + value.DestinationPath + "/" + value.FilePrefix + file.sourcename + "\n"
					} else {
						dupscript += "cp " + file.sourcepath + "/" + file.sourcename + " " + value.DestinationPath + "/" + value.FilePrefix + file.sourcename + "\n"
					}

				}
				//Create copy script
				if copyscript != "" {
					copyScriptName := fmt.Sprint(value.ScriptPath, "/", value.ScriptPrefix, constCopyTaskSuffix, i+1, ".sh")
					if err := os.WriteFile(copyScriptName, []byte(copyscript), 0777); err != nil {
						log.Error("Failed to create copy script. ", err)
					} else {
						runScriptCollection += fmt.Sprint("./", copyScriptName, "\n")
						log.Info("Created script: ", copyScriptName)
					}
				} else {
					log.Info("No files to copy.")
				}
				//Create duplicate script
				if dupscript != "" {
					dupScriptName := fmt.Sprint(value.ScriptPath, "/", value.ScriptPrefix, constDuplicateTaskSuffix, i+1, ".sh")
					if err := os.WriteFile(dupScriptName, []byte(dupscript), 0777); err != nil {
						log.Error("Failed to create duplicate script. ", err)
					} else {
						runScriptCollection += fmt.Sprint("./", dupScriptName, "\n")
						log.Info("Created script: ", dupScriptName)
					}
				} else {
					log.Info("No duplicate files to copy.")
				}
			}
			log.Info("Task complete.")
		}
	}

	//Summary and stats
	var filecount, duplicatecount int16
	for _, file := range filecache {
		if file.duplicate {
			duplicatecount++
		} else {
			filecount++
		}
	}

	log.Info()
	log.Info("---------------------------------------------")
	log.Info("Total cached files: ", len(filecache))
	log.Info("New: ", filecount)
	log.Info("Duplicates: ", duplicatecount)
	log.Info()
	log.Info("Finished file-org.")
	log.Info("*********************************************")
}

func buidFileCache(cache map[string]fileObj, startpath string, filetype []string) error {
	var throwAway int32 = 0
	err := filepath.WalkDir(startpath,
		func(pathstring string, info os.DirEntry, err error) error {

			if !info.IsDir() && contains(filetype, filepath.Ext(info.Name())) {
				//Generage hash string for the file
				hash := getSha1Checksum(pathstring)

				//Check if key already exists
				item, exists := cache[hash]
				if item.sourcepath == path.Dir(pathstring) {
					throwAway++
				} else {
					cache[hash] = fileObj{
						sha1hash:   hash,
						duplicate:  exists,
						sourcepath: cleanPath(path.Dir(pathstring)),
						sourcename: info.Name()}
				}
			}
			return nil
		})
	if err != nil {
		log.Error("Unexpected error building cache: ", err)
		return err
	}
	log.Info("Items thrown away: ", throwAway)
	return nil
}

//cleanPath cleans a path replacing spaces and invalid characters.
func cleanPath(path string) (cleanPath string) {
	path = strings.ReplaceAll(path, " ", "\\ ")
	path = strings.ReplaceAll(path, "(", ".")
	path = strings.ReplaceAll(path, ")", ".")
	return path
}

// contains
func contains(s []string, str string) bool {
	for _, v := range s {
		if strings.EqualFold(v, str) {
			return true
		}
	}
	return false
}

//getSha1Checksum returns the hash value of the file
func getSha1Checksum(filePath string) (hashstring string) {
	var returnSHA1String string
	file, err := os.Open(filePath)
	if err != nil {
		return returnSHA1String
	}
	defer file.Close()
	hash := sha1.New()
	if _, err := io.Copy(hash, file); err != nil {
		return returnSHA1String
	}
	hashInBytes := hash.Sum(nil)[:20]
	returnSHA1String = hex.EncodeToString(hashInBytes)
	return returnSHA1String
}

// loadTaskFile reads takes a taskfile path and loads the configuration into a TaskList struct.
func loadTaskFile(taskFile string) (exeTasks Tasklist, e error) {
	dat, _ := ioutil.ReadFile(taskFile)

	err := json.Unmarshal(dat, &exeTasks)
	if err != nil {
		return exeTasks, err
	}

	log.Info("Task file loaded: ", taskFile)
	return exeTasks, nil
}
