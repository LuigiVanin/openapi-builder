package openapi

import (
	"maps"
	"reflect"
	"slices"
	"strings"
)

type Builder struct {
	document *Document /* `json:"document"` */
	route    *RouteBuilder
}

func NewBuilder(title string, description string, version string) *Builder {
	return &Builder{
		route: nil,

		document: &Document{
			Openapi: "3.0.4",

			Info: Info{
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
	Header    any
	Body      any

	// String is the response status
	// The format of the response will always be application/json
	Responses map[string]any
}

type Route = RoutePayload

func (this *Builder) AddRoute(payload RoutePayload) *Builder {
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

func (this Builder) CreateResponse(responses map[string]any) map[string]Response {
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

func (this Builder) CreateBody(body any) Body {

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

func (this *Builder) CreateParameters(payload RoutePayload) []Parameter {

	// Copy of path
	parameters := []Parameter{}

	for _, in := range []string{"path", "query", "header"} {

		// Check if there is a parameter or query to be iterated
		if (in == "path" && payload.Parameter == nil) || (in == "query" && payload.Query == nil) || (in == "header" && payload.Header == nil) {
			continue
		}

		parameter := payload.Parameter

		if in == "query" {
			parameter = payload.Query
		} else if in == "header" {
			parameter = payload.Header
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

func (this *Builder) Route(method string, path string, opt ...Options) *RouteBuilder {

	if this.route != nil {
		return this.route
	}

	builder := NewRouteBuilder(path, method, opt...)
	return builder
}

func (this *Builder) Add(route *RouteBuilder) {
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

func (this *Builder) Build() *Document {
	return this.document
}
