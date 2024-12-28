# Название бинарника
BINARY_NAME=myapp

# Путь к директории с исходниками
SRC_DIR=./cmd

# Переменные для флагов компилятора
GO=go
GOFMT=gofmt
GOFLAGS=-v
LDFLAGS=-s -w

# Команда для компиляции бинарника
build:
	$(GO) build $(GOFLAGS) -o $(BINARY_NAME) $(SRC_DIR)

# Команда для запуска приложения
run:
	$(GO) run $(SRC_DIR)

# Команда для форматирования кода
fmt:
	$(GOFMT) -w $(SRC_DIR)

# Команда для тестирования
test:
	$(GO) test $(GOFLAGS) ./...

# Команда для проверки зависимостей
deps:
	$(GO) mod tidy

# Очистка скомпилированных файлов
clean:
	rm -f $(BINARY_NAME)

# Команда для сборки релизной версии
release:
	GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) -o $(BINARY_NAME)-linux-amd64 $(SRC_DIR)

# Выполнение всех шагов: fmt, test и build
all: fmt test build
