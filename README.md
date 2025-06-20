# JGame
An online quiz game built with Go and React. Play kahoot-like games with your friends or strangers by joining or creating your own rooms with unique set of quizzes.

Since the recent changes to cookie policy of most browsers, the dev environment must be ran with HTTPS. You can use `mkcert` application to generate local SSL certificates. The certificates are stored in `cert` symlink to a directory with the certificate (both private key and certificate) or `cert` directory. This decision will persist in all projects, since HTTPS is the environment where the application will run in prod.

## Backend
Written in Go 1.24.2 with `gorilla/mux` to route the incoming traffic and `jackc/pgx` to access the **PostgreSQL** database. `xeipuuv/gojsonschema` 1.2.0 is used to compare requests against json schemas in `/api` directory. `gorilla/websocket` is responsible for managing websocket connection for game rooms to update and synchronize game state in real-time.

> [!IMPORTANT]
> In order to run the application you need to create a `.env` file with the postgresql url set to DATABASE_URL environment variable. You should also set up authentication method to md5 inside `pg_hba.conf` file.

The backend relies on SSL certificate and key. You need to set `SSL_KEY_PATH` and `SSL_CERT_PATH` environment variables.


## Frontend
React with Typescript and `react-router-dom` for client-side routing, bundled with vite and served statically from the backend.
Simply do `npm install` and you will be fine.

The frontend relies on SSL certificate and key. You need to set `VITE_SSL_KEY_PATH` and `VITE_SSL_CERT_PATH` environment variables.

## For those who find this project
This is my first web application which I developed in one month. I'm not so proud of it, because the internal implementation is very dirty: there are a lot of under- and overfetch issues, and overall the backend implementation is not DRY. There's an issue with handling dead websocket connections which I didn't figure out how to solve, since I don't understand what's the root of the problem: the client or the server (create a room, enter it, refresh the page and the connection will never be ok - maybe someone will figure it out).

It's my first experience using React and Typescript in a such large project. I'm very satisfied with the library and its tools and I will definitely continue using it. When it comes to the backend, I was both disappointed and satisfied with what the Go programming language has to offer. Without a doubt, it's got a very strong easy-to-use standard `net/http` module which makes developing web applications very easy and fast. I struggeled a lot with the error handling in go, since it's too explicit: I'm 100% sure I've done some useless error handling during the course of development. Also, I know that I want to use ORMs to work with databases instead of using raw SQL queries. I didn't like working with the third party http and websocket libraries, even though they are simple to use.
