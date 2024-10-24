# Tracks

An open source web application for displaying and managing GPX tracks.

## Development setup

Install node dependencies first:

```bash
cd vendor
npm install
```

### tailwind

Tailwind is used to build the static CSS bundle in /static. It watches `static/{html,js}/**/*.{html,js}`, `ui/html/**/*.html` for changes and recompiles the CSS.

```bash
cd vendor
npx tailwindcss -i tracks.css -o ../static/css/tracks.css --watch
```

### rollup

Rollup is used to build the static javascript bundle in / static.

```
cd vendor
npx rollup -c -w
```

### air

Air is used to live reload the go server. It doesn't work if the current working path has symlinks.

```
air -d
```
