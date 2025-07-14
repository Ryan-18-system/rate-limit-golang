# Go Rate Limiter

Rate limiter em Go com suporte a limitação por IP e por token de acesso, utilizando Redis como backend de persistência. A configuração é feita via variáveis de ambiente ou `.env`, com suporte a Docker e Docker Compose.

---

## 🔧 Funcionalidades

* Middleware para limitar requisições HTTP.
* Suporte a dois modos:

  * **Por IP:** Limita número de requisições por segundo por IP.
  * **Por Token:** Limita por token informado no cabeçalho `API_KEY`.
* Tokens sobrepõem limites por IP.
* Resposta HTTP 429 com mensagem personalizada ao exceder o limite.
* Configuração via `.env`.
* Backend de persistência: Redis.
* Estratégia pluggável para futuras substituições do Redis.
* Tempo de bloqueio configurável.

---

## 📦 Estrutura do Projeto

```
go-rate-limiter/
├── cmd/server             # Entrada principal da aplicação
├── internal/config        # Leitura de variáveis de ambiente
├── internal/limiter       # Lógica de rate limiting e estratégia
├── internal/middleware    # Middleware HTTP
├── tests/                 # Testes automatizados
├── .env                   # Variáveis de ambiente
├── Dockerfile             # Build da aplicação Go
├── docker-compose.yml     # Redis + aplicação
└── README.md
```

---

## ⚙️ Variáveis de Ambiente

Defina no arquivo `.env`:

```dotenv
PORT=8080

# Limites por IP
RATE_LIMIT_IP=10
RATE_LIMIT_IP_DURATION=1s

# Limites por Token
RATE_LIMIT_TOKEN=100
RATE_LIMIT_TOKEN_DURATION=1s

# Tempo de bloqueio após ultrapassar limite
BLOCK_DURATION=300s

# Redis
REDIS_ADDR=redis:6379
REDIS_PASSWORD=
REDIS_DB=0
```

---

## ▶️ Como Rodar o Projeto

Certifique-se de ter o **Docker** e o **Docker Compose** instalados na sua máquina.

### 1. Clone o repositório

```bash
git clone https://github.com/seu-usuario/go-rate-limiter.git
cd go-rate-limiter
```

### 2. Configure as variáveis de ambiente

Edite o arquivo `.env` conforme necessário. Um exemplo já está incluso no projeto.

### 3. Suba os serviços com Docker Compose

```bash
docker-compose up --build
```

A aplicação estará disponível em:

```
http://localhost:8080
```

---

## ✅ Como Testar o Rate Limiter

Você pode testar usando **curl**, **Postman** ou outra ferramenta de requisição HTTP.

### 🔹 Testar Limite por IP

Faça múltiplas requisições sem cabeçalho de token:

```bash
for i in {1..15}; do curl -i http://localhost:8080; done
```

Se exceder o limite configurado no `.env`, você verá:

```
HTTP/1.1 429 Too Many Requests
you have reached the maximum number of requests or actions allowed within a certain time frame
```

### 🔹 Testar Limite por Token

Envie requisições com o cabeçalho `API_KEY`:

```bash
for i in {1..120}; do
  curl -i -H "API_KEY: token123" http://localhost:8080
done
```

Substitua `token123` por qualquer string. O rate limiter aplicará as regras definidas para `RATE_LIMIT_TOKEN`.

---

## 📊 Testes Automatizados

Se houver testes no diretório `tests/`, rode com:

```bash
go test ./tests/...
```

(Dentro do container ou com Go instalado localmente.)
