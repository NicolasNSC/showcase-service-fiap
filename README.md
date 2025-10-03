# showcase-service-fiap 

Microserviço para gerenciar o ciclo de vida de vendas na plataforma de revenda de automóveis, incluindo listagem de veículos, compra e confirmação de pagamento via webhook. O projeto segue os princípios da Clean Architecture, promovendo separação de responsabilidades, alta testabilidade e baixo acoplamento entre as camadas.

## Tecnologias Utilizadas

-   **Linguagem**: Go (v1.23+)
-   **Banco de Dados**: PostgreSQL
-   **Infraestrutura**: Docker & Docker Compose
-   **Roteador HTTP**: Chi
-   **Testes**: Testify & Gomock
-   **Documentação da API**: Swagger (OpenAPI)

## Como Executar

O projeto é totalmente containerizado, exigindo apenas **Docker** e **Docker Compose** instalados.

### 1. Clone o repositório

### 2. Configure as Variáveis de Ambiente

Crie um arquivo `.env` na raiz do projeto, baseado no `.env-sample`.

### 3. Suba os Containers

Use o Makefile para construir a imagem e iniciar a aplicação e o banco de dados:

```bash
make docker-up
```

A API estará disponível em [http://localhost:8081](http://localhost:8081).  
A documentação Swagger estará em [http://localhost:8081/swagger/index.html](http://localhost:8081/swagger/index.html).

---
### Comandos Úteis (Makefile)

- `make docker-up`: Inicia todo o ambiente containerizado.

- `make docker-down`: Para e remove os containers, redes e volumes.

- `make test`: Roda a suíte de testes unitários.

- `make cov`: Gera e abre o relatório de cobertura de testes no navegador.

---

## Endpoints da API

A documentação interativa completa está disponível em `/swagger/index.html`.

### Endpoints Públicos

- `GET /sales/available`: Lista todos os veículos disponíveis para venda.
- `GET /sales/sold`: Lista todos os veículos já vendidos.
- `POST /sales/{id}/purchase`: Inicia o processo de compra para uma venda específica.
- `POST /webhooks/payments`: Recebe a notificação de status de pagamento.
