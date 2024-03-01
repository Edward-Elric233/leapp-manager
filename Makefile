run:
	docker run --name leapp-manager -d --restart always -p 3000:3000 -v D:\Code\web\leapp-manager\data:/data edwardelric233/leapp-manager

list:
	docker ps

build:
	docker build -t edwardelric233/leapp-manager .

push:
	docker push edwardelric233/leapp-manager