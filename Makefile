# Definição de variáveis
DOCKER_COMPOSE = docker-compose -f ./docker/docker-compose.yml
ENV_FILE = configs/development.env
CONTAINER_NAME = container-golang

# Definição de cores
GREEN = \033[0;32m
NC = \033[0m

.PHONY: vendor

# Execução principal e em ordem
up: clean create-symlink start

# Limpar containers e imagens quebradas do docker
clean:
	@echo "$(GREEN)Limpando containers e imagens...$(NC)"
	./scripts/docker/cleanup.sh

# Construir e iniciar os containers
build:
	@echo "$(GREEN)Construindo e iniciando os containers...$(NC)"
	$(DOCKER_COMPOSE) build

# Iniciar os containers
start:
	@echo "$(GREEN)Iniciando os containers...$(NC)"
	$(DOCKER_COMPOSE) up

# Iniciar sem containers, apenas com go local
run: #create-symlink
	@echo "$(GREEN)Iniciando aplicação local...$(NC)"
	go run ./cmd/webserver/main.go

# Parar os containers
stop:
	@echo "$(GREEN)Parando os containers...$(NC)"
	$(DOCKER_COMPOSE) down

# Carregar as configurações de ambiente
reload-env:
	@echo "$(GREEN)Carregando as variáveis de ambiente do arquivo .env...$(NC)"
	source $(ENV_FILE)

# Configurar as variáveis de user e token do Git local: ~/.gitconfig
# Para usar o dagobah, são usadas no dockerfile
configure-git:
	@echo "$(GREEN)Configurando variáveis de usuário e token do Git...$(NC)"
    export GIT_USER=$$(git config --global user.name)
    export GIT_TOKEN=$$(git config --global github.token)

# Criar o link simbólico para o arquivo de configuração de ambiente
create-symlink:
	@echo "$(GREEN)Criando link simbólico para o arquivo de configuração de ambiente...$(NC)"
	rm -rf .env
	ln -s $(ENV_FILE) .env

generate-mocks:
	./scripts/generate_mocks.sh
# Criar a vendor do projeto com o docker golang

vendor:
	@if docker ps --filter name="$(CONTAINER_NAME)" -q >/dev/null; then \
		echo "Criando a pasta vendor..."; \
		if [ -d "vendor" ]; then \
			echo "A pasta vendor já existe. Removendo..."; \
			rm -rf "vendor"; \
		fi; \
		mkdir -p "vendor"; \
		if [ -f "go.mod" ]; then \
			echo "Baixando as dependências do projeto..."; \
			docker exec -it "$(CONTAINER_NAME)" sh -c "go mod vendor -v"; \
			docker exec -it "$(CONTAINER_NAME)" sh -c "chown -R 1000:1000 vendor"; \
		else \
			echo "Arquivo go.mod não encontrado. Certifique-se de que esteja no diretório do projeto Go."; \
		fi; \
	else \
		echo "O contêiner não está em execução."; \
	fi