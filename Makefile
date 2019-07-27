APP_NAME = go-jenkins-exporter

.PHONY: help
help:
	@echo "Simple golang devenv helper"
	@echo ""
	@echo "Usage:"
	@echo "  deps		to install and update dependencies"
	@echo "  install	to install the app in $$GOPATH/bin"
	@echo "  help		to show this help"
	@echo ""

.PHONY: deps
deps:
	export GO111MODULE=on
	go get -u -v

.PHONY: install
install:
	go build -o $$GOPATH/bin/$(APP_NAME) -i main.go
	@echo alias $(APP_NAME)="$$GOPATH/bin/$(APP_NAME)"

