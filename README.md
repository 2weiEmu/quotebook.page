# Quotebook of some Random AE Shitters

This website is going to be an automated quotebook of some guys from Aerospace Engineering at the TU Delft (and just one random computer science guy perhaps if he ever says something quote-worthy). Here you can find all the important quotes.

This website was written / is going to be writte in Golang - for the simple sake of spiting all the people using Javascript - also it performs better, and if these people ever decide that 10GB of quote data is to their pleasing well - well in that case that is honestly on the Database but yea...

## Functional Requirements List
- [ ] Page-based website, that loads 15 quotes at the same time, chronological order
- [ ] Easy Search interface that allows searching for the following parts of a quote:
  - [ ] Author
  - [ ] Contained Text
  - [ ] Date (day/month/year combo)
    - [ ] Year only
    - [ ] Year Month Combination
  - [ ] Date-Range 
- [ ] Add-quote Button, that allows a user to enter a new quote
  - [ ] With password authentication 
- [ ] Dark-mode detection, changing theme
- [ ] Next and Previous Page Buttons at the Bottom of the Page

## Feature List
- [ ] A Small light-weight database, perhaps H2 (but I don't know - maybe go works easiest with MySQL).
  - [ ] Good Indexing for decent performance
- [ ] The ability to POST to /quoteadd/ to add quotes
- [ ] The ability to interop with a discord webhook (should be included in above?)
- [ ] Quick-page based loading out of the DB
- [ ] An easy JSON schema to allow easy interaction with the website

# Technical Information

## Database Schema
```sql
TODO
```

