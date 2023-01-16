#!/bin/sh

cmd=$1
param1=$2

clean(){
	docker system prune -a -f
	docker image prune -a -f
	docker volume prune
}

up(){
	docker-compose up -d
}

down(){
	docker-compose down
	clean
	docker images
}

rebuild(){
	down
	up
}

run(){
	docker exec -it $param1 /bin/sh
}

case "$1" in
	clean) clean;;
	up) up;;
	down) down;;
	rebuild) rebuild;;
	run) run;;
esac
