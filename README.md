# Quotebook of some Random AE Shitters

This website is going to be an automated quotebook of some guys from Aerospace Engineering at the TU Delft (and just one random computer science guy perhaps if he ever says something quote-worthy). Here you can find all the important quotes.

This website was written / is going to be writte in Golang - for the simple sake of spiting all the people using Javascript - also it performs better, and if these people ever decide that 10GB of quote data is to their pleasing well - well in that case that is honestly on the Database but yea...

## Functional Requirements List
- [X] Page-based website, that loads 15 quotes at the same time, chronological order
- [X] Easy Search interface that allows searching for the following parts of a quote:
  - [ ] Author
  - [X] Contained Text
  - [ ] Date (day/month/year combo)
    - [ ] Year only
    - [ ] Year Month Combination
  - [ ] Date-Range 
- [X] Add-quote Button, that allows a user to enter a new quote
  - [ ] With password authentication 
- [ ] Dark-mode detection, changing theme
- [X] ~~Next and Previous Page Buttons at the Bottom of the Page~~ (Changed to Pagination)
- [X] Easy input verification, quotes shouldn't be longer than 512 chars, and people no longer than 50.

## Feature List
- [X] A Small light-weight database, now SQLite3
  - [ ] Good Indexing for decent performance
- [X] The ability to POST to / to add quotes
- [ ] The ability to interop with a discord webhook (should be included in above?)
- [ ] Quick-page based loading out of the DB
- [ ] An easy JSON schema to allow easy interaction with the website

# Technical Information

## Running

### Dependencies

- Go
- SQLite3

> Note that making this cooperate on windows maybe quite difficult / impossible. I did not attempt this yet. If you can make it work, feel free to submit a pull request with the required instructions.

### Installation & Running

Download the Git repo in whichever way you please.
Enter the first folder and run
```sh 
go mod init src
```
This is done to setup the Go environment. Then enter the main `src` folder, and create a file named `DATABASE`. This file should contain the database schema specified below (created using SQLite3).

#### Creating the SQL Schema

For this you require an installation of SQLite3 (`sqlite3`). Make sure that the `DATABASE` file is empty.

First enter sqlite3 with the command
```sh
sqlite3
```

at this point you should be in the sqlite3 environment.
Run the following commands to create the schema below:
```
.open DATABASE

CREATE TABLE quotes (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  quote VARCHAR(512) NOT NULL,
  date DATE NOT NULL,
  sayer VARCHAR(50) NOT NULL
);
```

Then quit out of sqlite3 by pressing `<C-c>` (Control-C) twice.

### Continue Installation

Return to the root folder, and run:
```sh 
go get github.com/mattn/go-sqlite3
```

Now you can run the project from the root folder using
```sh 
go run src/main.go
```
Note the first startup may take a second.

## Database Schema
```sql
CREATE TABLE quotes (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  quote VARCHAR(512) NOT NULL,
  date DATE NOT NULL,
  sayer VARCHAR(50) NOT NULL
);
```

