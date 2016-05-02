package main

import (
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha512"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
)

// type TokenManagerGet struct {
// 	User      string
// 	UserReply chan<- string
// }

// type TokenManagerPut struct {
// 	Token string
// 	User  string
// }

// func StartTokenManager() (chan<- TokenManagerGet, chan<- TokenManagerPut) {
// 	tokens := map[string]string{} // map[token]user
// 	getChan := make(chan TokenManagerGet)
// 	putChan := make(chan TokenManagerPut)
// 	go func() {
// 		for {
// 			select {
// 			case put := <-putChan:
// 				tokens[put.Token] = put.User
// 			case get := <-getChan:
// 				get.Reply <- tokens[get.Token]
// 			}
// 		}
// 	}()
// 	return getChan, putChan
// }

// func GetUserFromToken(get chan<- TokenManagerGet, token string) string {
// 	reply := make(chan string)
// 	get <- TokenManagerGet{Token: token, UserReply: reply}
// 	return <-reply
// }

// func PutToken(c chan<- TokenManagerPut, token, user string) {
// 	c <- TokenManagerPut{Token: token, User: user}
// }

func dosha1(pw string, salt string) string {
	h := sha1.New()
	io.WriteString(h, pw)
	return hex.EncodeToString((h.Sum(nil)))
}

func dosha512(pw, salt string) string {
	h := sha512.New()
	io.WriteString(h, pw)
	io.WriteString(h, salt)
	return hex.EncodeToString((h.Sum(nil)))
}

func dosha512_sha1(pw, salt string) string {
	return dosha512(dosha1(pw, salt), salt)
}

func GenerateToken() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func Authenticate(db *sql.DB, user, pass string) (bool, error) {
	encoding := ""
	hash := ""
	salt := ""
	err := db.QueryRow("select encoding, hash, salt from passwords where username = $1", user).Scan(&encoding, &hash, &salt)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil // if the user didn't exist, return false with a nil error.
		}
		return false, err
	}

	//	fmt.Printf("got user %s pass %s encoding %s hash %s salt %s\n", user, pass, encoding, hash, salt) // debug

	var hashFunc func(pw, salt string) string
	switch encoding {
	case "sha512":
		hashFunc = dosha512
	case "sha512_sha1":
		hashFunc = dosha512_sha1
	default:
		return false, fmt.Errorf("database returned unknown encoding: %s", encoding)
	}

	tryHash := hashFunc(pass, salt)
	if tryHash == hash {
		return true, nil
	}
	//	fmt.Printf("password fail: %s != %s\n", tryHash, hash) // debug
	return false, nil // password fail
}
