## How to setup the RPI

```bash

# install docker
curl -sSL https://get.docker.com | sh
sudo usermod -aG docker pi

# install python deps for docker-compose
sudo apt-get install libffi-dev libssl-dev
sudo apt install python3-dev
sudo apt-get install -y python3 python3-pip

# install+enable it
sudo pip3 install docker-compose
sudo systemctl enable docker

# copy the docker-compose file:
make copy-docker-compose

# start the docker-compose (on server)
docker-compose up -d


```