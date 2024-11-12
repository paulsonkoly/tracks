# Tracks

An open source web application for displaying and managing GPX tracks.

## Development setup

Install node dependencies first:

```bash
cd assets
npm install
```

### tailwind

Tailwind is used to build the static CSS bundle in /static. It watches `static/{html,js}/**/*.{html,js}`, `ui/html/**/*.html` for changes and recompiles the CSS.

```bash
cd assets
npx tailwindcss -i tracks.css -o ../static/css/tracks.css --watch
```

### rollup

Rollup is used to build the static javascript bundle in / static.

```
cd assets
npx rollup -c -w
```

### air

Air is used to live reload the go server. It doesn't work if the current working path has symlinks.

```
air -d
```

## database

dbmate is used for migrations and sqlc is used for database access. To set up the database the following steps are required:

  - set up a database role for the application and update .env file with the database url. 
  - give super user privileges to the app user. This is required because we can't create the postgis extension otherwise.
  - create the database with `dbmate create`
  - run `dbmate load` to load schema into postgres.
  - take away super user privileges from app user.

### first admin user set up

 Create a row in the `users` table with username `admin` and hashed_password `$2a$12$nmN0KvqozAaYm6CNOXDVQOJB1JDrQ0BjTgt2YTml/M6ebbggf48Ra`. Now you can log in with admin and password: 1234567. Navigate to `/users`, edit user and edit the username / password.
