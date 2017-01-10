package smawk

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"gopkg.in/telegram-bot-api.v4"
	"os/exec"
	"strconv"
)

// GenerateCertificate is used to create a self signed certificate for use with any
// instance of the bot
func GenerateCertificate(c string, st string, ct string, org string, dom string, key string, cert string) {
	// Generate our string for the certificate
	certstring := "/C="+c+"/ST="+st+"/L="+ct+"/O="+org+"/CN="+dom

	cmdname := "openssl"
	cmdargs := []string{"req","-newkey","rsa:2048","-sha256","-nodes","-keyout",key,"-x509","-days","365","-out",cert,"-subj",certstring}

	cmd := exec.Command(cmdname,cmdargs...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
	}
}

// ConnectDB takes care of opening a proper connection to the database to retrieve the scores that we need
func ConnectDB() (*sql.DB, error) {
	cfg := &mysql.Config {
		User: "smawk-bot",
		Passwd: "SM@WKisGR8",
		Net: "tcp",
		Addr: "107.170.45.12:3306",
		DBName: "smawk-bot",
	}
	return sql.Open("mysql", cfg.FormatDSN())
}

// EnterScore is responsible for updating the database with any upvote, downvote, bless, or curse commands
func EnterScore() {

}
