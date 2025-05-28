#!make
.PHONY: docker-build docker-tag docker-push docker runclient runserver
# Makefile for Rikon Scraper

include rikonscraper.env

docker-build:
	docker build -t rikon-scraper/go ./server/

docker-tag:
	docker tag rikon-scraper/go:latest 640592576015.dkr.ecr.us-east-2.amazonaws.com/rikon-scraper/go:latest

docker-push:
	docker push 640592576015.dkr.ecr.us-east-2.amazonaws.com/rikon-scraper/go:latest

docker: docker-build docker-tag docker-push

runclient:
	@cd web && npm run dev

runserver:
	@cd server && air