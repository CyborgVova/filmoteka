.PHONY: all
all:
	docker-compose up

.PHONY: clean
clean:
	docker-compose down

.PHONY: re
re: clean all