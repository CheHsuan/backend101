.PHONY: proto clean 

proto:
	./build.sh

run: proto clean server
	docker-compose up --build -d

clean:
	docker-compose down --remove-orphans
