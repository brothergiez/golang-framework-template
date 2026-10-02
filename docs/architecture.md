# Architecture

Dependency direction:

```text
cmd -> app/bootstrap -> interfaces + application
application -> domain interfaces
infrastructure implements domain interfaces
```

Service types share application/domain. Only transport adapters differ.
