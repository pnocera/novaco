# novaco

Nomad / Vault / Consul / Git http server in one service

echo "# testapp" >> README.md
git init
git add README.md
git commit -m "first commit"
git branch -M main
git remote add origin http://localhost:8888/admin/testapp.git
git push -u origin main

vault : http://192.168.1.145:8200
consul : http://localhost:8500
nomad : http://192.168.1.145:4646
git server : http://192.168.1.145:8888

openssl req -new -newkey rsa:4096 -days 3650 -nodes -x509 -subj "/C=US/ST=Denial/L=Springfield/O=Dis/CN=www.example.com" -keyout www.example.com.key  -out www.example.com.cert