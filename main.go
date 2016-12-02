package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"regexp"
	"strings"
)

const (
	url     = "http://www.mikrotik.com/client/ajax.php?action=getChangelog&id=24"
	pattern = `(?is)What's\s+new\s+in\s+([\d\.]+)[^\n]*\n(.*?)\nWhat's\s+new\s+in`
)

var mailTo = flag.String("to", "", "Mail 'To:' header")
var mailFrom = flag.String("from", "", "Mail 'From:' header")
var mailAddr = flag.String("addr", "smtp.yandex.ru:25", "SMTP server address")
var mailUsername = flag.String("username", "", "SMTP Auth username (if isn't set it takes the 'from' flag value.)")
var mailPassword = flag.String("password", "", "SMTP Auth password")
var versionStorage = flag.String("store", "$HOME/.mikrocheck", "File to store current version")

func init() {
	flag.Parse()
	if *versionStorage == "$HOME/.mikrocheck" {
		usr, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		*versionStorage = usr.HomeDir + "/.mikrocheck"
	}
}

func main() {
	m := newMail(*mailTo, *mailFrom, *mailAddr, *mailUsername, *mailPassword)

	ver, info := getLastVersion()
	if isNew(ver, *versionStorage) {
		m.Notify(ver, info)
		storeVersion(ver, *versionStorage)
		log.Printf("New version is available: %s", ver)
	} else {
		log.Printf("Version isn't changed: %s", ver)
	}
}

// storeVersion save version to the storage file.
func storeVersion(ver, storage string) {
	err := ioutil.WriteFile(storage, []byte(ver), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

// isNew checks whether the version is changed.
// It takes version and storage path params.
func isNew(ver, storage string) bool {
	if _, err := os.Stat(storage); err != nil {
		if os.IsNotExist(err) {
			return true
		} else {
			log.Fatal(err)
		}
	}

	data, err := ioutil.ReadFile(storage)
	if err != nil {
		log.Fatal(err)
	}
	return ver != string(data)
}

// getLastVersion gets the latest Mikrotik Firmware version.
func getLastVersion() (ver, info string) {
	re := regexp.MustCompile(pattern)
	match := re.FindAllStringSubmatch(getPage(), 2)
	if len(match) > 0 {
		return match[0][1], strings.TrimSpace(match[0][2])
	}
	return "", ""
}

// getPage loads Mikrotik change log.
func getPage() string {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	return string(body)
}
