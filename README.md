
# Open Api Builder for Swagger

## Concepts

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

### Concept 2 (Discarted)

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

### Concept 3 (Discarted)

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

## TODO

### Features
 - [x] Router Create Approch 
   - [x] Add Parameters
     - [x] Add Path
     - [x] Add Query
     - [x] Add Header
   - [x] Add Body
   - [x] Add Response
   - [x] Metadata (via options)
 - [x] Router Builder Approch
   - [x] Add Parameters
     - [x] Add Path
     - [x] Add Query
     - [x] Add Header
   - [x] Add Body
   - [x] Add Response
   - [x] Metadata
 - [ ] Multiple Media Types - curretly only `application/json`
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




