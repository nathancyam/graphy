# Graphy, opinionated Gqlgen/Neo4j Go boilerplate

This repository contains boilerplate code for projects that plan on utilising Neo4j and Gqlgen. As implied, this project relies on the following packages:

- Gqlgen (https://gqlgen.com/getting-started/)
- Neo4j (https://neo4j.com/developer/go/)
- Wire (https://github.com/google/wire)

Wire is specifically used to generate the dependency tree to create the application server.

## Project structure

```
.
├── cmd                         // Directory containing main binary generation and implementation binding
│ └── graphy
│     ├── neo
│     │ └── neo4j.go            // Neo4j storage implementation bootup
│     ├── main.go               // Main go entrypoint
│     ├── inject                // Dependency injection via wire code auto generation.
│     │ ├── wire_gen.go         // Generated dependency injection
│     │ └── wire.go             // Wire definitions
│     └── config                // Configuration package, exposes variables via viper
│         └── config.go         // Initialise viper, sets bindings and sets package variables
│
├── transport                   // Transport implementation and details
│ ├── http                      // HTTP package used to bind the GraphQL handler to server implementation
│ │ ├── utils.go                // Common HTTP utilities (e.g. request ID generation)
│ │ └── server.go               // Application server implementation, creates *http.Server implementation
│
│ └── graphql                   // GraphQL package, utilising gqlgen. Does not expose HTTP bindings
│     ├── schema.resolvers.go   // Root resolver schema resolvers. Part of gqlgen generation process
│     ├── resolver.go           // General resolver struct, used to connect resolver implications to domain services
│     ├── model                 // GraphQL model representations.
│     │ └── models_gen.go       // Gqlgen model generation stubs
│     └── generated             // Gqlgen GraphQL server implementation generated files
│         └── generated.go
│
├── storage                     // Storage implementations bound to the domain interfaces
│ └── graph
│     ├── repository.go
│     └── errors.go
│
├── pkg                         // Core domain logic implementations
│ ├── rounds
│ │ ├── update.go
│ │ ├── repository.go           // Interfaces for repository methods only
│ │ └── models.go
│ └── core
│     ├── core_test.go
│     └── core.go
│
├── gqlgen.yml                  // Gqlgen configuration file, used for autogen
├── api                         // API resources and definitions
│ └── schema.graphqls           // GraphQL schema defintion, used for gqlgen
└── Makefile
```
