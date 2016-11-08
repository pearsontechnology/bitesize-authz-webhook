NAME := authz-webhook
REPO := jetstackexperimental/authz-webhook

clean:
	rm -f $(NAME)-amd64

build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $(NAME)-amd64

test:
	go test . -v

image: build
	docker build -t $(REPO):temp .

push: SHELL=/bin/bash
push: image
	@set -e -x ; \
	if [[ $$TRAVIS_BRANCH == 'master' ]]; then \
		TAGS="canary" ; \
	fi ; \
	if [[ ! -z "$$TRAVIS_TAG" ]]; then \
		TAGS="$$TRAVIS_TAG latest" ; \
	fi ; \
	for tag in $${TAGS}; do \
		docker tag  $(REPO):temp   $(REPO):$${tag} ; \
		docker push $(REPO):$${tag} ; \
	done
