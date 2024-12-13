NAME=cronjob-sandbox

template:
	helm template \
		--debug \
		--release-name $(NAME) \
		--namespace $(NAME) \
		./helm \
		> out.yaml

install:
	helm upgrade --install $(NAME) \
		--namespace $(NAME) \
		--create-namespace \
		./helm

uninstall:
	helm uninstall $(NAME) \
		--namespace $(NAME)