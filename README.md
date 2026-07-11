# OpenAPI Builder

A small Go library to generate OpenAPI 3 specification files (`.json` / `.yaml`) directly from Go code — no comments, no annotations, no code generation step.

## Purpose

This library was created to fulfill a simple need: produce a `.json` file in the OpenAPI standard so a Swagger UI can serve the API documentation.

Instead of maintaining a spec by hand or annotating handlers with magic comments, you describe your routes with plain Go structs. The library uses reflection to convert your request/response types into OpenAPI schemas and writes the final document to disk, ready to be served by any Swagger UI handler.

## Installation

```bash
go get github.com/LuigiVanin/openapi-builder@latest
```

```go
import oas "github.com/LuigiVanin/openapi-builder/openapi"
```

## Quick start

```go
package main

import (
	oas "github.com/LuigiVanin/openapi-builder/openapi"
)

type CreateUserBody struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserResponse struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	builder := oas.NewBuilder(
		"My API",              // title
		"My API description",  // description
		"1.0.0",               // version
	)

	builder.Add(
		builder.Route("POST", "/users", oas.Options{Summary: "Create a user"}).
			AddTag("Users").
			AddBody(CreateUserBody{}, oas.Options{Required: true}).
			AddResponse(201, UserResponse{}),
	)

	document := builder.Build()

	// Writes specs/index.json and specs/index.yaml
	err := document.Write()
	if err != nil {
		panic(err)
	}
}
```

Running this program produces `specs/index.json` and `specs/index.yaml` containing the OpenAPI 3 document.

## Usage

There are two ways to declare routes. Both can be mixed freely on the same builder.

### Fluent route builder

Start a route with `builder.Route(method, path, options...)`, chain the parts you need, and register it with `builder.Add(...)`:

```go
builder.Add(
	builder.Route("GET", "/users/{id}", oas.Options{
		Summary:     "Fetch a user",
		Description: "Fetches a single user by its id",
	}).
		AddTag("Users").
		AddPathParam("id", "string", oas.Options{Required: true}).
		AddQueryParam("expand", "boolean").
		AddHeaderParam("X-Request-Id", "string").
		AddResponse(200, UserResponse{}).
		AddResponse(404, ErrorResponse{}),
)
```

Available chain methods:

| Method | Description |
| --- | --- |
| `AddTag(tag ...string)` | Adds one or more tags to the route |
| `AddPathParam(name, type, opts...)` | Declares a path parameter |
| `AddQueryParam(name, type, opts...)` | Declares a query parameter |
| `AddHeaderParam(name, type, opts...)` | Declares a header parameter |
| `AddPathParams(struct)` / `AddQueryParams(struct)` / `AddHeaderParams(struct)` | Declares parameters from a struct via reflection |
| `AddBody(struct, opts...)` | Sets the request body schema from a struct |
| `AddResponse(statusCode, struct, opts...)` | Adds a response schema for a status code |

### Declarative route

Describe the whole route in a single struct with `builder.AddRoute(...)`. Request/response schemas are inferred from the struct values via reflection:

```go
type PathParams struct {
	Id string `json:"id"`
}

type QueryParams struct {
	Category string `json:"category"`
	Page     int    `json:"page"`
}

builder.AddRoute(oas.Route{
	Method:    "PUT",
	Path:      "/products/{id}",
	Tags:      []string{"Products"},
	Summary:   "Update a product",
	Parameter: PathParams{},   // path parameters
	Query:     QueryParams{},  // query parameters
	Body:      UpdateBody{},   // request body (application/json)
	Responses: map[string]any{ // response schemas by status code
		"200": ProductResponse{},
		"404": ErrorResponse{},
	},
})
```

### Writing the document

`document.Write()` accepts optional `WriteOptions` to control the output:

```go
document.Write(oas.WriteOptions{
	Formats:    []string{"json"}, // default: []string{"yaml", "json"}
	FolderPath: "docs",           // default: "specs"
	FileName:   "openapi",        // default: "index"
})
```

You can also get the raw bytes without touching the filesystem:

```go
output, err := document.Output("json") // or "yaml"
```

### Serving with Swagger UI

The generated file plugs into any Swagger UI server. For example, with [flowchartsman/swaggerui](https://github.com/flowchartsman/swaggerui):

```go
package main

import (
	"net/http"
	"os"

	"github.com/flowchartsman/swaggerui"
)

func main() {
	spec, _ := os.ReadFile("specs/index.json")

	http.Handle("/docs/", http.StripPrefix("/docs", swaggerui.Handler(spec)))
	http.ListenAndServe(":8080", nil)
}
```

Then open `http://localhost:8080/docs/` to browse the documentation.

## Roadmap

### Features
 - [x] Router Create Approach
   - [x] Add Parameters
     - [x] Add Path
     - [x] Add Query
     - [x] Add Header
   - [x] Add Body
   - [x] Add Response
   - [x] Metadata (via options)
 - [x] Router Builder Approach
   - [x] Add Parameters
     - [x] Add Path
     - [x] Add Query
     - [x] Add Header
   - [x] Add Body
   - [x] Add Response
   - [x] Metadata
 - [ ] Multiple Media Types - currently only `application/json`
 - [ ] Ref notation on each Schema
 - [?] More Options Document
 - [ ] Tests
   - [ ] Document tests
   - [ ] Spec Document Builder tests
   - [ ] Route builder tests
   - [ ] Utils test
 - [x] Library publication

### Improvements
 - [x] Json format always have empty `items: {}` field

## Design history

<details>
<summary>API concepts considered during development</summary>

### Concept 1 - Declarative

Too verbose and it will require a lot of boiler plate that is not infered by any LSP, it would make harder to use. Also, extend information its pretty much impossible, to add a description to a path it would require a new attribute on the Route struct (?), a description to a bodyRequest? even harder.

```go
builder.AddRoute(
  swg.Route{
    Method:    "POST/{id}",
    Path:      "/test",
    Parameter: Parameter{},
    Query:     Query{},
    Body: map[string] { "aplication/json": RequestBody{} }
    Responses: map[string]map[string]Response: {
      "200": {
        "application/json": Response{}
      },
    },
  },
)

doc := builder.Build()
doc.Write()
```

### Concept 2 (Discarded)

A better ideia of DX, but I dont think it is practical to implement it. The ideia is to expose a RouteBuiler struct via SwaggerBuilder and then construct the document Route per Route, the problem lives that this RouteBuilder would be "fake", it would have been implemented in the SwaggerBuilder it self so the Push method could work. Another caveat is the fact that if the user doesnt Push a Route before skipping to the next route building a error could happen.

```go
builder := sb.NewSwaggerBuilder()
builder.PushRoute("GET", "test/{id}", Option {}).
  AddParameters(Parameter{}).
  AddQuery(Query{}).
  AddBody(RequestBody{}, "aplication/json", Option { description: string }).
  AddResponse(fiber.OK, Response {}).
  Pop()

doc := builder.Build()
doc.Write()
```

### Concept 3 (Discarded)

```go
builder := sb.NewSwaggerBuilder()
route := builder.Route("GET", "test/{id}", Option {}).
  AddParameters(Parameter{}).
  AddQuery(Query{}).
  AddBody(RequestBody{}, "aplication/json", Option { description: string }).
  AddResponse(fiber.OK, Response {}, "aplication/json").
  Build()

builder.Add(route)

doc := builder.Build()
doc.Write()
```

### Concept 4 - Fluent / Builder

```go
builder := sb.NewSwaggerBuilder()

builder.Add(
  sb.NewRoute("GET", "test/{id}", Option {}).
    AddParameters(Parameter{}).
    AddQuery(Query{}).
    AddBody(RequestBody{}, "aplication/json", Option { description: string }).
    AddResponse(fiber.OK, Response {}).
    Build()
)

doc := builder.Build()
doc.Write()
```

</details>