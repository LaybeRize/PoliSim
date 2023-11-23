# Roadmap

This is the offical Roadmap for the **PoliSim** Project.

## Goals for Beta

These are the goals which have to be reached to transfer the project from the Alpha-Phase to the Beta-Phase.

- tests for most logic (ca 60%)
- bug testing with users
- implementation of the Zwitscher-System
- basic general-purpose API for JSON
- updates discussion and vote on end
- advanced search for documents
- english localisation
- add css rule for buttons so their text isn't selectable anymore
- add information on new letters for user

### Tests

- [ ] finished

The 60% testcov are not project wide but inside the data folder for now. There the tests should cover most non-database
dependent logic. For that reason the project should decouple database-reliant and database-indendent validation/processing 
logic. If possible all database-independet logic should be covered.

### Bug testing

- [ ] finished

No specifics can be made currently to indicate when this phase is sufficently finished. But the general gist is, that 
if the product has been tested for a month, no obvious bugs should be floating around.

### Zwitscher System

- [ ] finished

The Zwitscher System should allow user to write short text messages on a public forum and let people comment under 
these short messages.   
These messages should be able to be viewed in general just in order of occurence and by user.  
In addition the user comments should be able to be seperated by general comments and answers to general comments.

### Basic JSON API

- [ ] finished

the website should have a rate limited JSON API for the most general things like list of organisations, titles
documents and so forth. Advanced queries for users with a token that must be sent with the request.

### Update Votes and Discussion when they end

- [x] finished

Votes and Discussion should automatically request an update when the run time expires.

### Advanced Search for documents

- [x] finished

Documents should not be only be available as sorted list by time, but also be able to search documents only 
from a specific organisation/person/type and also only documents before a specified date.  
Needs implementing being able to block documents too.

### Paging bug

- [x] finished

fix bug with paging occuring with three pages and more. (turns out this is pretty simple. You just stop being stupid and 
use the correct number.)

### English localisation

- [ ] finished

translation for all strings in DE.json into an appropriate english equivalent.

### Add special css rule for buttons

- [x] finished

add "user-select: none;" to all buttons.


### Add notification for new letters

- [x] finished

## Features for versions beyond the next

- notification information for letter and comments the account recevied.
- chat client between accounts and group chats for accounts