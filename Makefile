startdb-docker:
	 docker run -d --rm --name snippetboxdb -v snippetboxdb:/var/lib/mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=root mysql:latest

startdb:
	 brew services start mysql

stopdb:
	 brew services stop mysql