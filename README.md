# Entain-task

Test task for Entain, processes transaction records.

## Usage:

POST <endpoint_url>/process-record

Request body example

    "state": "lost",
    "amount": 1,
    "transaction_id": "6c21cc27-5be8-46ed-b898-d4f3ed6d58b2",
    "user_id": "0268c107-c4cf-45a5-8547-c16f86616d61"

"user_id" in example is default one. Odd records are processed every 1 minute

## Migrations

Applied automatically on app run


## How to run locally with docker
Create docker env .list file with DB credentials as environment variables

DB_USER=

DB_PASS=

DB_HOST=

DB_PORT=

DB_NAME=

Run docker container with arg ‘--env-file ./<filename>.list’

## How to run locally without 
Run entain-task/cmd/server/main.go with os env db variables, as above
