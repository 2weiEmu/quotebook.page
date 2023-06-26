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
	"strings"

	// External Packages
	"github.com/mattn/go-sqlite3"
)

// Creating the Main Database for Global Access
// TODO: @robert - this is perhaps not the best way to do this -> but it'll be fine
var db *sql.DB;
var sqlite3Conn sqlite3.SQLiteConn
var pageSearchStatement *sql.Stmt
var err error

/*
 * Useful Structs
 */
type QuoteQuery struct {
    ID int64
    Quote string 
    Date string
    Sayer string
}

type Pages struct {
    Page int
    Search string
}

type Data struct {
    TotalSearch string
    Quotes []QuoteQuery
    Pagination []Pages
}

type APIPost struct {
    Quote string `json:"Quote"`
    Date string `json:"Date"`
    Sayer string `json:"Sayer"`
}

type APIDelete struct {
    ID int `json:"ID"`
    Quote string `json:"Quote"`
    Date string `json:"Date"`
    Sayer string `json:"Sayer"`

}

type APIPut struct {
    Quote string `json:"Quote"`
    Date string `json:"Date"`
    Sayer string `json:"Sayer"`
    ChangedAttr string `json:"ChangedAttr"`
    NewValue string `json:"NewValue"`
}

func getQuotesPrepared(searchString string, pageNumber int) ([]QuoteQuery, error) {

    startNumber := pageNumber * 15;
    endNumber := startNumber + 15;
    searchString = "%" + searchString + "%"

    rows, err := pageSearchStatement.Query(searchString, startNumber, endNumber)
    defer rows.Close()

    if err != nil {
        fmt.Println("Prepared statement failed to execute with error:", err)
        return nil, err;
    }

    var returnQuotes []QuoteQuery;

    for rows.Next() {

        var quote QuoteQuery;
        err = rows.Scan(&quote.ID, &quote.Quote, &quote.Date, &quote.Sayer)

        if err != nil {
            fmt.Println("Failed to retrieve row:", rows, "With the error:", err)
            return nil, err
        }

        quote.Date = strings.TrimSuffix(quote.Date, "T00:00:00Z")
        returnQuotes = append(returnQuotes, quote)

    }
    return returnQuotes, nil
}


// Endpoint for sneding Form values
func updateHandling(w http.ResponseWriter, req *http.Request) {

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
func apiHandling(w http.ResponseWriter, req *http.Request) {

    var decoder *json.Decoder;

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

        fmt.Println(post);

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


func routeHandler(w http.ResponseWriter, req *http.Request) {

    pathRoute := req.URL.Path 

    if match, _ := regexp.MatchString(`^\/(\?((page=[0-9]+)|(search=[\w]+)))?$`, pathRoute); match {
        indexPage(w, req)

    } else if match, _ := regexp.MatchString(`^/css/`, pathRoute); match {
        fmt.Println("Serving Static file...")
        fs := http.FileServer(http.Dir("src/static/"))
        http.StripPrefix("static/", fs)
        fs.ServeHTTP(w, req)

    } else if match, _ := regexp.MatchString(`^/updates/`, pathRoute); match {
        updateHandling(w, req)

    } else if match, _ := regexp.MatchString(`^/api/`, pathRoute); match {
        apiHandling(w, req)
    }
}

func indexPage(w http.ResponseWriter, req *http.Request) {
    
    //searchAuthor := queryParams["author"] // TODO: add to regex
    //searchDate := queryParams["date"] // TODO: add to regex

    // Getting the query out of the Request
    queryParams := req.URL.Query()

    searchTotal := "Nothing"
    searchText := ""
    pageNumber := 0
    
    if queryParams.Has("page") {
        pageNumber, _ = strconv.Atoi(queryParams["page"][0])
    }

    if queryParams.Has("search") {
        searchText = queryParams["search"][0]
    }

    quotes, err := getQuotesPrepared(searchText, pageNumber)

    if err != nil {
        fmt.Println("Failed to get Quotes using prepared statement with error:", err)
        return
    }

    var pagesAround []Pages;

    for i := pageNumber - 2; i < pageNumber + 3; i++ {
        if i >= 0 {
            pagesAround = append(pagesAround, Pages{i, searchText});
        }
    }


    if searchText != "" {
        searchTotal = "\"" + searchText + "\""
    }

    data := Data {
        TotalSearch: searchTotal,
        Quotes: quotes,
        Pagination: pagesAround,
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

    if err = db.Ping(); err != nil {
        fmt.Println("DB Not Connected. Ping failed with error:", err)
    } else {
        fmt.Println("Connected to SQLite3 Database (src/DATABASE file)")
    }

    // Preparing Database Statement
    pageSearchStatement, err = db.Prepare(`SELECT * FROM quotes as q WHERE q.quote LIKE ? ORDER BY date DESC LIMIT ?, ?`)
    defer pageSearchStatement.Close()

    if err != nil {
        fmt.Println("Failed to prepare 'pageSearchStatement' with error:", err)
    }
    
    // HACK: we are actually just doing the jank shit now
    http.HandleFunc("/", routeHandler)

    fmt.Println("Hosted server on localhost port 8000.")

    http.ListenAndServe(":8000", nil)

}

