package main

import (
    "fmt"
    "database/sql"
    "strings"
)

func GetQuotesPrepared(preparedStatement *sql.Stmt, searchString string, pageNumber int, searchAuthor string) ([]QuoteQuery, bool, error) {

    startNumber := pageNumber * 15;
    endNumber := startNumber + 16;
    searchString = "%" + searchString + "%"
    searchAuthor = "%" + searchAuthor + "%"

    rows, err := preparedStatement.Query(searchString, searchAuthor, startNumber, endNumber)
    defer rows.Close()

    if err != nil {
        fmt.Println("Prepared statement failed to execute with error:", err)
        return nil, false, err;
    }

    var returnQuotes []QuoteQuery;

    for rows.Next() {

        var quote QuoteQuery;
        err = rows.Scan(&quote.ID, &quote.Quote, &quote.Date, &quote.Sayer)

        if err != nil {
            fmt.Println("Failed to retrieve row:", rows, "With the error:", err)
            return nil, false, err
        }

        quote.Date = strings.TrimSuffix(quote.Date, "T00:00:00Z")
        returnQuotes = append(returnQuotes, quote)

    }
    
    nextPage := len(returnQuotes) > 15

    return returnQuotes, nextPage, nil
}
