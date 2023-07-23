package main

import (
	// Go Library Packages
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"strconv"

	// External Packages
	"github.com/mattn/go-sqlite3"
)

// Creating the Main Database for Global Access
// TODO: @robert - this is perhaps not the best way to do this -> but it'll be fine
var db *sql.DB
var sqlite3Conn sqlite3.SQLiteConn
var pageSearchStatement *sql.Stmt
var err error

/*
 * Useful Structs
 */
type QuoteQuery struct {
	ID    int64
	Quote string
	Date  string
	Sayer string
}

type Pages struct {
    Page int
    Search string
    Author string
}

type Data struct {
    TotalSearch string
    AuthorSearch string
    Quotes []QuoteQuery
    Pagination []Pages
}

type APIPost struct {
	Quote string `json:"Quote"`
	Date  string `json:"Date"`
	Sayer string `json:"Sayer"`
}

type APIDelete struct {
	ID    int    `json:"ID"`
	Quote string `json:"Quote"`
	Date  string `json:"Date"`
	Sayer string `json:"Sayer"`
}

type APIPut struct {
	Quote       string `json:"Quote"`
	Date        string `json:"Date"`
	Sayer       string `json:"Sayer"`
	ChangedAttr string `json:"ChangedAttr"`
	NewValue    string `json:"NewValue"`
}

// Endpoint for sending Form values
func UpdateHandling(w http.ResponseWriter, req *http.Request) {

	if req.Method == http.MethodPost {
		fmt.Println("Received update...")

		quote := req.FormValue("Quote")
		date := req.FormValue("Date")
		sayer := req.FormValue("Sayer")

		db.Exec("INSERT INTO quotes (quote, date, sayer) VALUES ( ?, ?, ?)", quote, date, sayer)

		// Redirect to prevent form resubmission

	}

	http.Redirect(w, req, "/", http.StatusSeeOther)
}

// Endpoint for sending raw JSON
func ApiHandling(w http.ResponseWriter, req *http.Request) {

	var decoder *json.Decoder

	if req.Body != nil {
		decoder = json.NewDecoder(req.Body)
	}

	if req.Method == http.MethodPost {

		var post APIPost
		err := decoder.Decode(&post)

		if err != nil {
			fmt.Println("Post request failed. Error:", err)
			fmt.Fprintf(w, "%d", http.StatusInternalServerError)
		}

		fmt.Println(post)

		// WARNING: User input is not validated yet!

		// TODO: jaja prepared statements and all that
		db.Exec(`INSERT INTO quotes (quote, date, sayer) VALUES ( ?, ?, ?)`, post.Quote, post.Date, post.Sayer)

		// TODO: return the id to the post making the request -> they might need it
	} else if req.Method == http.MethodDelete {

		// TODO: delete does not really work - like - we have to figure out a bit how to do this... maybe we just don't...

		var apiDelete APIDelete

		err := decoder.Decode(&apiDelete)

		if err != nil {
			fmt.Println("Delete request failed. Error:", err)
			fmt.Fprintf(w, "%d", http.StatusInternalServerError)
		}

		fmt.Println(apiDelete)

	} else if req.Method == http.MethodPut {

		var apiPut APIPut

		err := decoder.Decode(&apiPut)

		if err != nil {
			fmt.Println("Put request failed. Error:", err)
			fmt.Fprintf(w, "%d", http.StatusInternalServerError)
		}

		// TODO: changing attribute listed

	}
}


func RouteHandler(w http.ResponseWriter, req *http.Request) {

	pathRoute := req.URL.Path

    if match, _ := regexp.MatchString(`^\/(\?((page=[0-9]+)|(search=[\w]+)|(author=[\w]+)))?$`, pathRoute); match {
        IndexPage(w, req)

	} else if match, _ := regexp.MatchString(`^/css/`, pathRoute); match {
		fmt.Println("Serving Static file...")
		fs := http.FileServer(http.Dir("src/static/"))
		http.StripPrefix("static/", fs)
		fs.ServeHTTP(w, req)

    } else if match, _ := regexp.MatchString(`^/updates/`, pathRoute); match {
        UpdateHandling(w, req)

    } else if match, _ := regexp.MatchString(`^/api/`, pathRoute); match {
        ApiHandling(w, req)
    }
}

func IndexPage(w http.ResponseWriter, req *http.Request) {
    
    //searchDate := queryParams["date"] // TODO: add to regex

	//searchAuthor := queryParams["author"] // TODO: add to regex
	//searchDate := queryParams["date"] // TODO: add to regex

	// Getting the query out of the Request
	queryParams := req.URL.Query()

    searchTotal := "Nothing"
    searchText := ""
    searchAuthor := ""
    pageNumber := 0
    
    if queryParams.Has("page") {
        pageNumber, _ = strconv.Atoi(queryParams["page"][0])
    }

    if queryParams.Has("author") {
        searchAuthor = queryParams["author"][0]
    }

    quotes, nextPageMarker, err := GetQuotesPrepared(pageSearchStatement, searchText, pageNumber, searchAuthor)

	if queryParams.Has("search") {
		searchText = queryParams["search"][0]
	}


	var pagesAround []Pages

    if pageNumber > 0 {
        pagesAround = append(pagesAround, Pages{pageNumber - 1, searchText, searchAuthor})
    }

    pagesAround = append(pagesAround, Pages{pageNumber, searchText, searchAuthor})

    if nextPageMarker {
        pagesAround = append(pagesAround, Pages{pageNumber + 1, searchText, searchAuthor})
    }
	if err != nil {
		fmt.Println("Failed to get Quotes using prepared statement with error:", err)
		return
	}

    AuthorSearch := "Nobody"
    if searchAuthor != "" {
        AuthorSearch = "\"" + searchAuthor + "\""
    }

    data := Data {
        TotalSearch: searchTotal,
        AuthorSearch: AuthorSearch,
        Quotes: quotes,
        Pagination: pagesAround,
    }

	if searchText != "" {
		searchTotal = "\"" + searchText + "\""
	}

	indexPage, _ := template.ParseFiles("src/static/templates/index.html")

	err = indexPage.Execute(w, data)

	if err != nil {
		fmt.Fprintf(w, "Something went wrong: %s", err)
	}
}

/*
 * Main
 */
func main() {

    db, err = sql.Open("sqlite3", "file:src/DATABASE?cache=shared")

	defer db.Close()

    if err != nil {
        fmt.Println("Failed to connect to the database (src/DATABASE) with error:", err)
    }

    // Preparing Database Statement
    pageSearchStatement, err = db.Prepare(`SELECT * FROM quotes as q WHERE q.quote LIKE ? AND sayer LIKE ? ORDER BY date DESC LIMIT ?, ?`)

    defer pageSearchStatement.Close()
	if e := db.Ping(); e != nil {
		fmt.Println("DB Not Connected. Ping failed with error:", e)
	} else {
		fmt.Println("Connected to SQLite3 Database (src/DATABASE file)")
	}
    if err != nil {
        fmt.Println("Failed to prepare 'pageSearchStatement' with error:", err)
    }
    
    // HACK: we are actually just doing the jank shit now
    // NOTE: apparently this is the accepted way of doing this?
    http.HandleFunc("/", RouteHandler)

	// Preparing Database Statement


    fmt.Println("Hosted server on localhost port 8000.")

    http.ListenAndServe(":8000", nil)

}
