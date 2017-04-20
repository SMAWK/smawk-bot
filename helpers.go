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

// GenerateUpdate is a helper function that generates a test update
// (see Update sruct in tgbotapi/types). This function can be called
// from the test files of programs that implement this library
func (bot *SmawkBot) GenerateUpdate(cmd string) (tgbotapi.Update, error) {
    // Create our Update Var
    var upd tgbotapi.Update

    // Create our JSON blob
    var updjson = []byte(`{
        "update_id":322176086,
        "message":{
            "message_id":178,
            "from":{
                "id":55997207,
                "first_name":"Benjamin",
                "last_name":"Matthews",
                "username":"bnmtthews"
            },
            "chat":{
                "id":55997207,
                "first_name":"Benjamin",
                "last_name":"Matthews",
                "username":"bnmtthews",
                "type":"private"
            },
            "date":1468013062,
            "text":"`+cmd+`",
            "entities":[{
                "type":"bot_command",
                "offset":0,
                "length":`+strconv.Itoa(len(cmd))+`
            }]
        }
    }`)

    // Create our update
    json.Unmarshal(updjson, &upd)

    // Return our update
    return upd, nil
}

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
func ConnectDB(password string) (*sql.DB, error) {
    cfg := &mysql.Config {
        User: "smawk-bot",
        Passwd: password,
        Net: "tcp",
        Addr: "107.170.45.12:3306",
        DBName: "smawk-bot",
    }
    return sql.Open("mysql", cfg.FormatDSN())
}
