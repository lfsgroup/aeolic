.PHONY: all

.DEFAULT_GOAL := help

SHELL=bash
PROJECT_ROOT:=$(shell git rev-parse --show-toplevel)
COMMIT := $(shell git rev-parse --short HEAD)
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
REPO := $(shell basename `git rev-parse --show-toplevel`)
DATE := $(shell date +%Y-%m-%d-%H-%M-%S)
APP_NAME := aeolic

# Load env properties , db name, port, etc...
# nb: You can change the default config with `make ENV_CONTEXT=".env.uat" `
ENV_CONTEXT ?= $(PROJECT_ROOT)/.env.local
LOCAL_ENV_MINE=$(PROJECT_ROOT)/.env.local.mine
## Override any default values in the parent .local.env, with your own
-include $(ENV_CONTEXT) $(LOCAL_ENV_MINE)


test: ## Run unit tests
	go test --short -cover -failfast ./...

test_integration: test_build ## Integration test, post to slack channel
	SLACK_API_TOKEN=$(SLACK_API_TOKEN) TEST_SLACK_CHANNEL=$(TEST_SLACK_CHANNEL)  ./$(APP_NAME)

test_build: ## Build test binary
	go build -o $(APP_NAME) ./cmd/slack


# HELP
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ## This help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)



#####################
#  Private Targets  #
#####################

log: # log env vars
	@echo "\n"
	@echo "COMMIT               $(COMMIT)"
	@echo "BRANCH               $(BRANCH)"
	@echo "APP_NAME             $(APP_NAME)"
	@echo "DATE                 $(DATE)"
	@echo "\n"
