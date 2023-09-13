.PHONY: update-go-deps
build-cli:
	@echo "Building CLI"
	GOOS=linux GOARCH=amd64 go build -o bin/move-linux-amd64 cmd/move/move.go
	# GOOS=linux GOARCH=arm64 go build -o bin/move-linux-arm64 cmd/move/move.go
build-gui:
	@echo "Building GUI"
	fyne-cross linux cmd/gui
update-go-deps:
	@echo ">> updating Go dependencies"
	@for m in $$(go list -mod=readonly -m -f '{{ if and (not .Indirect) (not .Main)}}{{.Path}}{{end}}' all); do \
		go get $$m; \
	done
	go mod tidy
ifneq (,$(wildcard vendor))
	go mod vendor
endif