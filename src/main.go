package main

import (
    "fmt"
    "net/http"
    "html/template"
    "strconv"
    "github.com/mattn/go-sqlite3"
    "database/sql"
)



// Creating the Main Database for Global Access
// TODO @robert - this is perhaps not the best way to do this -> but it'll be fine
var db sql.DB;
var sqlite3Conn sqlite3.SQLiteConn;

func rootPage(w http.ResponseWriter, req *http.Request) {
    
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


    // Get index.html and render to client (reponseWriter)
    indexPage, _ := template.ParseFiles("static/templates/index.html")
    err = indexPage.Execute(w, "")

    if err != nil {
        fmt.Fprintf(w, "Something went wrong...")
    }
}

type QuoteQuery struct {
    ID int64
    Quote string 
    Date string
    Sayer string
}

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

func addQuoteToDb(quote QuoteQuery) {

}


func main() {


    /**
     * DATABASE SETUP
     */
    
    // file:test.db?cache=shared&mode=memory <- an in-memory database DSN you can use (but not on the Raspi lol)
    db, connectionError := sql.Open("sqlite3", "file:src/DATABASE?cache=shared")

    if e := db.Ping(); e != nil || connectionError != nil {
        fmt.Println("Failed to Start the DB")
    } else {
        fmt.Println("Connected to SQLite3 Database (DATABASE file)")
    }

    // Close the database once Main exits
    defer db.Close()
    


    /**
     * HTTP HANDLING
     */
    fs := http.FileServer(http.Dir("static/"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

    http.HandleFunc("/", rootPage)


    http.ListenAndServe(":8000", nil)
}
