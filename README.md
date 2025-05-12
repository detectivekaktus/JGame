# JGame
An online quiz game built with Go and React. Play jeopardy-like games with your friends or strangers by joining or creating your own rooms with unique set of quizzes.


## Backend
Written in Go 1.24.2 with `gorilla/mux` to route the incoming traffic and `jackc/pgx` to access the **PostgreSQL** database.

> [!IMPORTANT]
> In order to run the application you need to create a `.env` file with the postgresql url set to DATABASE_URL environment variable. You should also set up authentication method to md5 inside `pg_hba.conf` file.


## Frontend
React with Typescript and `react-router-dom` for client-side routing, bundled with vite and served statically from the backend.
Simply do `npm install` and you will be fine.
