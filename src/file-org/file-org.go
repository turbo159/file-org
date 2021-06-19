package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.Info("In init")
}

func main() {
	log.Info("In main")

	cf, e := loadConfigFile("./config.json")
	if e != nil {
		log.Error("Errors", e)
	}

	tl, e := loadTaskFile(cf.Taskfile)
	if e != nil {
		log.Error("Errors")
	}

	if len(tl.Tasks) > 0 {
		// log.Info("Tasklist: ")
		log.Printf("Tasklist: "+"%+v\n", tl)
	} else {
		log.Info("Tasklist totally empty")
	}

	var fileObjList []fileObj

	for _, value := range tl.Tasks {
		if value.IsEnabled {
			x := getFileList(value.Sourcepath, value.Filetype, cf.SafeMove)

			for _, obj := range x {
				fileObjList = append(fileObjList, obj)
			}
		}
	}

	for i, file := range fileObjList {
		fmt.Println(i, file.sha1hash, file.sourcename, file.sourcepath)
	}

}

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

// getFileList returns a list of files from the received path
func getFileList(startpath string, filetype []string, checksum bool) (files []fileObj) {
	err := filepath.Walk(startpath,
		func(pathstring string, info os.FileInfo, err error) error {
			// if err != nil {
			// 	//return err
			// }

			if !info.IsDir() && contains(filetype, filepath.Ext(info.Name())) {
				currentpath, _ := os.Getwd()
				var hash string
				if checksum {
					hash = getSha1Checksum(pathstring)
				} else {
					hash = ""
				}
				files = append(files, fileObj{sourcepath: currentpath, sourcename: info.Name(), sha1hash: hash})
			}

			return nil
		})
	if err != nil {
		log.Println(err)
	}
	return files
}

// loadTaskFile reads takes a taskfile path and loads the configuration into a TaskList struct.
func loadTaskFile(taskFile string) (exeTasks Tasklist, e error) {
	dat, _ := ioutil.ReadFile(taskFile)

	err := json.Unmarshal(dat, &exeTasks)
	if err != nil {
		log.Error("JSON import error: ", err)
		return exeTasks, err
	}

	log.Info("loadTaskFile complete")
	return exeTasks, nil
}

// loadConfigFile reads a configuration file toc initialize the program.
func loadConfigFile(configFile string) (c config, e error) {
	dat, _ := ioutil.ReadFile(configFile)

	err := json.Unmarshal(dat, &c)
	if err != nil {
		log.Error("Config file import error: ", err)
		return c, err
	}
	log.Info("loadConfigFile complete: ", c)

	return c, nil
}
