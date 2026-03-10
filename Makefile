.PHONY: build clean dev test sync-skill build-web dev-web

sync-skill:
	cp .claude/commands/prm.md cmd/skill_prm.md

build-web:
	cd web && npm run build
	rm -rf internal/web/static
	mkdir -p internal/web/static
	cp -r web/dist/* internal/web/static/

build: sync-skill build-web
	go build -o bin/prm .

clean:
	rm -f bin/prm
	rm -rf web/dist internal/web/static/*

dev:
	go run . $(ARGS)

dev-web:
	cd web && npm run dev

test:
	go test ./...
