sudo docker run -p 3307:3306 --name mysql-instance -v $PWD/conf:/etc/mysql/conf.d -v $PWD/logs:/logs -v $PWD/data:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=123456 -d mysql
sudo docker run --name beego-instance -p 8081:8080 -d beego-vue-app
