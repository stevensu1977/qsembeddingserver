package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/url"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"stevensu1977/quicksight/utils"
)

//Init Session Store
var store = sessions.NewCookieStore([]byte("0123456789ABCDEFGHIJKLabcdefghijkl"))

//user data was loaded from yaml file , use RDBMS/NoSQL for production environment
var users = []utils.User{}

var region = ""

var roleName = ""

//sessionCheckMiddleware is mux middleware , check login session
func sessionCheckMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "QuichSightPortal")
		if session.Values["username"] != nil {
			next.ServeHTTP(w, r)
			log.Printf("Already login, %s", session.Values["username"])
		} else {
			if r.URL.Path != "/static/login.html" && r.URL.Path != "/login" && r.URL.Path != "/logout" {
				NeedAuthHandler(w, r)
				return
			}
			next.ServeHTTP(w, r)
		}
	})
}

//loadDashboardIDs , load quicksight dashboard ids by user email
func loadDashboardIDs(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, "QuichSightPortal")
	jsonResponse, jsonError := json.Marshal(session.Values["dashboards"])
	if jsonError != nil {
		http.Error(w, jsonError.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

}

//login check useremail/password
func login(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fmt.Println(params)
	r.ParseForm()

	if r.FormValue("username") != "" && r.FormValue("password") != "" {
		fmt.Println(r.FormValue("username"), r.FormValue("password"), utils.Sha1(r.FormValue("password")))
		for _, v := range users {

			if r.FormValue("username") == v.Email && utils.Sha1(r.FormValue("password")) == v.Password {
				fmt.Println("login successful")
				session, _ := store.Get(r, "QuichSightPortal")
				session.Values["username"] = v.Email
				session.Values["dashboards"] = v.DashBoards
				session.Save(r, w)

				//redirect index.html
				http.Redirect(w, r, "/static/index.html", http.StatusTemporaryRedirect)
				return
			}
		}

	}
	//if not login redirect login.html
	http.Redirect(w, r, "/static/login.html", http.StatusTemporaryRedirect)
}

//logout destory login session
func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "QuichSightPortal")

	session.Values["username"] = nil
	session.Values["dashboards"] = nil

	session.Save(r, w)

	http.Redirect(w, r, "/static/login.html", http.StatusTemporaryRedirect)

}

//loadUsers load user data from yaml file
func loadUsers(dataFile string) {
	var err error
	users, err = utils.LoadUsersFromFile(dataFile)
	if err != nil {
		panic(err)
	}
	for _, v := range users {
		log.Println(fmt.Sprintf("User: %s, %s", v.Email, v.DashBoards))
	}
}

//buildEmbedURL invoke AWS Golang SDK get quicksight embedding URL
func buildEmbedURL(w http.ResponseWriter, r *http.Request) {

	awsAccountId, err := utils.GetAccountId()

	if err != nil {
		panic(err)
	}

	session, _ := store.Get(r, "QuichSightPortal")

	userEmail := session.Values["username"].(string)

	urlParams := r.URL.Query()
	dashboardID := urlParams.Get("id")
	dashboardURL := utils.GetEmbedUrl(*awsAccountId, region, roleName, dashboardID, userEmail)

	log.Println(*dashboardURL)
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.Encode(map[string]string{"url": *dashboardURL})
}

func NeedAuthHandler(w http.ResponseWriter, r *http.Request) {
	// Get a session. We're ignoring the error resulted from decoding an
	// existing session: Get() always returns a session, even if empty.
	http.Redirect(w, r, "/static/login.html", 401)

}

func main() {

	addr := flag.String("addr", ":8000", "listen address")
	dir := flag.String("dir", "./static", "the directory to serve files from. Defaults to the current dir")
	dataFile := flag.String("data", "users.yaml", "data file")
	flag.StringVar(&region, "region", "us-west-2", "quicksight region")
	flag.StringVar(&roleName, "role", "QSDER", "quicksight role ")

	flag.Parse()

	fmt.Println(*addr, *dir, region)
	log.Println("================QuickSight Embedding Server===============")
	log.Printf("Region: %s, Role: %s", region, roleName)

	//load users from data file
	loadUsers(*dataFile)

	//init mux Roouter
	r := mux.NewRouter()

	//static content
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(*dir))))

	//dynamic content
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/static/login.html", 302)
	})

	r.HandleFunc("/login", login).Methods("POST")
	r.HandleFunc("/logout", logout).Methods("GET")

	r.HandleFunc("/api/dashboard/ids", loadDashboardIDs).Methods("GET")
	r.HandleFunc("/api/dashboard/embedURL", buildEmbedURL).Methods("GET")

	//enable loginMiddleware
	r.Use(sessionCheckMiddleware)

	srv := &http.Server{
		Handler: r,
		Addr:    *addr,

		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Printf("server start successfully, listen on %s ", *addr)
	log.Fatal(srv.ListenAndServe())

}
