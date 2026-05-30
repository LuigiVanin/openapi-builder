package swagger_builder

import (
	"reflect"
	"slices"
	"strconv"
	"strings"
)

type Options struct {
	Summary     string
	Description string
	Required    bool
	Format      string
	MediaType   string
	Tags        []string
}

type RouteBuilder struct {
	path   string
	method string

	summary     string
	description string
	tags        []string

	pathParameters  map[string]Parameter
	queryParameters map[string]Parameter

	body      Body
	responses map[string]Response
}

func NewRouteBuilder(path string, method string, options ...Options) *RouteBuilder {
	method = strings.ToLower(method)
	path = FormatRoutePath(path)
	httpMethods := []string{"get", "put", "delete", "options", "patch", "post"}

	opt := Options{}
	if len(options) > 0 {
		opt = options[0]
	}

	// validation step
	if !slices.Contains(httpMethods, method) || path == "" {
		method = "get"
	}

	return &RouteBuilder{
		path:   path,
		method: method,

		summary:     opt.Summary,
		description: opt.Description,

		pathParameters:  map[string]Parameter{},
		queryParameters: map[string]Parameter{},
		body: Body{
			Content: map[string]MediaTypeObject{},
		},
		responses: map[string]Response{},
		tags:      opt.Tags,
	}
}

func (this *RouteBuilder) AddTag(tag ...string) *RouteBuilder {
	this.tags = append(this.tags, tag...)

	return this
}

func (this *RouteBuilder) AddPathParam(name string, typename string, options ...Options) *RouteBuilder {
	opt := Options{}
	if len(options) > 0 {
		opt = options[0]
	}

	this.pathParameters[name] = Parameter{
		Name: name,
		Schema: Schema{
			Type:   typename,
			Format: opt.Format,
		},
		In:          "path",
		Required:    opt.Required,
		Description: opt.Description,
	}

	return this
}

func (this *RouteBuilder) AddQueryParam(name string, typename string, options ...Options) *RouteBuilder {
	opt := Options{}
	if len(options) > 0 {
		opt = options[0]
	}

	this.queryParameters[name] = Parameter{
		Name: name,
		Schema: Schema{
			Type:   typename,
			Format: opt.Format,
		},
		In:          "query",
		Required:    opt.Required,
		Description: opt.Description,
	}

	return this
}

func (this *RouteBuilder) AddQueryParams(payload any) *RouteBuilder {
	t := reflect.TypeOf(payload)

	parameters := TypeToParam(t)

	for _, param := range parameters {
		param.In = "query"
		this.pathParameters[param.Name] = param
	}

	return this
}

func (this *RouteBuilder) AddPathParams(payload any) *RouteBuilder {

	t := reflect.TypeOf(payload)

	parameters := TypeToParam(t)

	for _, param := range parameters {
		param.In = "path"
		this.pathParameters[param.Name] = param
	}

	return this
}

func (this *RouteBuilder) mergePathQuery() []Parameter {
	params := []Parameter{}

	for _, param := range this.pathParameters {
		params = append(params, param)
	}

	for _, query := range this.queryParameters {
		params = append(params, query)
	}

	return params
}

func (this *RouteBuilder) AddBody(payload any, options ...Options) *RouteBuilder {
	opt := Options{}
	if len(options) > 0 {
		opt = options[0]
	}

	mediaType := "application/json"

	if opt.MediaType != "" {
		mediaType = opt.MediaType
	}

	if opt.Description != "" {
		this.body.Description = opt.Description
	}

	if opt.Required {
		this.body.Required = opt.Required
	}

	t := reflect.TypeOf(payload)

	if this.body.Content == nil {
		this.body.Content = map[string]MediaTypeObject{}
	}

	// The Content map is defined as a value and  not as pointer, with that definitio we cannot change the
	// value in a map this.body.Content[mediaType].Schema = xxxx, because it is a fixed value not a pointer, so
	// We need to copy it, change the field we want and then substitute the whole field value
	mt := this.body.Content[mediaType]
	mt.Schema = StructToSchema(t)
	this.body.Content[mediaType] = mt

	return this
}

func (this *RouteBuilder) AddResponse(statusCode int, payload any, options ...Options) *RouteBuilder {
	opt := Options{}
	if len(options) > 0 {
		opt = options[0]
	}

	mediaType := "application/json"
	if opt.MediaType != "" {
		mediaType = opt.MediaType
	}

	statusKey := strconv.Itoa(statusCode)
	resp := this.responses[statusKey]

	if resp.Content == nil {
		resp.Content = map[string]MediaTypeObject{}
	}

	if opt.Description != "" {
		resp.Description = opt.Description
	}

	if payload != nil {
		t := reflect.TypeOf(payload)
		resp.Content[mediaType] = MediaTypeObject{
			Schema: StructToSchema(t),
		}
	}

	this.responses[statusKey] = resp

	return this
}

func (this *RouteBuilder) Build() map[string]map[string]Path {
	route := Path{}

	params := this.mergePathQuery()

	route.Parameters = params
	route.RequestBody = this.body
	route.Responses = this.responses
	route.Description = this.description
	route.Summary = this.summary
	route.Tags = this.tags
	return map[string]map[string]Path{
		this.path: {
			this.method: route,
		},
	}
}
