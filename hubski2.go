package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

// TODO(check and log errors from io.Write/fmt.Fprintf(w )

// HandleAuthErr writes the appropriate error response, on an auth error from dbGetTokenUser
func HandleAuthErr(userFound bool, err error, w http.ResponseWriter) {
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else if !userFound {
		log.Println("User not found: returning Unauthorized")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
}

// \todo verify params exist in map; log; make tokens expire
func Login(db *sql.DB, params map[string]string, w http.ResponseWriter, r *http.Request) {
	user := params["user"]
	pass := params["pass"]
	authenticated, err := Authenticate(db, user, pass)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if !authenticated {
		log.Println("Authentication failed: returning Unauthorized")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	token, err := GenerateToken()
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = dbInsertToken(token, user, db)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s", token)
}

func UserPageData(db *sql.DB, params map[string]string, w http.ResponseWriter, r *http.Request) {
	user, userFound, err := dbGetTokenUser(params["token"], db)
	if !userFound || err != nil {
		HandleAuthErr(userFound, err, w)
		return
	}

	userPageData, userFound, err := dbGetUserPageUserData(user, db)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	if !userFound {
		// TODO(501? It shouldn't be possible to get a token for a user that doesn't exist, right?)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	jsonUserPageData, err := json.Marshal(userPageData)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s", jsonUserPageData)
}

func UserPostsAndShares(db *sql.DB, params map[string]string, w http.ResponseWriter, r *http.Request) {
	user, userFound, err := dbGetTokenUser(params["token"], db)
	if !userFound || err != nil {
		HandleAuthErr(userFound, err, w)
		return
	}

	numPostsAndShares := 100

	publications, err := dbGetUserPagePostsAndShares(user, numPostsAndShares, db)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	jsonPublications, err := json.Marshal(publications)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s", jsonPublications)
}

// TODO(log)
func Donations(db *sql.DB, params map[string]string, w http.ResponseWriter, r *http.Request) {
	user, userFound, err := dbGetTokenUser(params["token"], db)
	if !userFound || err != nil {
		HandleAuthErr(userFound, err, w)
		return
	}

	donations, err := dbGetDonations(user, db)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%d", donations)
}

func UserPageComments(db *sql.DB, params map[string]string, w http.ResponseWriter, r *http.Request) {
	numUserPageComments := 12

	user, userFound, err := dbGetTokenUser(params["token"], db)
	if !userFound || err != nil {
		HandleAuthErr(userFound, err, w)
		return
	}

	comments, err := dbGetUserPageComments(user, numUserPageComments, db)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	jsonComments, err := json.Marshal(comments)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s", jsonComments)
}

func UserPage(db *sql.DB, params map[string]string, w http.ResponseWriter, r *http.Request) {
	// TODO load on startup and pass as parameter
	htmlBytes, err := ioutil.ReadFile("user.html")
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	fmt.Fprintf(w, "%s", htmlBytes)
}

func GetRoutes() []Route {
	return []Route{
		Route{"donations/{token}", Donations},
		Route{"login/{user}/{pass}", Login}, // TODO(put pass in body?)
		Route{"user-page-data/{token}", UserPageData},
		Route{"user-page-comments/{token}", UserPageComments},
		Route{"user-page-posts-and-shares/{token}", UserPostsAndShares},
		Route{"user-page", UserPage},
	}
}

func RegisterRouter(routes []Route, db *sql.DB) {
	compiledRoutes := CompileRoutes(routes)

	wrapHandler := func(w http.ResponseWriter, r *http.Request) {
		Handler(db, compiledRoutes, w, r)
	}

	http.HandleFunc("/", WrapLogRequests(WrapCORS(wrapHandler)))
}

func WrapCORS(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		f(w, r)
	}
}

func WrapLogRequests(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s requested %s\n", r.RemoteAddr, r.RequestURI)
		f(w, r)
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("To run, pass PostgreSQL connection URI as the first argument, and port to serve on as the second argument")
		return
	}
	connUri := os.Args[1]
	port, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Println("port (second argument) is not an integer")
		return
	}

	db, err := dbGet(connUri)
	if err != nil {
		log.Println(err)
		return
	}

	RegisterRouter(GetRoutes(), db)
	log.Printf("Listening on %d\n", port)
	err = http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		log.Println(err)
	}
}
