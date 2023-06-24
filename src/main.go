package main

import (
    "fmt"
    "net/http"
    "html/template"
    "strconv"
    "github.com/mattn/go-sqlite3"
    "database/sql"
    "regexp"
)



// Creating the Main Database for Global Access
// TODO: @robert - this is perhaps not the best way to do this -> but it'll be fine
var db *sql.DB;
var sqlite3Conn sqlite3.SQLiteConn;


type QuoteQuery struct {
    ID int64
    Quote string 
    Date string
    Sayer string
}

type Data struct {
    Quotes []QuoteQuery
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

func getQuotePageFromDB(pageNumber int) ([]QuoteQuery) {
    rows, err := db.Query(getQueryForPage(pageNumber))

    var returnQuotes []QuoteQuery;

    if err != nil { fmt.Println("Query failed...", err) }
    defer rows.Close()

    for rows.Next() {

        var quote QuoteQuery;

        if err := rows.Scan(&quote.ID, &quote.Quote, &quote.Date, &quote.Sayer); err != nil {
            fmt.Println("Failed to format a result...", err)
        }
        returnQuotes = append(returnQuotes, quote)
    }

    return returnQuotes
}

func updateHandling(w http.ResponseWriter, req *http.Request) {


    //r.FormValue
}

func routeHandler(w http.ResponseWriter, req *http.Request) {

    pathRoute := req.URL.Path 

    if match, _ := regexp.MatchString(`^/(\?page=[0-9]+)?$`, pathRoute); match {
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
    
    if req.Method == http.MethodPost {
        fmt.Println("Received update...")

        fmt.Println(req.FormValue("Quote"))
        quote := req.FormValue("Quote")
        fmt.Println(req.FormValue("Date"))
        date := req.FormValue("Date")
        fmt.Println(req.FormValue("Sayer"))
        sayer := req.FormValue("Sayer")

        db.Exec("INSERT INTO quotes (quote, date, sayer) VALUES ( ?, ?, ?)", quote, date, sayer)

    }

    // Getting the query out of the Request
    queryParams := req.URL.Query()
    pageList := queryParams["page"]

    var pageQuery string;

    if len(pageList) > 0 {
        pageQuery = pageList[0]
    }
    
    // Getting the Page number that was queried for
    pageNumber, err := strconv.Atoi(pageQuery) // TODO: @robert: actually handle this error please...
    fmt.Println(pageNumber)

    // Getting relevant rows out of the DB
    quotes := getQuotePageFromDB(pageNumber)

    data := Data {
        Quotes: quotes,
    }

    // Get index.html and render to client (reponseWriter)
    indexPage, _ := template.ParseFiles("src/static/templates/index.html")

    err = indexPage.Execute(w, data)

    if err != nil {
        fmt.Fprintf(w, "Something went wrong...", err)
    }
}

func main() {

    var connectionError error;
    db, connectionError = sql.Open("sqlite3", "file:src/DATABASE?cache=shared")

    if e := db.Ping(); e != nil || connectionError != nil {
        fmt.Println("Failed to Start the DB")
    } else {
        fmt.Println("Connected to SQLite3 Database (DATABASE file)")
    }
    
    // HACK: we are actually just doing the jank shit now
    http.HandleFunc("/", routeHandler)

    http.ListenAndServe(":8000", nil)
}
