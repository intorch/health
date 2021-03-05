# Simple Health check Library

[![Codacy Badge](https://api.codacy.com/project/badge/Grade/92f3375a8d1045fcb631fb027a8468eb)](https://app.codacy.com/gh/intorch/health?utm_source=github.com&utm_medium=referral&utm_content=intorch/health&utm_campaign=Badge_Grade_Settings)

A simple health check library that starts a goroutine providing a endpoint to check the helth of application. The endpoint path and port are passed by configuration object.

## Usage

Use the code below to start health check goroutine.

```go
conf := &config.Configuration{
    Health: &config.Health{
        Addr:  ":3000",
        Route: "/health",
    },
}

health = New()
health.Start(conf)
```

To add any error to the health check object use the code below, that add http status 500.

```go
h.AddError(http.StatusInternalServerError, "Some Erro")
```
