package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/baofengcloud/go-sdk/src/baofengcloud"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"strings"
)

var serviceType = flag.Int("service", 0, "service type, 0 Paas, 1 Saas")
var isPrivate = flag.Bool("private", false, "private file")

var commands = []string{"config", "query", "delete", "upload"}

var fileName, localFilePath string

type ConfFile struct {
	AccessKey string
	SecretKey string
}

var confFile ConfFile

var confFilePath string

func main() {

	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("please specify an command: ", strings.Join(commands, ", "))
		return
	}

	action := strings.ToLower(args[0])

	validCommand := false
	for _, c := range commands {
		if c == action {
			validCommand = true
			break
		}
	}

	if !validCommand {
		fmt.Println("unknown command ", action)
		return
	}

	user, _ := user.Current()

	confFilePath = path.Join(user.HomeDir, ".bfcloud")

	if fd, err := os.Open(confFilePath); err == nil {

		jsonStr, _ := ioutil.ReadAll(fd)

		json.Unmarshal(jsonStr, &confFile)
	}

	if action == "config" {

		fmt.Print("set access key:")
		fmt.Scanln(&confFile.AccessKey)
		fmt.Print("set secret key:")
		fmt.Scanln(&confFile.SecretKey)

		fd, err := os.Create(confFilePath)
		if err != nil {
			fmt.Println(err)
			return
		}

		defer fd.Close()

		jsonStr, _ := json.Marshal(&confFile)

		fd.WriteString(string(jsonStr))

		return
	}

	if len(confFile.AccessKey) == 0 || len(confFile.SecretKey) == 0 {
		fmt.Println("access key and secret key not set, run bfcloud config first!")
		return
	}

	conf := &baofengcloud.Configure{
		AccessKey: confFile.AccessKey,
		SecretKey: confFile.SecretKey,
	}

	var fileInfo *baofengcloud.FileInfo
	var result *baofengcloud.Result
	var err error

	if action == "query" {

		if len(args) < 2 {
			fmt.Println("please specify a name")
			return
		}

		fileInfo, err = baofengcloud.QueryFile(conf, baofengcloud.ServiceType(*serviceType), args[1], "")

	} else if action == "delete" {

		if len(args) < 2 {
			fmt.Println("please specify a name")
			return
		}

		result, err = baofengcloud.DeleteFile(conf, baofengcloud.ServiceType(*serviceType), args[1], "", "")

	} else if action == "upload" {

		if len(args) < 3 {
			fmt.Println("please specify a name and local file path: upload {dstname} {localfilepath}")
			return
		}

		var fileType baofengcloud.FileType

		if *isPrivate {
			fileType = baofengcloud.Private
		} else {
			fileType = baofengcloud.Public
		}

		err = baofengcloud.UploadFile2(conf, baofengcloud.ServiceType(*serviceType),
			fileType, args[2], args[1], "", "")
	}

	if err != nil {
		fmt.Println(err)
		return
	}

	if result != nil {
		fmt.Printf("%+v\n", *result)
	}

	if fileInfo != nil {
		fmt.Printf("%+v\n", *fileInfo)
	}
}
