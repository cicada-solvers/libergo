echo "cleaning up the containers"
podman stop postgres
podman stop mysql1
podman stop mysql2
podman stop mysql3

podman rm -f postgres
podman rm -f mysql1
podman rm -f mysql2
podman rm -f mysql3

podman rmi -f docker.io/library/postgres
podman rmi -f docker.io/library/mysql

podman volume rm --all
podman volume prune --force

echo "Pulling the images"
podman pull docker.io/library/postgres
podman pull docker.io/library/mysql

echo "Bringing up the db containers"
podman run --name postgres -e POSTGRES_PASSWORD=lppasswd -dt -p 5432:5432 docker.io/library/postgres
podman run --name mysql1 -e MYSQL_ROOT_PASSWORD=dpasswd -e MYSQL_DATABASE=wordsdb -e MYSQL_USER=runedonkey -e MYSQL_PASSWORD=dpasswd -dt -p 3306:3306 docker.io/library/mysql
podman run --name mysql2 -e MYSQL_ROOT_PASSWORD=dpasswd -e MYSQL_DATABASE=wordsdb -e MYSQL_USER=runedonkey -e MYSQL_PASSWORD=dpasswd -dt -p 3307:3306 docker.io/library/mysql
podman run --name mysql3 -e MYSQL_ROOT_PASSWORD=dpasswd -e MYSQL_DATABASE=wordsdb -e MYSQL_USER=runedonkey -e MYSQL_PASSWORD=dpasswd -dt -p 3308:3306 docker.io/library/mysql

echo "seeing processes"
podman ps -a