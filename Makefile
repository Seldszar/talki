all: build-public binary-windows-amd64 binary-linux-amd64 binary-linux-arm64

build-public:
	cd public && npm install && npm run build

binary-windows-amd64:
	CGO_ENABLED=0 GOGC=off GOOS=windows GOARCH=amd64 go build -o "./dist/talki-windows-amd64.exe" .

binary-linux-amd64:
	CGO_ENABLED=0 GOGC=off GOOS=linux GOARCH=amd64 go build -o "./dist/talki-linux-amd64" .

binary-linux-arm64:
	CGO_ENABLED=0 GOGC=off GOOS=linux GOARCH=arm64 go build -o "./dist/talki-linux-arm64" .
