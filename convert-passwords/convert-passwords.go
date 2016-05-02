package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"io/ioutil"
	"os"
	"strings"
)

type UserPass struct {
	User     string
	Encoding string
	Hash     string
	Salt     string
}

func parsePasswords(body string) []UserPass {
	startI := strings.IndexRune(body, '(') + 1
	body = body[startI:]

	var passes []UserPass
	for i := strings.IndexRune(body, '('); i >= 0; i = strings.IndexRune(body, '(') {
		body = body[i+2:] // skip to the first letter of the username
		endUsernameQuote := strings.IndexRune(body, '"')

		username := body[:endUsernameQuote]

		body = body[endUsernameQuote+3:] // skip to first letter of the encoding
		spaceAfterEncoding := strings.IndexRune(body, ' ')

		encoding := body[:spaceAfterEncoding]

		body = body[spaceAfterEncoding+2:] // skip to first character of hash
		endHashQuote := strings.IndexRune(body, '"')

		hashText := body[:endHashQuote]

		body = body[endHashQuote+3:] // skip to beginning of salt
		endSaltQuote := strings.IndexRune(body, '"')

		saltText := body[:endSaltQuote]

		body = body[endSaltQuote+3:] // skip to right before open paren of next user

		passes = append(passes, UserPass{User: username, Encoding: encoding, Hash: hashText, Salt: saltText})
	}

	return passes
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("To convert passwords, pass the PostgreSQL connection URI as the first argument, and the password file as the second argument")
		return
	}
	connUri := os.Args[1]
	passwordFile := os.Args[2]

	passwordFileText, err := ioutil.ReadFile(passwordFile)
	if err != nil {
		fmt.Println("Error reading password file: " + err.Error())
		return
	}

	db, err := sql.Open("postgres", connUri)
	if err != nil {
		fmt.Println("Error opening database: " + err.Error())
		return
	}

	tx, err := db.Begin()
	if err != nil {
		fmt.Println("Error creating transaction: " + err.Error())
		return
	}

	stmt, err := tx.Prepare("insert into passwords (username, encoding, hash, salt) values ($1, $2, $3, $4);")
	if err != nil {
		fmt.Println("Error preparing statement: " + err.Error())
		return
	}

	fmt.Printf("parsing password file...\n")
	passes := parsePasswords(string(passwordFileText))

	fmt.Printf("inserting %d passwords...\n", len(passes))
	for _, p := range passes {
		_, err := stmt.Exec(p.User, p.Encoding, p.Hash, p.Salt)
		if err != nil {
			fmt.Println("Error inserting " + p.User + ": " + err.Error())
			return
		}
	}

	fmt.Printf("committing transaction...")
	err = tx.Commit()
	if err != nil {
		fmt.Println("Error committing transaction: " + err.Error())
		return
	}
}
