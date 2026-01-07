[![project](https://img.shields.io/badge/github-psaraiva%2Fgo--horse--racing--by--cobra-blue)](https://img.shields.io/badge/github-psaraiva%2Fgo--horse--racing--by--cobra-blue)
[![License](https://img.shields.io/badge/license-MIT-%233DA639.svg)](https://opensource.org/licenses/MIT)

[![Go Report Card](https://goreportcard.com/badge/github.com/psaraiva/go-horse-racing-by-cobra)](https://goreportcard.com/report/github.com/psaraiva/go-horse-racing-by-cobra)
![Codecov](https://img.shields.io/codecov/c/github/psaraiva/go-horse-racing-by-cobra)

[![Idioma: Português](https://img.shields.io/badge/Idioma-Português-green?style=flat-square)](README_pt_br.md)

# 🐎 Horse Racing by Cobra 🐍

## 🎯 Objective
This game aims to demonstrate the use of Goroutines in a simple, practical, and fun way using Cobra.

## ⚙️ How it works?
The horses run until the first one crosses the finish line.

## � Quick Start
```bash
# Clone the repository
git clone https://github.com/psaraiva/go-horse-racing-by-cobra.git
cd go-horse-racing-by-cobra

# Run with Docker
make docker-build
make docker-run
```

## Preview
![Preview](./asset/horse_race.gif)

## 🛠️ Makefile
The project includes a Makefile with useful commands for development and execution:

### Development Commands
```bash
make help            # Display all available commands
make test            # Run all tests
make test-race       # Run tests with race condition detection
make test-coverage   # Run tests with coverage and generate HTML report
make build           # Build the project
make clean           # Remove generated files
```

### Docker Commands
```bash
make docker-build    # Build Docker image
make docker-run      # Run container interactively (default configuration)
make docker-stop     # Stop and remove Docker container
make docker-deploy   # Stop container, build image and prepare for execution
make docker-clean    # Stop container and remove Docker image
make docker-rebuild  # Rebuild Docker image from scratch
```

To run with custom parameters, use the ARGS variable:
```bash
# With 5 horses and score target of 50
make docker-run ARGS="--horses-quantity 5 --score-target 50"

# With label 'C' and timeout of 15 seconds
make docker-run ARGS="--horse-label C --game-timeout 15s"

# With 20 horses, label 'P', target 75 points and timeout of 90 seconds
make docker-run ARGS="--horses-quantity 20 --horse-label P --score-target 75 --game-timeout 90s"
```

## 🔧 Parameters
- `--horse-label`
  - default value `H`
  - valid value `char(1)`
- `--horses-quantity`
  - default value `2`
  - valid value `int 2..99`
- `--score-target`
  - default value `75`
  - valid value `int 15..100`
- `--game-timeout`
  - default value `10s`
  - valid value `string 10s..90s`

## Example
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

Timeout message
```bash
   +---------|---------|---------|---------|---------|---------|---------|---------|--+
H01|..................H01                                                             |
H02|...............H02                                                                |
H03|.....................H03                                                          |
   +---------|---------|---------|---------|---------|---------|---------|---------|--+

Today is a very hot day, the horses are tired!
```
