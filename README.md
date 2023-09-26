## Corpus Reader

### Task Description

Please implement a program that scans through a provided repository and counts how many times a specific word appears in it. The challenge consists of two parts:

#### Part 1 - Command Line Interface (CLI)

Implement a CLI program in Go that accepts a path to a corpus repository and the word that we are interested in:

**Example of Execution:**
```shell
counter --dir ./corpus --word john
```

**Output:**
```
count: 58
```

#### Part 2 - HTTP Server

Expose the functionality from part 1 via an HTTP server.

**Example of Execution:**
```shell
curl -X POST -H "Content-Type: application/json" -d '{"directory":"corpus", "word":"john"}' http://localhost:8080/counter
```

**Output:**
```json
{ "count": 58 }
```

### Remarks

- The solution can be as simple or as complex as you want.
- Aim for GO idiomatic code.
- Prioritize code maintainability.
- Faster is better.
- Include tests, documentation, and benchmarks if necessary.
- Follow Git best practices.
- You are welcome to enhance the challenge in any manner you prefer.