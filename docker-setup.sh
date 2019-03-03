curl -sSL https://get.daocloud.io/daotools/set_mirror.sh | sh -s http://f1361db2.m.daocloud.io
sudo systemctl restart docker.service

echo "Restarted docker."
sudo docker pull mysql
sudo docker build -t beego-vue-app .
