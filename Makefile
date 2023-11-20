BUILD_DIR=build
CMD_DIR=cmd
CMDS=$(patsubst $(CMD_DIR)/%,%,$(wildcard $(CMD_DIR)/*))
LOCAL_MOD_NAME=$(shell go list -m)

.PHONY: fmt check test

all: fmt check test bin

fmt:
	gofmt -s -w -l .
	@echo 'goimports' && goimports -w -local $(LOCAL_MOD_NAME) $(shell find . -type f -name '*.go' -not -path "./internal/*")
	gci write -s standard -s default -s "Prefix($(LOCAL_MOD_NAME))" --skip-generated .
	@files=$$(git diff --diff-filter=AM --name-only HEAD^ | grep '.go$$'); \
	if [ -n "$$files" ]; then \
	  echo 'golines' $$files; \
		golines --ignore-generated -m 80 --reformat-tags --shorten-comments -w $$files; \
	fi
	go mod tidy

check:
	revive -exclude pkg/... -formatter friendly -config tools/revive.toml  ./...
	find . -name "*.json" | xargs -n 1 -t gojq . >/dev/null
	go vet -all -printfuncs=unexpectederrorf,paramerrorf,wrapf,warnf,warningf,infof,debugf,errorf,fatalf,panicf,debug,info,warning,error,fatal,panic ./...
	golangci-lint run
	misspell -error */**
	@echo 'staticcheck' && staticcheck $(shell go list ./... | grep -v internal)

vet:
	./vet.sh

test:
	go test ./...

bin:
	./build_docker.sh --bin

docker: $(CMDS)
	./build_docker.sh --docker $^

$(CMDS):
	./build_docker.sh --bin $@

proto:
		protoc --go_out=internal/rpc/hello --go_opt=module=protos --go-grpc_out=internal/rpc/hello --go-grpc_opt=module=protos protos/HelloService.proto

setup: setup-postgres setup-redis setup-stan

# add this in /etc/fstab and run `sudo mount -a`
# tmpfs /mtmp tmpfs size=2048m,mode=1777 0 0
setup-postgres:
	@if ! docker ps | /bin/grep postgres-localdev; then \
		docker run --name postgres-localdev \
			-p 5432:5432 \
			-d --tmpfs /var/lib/postgresql/data:rw,noexec,nosuid,size=4096m \
			-d --tmpfs /run:rw,noexec,nosuid,size=4096m \
			-e POSTGRES_DB=testing \
			-e POSTGRES_USER=tester \
			-e POSTGRES_PASSWORD=aaaa1234 \
			--restart always \
			-d postgres:12 \
			-c fsync=off -c full_page_writes=off; \
		docker run --rm --link postgres-localdev:postgres-localdev aanand/wait; \
	fi

setup-redis:
	@if ! docker ps | /bin/grep redis-localdev; then \
		docker run --name redis-localdev -p 6379:6379 \
			--restart always \
			-d redis:alpine; \
		docker run --rm --link redis-localdev:redis-localdev aanand/wait; \
	fi

setup-stan:
	@if ! docker ps | /bin/grep nats-localdev; then \
		docker run --restart=always -d --name nats-localdev -p 5222:4222 nats-streaming; \
	fi

remove:
	docker rm -f postgres-localdev redis-localdev nats-localdev
