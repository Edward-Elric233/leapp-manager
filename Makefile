build:
	docker compose build --no-cache
	docker build -t edwardelric233/leapp-manager .

run:
	docker compose up -d
	docker run --name leapp-manager -d --restart always -p 3000:3000 -v D:\Code\web\leapp-manager\data:/data edwardelric233/leapp-manager

list:
	docker ps

push:
	docker push edwardelric233/leapp-manager