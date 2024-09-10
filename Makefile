TARGET=pi@domestication.fritz.box

include Makefile.secrets

build-rpi-cross-docker:
	docker build -f Dockerfile.cross . -t rpi-cross

build-rpi:
	docker run --rm -v "${PWD}":/usr/src/myapp -w /usr/src/myapp rpi-cross /bin/sh -c 'go build -o build/bot cmd/bot/main.go'

.PHONY: deploy

run-local:
	LINKBOT_TOKEN=${LINKBOT_TOKEN} LINKBOT_INITIALUSERID=${LINKBOT_INITIALUSERID} LINKBOT_INITIALUSERNAME=${LINKBOT_INITIALUSERNAME} go run cmd/bot/main.go

copy-secrets:
	scp Makefile.secrets ${TARGET}:/tmp/Linkbot.secrets
	ssh ${TARGET} 'sudo mv /tmp/Linkbot.secrets /var/opt/'
	# create data dir
	ssh ${TARGET} 'sudo mkdir -p /var/opt/linkbot'

deploy: build-rpi

	ssh ${TARGET} 'mkdir -p /tmp/linkbot'
	scp build/bot ${TARGET}:/tmp/linkbot/
	scp Dockerfile ${TARGET}:/tmp/linkbot/

	ssh ${TARGET} 'cd /tmp/linkbot && docker build . -t linkbot'
	-ssh ${TARGET} 'docker stop linkbot_prod'
	-ssh ${TARGET} 'docker rm linkbot_prod'
	ssh ${TARGET} 'docker run --name linkbot_prod \
		--env-file /var/opt/Domestication.secrets \
		-p 9099:9090 \
		-v /var/opt/linkbot:/data \
		--env DOMESTICATION_DATABASE=/data/linkbot_data.db \
		-d --restart=always \
		linkbot'

