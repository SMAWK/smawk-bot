package smawk_test

import (
	//"bytes"
	//"database/sql"
	"encoding/json"
	"fmt"
	"github.com/bmatt468/smawk-bot"
	//"github.com/go-sql-driver/mysql"
	"gopkg.in/telegram-bot-api.v4"
	"log"
	//"os"
	//"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"
)

// Create our constants for use throughout the testing functions
const (
	SMAWKToken              = "249930361:AAHz1Gksb-eT0SQG47lDb7WbJxujr7kGCkU"
	ChatID                 	= 55997207
)

/* ================================================ */
/*                 Helper functions                 */
/* ================================================ */

// generateUpdate is a helper function that generates a test update
// (see Update sruct in tgbotapi/types)
func GenerateUpdate(cmd string) (tgbotapi.Update) {
	// Create our Update Var
	var upd tgbotapi.Update

	// Create our JSON blob
	var updjson = []byte(`{
		"update_id":322176086,
		"message":{
			"message_id":178,
			"from":{
				"id":`+strconv.Itoa(ChatID)+`,
				"first_name":"Benjamin",
				"last_name":"Matthews",
				"username":"bnmtthews"
			},
			"chat":{
				"id":`+strconv.Itoa(ChatID)+`,
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
	return upd
}

func timestamp() string {
	t := time.Now()
	return t.Format("2006-01-02 15:04:05.000 ");
}

/* ================================================ */
/*                Testing functions                 */
/* ================================================ */
func TestBot(t *testing.T) {
	// ==================== //
	// Start Helper Tests   //
	// ==================== //
	fmt.Println("======= Starting Helper Tests =======")

	/** === Loading Bot === **/
	fmt.Print(timestamp()+"Loading SMÄWK_bot.... ")

	// Fetch our bot using the helper function
	bot, err := smawk.Connect(SMAWKToken,false)
	bot.Testing = true;

	// Check to see if something bad happened and break if need be
	if err != nil {
		log.Fatal(err)
		t.FailNow()
	}

	// Otherwise, log to the console that we authenticated properly
	fmt.Println("done")

	/** === Database connection === **/
	fmt.Print(timestamp()+"Connecting to database.... ")
	db, err := smawk.ConnectDB()
	if err != nil {
		log.Fatal(err)
		t.FailNow()
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
		t.FailNow()
	}

	fmt.Println("done")
	fmt.Println("======= Helper Tests Succeeded =======\n")

	// ==================== //
	// Start Command Tests  //
	// ==================== //
	fmt.Println("======= Starting Command Tests =======")

	/** === Start Command === **/
	fmt.Print(timestamp()+"Running /start tests.... ")
	upd := GenerateUpdate("/start")
	_,err = bot.ParseAndExecuteUpdate(upd)

	if err != nil {
		log.Fatal(err)
		t.FailNow();
	}

	fmt.Println("done")

	/** === ID Command === **/
	fmt.Print(timestamp()+"Running /id tests.... ")
	upd = GenerateUpdate("/id")
	msg,err := bot.ParseAndExecuteUpdate(upd)

	if err != nil {
		log.Fatal(err)
		t.FailNow();
	}

	strText := msg.(tgbotapi.Message).Text
	if strings.Replace(strText, "Your chat ID is: ","",-1) != strconv.Itoa(ChatID) {
		log.Fatal("id mismatch")
		t.FailNow();
	}

	fmt.Println("done")

	/** === SMAWK Command === **/
	fmt.Print(timestamp()+"Running /smawk tests.... ")
	// Check to see if message is returned with proper user
	upd = GenerateUpdate("/smawk wheeee")
	msg,err = bot.ParseAndExecuteUpdate(upd)

	strText = msg.(tgbotapi.Message).Text
	if strings.Replace(strText, upd.Message.From.UserName+" ","",-1) != "wheeee" {
		log.Fatalf("/smawk string mismatch. Expected wheeee - got %s",strText)
		t.FailNow();
	}

	// Check to see if response is empty if not a user
	upd = GenerateUpdate("/smawk fail_plox")
	upd.Message.From.UserName = "not_bnmtthews"
	msg,err = bot.ParseAndExecuteUpdate(upd)

	if msg != (tgbotapi.Message{}) {
		log.Fatalf("/smawk not failing on bad user")
		t.FailNow();
	}

	// Make sure phrases work
	upd = GenerateUpdate("/smawk abc def ghi")
	msg,err = bot.ParseAndExecuteUpdate(upd)

	strText = msg.(tgbotapi.Message).Text
	if strings.Replace(strText, upd.Message.From.UserName+" ","",-1) != "abc def ghi" {
		log.Fatalf("/smawk string mismatch. Expected wheeee - got %s",strText)
		t.FailNow();
	}

	fmt.Println("done")

	/** === Me Command === **/
	fmt.Print(timestamp()+"Running /me tests.... ")
	// Check to see if message is returned with proper user
	upd = GenerateUpdate("/me wheeee")
	msg,err = bot.ParseAndExecuteUpdate(upd)

	strText = msg.(tgbotapi.Message).Text
	if strings.Replace(strText, upd.Message.From.UserName+" ","",-1) != "wheeee" {
		log.Fatalf("/me string mismatch. Expected wheeee - got %s",strText)
		t.FailNow();
	}

	// Check to see if response is empty if not a user
	upd = GenerateUpdate("/me fail_plox")
	upd.Message.From.UserName = "not_bnmtthews"
	msg,err = bot.ParseAndExecuteUpdate(upd)

	if msg != (tgbotapi.Message{}) {
		log.Fatalf("/me not failing on bad user")
		t.FailNow();
	}

	// Make sure phrases work
	upd = GenerateUpdate("/me abc def ghi")
	msg,err = bot.ParseAndExecuteUpdate(upd)

	strText = msg.(tgbotapi.Message).Text
	if strings.Replace(strText, upd.Message.From.UserName+" ","",-1) != "abc def ghi" {
		log.Fatalf("/me string mismatch. Expected wheeee - got %s",strText)
		t.FailNow();
	}

	fmt.Println("done")

	/** === All Command === **/
	fmt.Print(timestamp()+"Running /all tests.... ")

	upd = GenerateUpdate("/all")
	msg,err = bot.ParseAndExecuteUpdate(upd)

	allExpectedString := "@bnmtthews @cyberbatman @CMoneys @taborneighbor @wiglz @ReverendRecker @izaabsharp @PGB_Almighty @smawk_bot"
	if msg.(tgbotapi.Message).Text != allExpectedString {
		log.Fatalf("/all fail. Expected %s - got %s",allExpectedString,msg.(tgbotapi.Message).Text)
	}

	fmt.Println("done")

	/** === Hype Command === **/
	fmt.Print(timestamp()+"Running /hype tests.... ")
	upd = GenerateUpdate("/hype")
	msg,err = bot.ParseAndExecuteUpdate(upd)
	fmt.Println("done")

	/** === Label Command === **/
	fmt.Print(timestamp()+"Running /label tests.... ")
	fmt.Println("done")

	/** === Whois Command === **/
	fmt.Print(timestamp()+"Running /whois tests.... ")
	upd = GenerateUpdate("/whois")
	msg,err = bot.ParseAndExecuteUpdate(upd)

	strText = msg.(tgbotapi.Message).Text
	if strText != "Correct Usage: /whois @username" {
		log.Fatalf("/whois string mismatch. Expected Correct Usage: /whois @username - got %s",strText)
		t.FailNow();
	}

	fmt.Println("done")

	/** === Score Command === **/
	fmt.Print(timestamp()+"Running /score tests.... ")
	fmt.Println("done")

	/** === Upvote Command === **/
	fmt.Print(timestamp()+"Running /upvote tests.... ")
	fmt.Println("done")

	/** === Downvote Command === **/
	fmt.Print(timestamp()+"Running /downvote tests.... ")
	fmt.Println("done")

	/** === Bless Command === **/
	fmt.Print(timestamp()+"Running /bless tests.... ")
	fmt.Println("done")

	/** === Curse Command === **/
	fmt.Print(timestamp()+"Running /curse tests.... ")
	fmt.Println("done")

	fmt.Println("======= Command Tests Succeeded =======")
}
