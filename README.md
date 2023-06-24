# Quotebook of some Random AE Shitters

This website is going to be an automated quotebook of some guys from Aerospace Engineering at the TU Delft (and just one random computer science guy perhaps if he ever says something quote-worthy). Here you can find all the important quotes.

This website was written / is going to be writte in Golang - for the simple sake of spiting all the people using Javascript - also it performs better, and if these people ever decide that 10GB of quote data is to their pleasing well - well in that case that is honestly on the Database but yea...

## Functional Requirements List
- [X] Page-based website, that loads 15 quotes at the same time, chronological order
- [ ] Easy Search interface that allows searching for the following parts of a quote:
  - [ ] Author
  - [ ] Contained Text
  - [ ] Date (day/month/year combo)
    - [ ] Year only
    - [ ] Year Month Combination
  - [ ] Date-Range 
- [X] Add-quote Button, that allows a user to enter a new quote
  - [ ] With password authentication 
- [ ] Dark-mode detection, changing theme
- [ ] Next and Previous Page Buttons at the Bottom of the Page
- [ ] Easy input verification, quotes shouldn't be longer than 512 chars, and people no longer than 50.

## Feature List
- [X] A Small light-weight database, now SQLite3
  - [ ] Good Indexing for decent performance
- [X] The ability to POST to / to add quotes
- [ ] The ability to interop with a discord webhook (should be included in above?)
- [ ] Quick-page based loading out of the DB
- [ ] An easy JSON schema to allow easy interaction with the website

# Technical Information

## Running

You are going to have to go to the primary directory, and run the following command: `go mod init src/`. In the `src` directory you are also going to have to create a file called `DATABASE` (yes I know, great naming), in which you should save the schema specified in the below section. Then you can once again go to the primary directory, and run `go run src/main.go`. 
Some errors may be thrown, these will tell you to install certain git packages, using various `go mod` commands.

## Database Schema
```sql
TABLE quotes (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  quote VARCHAR(512) NOT NULL,
  date DATE NOT NULL,
  sayer VARCHAR(50) NOT NULL
);
```

