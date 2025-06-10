# JGame
An online quiz game built with Go and React. Play jeopardy-like games with your friends or strangers by joining or creating your own rooms with unique set of quizzes.

Since the recent changes to cookie policy of most browsers, the dev environment must be ran with HTTPS. You can use `mkcert` application to generate local SSL certificates. The certificates are stored in `cert` symlink to a directory with the certificate (both private key and certificate) or `cert` directory. This decision will persist in all projects, since HTTPS is the environment where the application will run in prod.

## Backend
Written in Go 1.24.2 with `gorilla/mux` to route the incoming traffic and `jackc/pgx` to access the **PostgreSQL** database. `xeipuuv/gojsonschema` 1.2.0 is used to compare requests against json schemas in `/api` directory.

> [!IMPORTANT]
> In order to run the application you need to create a `.env` file with the postgresql url set to DATABASE_URL environment variable. You should also set up authentication method to md5 inside `pg_hba.conf` file.

The backend relies on SSL certificate and key. You need to set `SSL_KEY_PATH` and `SSL_CERT_PATH` environment variables.


## Frontend
React with Typescript and `react-router-dom` for client-side routing, bundled with vite and served statically from the backend.
Simply do `npm install` and you will be fine.

The frontend relies on SSL certificate and key. You need to set `SSL_KEY_PATH` and `SSL_CERT_PATH` environment variables.
