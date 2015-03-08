.PHONY: build
build:
	echo building $(PLAYBOOK)
	go build -a $(PLAYBOOK)/main.go
