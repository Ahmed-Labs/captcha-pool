# Captcha Pool

A tool that asynchronously solves and stores captchas for immediate retrieval when needed. 
It continuously regenerates captchas in the backgorund and discards expired captchas given a TTL (varies depending on captcha type).

## Installation

```bash
go get github.com/Ahmed-Labs/captcha-pool
```

## Usage

```go
    solve := func() (string, error) {
        // Some third party captcha solver
        return "token", err
    }

    pool = captchapool.New(solve,
        &captchapool.Options{
            Count:           5,
            Refresh:         true,
            RefreshInterval: time.Second * 40,
            TTL:             time.Second * 60,
        },
    )

    pool.Start()
    defer pool.Stop()

    pool.GetToken()
```

See `./example` for more detailed usage.