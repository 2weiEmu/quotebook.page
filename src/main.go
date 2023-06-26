package main

import (
    // Go Library Packages
    "database/sql"
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
var sqlite3Conn sqlite3.SQLiteConn;
var pageSearchStatement *sql.Stmt

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
}

type Data struct {
    Quotes []QuoteQuery
    Pagination []Pages
}

// TODO: @robert one day prepared statements pls
func getQueryForPage(pageNumber int) string {
    sqlStmt := `SELECT id, quote, date, sayer FROM quotes as q
    ORDER BY date DESC
    LIMIT `

    startOffset := pageNumber * 15
    endOffset := startOffset + 15

    finalStatement := fmt.Sprintf("%s%d,%d", sqlStmt, startOffset, endOffset)

    return finalStatement
}


// TODO: add returning the error
func getQuotesPrepared(searchString string, pageNumber int) ([]QuoteQuery, error) {

    // pageSearchStatement.Query takes "searchString" "startNumber" "endNumber"
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

        err := rows.Scan(&quote.ID, &quote.Quote, &quote.Date, &quote.Sayer)

        if err != nil {
            fmt.Println("Failed to retrieve row:", rows, "With the error:", err)
        }

        quote.Date = strings.TrimSuffix(quote.Date, "T00:00:00Z")
        returnQuotes = append(returnQuotes, quote)

    }

    fmt.Println("Returned Quotes:", returnQuotes)

    return returnQuotes, nil

}

func updateHandling(w http.ResponseWriter, req *http.Request) {


    if req.Method == http.MethodPost {
        fmt.Println("Received update...")

        quote := req.FormValue("Quote")
        date := req.FormValue("Date")
        sayer := req.FormValue("Sayer")

        db.Exec("INSERT INTO quotes (quote, date, sayer) VALUES ( ?, ?, ?)", quote, date, sayer)

        // Redirect to prevent form resubmission

    } else if req.Method == http.MethodDelete {

    }


    http.Redirect(w, req, "/", http.StatusSeeOther)

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

    }
}

func indexPage(w http.ResponseWriter, req *http.Request) {
    
    //searchAuthor := queryParams["author"] // TODO: add to regex
    //searchDate := queryParams["date"] // TODO: add to regex

    // Getting the query out of the Request
    queryParams := req.URL.Query()

    searchText := ""
    pageNumber := 0
    
    if queryParams.Has("page") {
        pageNumber, _ = strconv.Atoi(queryParams["page"][0])
    }

    if queryParams.Has("search") {
        searchText = queryParams["search"][0]
    }

    quotes, _ := getQuotesPrepared(searchText, pageNumber)

    var pagesAround []Pages;

    for i := pageNumber - 2; i < pageNumber + 3; i++ {
        if i >= 0 {
            pagesAround = append(pagesAround, Pages{i});
        }
    }

    data := Data {
        Quotes: quotes,
        Pagination: pagesAround,
    }

    indexPage, _ := template.ParseFiles("src/static/templates/index.html")

    err := indexPage.Execute(w, data)

    if err != nil {
        fmt.Fprintf(w, "Something went wrong: %s", err)
    }
}


/*
 * Main
 */
func main() {

    var connectionError error;
    db, connectionError = sql.Open("sqlite3", "file:src/DATABASE?cache=shared")

    defer db.Close()

    if connectionError != nil {
        fmt.Println("Failed to connect to the database (src/DATABASE) with error:", connectionError)
    }

    if e := db.Ping(); e != nil {
        fmt.Println("DB Not Connected. Ping failed with error:", e)
    } else {
        fmt.Println("Connected to SQLite3 Database (src/DATABASE file)")
    }

    // Preparing Database Statement
    var err error
    pageSearchStatement, err = db.Prepare(`SELECT * FROM quotes as q WHERE q.quote LIKE ? ORDER BY date DESC LIMIT ?, ?`)
    defer pageSearchStatement.Close()

    if err != nil {
        fmt.Println("Failed to prepare 'pageSearchStatement' with error:", err)
    }
    
    // HACK: we are actually just doing the jank shit now
    http.HandleFunc("/", routeHandler)

    http.ListenAndServe(":8000", nil)
}

