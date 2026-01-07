[![project](https://img.shields.io/badge/github-psaraiva%2Fgo--horse--racing--by--cobra-blue)](https://img.shields.io/badge/github-psaraiva%2Fgo--horse--racing--by--cobra-blue)
[![License](https://img.shields.io/badge/license-MIT-%233DA639.svg)](https://opensource.org/licenses/MIT)

[![Go Report Card](https://goreportcard.com/badge/github.com/psaraiva/go-horse-racing-by-cobra)](https://goreportcard.com/report/github.com/psaraiva/go-horse-racing-by-cobra)
![Codecov](https://img.shields.io/codecov/c/github/psaraiva/go-horse-racing-by-cobra)

[![Idioma: English](https://img.shields.io/badge/Idioma-English-blue?style=flat-square)](README.md)

# Corrida de Cavalos por Cobra 🐎🐍

## 🎯 Objetivo
Este jogo demonstra o uso de Goroutines (concorrência) de uma forma simples, prática e divertida, utilizando a biblioteca Cobra.

## ⚙️ Como isso funciona?
Os cavalos correm até o primeiro cruzar a linha de chegada.

## 🚀 Início Rápido
```bash
# Clone o repositório
git clone https://github.com/psaraiva/go-horse-racing-by-cobra.git
cd go-horse-racing-by-cobra

# Execute com Docker
make docker-build
make docker-run
```

## Preview
![Preview](./asset/horse_race.gif)

## 🛠️ Makefile
O projeto inclui um Makefile com comandos úteis para desenvolvimento e execução:

### Comandos de Desenvolvimento
```bash
make help            # Exibe todos os comandos disponíveis
make test            # Executa todos os testes
make test-race       # Executa testes com detecção de race conditions
make test-coverage   # Executa testes com cobertura e gera relatório HTML
make build           # Compila o projeto
make clean           # Remove arquivos gerados
```

### Comandos Docker
```bash
make docker-build    # Constrói a imagem Docker
make docker-run      # Executa o container interativamente (configuração padrão)
make docker-stop     # Para e remove o container Docker
make docker-deploy   # Para container, constrói imagem e prepara para execução
make docker-clean    # Para container e remove a imagem Docker
make docker-rebuild  # Reconstrói a imagem do zero
```

Para executar com parâmetros customizados, use a variável ARGS:
```bash
# Com 5 cavalos e alvo de 50 pontos
make docker-run ARGS="--horses-quantity 5 --score-target 50"

# Com label 'C' e timeout de 15 segundos
make docker-run ARGS="--horse-label C --game-timeout 15s"

# Com 20 cavalos, label 'P', alvo 75 pontos e timeout de 90 segundos
make docker-run ARGS="--horses-quantity 20 --horse-label P --score-target 75 --game-timeout 90s"
```

## 🔧 Parâmetros
- `--horse-label`
  - valor padrão `H`
  - valor válido `char(1)`
- `--horses-quantity`
  - valor padrão `2`
  - valor válido `int 2..99`
- `--score-target`
  - valor padrão `75`
  - valor válido `int 15..100`
- `--game-timeout`
  - valor padrão `10s`
  - valor válido `string 10s..90s`

## Exemplo
```bash
   +---------|---------|---------|---------|---------|---------|---------|---------|--+
H01|................................................................................H01|
H02|........................................................................H02       |
H03|..............................................................................H03 |
H04|............................................................................H04   |
H05|...............................................................................H05|
H06|..............................................................................H06 |
H07|.............................................................................H07  |
H08|..............................................................................H08 |
H09|.........................................................................H09      |
H10|.........................................................................H10      |
   +---------|---------|---------|---------|---------|---------|---------|---------|--+
```

Mensagem de tempo esgotado
```bash
   +---------|---------|---------|---------|---------|---------|---------|---------|--+
H01|..................H01                                                             |
H02|...............H02                                                                |
H03|.....................H03                                                          |
   +---------|---------|---------|---------|---------|---------|---------|---------|--+

Today is a very hot day, the horses are tired!
```
