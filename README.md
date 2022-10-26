# Rigid ENV

A no-dependency Go package that allows safe .env files loading.

## Installation

```sh
go get github.com/terawatthour/rigid-env
```

## Usage

```go
import (
    env "github.com/terawatthour/rigid-env"
)

type Environment struct {
	ENV         string
	TARGET_PROD string
	TARGET_DEV  string
	VERSION     int
}

func main() {
    var environment Environment

    if err := env.Load(&environment, ".env"); err != nil {
        panic(err)
    }

    // here you have access to the loaded values, last file in order takes precedence
    fmt.Println(environment.VERSION)
}
```

## Compatibility

The parser will handle:

-   multiline comments
-   string interpolation
-   blank lines
-   comments that start in the beginning of the line (with hash)
-   type casting

The parser won't handle:

-   comments after assignment
-   non-standard variable names

## License

© TeraWattHour 2022-  
Published under the MIT License
