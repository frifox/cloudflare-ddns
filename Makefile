default:
	rm -f bootstrap cloudflare-ddns.zip
	GOOS=linux GOARCH=arm64 go build -o bootstrap -tags lambda.norpc -ldflags "-s -w"
	zip cloudflare-ddns.zip bootstrap
	rm bootstrap