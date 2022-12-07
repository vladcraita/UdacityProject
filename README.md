# Udacity CRM Backend

## Installation
Make sure you have go installed locally.
Open a terminal window and type go env GOPATH, you should see something like /Users/vladcraita/go, depending on your OS.
Create a new package, ex /Users/vladcraita/go/src/github.com/vladcraita
Inside the newly created directory type:
git clone https://github.com/vladcraita/UdacityProject.git

The project uses external modules, to make them available for import enter:
go get github.com/gorilla/mux
go get github.com/jinzhu/copier

Go inside the UdacityProject folder then type:
go build

If everything is successful enter go run main.go. Server should start.

API operations are
GET /customers - retrieves all customers from db
POST /customers - creates a new customer
GET /customers/{id} - get a customer by id
DELETE/customers/{id} - removes a customer by id
PUT /customers/{id} - replaces a customer at id 
PATCH /customers/{id} - updates a customer by id

