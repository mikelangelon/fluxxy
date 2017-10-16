NO_VENDOR := `go list bol/... | grep -v 'bol/vendor/'`

build:
	go build -ldflags "-s $(LDFLAGS)" game/cmd/fluxxy

sanitize:
	go fmt $(NO_VENDOR)
	go vet $(NO_VENDOR)

clean:
	rm -f pdfMaker
