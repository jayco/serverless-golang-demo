.PHONY: build build-debug clean deploy gomod local test-unit test-integration

# Run system build
build:
	GO111MODULE=on GOARCH=amd64 GOOS=linux go build -ldflags="-w -s" -o bin/health internal/routes/health/main.go

# Run local unit tests from underlying system
test-unit:
	GO111MODULE=on go test -count=1 -timeout 30s github.com/jayco/serverless-golang-demo/internal/routes/health -run TestHandlerUnit

# integration tests require sam-template and debug build
test-integration: clean build-debug generate-template
	GO111MODULE=on go test -v -timeout 30s github.com/jayco/serverless-golang-demo/internal/routes/health -run TestHandlerIntegration

build-debug:
ifneq ($(shell test -a $(PWD)/.debug/dlv && echo 1), 1)
	echo "Building delve debugger"
	@docker run -t -i --rm \
	-e "GO111MODULE=on" \
	-e "CGO_ENABLED=1" \
	-e "GOARCH=amd64" \
	-e "GOOS=linux" \
	-v $(PWD):/src \
	-w /src/ lambci/lambda:build-go1.x \
	go build -o .debug/dlv github.com/derekparker/delve/cmd/dlv
endif
	@echo "Building health route"
	@docker run \
	-e "GO111MODULE=on" \
	-e "CGO_ENABLED=1" \
	-e "GOARCH=amd64" \
	-e "GOOS=linux" \
	-v $(PWD):/src \
	-w /src/ lambci/lambda:build-go1.x \
	go build -gcflags='-N -l' -o bin/health internal/routes/health/main.go

generate-template:
	@echo "Generating local template for SAM"
	serverless plugin install --name serverless-sam
	serverless sam export --output ./sam-local-template.yml

local-debug: clean build-debug generate-template
	@echo "Bringing stack up"
	sam local start-api --template ./sam-local-template.yml -d 5986 --debugger-path ./.debug

local: clean build
	@echo "Generating local template for SAM"
	serverless plugin install --name serverless-sam
	serverless sam export --output ./sam-local-template.yml

	@echo "Bringing stack up"
	sam local start-api --skip-pull-image --template ./sam-local-template.yml

clean:
	rm -rf ./bin sam-local-template.yml node_modules

deploy: clean build
	sls deploy --verbose

gomod:
	chmod u+x gomod.sh
	./gomod.sh
