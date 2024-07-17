win32:
	GOOS=windows go build -ldflags "-H=windowsgui" -o dedawn.exe dedawn/cmd/dedawn

linux:
	go build -o dedawn dedawn/cmd/dedawn

server:
	go build -o dedawn_server dedawn/cmd/dedawn_server