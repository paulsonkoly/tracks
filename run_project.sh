#!/bin/bash

set -e

PROJECT_DIR="$(cd "$(dirname "$0")"; pwd)"

REQUIRED_TOOLS=("npm" "npx" "node" "psql" "dbmate")

echo "🔍 Checking for required tools..."

for tool in "${REQUIRED_TOOLS[@]}"; do
  if ! command -v "$tool" >/dev/null 2>&1; then
    echo "❌ Missing required tool: $tool"
    echo "💡 Please install $tool before continuing."
    exit 1
  fi
done

# Special check for 'air'
if command -v air >/dev/null 2>&1; then
  :
elif [ -x "$HOME/go/bin/air" ]; then
  export PATH="$HOME/go/bin:$PATH"
else
  echo "❌ Missing required tool: air"
  echo "💡 You can install air with:"
  echo "   go install github.com/cosmtrek/air@latest"
  exit 1
fi

echo "✅ All required tools are installed."

osascript <<EOF
tell application "iTerm"
  -- Tab 1: init-db + Tailwind watcher
  set newWindow to (create window with default profile)
  tell current session of newWindow
      write text "./init_db.sh && cd \"$PROJECT_DIR/assets\" && npm install && npx tailwindcss -i tracks.css -o ../static/css/tracks.css --watch"
  end tell

  delay 3

  -- Tab 2: Rollup watcher
  tell current window
      create tab with default profile
      tell current session
          write text "cd \"$PROJECT_DIR/assets\" && npx rollup -c -w"
      end tell
  end tell

  -- Tab 3: air dev server
  tell current window
      create tab with default profile
      tell current session
          write text "cd \"$PROJECT_DIR\" && air -d"
      end tell
  end tell
end tell
EOF

sleep 5s
open http://localhost:9999/
