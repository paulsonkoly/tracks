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

dbmate is used for migrations and sqlc is used for database access.
