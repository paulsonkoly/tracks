.PHONY:
deploy: copy-up migrate
	@echo "Restart the prod server"

.PHONY:
copy-up:
	env GOOS=linux GOARCH=amd64 go build -o tracks cmd/tracks/tracks.go
	ssh tracks@ssh-tracks.alwaysdata.net RemoteCommand="\"rm -rf tracks\""
	scp -r tracks tracks@ssh-tracks.alwaysdata.net:
	ssh tracks@ssh-tracks.alwaysdata.net RemoteCommand="\"rm -rf ui\""
	scp -r ui tracks@ssh-tracks.alwaysdata.net:
	ssh tracks@ssh-tracks.alwaysdata.net RemoteCommand="\"rm -rf static\""
	scp -r static tracks@ssh-tracks.alwaysdata.net:

.PHONY:
migrate:
	@echo "!!!!"
	@echo "!!!! Make sure DATABASE_URL is set to production"
	@echo "!!!!"
	dbmate status
	dbmate up
