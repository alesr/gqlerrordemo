# gqlerrordemo

## Overview
**gqlerrordemo** is a Go-based demonstration repository that illustrates handling GraphQL errors in a layered architecture. This project focuses on providing clear error handling within GraphQL operations.

## Structure
- **main.go**: Initializes the server and integrates the GraphQL setup with custom error handling.
- **gqlgen.yml**: Configuration for gqlgen, a tool for building GraphQL servers in Go.
- **schema.resolvers.go**: Implements resolvers for the GraphQL schema.
- **errors.go**: Custom GraphQL error definitions and handling logic.
- **resolver.go**: Root resolver setup for the GraphQL API, integrating services and error handling.
- **foo.go**: Demonstrates domain logic with the Foo service.
- **fakedb.go**: Simulates database functionality, aiding in error generation and handling.


## Features
- **GraphQL Server**: Utilizes gqlgen for server setup and operations.
- **Custom Error Handling**: Showcases creation and management of custom GraphQL errors.
- **Layered Architecture**: Demonstrates separation of concerns in service, resolver, and database simulation layers.
- **Error Presenter**: Customizes how errors are formatted and presented in GraphQL responses.

## Using gqlerrordemo
Start the server (go run main.go) and access the GraphQL playground (http://localhost:8080).

### Example Queries and Mutations

Create a Foo with Missing ID

```graphql
mutation {
  createFoo(id: "") {
    id
  }
}
```

Response:

```json
{
  "errors": [
    {
      "message": "Foo missing id",
      "path": ["createFoo"],
      "extensions": { "code": "FOO_MISSING_ID" }
    }
  ],
  "data": { "createFoo": null }
}
````

Create a Foo with an ID that Already Exists

```graphql
mutation {
  createFoo(id: "alreadyexists") {
    id
  }
}
````

Response:

```json
{
  "errors": [
    {
      "message": "Foo already exists",
      "path": ["createFoo"],
      "extensions": { "code": "FOO_ALREADY_EXISTS" }
    }
  ],
  "data": { "createFoo": null }
}
```

Successful Creation of a Foo

```graphql
mutation {
  createFoo(id: "123") {
    id
  }
}
```

Response:

```json
{
  "data": { "createFoo": { "id": "123" } }
}
```

Retrieving a Foo

```graphql
query {
  getFoo(id: "123") {
    id
  }
}
```

Response:

```json
{
  "data": { "getFoo": { "id": "foo" } }
}
```

Retrieving a Non-Existent Foo

```graphql
query {
  getFoo(id: "notfound") {
    id
  }
}
```

Response:

```json
{
  "errors": [
    {
      "message": "Foo not found",
      "path": ["getFoo"],
      "extensions": { "code": "FOO_NOT_FOUND" }
    }
  ],
  "data": { "getFoo": null }
}
```

Retrieving a Foo with Missing ID

```graphql
query {
  getFoo(id: "") {
    id
  }
}
```

Response:

```json
{
  "errors": [
    {
      "message": "Foo missing id",
      "path": ["getFoo"],
      "extensions": { "code": "FOO_MISSING_ID" }
    }
  ],
  "data": { "getFoo": null }
}
```

### Conclusion
**gqlerrordemo** provides a hands-on experience with effective error handling in a GraphQL server using a layered architecture in Go. The included examples in the GraphQL playground demonstrate the handling of various scenarios, offering insights into implementing similar structures in your own projects.

