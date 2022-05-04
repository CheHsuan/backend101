.PHONY: run clean 

run: clean
	docker-compose up --build -d
	docker-compose logs -f normal_app

clean:
	docker-compose down --remove-orphans
