# Wise API Go Client

Go client for the [Wise API](https://docs.wise.com/api-reference).

## Installation

```bash
go get github.com/joeblew999/investment/wise-api
```

## Usage

```go
package main

import (
    "context"
    "fmt"
    "os"

    wise "github.com/joeblew999/investment/wise-api"
)

func main() {
    // Create client (uses production by default)
    client := wise.NewClient(os.Getenv("WISE_API_TOKEN"))

    // Or use sandbox
    // client := wise.NewClient(os.Getenv("WISE_API_TOKEN"), wise.WithSandbox())

    ctx := context.Background()

    // Get profiles
    profiles, _ := client.Profiles.List(ctx)
    fmt.Printf("Profiles: %+v\n", profiles)

    // Get exchange rate
    rate, _ := client.ExchangeRates.Get(ctx, wise.USD, wise.EUR)
    fmt.Printf("USD/EUR: %f\n", rate.Rate)

    // Get balances
    balances, _ := client.Balances.List(ctx, profiles[0].ID, nil)
    fmt.Printf("Balances: %+v\n", balances)
}
```

## Services

- **Profiles** - List, get, create personal/business profiles
- **Quotes** - Create and manage transfer quotes
- **Recipients** - Manage recipient accounts
- **Transfers** - Create, list, cancel transfers
- **ExchangeRates** - Get live and historical exchange rates
- **Balances** - Multi-currency balance management

## Environment Variables

```bash
WISE_API_TOKEN=your-api-token
```

## API Reference

https://docs.wise.com/api-reference
