package lib_test

import (
	"encoding/json"
	"maps"
	"reflect"
	"slices"
	"testing"
	"time"

	openapi "github.com/LuigiVanin/openapi-builder/openapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// TaggedPayload cobre cada forma de tag json e cada tipo com tratamento próprio.
type TaggedPayload struct {
	Name     string `json:"name"`
	Optional string `json:"optional,omitempty"`
	Hidden   string `json:"-"`
	Dash     string `json:"-,"`
	NoTag    string
	When     time.Time       `json:"when"`
	Raw      json.RawMessage `json:"raw"`
	Bytes    []byte          `json:"bytes"`
	Child    *TaggedPayload  `json:"child"`

	internal string //nolint:unused // existe só para provar que campo não exportado é ignorado
}

// NodeA e NodeB se referenciam mutuamente: sem a proteção de ciclo a recursão
// não termina.
type NodeA struct {
	Name string `json:"name"`
	B    *NodeB `json:"b"`
}

type NodeB struct {
	Name string `json:"name"`
	A    *NodeA `json:"a"`
}

// Promoted é embedado sem tag: os campos dele sobem um nível.
type Promoted struct {
	Inherited string `json:"inherited"`
	Shadowed  string `json:"shadowed"`
}

// AmbiguousLeft e AmbiguousRight trazem o mesmo nome na mesma profundidade e
// nenhum tem tag: o encoding/json descarta os dois.
type AmbiguousLeft struct {
	Ambiguous string
}

type AmbiguousRight struct {
	Ambiguous string
}

// TaggedWinner e UntaggedLoser também empatam em profundidade, mas só um tem
// nome na tag — esse vence. O tipo diferente é o que denuncia qual sobreviveu.
type TaggedWinner struct {
	Value int `json:"Winner"`
}

type UntaggedLoser struct {
	Winner string
}

// unexportedEmbed não é exportado, mas os campos exportados de dentro dele
// continuam sendo serializados.
type unexportedEmbed struct {
	Visible string `json:"visible"`
}

type Slug string

type EmbedHost struct {
	Promoted
	AmbiguousLeft
	AmbiguousRight

	Shadowed string `json:"shadowed"`
	Own      string `json:"own"`
}

type TieBreakHost struct {
	TaggedWinner
	UntaggedLoser
}

type PointerEmbedHost struct {
	*Promoted

	Own string `json:"own"`
}

type TaggedEmbedHost struct {
	Promoted `json:"promoted"`

	Own string `json:"own"`
}

type HiddenEmbedHost struct {
	Promoted `json:"-"`

	Own string `json:"own"`
}

type ScalarEmbedHost struct {
	Slug

	Own string `json:"own"`
}

type UnexportedEmbedHost struct {
	unexportedEmbed

	Own string `json:"own"`
}

type UtilsTestSuite struct {
	suite.Suite
}

func (this *UtilsTestSuite) SetupTest() {

}

func TestUtilsTestSuite(t *testing.T) {
	suite.Run(t, new(UtilsTestSuite))
}

func (this *UtilsTestSuite) TestTypeToSwagger_Success() {
	assert.Equal(this.T(), "string", openapi.TypeToSwagger(reflect.String))
	assert.Equal(this.T(), "integer", openapi.TypeToSwagger(reflect.Int))
	assert.Equal(this.T(), "integer", openapi.TypeToSwagger(reflect.Int64))
	assert.Equal(this.T(), "boolean", openapi.TypeToSwagger(reflect.Bool))
	assert.Equal(this.T(), "number", openapi.TypeToSwagger(reflect.Float32))
	assert.Equal(this.T(), "number", openapi.TypeToSwagger(reflect.Float64))
	assert.Equal(this.T(), "object", openapi.TypeToSwagger(reflect.Struct))
	assert.Equal(this.T(), "object", openapi.TypeToSwagger(reflect.Map))
	assert.Equal(this.T(), "array", openapi.TypeToSwagger(reflect.Slice))
	assert.Equal(this.T(), "array", openapi.TypeToSwagger(reflect.Array))
}

func (this *UtilsTestSuite) TestTypeToSwaggerDefault_Success() {
	assert.Equal(this.T(), "string", openapi.TypeToSwagger(reflect.Chan))
}

func (this *UtilsTestSuite) TestTypeToSchemaPrimitive_Success() {
	schema := openapi.TypeToSchema(reflect.TypeOf("some string"))

	assert.Equal(this.T(), "string", schema.Type)
	assert.Empty(this.T(), schema.Properties)
}

func (this *UtilsTestSuite) TestTypeToSchemaStruct_Success() {
	schema := openapi.TypeToSchema(reflect.TypeOf(UserPayload{}))

	assert.Equal(this.T(), "object", schema.Type)
	assert.Contains(this.T(), schema.Properties, "id")
	assert.Contains(this.T(), schema.Properties, "name")
	assert.Equal(this.T(), "string", schema.Properties["id"].Type)
}

func (this *UtilsTestSuite) TestTypeToSchemaSlice_Success() {
	schema := openapi.TypeToSchema(reflect.TypeOf([]UserPayload{}))

	assert.Equal(this.T(), "array", schema.Type)
	assert.Equal(this.T(), "object", schema.Items.Type)
	assert.Contains(this.T(), schema.Items.Properties, "id")
	assert.Contains(this.T(), schema.Items.Properties, "name")
}

func (this *UtilsTestSuite) TestTypeToSchemaPointer_Success() {
	schema := openapi.TypeToSchema(reflect.TypeOf(&UserPayload{}))

	assert.Equal(this.T(), "object", schema.Type)
	assert.Contains(this.T(), schema.Properties, "id")
}

func (this *UtilsTestSuite) TestTypeToParam_Success() {
	params := openapi.TypeToParam(reflect.TypeOf(Query{}))

	assert.NotEmpty(this.T(), params)

	names := map[string]bool{}
	for _, param := range params {
		names[param.Name] = true
		assert.True(this.T(), param.Required)
	}

	assert.True(this.T(), names["category"])
	assert.True(this.T(), names["Jump"])
}

func (this *UtilsTestSuite) TestTypeToParamNonStruct_Success() {
	params := openapi.TypeToParam(reflect.TypeOf(""))

	assert.Empty(this.T(), params)
}

func (this *UtilsTestSuite) TestFormatRoutePath_Success() {
	assert.Equal(this.T(), "/test", openapi.FormatRoutePath("test"))
	assert.Equal(this.T(), "/test", openapi.FormatRoutePath("/test"))
	assert.Equal(this.T(), "/test/id", openapi.FormatRoutePath("test/id"))
	assert.Equal(this.T(), "/test", openapi.FormatRoutePath("/test/"))
}

func (this *UtilsTestSuite) TestMerge_Success() {
	base := openapi.WriteOptions{
		Formats:    []string{"yaml", "json"},
		FolderPath: "specs",
		FileName:   "index",
	}
	override := openapi.WriteOptions{
		FileName: "custom",
	}

	result := openapi.Merge(base, override)

	assert.Equal(this.T(), "custom", result.FileName)
	assert.Equal(this.T(), "specs", result.FolderPath)
	assert.Equal(this.T(), []string{"yaml", "json"}, result.Formats)
}

func (this *UtilsTestSuite) TestTagWithOptionIsTrimmed_Success() {
	schema := openapi.TypeToSchema(reflect.TypeOf(TaggedPayload{}))

	assert.Contains(this.T(), schema.Properties, "optional")
	assert.NotContains(this.T(), schema.Properties, "optional,omitempty")
	assert.Equal(this.T(), "string", schema.Properties["optional"].Type)
}

func (this *UtilsTestSuite) TestDashTagIsSkipped_Success() {
	schema := openapi.TypeToSchema(reflect.TypeOf(TaggedPayload{}))

	assert.NotContains(this.T(), schema.Properties, "Hidden")
	assert.NotContains(this.T(), schema.Properties, "hidden")
}

func (this *UtilsTestSuite) TestDashCommaIsAFieldNamedDash_Success() {
	schema := openapi.TypeToSchema(reflect.TypeOf(TaggedPayload{}))

	assert.Contains(this.T(), schema.Properties, "-")
	assert.Equal(this.T(), "string", schema.Properties["-"].Type)
}

func (this *UtilsTestSuite) TestUntaggedFieldKeepsGoName_Success() {
	schema := openapi.TypeToSchema(reflect.TypeOf(TaggedPayload{}))

	assert.Contains(this.T(), schema.Properties, "NoTag")
}

func (this *UtilsTestSuite) TestUnexportedFieldIsSkipped_Success() {
	schema := openapi.TypeToSchema(reflect.TypeOf(TaggedPayload{}))

	assert.NotContains(this.T(), schema.Properties, "internal")
}

func (this *UtilsTestSuite) TestTimeIsDateTimeString_Success() {
	schema := openapi.TypeToSchema(reflect.TypeOf(TaggedPayload{}))
	when := schema.Properties["when"]

	assert.Equal(this.T(), "string", when.Type)
	assert.Equal(this.T(), "date-time", when.Format)
	assert.Empty(this.T(), when.Properties)
}

func (this *UtilsTestSuite) TestTimePointerIsDateTimeString_Success() {
	schema := openapi.TypeToSchema(reflect.TypeOf(&time.Time{}))

	assert.Equal(this.T(), "string", schema.Type)
	assert.Equal(this.T(), "date-time", schema.Format)
}

func (this *UtilsTestSuite) TestRawMessageIsObject_Success() {
	schema := openapi.TypeToSchema(reflect.TypeOf(TaggedPayload{}))
	raw := schema.Properties["raw"]

	assert.Equal(this.T(), "object", raw.Type)
	assert.Empty(this.T(), raw.Items.Type)
}

func (this *UtilsTestSuite) TestByteSliceIsBase64String_Success() {
	schema := openapi.TypeToSchema(reflect.TypeOf(TaggedPayload{}))
	bytes := schema.Properties["bytes"]

	assert.Equal(this.T(), "string", bytes.Type)
	assert.Equal(this.T(), "byte", bytes.Format)
}

func (this *UtilsTestSuite) TestPointerIsRecomposed_Success() {
	schema := openapi.TypeToSchema(reflect.TypeOf(TaggedPayload{}))
	child := schema.Properties["child"]

	assert.Equal(this.T(), "object", child.Type)
}

func (this *UtilsTestSuite) TestSelfReferenceDoesNotRecurse_Success() {
	schema := openapi.TypeToSchema(reflect.TypeOf(TaggedPayload{}))
	child := schema.Properties["child"]

	// O ciclo é cortado: o filho vira um objeto sem propriedades.
	assert.Equal(this.T(), "object", child.Type)
	assert.Empty(this.T(), child.Properties)
}

func (this *UtilsTestSuite) TestMutualReferenceDoesNotRecurse_Success() {
	schema := openapi.TypeToSchema(reflect.TypeOf(NodeA{}))

	b := schema.Properties["b"]
	assert.Equal(this.T(), "object", b.Type)
	assert.Contains(this.T(), b.Properties, "name")

	a := b.Properties["a"]
	assert.Equal(this.T(), "object", a.Type)
	assert.Empty(this.T(), a.Properties)
}

func (this *UtilsTestSuite) TestSiblingBranchesAreNotCutByTheCycleGuard_Success() {
	type Leaf struct {
		Value string `json:"value"`
	}

	type Root struct {
		Left  Leaf `json:"left"`
		Right Leaf `json:"right"`
	}

	schema := openapi.TypeToSchema(reflect.TypeOf(Root{}))

	assert.Contains(this.T(), schema.Properties["left"].Properties, "value")
	assert.Contains(this.T(), schema.Properties["right"].Properties, "value")
}

func (this *UtilsTestSuite) TestScalarSchemaHasNoProperties_Success() {
	schema := openapi.TypeToSchema(reflect.TypeOf(0))

	assert.Equal(this.T(), "integer", schema.Type)
	assert.Nil(this.T(), schema.Properties)
}

func (this *UtilsTestSuite) TestTypeToParamTagIsTrimmedAndDashSkipped_Success() {
	params := openapi.TypeToParam(reflect.TypeOf(TaggedPayload{}))

	names := map[string]bool{}
	for _, param := range params {
		names[param.Name] = true
	}

	assert.True(this.T(), names["optional"])
	assert.False(this.T(), names["optional,omitempty"])
	assert.True(this.T(), names["-"]) // o campo `json:"-,"`
	assert.False(this.T(), names["internal"])
	assert.False(this.T(), names["Hidden"])

	// `json:"-"` aparecia duas vezes com o mesmo nome "-"; agora só sobra o `-,`.
	dashes := 0
	for _, param := range params {
		if param.Name == "-" {
			dashes++
		}
	}
	assert.Equal(this.T(), 1, dashes)
}

func (this *UtilsTestSuite) TestTypeToParamResolvesFieldTypes_Success() {
	params := openapi.TypeToParam(reflect.TypeOf(TaggedPayload{}))

	schemas := map[string]openapi.Schema{}
	for _, param := range params {
		schemas[param.Name] = param.Schema
	}

	assert.Equal(this.T(), "string", schemas["when"].Type)
	assert.Equal(this.T(), "date-time", schemas["when"].Format)
	assert.Equal(this.T(), "object", schemas["child"].Type)
	assert.Equal(this.T(), "string", schemas["bytes"].Type)
}

// schemaKeys devolve os nomes das propriedades do schema, ordenados.
func schemaKeys(t reflect.Type) []string {
	return slices.Sorted(maps.Keys(openapi.TypeToSchema(t).Properties))
}

// marshalKeys devolve as chaves que o encoding/json realmente emite para o
// valor, ordenadas. É o gabarito contra o qual o schema é comparado.
func (this *UtilsTestSuite) marshalKeys(value any) []string {
	raw, err := json.Marshal(value)
	assert.NoError(this.T(), err)

	decoded := map[string]any{}
	assert.NoError(this.T(), json.Unmarshal(raw, &decoded))

	return slices.Sorted(maps.Keys(decoded))
}

// --- ponto 1: todos os tamanhos de inteiro ---

func (this *UtilsTestSuite) TestTypeToSwaggerEveryIntegerKind_Success() {
	kinds := []reflect.Kind{
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
	}

	for _, kind := range kinds {
		assert.Equal(this.T(), "integer", openapi.TypeToSwagger(kind), kind.String())
	}
}

func (this *UtilsTestSuite) TestSizedIntegerFieldsAreIntegers_Success() {
	type Numbers struct {
		Tiny   int8   `json:"tiny"`
		Small  int16  `json:"small"`
		Medium int32  `json:"medium"`
		Big    int64  `json:"big"`
		Count  uint   `json:"count"`
		Total  uint64 `json:"total"`
		Port   uint16 `json:"port"`
	}

	schema := openapi.TypeToSchema(reflect.TypeOf(Numbers{}))

	for name := range schema.Properties {
		assert.Equal(this.T(), "integer", schema.Properties[name].Type, name)
	}

	assert.Len(this.T(), schema.Properties, 7)
}

func (this *UtilsTestSuite) TestByteArrayIsAListOfIntegers_Success() {
	// Array de byte (tamanho fixo) sai como lista de números no encoding/json,
	// diferente do slice, que sai em base64.
	schema := openapi.TypeToSchema(reflect.TypeOf([4]byte{}))

	assert.Equal(this.T(), "array", schema.Type)
	assert.Equal(this.T(), "integer", schema.Items.Type)
}

// --- ponto 3: structs embedados ---

func (this *UtilsTestSuite) TestEmbeddedStructIsPromoted_Success() {
	schema := openapi.TypeToSchema(reflect.TypeOf(EmbedHost{}))

	assert.Contains(this.T(), schema.Properties, "inherited")
	assert.NotContains(this.T(), schema.Properties, "Promoted")
}

func (this *UtilsTestSuite) TestShallowFieldWinsOverPromoted_Success() {
	schema := openapi.TypeToSchema(reflect.TypeOf(EmbedHost{}))

	assert.Contains(this.T(), schema.Properties, "shadowed")
	assert.Len(this.T(), schema.Properties, 3) // inherited, shadowed, own — e nada de ambiguous
}

func (this *UtilsTestSuite) TestAmbiguousPromotedFieldIsDropped_Success() {
	schema := openapi.TypeToSchema(reflect.TypeOf(EmbedHost{}))

	assert.NotContains(this.T(), schema.Properties, "Ambiguous")
}

// Empatando na profundidade, o único campo com nome na tag leva o nome.
func (this *UtilsTestSuite) TestTaggedFieldWinsTheTieBreak_Success() {
	schema := openapi.TypeToSchema(reflect.TypeOf(TieBreakHost{}))

	assert.Len(this.T(), schema.Properties, 1)
	assert.Equal(this.T(), "integer", schema.Properties["Winner"].Type)
}

func (this *UtilsTestSuite) TestEmbeddedPointerStructIsPromoted_Success() {
	schema := openapi.TypeToSchema(reflect.TypeOf(PointerEmbedHost{}))

	assert.Contains(this.T(), schema.Properties, "inherited")
	assert.Contains(this.T(), schema.Properties, "own")
}

func (this *UtilsTestSuite) TestTaggedEmbeddedStructIsNotPromoted_Success() {
	schema := openapi.TypeToSchema(reflect.TypeOf(TaggedEmbedHost{}))

	assert.NotContains(this.T(), schema.Properties, "inherited")
	assert.Contains(this.T(), schema.Properties, "promoted")
	assert.Contains(this.T(), schema.Properties["promoted"].Properties, "inherited")
}

func (this *UtilsTestSuite) TestEmbeddedStructWithDashTagIsSkipped_Success() {
	schema := openapi.TypeToSchema(reflect.TypeOf(HiddenEmbedHost{}))

	assert.Equal(this.T(), []string{"own"}, schemaKeys(reflect.TypeOf(HiddenEmbedHost{})))
	assert.NotContains(this.T(), schema.Properties, "inherited")
}

func (this *UtilsTestSuite) TestEmbeddedScalarKeepsTypeName_Success() {
	schema := openapi.TypeToSchema(reflect.TypeOf(ScalarEmbedHost{}))

	assert.Contains(this.T(), schema.Properties, "Slug")
	assert.Equal(this.T(), "string", schema.Properties["Slug"].Type)
}

func (this *UtilsTestSuite) TestUnexportedEmbeddedStructIsStillPromoted_Success() {
	schema := openapi.TypeToSchema(reflect.TypeOf(UnexportedEmbedHost{}))

	assert.Contains(this.T(), schema.Properties, "visible")
}

// TestPromotionMatchesEncodingJson compara o schema com o que o encoding/json
// realmente emite — é o gabarito de todas as regras de promoção juntas.
func (this *UtilsTestSuite) TestPromotionMatchesEncodingJson_Success() {
	values := []any{
		EmbedHost{},
		PointerEmbedHost{Promoted: &Promoted{}},
		TaggedEmbedHost{},
		HiddenEmbedHost{},
		ScalarEmbedHost{},
		UnexportedEmbedHost{},
		TieBreakHost{},
	}

	for _, value := range values {
		t := reflect.TypeOf(value)

		assert.Equal(this.T(), this.marshalKeys(value), schemaKeys(t), t.Name())
	}
}

func (this *UtilsTestSuite) TestTypeToParamPromotesEmbeddedFields_Success() {
	params := openapi.TypeToParam(reflect.TypeOf(EmbedHost{}))

	names := []string{}
	for _, param := range params {
		names = append(names, param.Name)
	}

	assert.Contains(this.T(), names, "inherited")
	assert.Contains(this.T(), names, "own")
	assert.NotContains(this.T(), names, "Ambiguous")
	assert.NotContains(this.T(), names, "Promoted")
}

// --- ponto 4: items aninhado ---

func (this *UtilsTestSuite) TestNestedSliceKeepsInnerItems_Success() {
	schema := openapi.TypeToSchema(reflect.TypeOf([][]string{}))

	assert.Equal(this.T(), "array", schema.Type)
	assert.Equal(this.T(), "array", schema.Items.Type)

	if assert.NotNil(this.T(), schema.Items.Items) {
		assert.Equal(this.T(), "string", schema.Items.Items.Type)
	}
}

func (this *UtilsTestSuite) TestNestedSliceOfStructKeepsInnerProperties_Success() {
	schema := openapi.TypeToSchema(reflect.TypeOf([][]UserPayload{}))

	if assert.NotNil(this.T(), schema.Items.Items) {
		assert.Equal(this.T(), "object", schema.Items.Items.Type)
		assert.Contains(this.T(), schema.Items.Items.Properties, "id")
	}
}

func (this *UtilsTestSuite) TestFlatSliceHasNoInnerItems_Success() {
	schema := openapi.TypeToSchema(reflect.TypeOf([]string{}))

	assert.Equal(this.T(), "string", schema.Items.Type)
	assert.Nil(this.T(), schema.Items.Items)
}

// --- ponto 6: required derivado do omitempty ---

func (this *UtilsTestSuite) TestOmitemptyMakesParamOptional_Success() {
	params := openapi.TypeToParam(reflect.TypeOf(TaggedPayload{}))

	required := map[string]bool{}
	for _, param := range params {
		required[param.Name] = param.Required
	}

	assert.False(this.T(), required["optional"])
	assert.True(this.T(), required["name"])
	assert.True(this.T(), required["NoTag"])
	assert.True(this.T(), required["-"])
}

func (this *UtilsTestSuite) TestOtherTagOptionsDoNotMakeParamOptional_Success() {
	type Payload struct {
		Number string `json:"number,string"`
		Both   string `json:"both,string,omitempty"`
	}

	params := openapi.TypeToParam(reflect.TypeOf(Payload{}))

	required := map[string]bool{}
	for _, param := range params {
		required[param.Name] = param.Required
	}

	assert.True(this.T(), required["number"])
	assert.False(this.T(), required["both"])
}
