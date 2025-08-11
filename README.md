# Sistema CEP + Clima com OTEL e Zipkin

Este projeto é composto por **dois microsserviços** que trabalham juntos para fornecer informações climáticas a partir de um CEP informado.  
Além disso, o sistema implementa **observabilidade** com **OpenTelemetry** e visualização de traces via **Zipkin**.

---

## 📌 Arquitetura

- **Service A**: Recebe o CEP via HTTP POST, valida, consulta o Service B e retorna:
  - Cidade
  - Temperatura em Celsius, Fahrenheit e Kelvin

- **Service B**: Recebe um CEP, consulta a **Weather API** e retorna os dados de clima.

- **Zipkin**: Ferramenta de visualização de traces para acompanhar as chamadas distribuídas entre os serviços.

---

## 🚀 Tecnologias Utilizadas

- **Go** (>= 1.23)
- **Docker** & **Docker Compose**
- **OpenTelemetry**
- **Zipkin**
- **Weather API** (para dados climáticos)

---

## 📂 Estrutura do Projeto

```

.
├── service-a/       # Serviço que recebe CEP e retorna clima
├── service-b/       # Serviço que consulta a Weather API
├── docker-compose.yml
├── .env
└── README.md
```

## ⚙️ Configuração do Ambiente

1. **Clonar o repositório**
   ```bash
   git clone https://github.com/seu-usuario/sistema-cep-clima.git
   cd sistema-cep-clima```

2. **Subir os serviços com Docker**

   ```bash
   docker compose up --build
   ```

   Isso irá:

   * Compilar e iniciar o **Service A** na porta **8080**
   * Compilar e iniciar o **Service B** na porta **8081**
   * Iniciar o **Zipkin** na porta **9411**

---

## 📬 Como Utilizar

### 1. Enviar uma requisição para o **Service A**

```bash
curl -X POST http://localhost:8080/cep \
  -H "Content-Type: application/json" \
  -d '{"cep":"01001000"}'
```

📌 **Resposta esperada:**

```json
{
  "city": "São Paulo",
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.65
}
```

---

### 2. Testar o **Service B** diretamente (opcional)

```bash
curl -X POST http://localhost:8081/weather \
  -H "Content-Type: application/json" \
  -d '{"cep":"01001000"}'
```

---

## 🔍 Observabilidade com Zipkin

1. Com os serviços rodando, acesse:

   ```
   http://localhost:9411
   ```

2. Clique em **"Run Query"** para listar as traces.

3. Selecione uma trace para visualizar o fluxo completo da requisição:

Isso permite verificar:

* Latência total da requisição
* Tempo gasto em cada serviço
* Possíveis gargalos

---

## 📄 Licença

Este projeto é distribuído sob a licença MIT. Consulte o arquivo `LICENSE` para mais detalhes.

```

