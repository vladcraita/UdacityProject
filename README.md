# Udacity CRM Backend

## Installation
Make sure you have go installed locally.

Clone the repo:

git clone https://github.com/vladcraita/UdacityProject.git

The project uses external modules, to make them available for import enter:

Go inside the UdacityProject folder then type:

go run main.go

If everything is successful enter go run main.go. Server should start.


API operations are
GET /customers - retrieves all customers from db

POST /customers - creates a new customer

payload: 
{
    "name": "Al Bundy",
    "role": "Shoe Salesman",
    "email": "al.bundy@garys.com",
    "phone": "1078212232",
    "contacted": false
}

GET /customers/{id} - get a customer by id

ex: GET /customers/1

DELETE/customers/{id} - removes a customer by id

ex: DELETE /customers/{id}

PUT /customers/{id} - replaces a customer at id 

ex: PUT /customers/1

payload 
{
    "name": "Al Bundy",
    "role": "Shoe Salesman",
    "email": "al.bundy@garys.com",
    "phone": "1078212232",
    "contacted": false
}

PATCH /customers/{id} - updates a customer by id

ex: PATCH /customers/1

payload
{
    "name": "Al Bundy",
    "role": "Shoe Salesman",
    "email": "al.bundy@garys.com",
    "phone": "111",
    "contacted": false
}

