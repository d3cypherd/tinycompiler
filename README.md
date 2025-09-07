<h1 align=center>
    <img src='imgs/logo.png'>
</h1>


<p align="center">
  <i align="center">A lightweight recursive descent compiler implemented in Go.</i>
</p>

![GIF demo](imgs/preview-clip.gif)

> _Built with [Fyne](fyne.io) for a cross-platform GUI._

## Features

- **Lexical analysis** – Scans and tokenizes source code.
- **Parsing** – Builds an abstract syntax tree (AST) using recursive descent.
- **Interactive editor** – Simple UI to write and test code.
- **Error reporting** – Prints errors with line and word numbers.

## Usage
1. **Clone the repository**
   ```bash
   git clone https://github.com/yourusername/tinycompiler.git
   cd tinycompiler
    ```

#### Dependencies
This project uses Go modules for dependency management.  
Required: **Go 1.22+** (or latest stable release).  

All external libraries are listed in [`go.mod`](./go.mod) and installed automatically when you build or run:

2. **Run the application**
```bash
go run .
```

3. **Write and compile code**
- Use the **built-in editor** to write your program, or click **Upload File** to import a source file.
- Click **Scan** to tokenize the code, then **Parse** to generate the syntax tree.
- Any errors will be displayed with their line and word numbers.

## Build
To build a standalone executable, run:

```bash
go build -o tinycompiler
```

This will create a binary named `tinycompiler` in the current directory.

#### Cross-compiling

Go makes it easy to build for other platforms. For example:

- **Linux (64-bit):**
```bash
GOOS=linux GOARCH=amd64 go build -o tinycompiler-linux
```

- **Windows (64-bit):**
```bash
GOOS=windows GOARCH=amd64 go build -o tinycompiler.exe
```

- **macOS (Applic Silicon):**
```bash
GOOS=darwin GOARCH=arm64 go build -o tinycompiler-macos
```

## TINY Language Tokens
| TokenType      | Value / Example |
|----------------|-----------------|
| SEMICOLON      | `;`             |
| IF             | `if`            |
| THEN           | `then`          |
| END            | `end`           |
| REPEAT         | `repeat`        |
| UNTIL          | `until`         |
| IDENTIFIER     | `x`, `abc`, `xyz` |
| ASSIGN         | `:=`            |
| READ           | `read`          |
| WRITE          | `write`         |
| LESSTHAN       | `<`             |
| EQUAL          | `=`             |
| PLUS           | `+`             |
| MINUS          | `-`             |
| MULT           | `*`             |
| DIV            | `/`             |
| OPENBRACKET    | `(`             |
| CLOSEDBRACKET  | `)`             |
| NUMBER         | `12`, `289`     |
