# SnippetBox Golang

This is a **modified version** of the [Let's Go e-book by Alex Edwards](https://lets-go.alexedwards.net/). It's a Go web application called "Snippetbox" that lets users CRUD text snippets. Developed without using web framework.

<img width="500" src="https://github.com/cullenjett/snippetbox/raw/master/lets-go-screenshot.png">

### Features

 0. Basic routing :white_check_mark:
 1. Configuration management and error handling :white_check_mark:
 2. Database connection management :white_check_mark: 
 3.  Dynamic template rendering :white_check_mark:
 4.  Middleware chaining using [justinas/alice](https://github.com/justinas/alice) :white_check_mark:
 5.  RESTful Routing using [bmizerany/pat](https://github.com/bmizerany/pat) router :white_check_mark:
 6.  Form processing and validation :white_check_mark:
 7.  Stateful HTTP (Session Manager) using [golangcollege/sessions](https://github.com/golangcollege/sessions) 
 :white_check_mark:
 8.  Self-signed TLS, HTTPS Server :white_check_mark:
 9.  User authentication and authorization :white_check_mark:
 - Password Encryption :white_check_mark:
 - Login/Logout & Authorization :white_check_mark:
 -  [CSRF](https://www.gnucitizen.org/blog/csrf-demystified/) Protection using [justinas/nosurf](https://github.com/justinas/nosurf) :building_construction:
 10.    Golang `context` to carry request-scoped data :white_check_mark:
 11.    Testing:
 - Simple unit testing :white_check_mark:
 - Testing HTTP handlers :white_check_mark:
 - End-to-end testing :white_check_mark:
 - Mocking Database Dependencies :white_check_mark:
 - Testing HTML Forms :building_construction:
 - Integration Test :white_check_mark:

#### To do
- :white_check_mark: Configuration management from `.env` using [envdecode](https://github.com/joeshaw/envdecode)
- :building_construction: Dockerize app + MySQL 
- :building_construction: Project desc
- :building_construction: API Endpoint


### Development


##### Configuration
Create an `.env` file containing:
```
HTTP_SERVER_PORT=4000          //Listening port
HTTP_SERVER_TIMEOUT_IDLE=60s  
HTTP_SERVER_TIMEOUT_READ=5s
HTTP_SERVER_TIMEOUT_WRITE=10s
STATIC_PATH="./ui/static"
DEBUG=True
DB_HOST=db
DB_PORT=3306
DB_USER=#database user 
DB_PASS=#database password
DB_NAME=#database name 
TLS_CERT_PATH=#path to cert.pem"
TLS_PKEY_PATH=#path to key.pem"
COOKIE_PKEY="s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge"

```
Export `.env` before:

##### `go run cmd/web/*`

Starts the local web server with HTTPS on designated port ( e.g. [https://localhost:4000](https://localhost:4000)).
The port can be specified also with `--addr` param:
##### `go run cmd/web/* --addr :4000`

