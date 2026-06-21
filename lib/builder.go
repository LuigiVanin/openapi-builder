package swagger_builder

import (
	"maps"
	"reflect"
	"slices"
	"strings"
)

type SwaggerBuilder struct {
	document *SwaggerDocument /* `json:"document"` */
	route    *RouteBuilder
}

func NewSwaggerBuilder(title string, description string, version string) *SwaggerBuilder {
	return &SwaggerBuilder{
		route: nil,

		document: &SwaggerDocument{
			Openapi: "3.0.4",

			Info: SwaggerInfo{
				Title:       title,
				Description: description,
				Version:     version,
			},

			Paths: map[string]map[string]Path{},
		},
	}
}

type RoutePayload struct {
	Method      string   `json:"method"`
	Path        string   `json:"path"`
	Tags        []string `json:"tags"`
	Summary     string   `json:"summary"`
	Description string   `json:"description"`

	Parameter any
	Query     any
	Body      any

	// String is the response status
	// The format of the response will always be application/json
	Responses map[string]any
}

type Route = RoutePayload

func (this *SwaggerBuilder) AddRoute(payload RoutePayload) *SwaggerBuilder {
	method := strings.ToLower(payload.Method)
	path := FormatRoutePath(payload.Path)
	httpMethods := []string{"get", "put", "delete", "options", "patch", "post"}

	// validation step
	if !slices.Contains(httpMethods, method) || payload.Path == "" {
		return this
	}

	if this.document.Paths == nil {
		this.document.Paths = map[string]map[string]Path{}
	}

	if this.document.Paths[path] == nil {
		this.document.Paths[path] = map[string]Path{}
	}

	this.document.Paths[path][method] = Path{
		Summary:     payload.Summary,
		Description: payload.Description,
		Tags:        payload.Tags,
		Parameters:  this.CreateParameters(payload),
		RequestBody: this.CreateBody(payload.Body),
		Responses:   this.CreateResponse(payload.Responses),
	}

	return this
}

func (this SwaggerBuilder) CreateResponse(responses map[string]any) map[string]Response {
	r := map[string]Response{}

	for key, value := range maps.All(responses) {
		content := map[string]MediaTypeObject{}
		t := reflect.TypeOf(value)

		content["application/json"] = MediaTypeObject{
			Schema: TypeToSchema(t),
		}

		r[key] = Response{
			Content: content,
		}
	}

	return r
}

func (this SwaggerBuilder) CreateBody(body any) Body {

	if body == nil {
		return Body{}
	}

	t := reflect.TypeOf(body)

	schema := TypeToSchema(t)

	return Body{
		Description: "",
		Required:    true,
		Content:     map[string]MediaTypeObject{"application/json": {Schema: schema}},
	}
}

func (this *SwaggerBuilder) CreateParameters(payload RoutePayload) []Parameter {

	// Copy of path
	parameters := []Parameter{}

	for _, in := range []string{"path", "query"} {

		// Check if there is a parameter or query to be iterated
		if (in == "path" && payload.Parameter == nil) || (in == "query" && payload.Query == nil) {
			continue
		}

		parameter := payload.Parameter

		if in == "query" {
			parameter = payload.Query
		}

		t := reflect.TypeOf(parameter)

		tempParams := TypeToParam(t)
		for index := range tempParams {
			tempParams[index].In = in
		}

		parameters = append(parameters, tempParams...)

	}

	return parameters
}

func (this *SwaggerBuilder) Route(method string, path string, opt ...Options) *RouteBuilder {

	if this.route != nil {
		return this.route
	}

	builder := NewRouteBuilder(path, method, opt...)
	return builder
}

func (this *SwaggerBuilder) Add(route *RouteBuilder) {
	path := route.Build()

	for key, p := range path {
		if this.document.Paths[key] == nil {
			this.document.Paths[key] = map[string]Path{}
		}

		for method, r := range p {
			this.document.Paths[key][method] = r
		}
	}

}

func (this *SwaggerBuilder) Build() *SwaggerDocument {
	return this.document
}
