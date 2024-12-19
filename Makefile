include .env
export $(shell sed 's/=.*//' .env)
TAG=v8
NAME=cronjob-sandbox

template:
	helm template \
		--debug \
		--release-name $(NAME) \
		--namespace $(NAME) \
		--set image=$(FQIN) \
		--set tag=$(TAG) \
		./helm \
		> out.yaml

install:
	@echo "FQIN is $(FQIN)"
	helm upgrade --install $(NAME) \
		--namespace $(NAME) \
		--set image=$(FQIN) \
		--set tag=$(TAG) \
		--create-namespace \
		./helm

uninstall:
	helm uninstall $(NAME) \
		--namespace $(NAME)