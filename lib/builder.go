package swagger_builder

import (
	"reflect"
	"slices"
	"strings"
)

type SwaggerBuilder struct {
	document *SwaggerDocument /* `json:"document"` */
}

func NewSwaggerBuilder(title string, description string, version string) *SwaggerBuilder {
	return &SwaggerBuilder{
		document: &SwaggerDocument{
			Openapi: "3.0.4",

			Info: SwaggerInfo{
				Title:        title,
				Descriptiotn: description,
				Version:      version,
			},
		},
	}
}

type DocumentRoutePayload struct {
	Method      string   `json:"method"`
	Path        string   `json:"path"`
	Tags        []string `json:"tags"`
	Summary     string   `json:"summary"`
	Description string   `json:"description"`

	Parameter any
	Query     any
	Body      any
}

type Route = DocumentRoutePayload

func (this *SwaggerBuilder) AddRoute(payload DocumentRoutePayload) *SwaggerBuilder {
	method := strings.ToLower(payload.Method)
	path := FormatRoutePath(payload.Path)

	// validation step
	if !slices.Contains([]string{"get", "put", "delete", "options", "patch", "post"}, method) || payload.Path == "" {
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
	}

	return this
}

func (this SwaggerBuilder) CreateBody(body any) Body {

	if body == nil {
		return Body{}
	}

	t := reflect.TypeOf(body)

	schema := StructToSchema(t)

	return Body{
		Description: "",
		Required:    true,
		Content:     map[string]MediaTypeObject{"application/json": {Schema: schema}},
	}
}

func (this *SwaggerBuilder) CreateParameters(payload DocumentRoutePayload) []Parameter {

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

		if t.Kind() == reflect.Pointer {
			t = t.Elem()
		}

		if t.Kind() != reflect.Struct {
			continue
		}

		for index := range t.NumField() {
			field := t.Field(index)

			name := field.Tag.Get("json")

			if name == "" {
				name = field.Name
			}

			parameter := Parameter{
				In:       in,
				Name:     name,
				Required: true,
				Schema: map[string]any{
					"type": TypeToSwagger(field.Type.Kind()),
				},
			}
			parameters = append(parameters, parameter)
		}

	}

	return parameters
}

func (this *SwaggerBuilder) Build() *SwaggerDocument {
	return this.document
}
