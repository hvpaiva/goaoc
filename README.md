# Go AOC (Advent of Code)

**Go AOC** is a Go library designed to simplify the process of running **Advent of Code** challenges. It streamlines 
input/output handling and lets you manage challenge execution with ease.

## Table of Contents
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Usage](#usage)
  - [Basic Example](#basic-example)
  - [Defining Custom Challenges](#defining-custom-challenges)
  - [Providing the Part Parameter](#providing-the-part-parameter)
  - [Configuration Options](#configuration-options)
  - [Clipboard Support](#clipboard-support)
- [IO Manager](#io-manager)
  - [Environment](#environment)
- [Error Handling](#error-handling)
- [Troubleshooting](#troubleshooting)
- [Contributing](#contributing)
- [License](#license)

## Installation

To install the library, you can use the following command:

```bash
go get -u github.com/hvpaiva/goaoc
```

## Quick Start

To quickly integrate Go AOC in your workflow, execute a simple challenge:

```go
package main

import (
   "log"

   "github.com/hvpaiva/goaoc"
)

func main() {
   err := goaoc.Run("yourInputData", partOne, partTwo, goaoc.WithPart(1))
   if err != nil {
      log.Fatalf("Run failed: %v", err)
   }
}

func partOne(input string) int {
   // Implement your algorithm for part one here
}

func partTwo(input string) int {
   // Implement your algorithm for part two here
}
```

## Usage

### Basic Example

Here's how to use Go AOC in a project:

```go
package main

import (
   "log"

   "github.com/hvpaiva/goaoc"
)

func main() {
   err := goaoc.Run("example input", partOne, partTwo)
   if err != nil {
      log.Fatalf("Error running challenge: %v", err)
   }
}

func partOne(input string) int {
   // Logic for part one
   return len(input)
}

func partTwo(input string) int {
   // Logic for part two (e.g., double length)
   return len(input) * 2
}
```

### Defining Custom Challenges

Challenge functions should receive a `string` input and return an `int`. Design purposes or parsing can be done within 
these functions.

### Providing the `part` Parameter

Multiple strategies exist for specifying the challenge part:

1. **Using a Flag**: You can pass the `--part` flag when running the challenge. Valid values are `1` or `2`.
   ```bash
   go run main.go --part=1
   ```

2. **Using an Environment Variable**: Set the `GOAOC_CHALLENGE_PART` environment variable to `1` or `2`.
   ```bash
   export GOAOC_CHALLENGE_PART=2
   go run main.go
   ```

3. **Through Standard Input**: If neither a flag nor an environment variable is provided, the program will prompt you to 
input the part number via the console.
   ```bash
   Which part do you want to run? (1/2)
   > 1
   ```

4. **Using a Function Parameter**: Directly specify the part by using the `goaoc.WithPart(part)` option when calling `goaoc.Run`.

```go
goaoc.Run(input, partOne, partTwo, goaoc.WithPart(1))
```

### Configuration Options

`goaoc.Run` supports configurations via options like:

- **WithPart(part challenge.Part)**: Specifies the part of the challenge to run (1 or 2).
- **WithManager(env io.Env)**: Sets up custom [IO Manager](#io-manager).

### Clipboard Support

Auto-copies results to clipboard—useful for quick submission.

> Disable using `GOAOC_DISABLE_COPY_CLIPBOARD=true`.

## IO Manager

Implement custom input/output handling using your own `IOManager`:

```go
type customManager struct {}

func (m *customManager) Read(arg string) (string, error) {
    // Custom input logic
}

func (m *customManager) Write(output string) error {
    // Custom output logic
}

customManager := &customManager{}
goaoc.Run(input, do, doAgain, goaoc.WithManager(customManager))
```

### Environment

Alter the default environment setting for `DefaultConsoleManager`:

```go
var customEnv = goaoc.Env{
   Stdin:  bytes.NewBufferString(""),
   Stdout: new(bytes.Buffer),
   Args:   []string{},
}

goaoc.Run(input, do, doAgain, goaoc.WithManager(goaoc.DefaultConsoleManager{Env: customEnv}))

```

## Error Handling

The `Run` function propagates errors for handling:

```go
if err := goaoc.Run(input, do, doAgain); err != nil {
	log.Fatal(err)
}
```

> Note: All errors in internal flow are returned in goaoc.Run functions. Except for copying to clipboard, which just
> logs the error, but does not break the execution. The errors are also all typed, so you can check the type of the error.

## Troubleshooting

If you encounter issues, consider:
- Checking file permissions for clipboard commands.
- Validating environment paths.
- Inspecting error messages for guidance.

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue on GitHub.

## License

This project is licensed under the Apache 2.0 License. See the [LICENSE](./LICENSE) file for details.
