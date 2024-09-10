TARGET=pi@domestication.fritz.box

copy-prometheus:
	scp -r prometheus ${TARGET}:/tmp/prometheus
	ssh ${TARGET} 'sudo mv /tmp/prometheus /var/opt/'

copy-docker-compose:
	scp docker-compose.yml ${TARGET}:/tmp/docker-compose.yml
	ssh ${TARGET} 'sudo mv /tmp/docker-compose.yml /var/opt/'