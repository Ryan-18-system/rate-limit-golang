# Go Rate Limiter

Rate limiter em Go com suporte a limitação por IP e por token de acesso, utilizando Redis como backend de persistência. A configuração é feita via variáveis de ambiente ou `.env`, com suporte a Docker e Docker Compose.

---

## 🔧 Funcionalidades

- Middleware para limitar requisições HTTP.
- Suporte a dois modos:
  - **Por IP:** Limita número de requisições por segundo por IP.
  - **Por Token:** Limita por token informado no cabeçalho `API_KEY`.
- Tokens sobrepõem limites por IP.
- Resposta HTTP 429 com mensagem personalizada ao exceder o limite.
- Configuração via `.env`.
- Backend de persistência: Redis.
- Estratégia pluggável para futuras substituições do Redis.
- Tempo de bloqueio configurável.

---

## 📦 Estrutura do Projeto

go-rate-limiter/
├── cmd/server # Entrada principal da aplicação
├── internal/config # Leitura de variáveis de ambiente
├── internal/limiter # Lógica de rate limiting e estratégia
├── internal/middleware # Middleware HTTP
├── tests/ # Testes automatizados
├── .env # Variáveis de ambiente
├── Dockerfile # Build da aplicação Go
├── docker-compose.yml # Redis + aplicação
└── README.md

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