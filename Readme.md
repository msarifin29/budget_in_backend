# Required 
[go version 1.21.5](https://go.dev/dl/)
[postgresql](https://www.postgresql.org/)

## How to Run

* first
 make schema database in folder database/migration
* second 
create `dev.env` in root project
* third
copy this into your `dev.env` file
```
ENVIRONMENT=development
APP=budget_in
DB_POSTGRES_SOURCE=postgresql://postgres:{PASSWORD}@localhost:5432/{DB_NAME}?sslmode=disable
DB_POSTGRES_DRIVER=postgres
ACCESS_TOKEN_DURATION=131400m
TOKEN_SYMMETRIC_KEY=123456789ABCDEFGHIJ987654321ACBD	
SERVER_ADDRESS=0.0.0.0
SET_MAX_IDLE_CONNS=64
SET_MAX_OPEN_CONNS=64
SET_CONN_MAX_LIFE_TIME=60
SET_CONN_MAX_IDLE_TIME=10
SENDER_NAME =Budget In {your email}
AUTH_EMAIL={your email}
AUTH_PASSWORD={your auth password}
```
following this article to generated auth password [article](https://knowledge.workspace.google.com/kb/how-to-create-app-passwords-000009237)
