package build

import (
	"encoding"
	"errors"
	"fmt"
	"math/bits"
	"net"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/talav/openapi/config"
	"github.com/talav/openapi/hook"
	"github.com/talav/openapi/internal/model"
	"github.com/talav/openapi/metadata"
	"github.com/talav/schema"
)

const (
	// JSON Schema type constants.
	TypeString  = "string"
	TypeArray   = "array"
	TypeObject  = "object"
	TypeBoolean = "boolean"
	TypeInteger = "integer"
	TypeNumber  = "number"

	formatInt32           = "int32"
	formatInt64           = "int64"
	contentEncodingBase64 = "base64"
)

var (
	// Interface types for efficient implementation checks without allocation.
	schemaTransformerType = reflect.TypeOf((*hook.SchemaTransformer)(nil)).Elem()
	schemaProviderType    = reflect.TypeOf((*hook.SchemaProvider)(nil)).Elem()
	textUnmarshalerType   = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()

	// Standard library types for schema generation.
	timeType   = reflect.TypeOf(time.Time{})
	urlType    = reflect.TypeOf(url.URL{})
	ipType     = reflect.TypeOf(net.IP{})
	ipAddrType = reflect.TypeOf(net.IPAddr{})
)

type schemaNamerFunc func(t reflect.Type, hint string) string

// SchemaGenerator generates and caches OpenAPI schemas from Go types.
// It handles schema generation, caching, reference management, and type aliases.
type SchemaGenerator struct {
	// Configuration
	prefix   string
	namer    schemaNamerFunc
	metadata *schema.Metadata
	tagCfg   config.TagConfig

	// Cache
	schemas map[string]*model.Schema
	types   map[string]reflect.Type
	seen    map[reflect.Type]string // type -> name mapping for deduplication

	// Options
	inlineOnly map[string]bool               // Schemas excluded from components
	aliases    map[reflect.Type]reflect.Type // Type aliases
}

// NewSchemaGenerator creates a new schema generator with the given configuration.
func NewSchemaGenerator(prefix string, m *schema.Metadata, tagCfg config.TagConfig) *SchemaGenerator {
	return &SchemaGenerator{
		prefix:     prefix,
		namer:      schemaNamer,
		metadata:   m,
		tagCfg:     tagCfg,
		schemas:    make(map[string]*model.Schema),
		types:      make(map[string]reflect.Type),
		seen:       make(map[reflect.Type]string),
		inlineOnly: make(map[string]bool),
		aliases:    make(map[reflect.Type]reflect.Type),
	}
}

// Schema generates a schema for the given type. It handles caching, references,
// and type aliases automatically. For most use cases, this is the only method needed.
func (g *SchemaGenerator) Schema(t reflect.Type) *model.Schema {
	return g.schema(t, true, "")
}

// Schemas returns all generated schemas as a map, suitable for OpenAPI components/schemas.
// Inline-only schemas (marked via MarkInlineOnly) are excluded.
func (g *SchemaGenerator) Schemas() map[string]*model.Schema {
	result := make(map[string]*model.Schema, len(g.schemas))
	for name, schema := range g.schemas {
		if !g.inlineOnly[name] {
			result[name] = schema
		}
	}

	return result
}

// markInlineOnly marks a type to be excluded from the Schemas() map.
// The schema will still be generated and can be referenced, but won't appear
// in components/schemas. Useful for types that are only used inline.
// The hint parameter should match the hint used when generating the schema.
func (g *SchemaGenerator) markInlineOnly(t reflect.Type, hint string) {
	t = deref(t)
	name := g.namer(t, hint)
	g.inlineOnly[name] = true
}

// schema is the internal method that handles the full schema generation logic.
// allowRef controls whether to return a $ref or inline schema.
// hint is used for naming unnamed types.
//
//nolint:cyclop // exclude
func (g *SchemaGenerator) schema(t reflect.Type, allowRef bool, hint string) *model.Schema {
	origType := t
	t = deref(t)

	// Pointer to array should decay to array
	if t.Kind() == reflect.Array || t.Kind() == reflect.Slice {
		origType = t
	}

	// Handle type aliases
	if alias, ok := g.aliases[t]; ok {
		return g.schema(alias, allowRef, hint)
	}

	// Determine if this type should get a reference
	getsRef := g.shouldGetRef(t)
	name := g.namer(origType, hint)

	// Check cache if it gets a ref
	//nolint:nestif // Complex nested logic for reference handling - acceptable complexity
	if getsRef {
		if s, ok := g.schemas[name]; ok {
			// Verify type consistency
			if seenName, exists := g.seen[t]; !exists || seenName != name {
				// Name matches but type is different, so we have a dupe.
				panic(fmt.Errorf("duplicate name: %s, new type: %s, existing type: %s", name, t, g.types[name]))
			}
			if allowRef {
				return &model.Schema{Ref: g.prefix + name}
			}

			return s
		}
	}

	// Register placeholder for recursive types
	if getsRef {
		g.schemas[name] = &model.Schema{}
		g.types[name] = t
		g.seen[t] = name
	}

	// Generate the schema
	s, err := g.generate(origType)
	if err != nil {
		panic(fmt.Errorf("failed to generate schema for type %s: %w", origType, err))
	}

	// Store if it gets a ref
	if getsRef {
		g.schemas[name] = s
	}

	// Return ref or inline
	if getsRef && allowRef {
		return &model.Schema{Ref: g.prefix + name}
	}

	return s
}

// shouldGetRef determines if a type should be stored with a reference.
func (g *SchemaGenerator) shouldGetRef(t reflect.Type) bool {
	if t.Kind() != reflect.Struct {
		return false
	}

	// Special case: time.Time is always a string.
	if t == timeType {
		return false
	}

	// Check for special interfaces
	v := reflect.New(t).Interface()
	if _, ok := v.(hook.SchemaProvider); ok {
		return false
	}
	if _, ok := v.(encoding.TextUnmarshaler); ok {
		return false
	}

	return true
}

// generate creates a schema for a type (internal, no caching or refs).
func (g *SchemaGenerator) generate(t reflect.Type) (*model.Schema, error) {
	isPointer := t.Kind() == reflect.Pointer
	t = deref(t)

	// Check for interface implementations that override schema generation
	if schema, err := g.schemaFromInterface(t, isPointer); schema != nil || err != nil {
		return schema, err
	}

	// Lookup in maps (type first, then kind)
	if s := g.schemaForSimpleType(t, isPointer); s != nil {
		return s, nil
	}

	//nolint:exhaustive // Only handling supported Go types for OpenAPI schema generation
	switch t.Kind() {
	case reflect.Slice, reflect.Array:
		return g.generateArray(t, isPointer)
	case reflect.Map:
		return g.generateMap(t)
	case reflect.Struct:
		return g.generateStruct(t)
	case reflect.Interface:
		// Interfaces mean any object.
		return &model.Schema{}, nil
	default:
		//nolint:nilnil // Returning nil schema for unsupported types is intentional
		return nil, nil
	}
}

// schemaFromInterface checks if the type implements SchemaProvider or TextUnmarshaler.
func (g *SchemaGenerator) schemaFromInterface(t reflect.Type, isPointer bool) (*model.Schema, error) {
	// Check SchemaProvider without allocation first
	if t.Implements(schemaProviderType) || reflect.PointerTo(t).Implements(schemaProviderType) {
		// Special case: type provides its own schema. Do not try to generate.
		v := reflect.New(t).Interface()
		sp, ok := v.(hook.SchemaProvider)
		if !ok {
			return nil, fmt.Errorf("type does not implement SchemaProvider")
		}

		return sp.Schema(g), nil
	}

	// Check TextUnmarshaler without allocation
	if t.Implements(textUnmarshalerType) || reflect.PointerTo(t).Implements(textUnmarshalerType) {
		// Special case: types that implement encoding.TextUnmarshaler are able to
		// be loaded from plain text, and so should be treated as strings.
		return &model.Schema{Type: TypeString, Nullable: isPointer}, nil
	}

	//nolint:nilnil // Returning (nil, nil) signals that no interface implementation was found
	return nil, nil
}

var (
	lookUpByType = map[reflect.Type]*model.Schema{
		timeType:   {Type: TypeString, Format: "date-time"},
		urlType:    {Type: TypeString, Format: "uri"},
		ipType:     {Type: TypeString, Format: "ipv4"},
		ipAddrType: {Type: TypeString, Format: "ipv4"},
	}

	lookUpByKind = map[reflect.Kind]*model.Schema{
		reflect.Bool:    {Type: TypeBoolean},
		reflect.Int8:    {Type: TypeInteger, Format: formatInt32},
		reflect.Int16:   {Type: TypeInteger, Format: formatInt32},
		reflect.Int32:   {Type: TypeInteger, Format: formatInt32},
		reflect.Int64:   {Type: TypeInteger, Format: formatInt64},
		reflect.Uint8:   {Type: TypeInteger, Format: formatInt32, Minimum: &model.Bound{Value: 0}},
		reflect.Uint16:  {Type: TypeInteger, Format: formatInt32, Minimum: &model.Bound{Value: 0}},
		reflect.Uint32:  {Type: TypeInteger, Format: formatInt32, Minimum: &model.Bound{Value: 0}},
		reflect.Uint64:  {Type: TypeInteger, Format: formatInt64, Minimum: &model.Bound{Value: 0}},
		reflect.Float32: {Type: TypeNumber, Format: "float"},
		reflect.Float64: {Type: TypeNumber, Format: "double"},
		reflect.String:  {Type: TypeString},
	}
)

// schemaForSimpleType looks up schema information by type first, then by kind.
func (g *SchemaGenerator) schemaForSimpleType(t reflect.Type, isPointer bool) *model.Schema {
	// Try type lookup first (for stdlib types)
	if found, ok := lookUpByType[t]; ok {
		s := *found
		applyNullableForScalar(&s, isPointer)

		return &s
	}

	// Try kind lookup
	kind := t.Kind()
	if kind == reflect.Int || kind == reflect.Uint {
		s := &model.Schema{Type: TypeInteger}
		if bits.UintSize == 32 {
			s.Format = formatInt32
		} else {
			s.Format = formatInt64
		}
		if kind == reflect.Uint {
			s.Minimum = &model.Bound{Value: 0}
		}
		applyNullableForScalar(s, isPointer)

		return s
	}

	if found, ok := lookUpByKind[kind]; ok {
		s := *found
		applyNullableForScalar(&s, isPointer)

		return &s
	}

	return nil
}

// generateArray generates a schema for slice or array types.
func (g *SchemaGenerator) generateArray(t reflect.Type, isPointer bool) (*model.Schema, error) {
	s := model.Schema{}

	if t.Elem().Kind() == reflect.Uint8 {
		// Special case: []byte will be serialized as a base64 string.
		s.Type = TypeString
		s.ContentEncoding = contentEncodingBase64
		s.ContentMediaType = "application/octet-stream"
		s.Nullable = isPointer
	} else {
		s.Type = TypeArray
		s.Nullable = false
		s.Items = g.schema(t.Elem(), true, t.Name()+"Item")

		if t.Kind() == reflect.Array {
			l := t.Len()
			s.MinItems = &l
			s.MaxItems = &l
		}
	}

	return &s, nil
}

// generateMap generates a schema for map types.
func (g *SchemaGenerator) generateMap(t reflect.Type) (*model.Schema, error) {
	s := model.Schema{Type: TypeObject}
	valueSchema := g.schema(t.Elem(), true, t.Name()+"Value")
	s.Additional = &model.Additional{Schema: valueSchema}

	return &s, nil
}

// structFieldsResult contains the results of processing struct fields.
type structFieldsResult struct {
	// props maps property names to their OpenAPI schemas.
	// These become the "properties" field in the generated object schema.
	props map[string]*model.Schema

	// required lists property names that must be present in the object.
	// These become the "required" array in the generated schema.
	required []string

	// dependentRequired maps a field name to a list of other fields that must be present
	// when the mapped field is present. This implements JSON Schema 2019-09 / OpenAPI 3.1
	// dependentRequired feature for conditional required fields.
	dependentRequired map[string][]string
}

// generateStruct generates a schema for struct types.
func (g *SchemaGenerator) generateStruct(t reflect.Type) (*model.Schema, error) {
	// Get struct metadata
	structMeta, err := g.metadata.GetStructMetadata(t)
	if err != nil {
		return nil, fmt.Errorf("failed to get struct metadata for type %s: %w", t, err)
	}

	s := model.Schema{Type: TypeObject}

	// Process each field and build properties
	result := g.processStructFields(t, *structMeta)

	// Validate dependent required fields
	if err := validateDependentRequired(result.dependentRequired, result.props); err != nil {
		return nil, err
	}

	// Store dependentRequired (JSON Schema 2019-09 / OpenAPI 3.1 feature)
	if len(result.dependentRequired) > 0 {
		s.DependentRequired = result.dependentRequired
	}

	// Handle struct-level metadata (_ field)
	g.applyStructLevelMetadata(&s, structMeta)

	// Apply SchemaTransformer if implemented
	if t.Implements(schemaTransformerType) || reflect.PointerTo(t).Implements(schemaTransformerType) {
		v := reflect.New(t).Interface()
		if st, ok := v.(hook.SchemaTransformer); ok {
			s = *st.TransformSchema(g, &s)
		}
	}

	s.Properties = result.props
	s.Required = result.required

	return &s, nil
}

// processStructFields iterates through struct fields and builds property schemas.
func (g *SchemaGenerator) processStructFields(t reflect.Type, structMeta schema.StructMetadata) structFieldsResult {
	result := structFieldsResult{
		props:             make(map[string]*model.Schema),
		dependentRequired: make(map[string][]string),
	}

	// Iterate through metadata fields
	for _, fieldMeta := range structMeta.Fields {
		if g.isHidden(fieldMeta) {
			continue
		}

		reflectField := t.Field(fieldMeta.Index)
		fs := g.schema(reflectField.Type, true, t.Name()+fieldMeta.StructFieldName+"Struct")
		if fs == nil {
			continue
		}
		// Extract field name from metadata (respects JSON tags)
		name := g.defineFieldName(reflectField, fieldMeta)

		// Determine required status from metadata
		fieldRequired := isRequiredFromMetadata(&fieldMeta, g.tagCfg)

		// Apply OpenAPI metadata
		g.applyOpenAPIMetadata(fs, fieldMeta)

		// Apply validation metadata
		g.applyValidateMetadata(fs, fieldMeta)

		// If field is required, it cannot be null
		if fieldRequired {
			fs.Nullable = false
		}

		// Apply default value from default tag
		g.applyDefaultValue(fs, fieldMeta)

		// Apply dependent required metadata (on object schema, not field schema)
		g.applyDependentRequired(result.dependentRequired, fieldMeta, name)

		// Add to properties
		result.props[name] = fs

		if fieldRequired {
			result.required = append(result.required, name)
		}
	}

	return result
}

// validateDependentRequired validates that all dependent required fields exist.
func validateDependentRequired(dependentRequired map[string][]string, props map[string]*model.Schema) error {
	var errs []error
	for field, dependents := range dependentRequired {
		for _, dependent := range dependents {
			if _, ok := props[dependent]; !ok {
				errs = append(errs, fmt.Errorf("dependent field '%s' for field '%s' does not exist", dependent, field))
			}
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("dependent required validation failed: %w", errors.Join(errs...))
	}

	return nil
}

// defineFieldName extracts the field name from metadata, respecting JSON tags.
// Priority: JSON tag > explicit schema tag > struct field name.
func (g *SchemaGenerator) defineFieldName(field reflect.StructField, fieldMeta schema.FieldMetadata) string {
	// First, check JSON tag for field name (most common case for OpenAPI schemas)
	if jsonTag, ok := field.Tag.Lookup("json"); ok {
		// Parse JSON tag (format: "name,omitempty,string")
		parts := strings.Split(jsonTag, ",")
		if len(parts) > 0 && parts[0] != "" && parts[0] != "-" {
			return parts[0]
		}
	}

	// Second, check schema tag for explicit parameter name
	if schemaMeta, ok := schema.GetTagMetadata[*schema.SchemaMetadata](&fieldMeta, g.tagCfg.Schema); ok {
		if schemaMeta.ParamName != "" {
			return schemaMeta.ParamName
		}
	}

	// Fall back to struct field name
	return fieldMeta.StructFieldName
}

// isHidden determines if a field is hidden based on metadata.
func (g *SchemaGenerator) isHidden(fieldMeta schema.FieldMetadata) bool {
	if openAPIMeta, ok := schema.GetTagMetadata[*metadata.OpenAPIMetadata](&fieldMeta, g.tagCfg.OpenAPI); ok {
		return toBool(openAPIMeta.Hidden)
	}

	return false
}

// applyOpenAPIMetadata applies OpenAPI metadata to a schema.
func (g *SchemaGenerator) applyOpenAPIMetadata(fs *model.Schema, fieldMeta schema.FieldMetadata) {
	openAPIMeta, ok := schema.GetTagMetadata[*metadata.OpenAPIMetadata](&fieldMeta, g.tagCfg.OpenAPI)
	if !ok {
		return
	}

	fs.Title = openAPIMeta.Title
	fs.Description = openAPIMeta.Description
	fs.Format = openAPIMeta.Format
	fs.Examples = openAPIMeta.Examples
	fs.ReadOnly = toBool(openAPIMeta.ReadOnly)
	fs.WriteOnly = toBool(openAPIMeta.WriteOnly)
	fs.Deprecated = toBool(openAPIMeta.Deprecated)
	fs.Extensions = openAPIMeta.Extensions
}

// applyStructLevelMetadata extracts struct-level metadata from the _ field.
func (g *SchemaGenerator) applyStructLevelMetadata(s *model.Schema, structMeta *schema.StructMetadata) {
	fieldMeta, ok := structMeta.Field("_")
	if !ok {
		return
	}

	openAPIMeta, ok := schema.GetTagMetadata[*metadata.OpenAPIMetadata](fieldMeta, g.tagCfg.OpenAPI)
	if !ok {
		return
	}

	// Apply struct-level options from parsed metadata (only valid when used on _ field)
	if openAPIMeta.AdditionalProperties != nil {
		// Convert to model.Additional format
		allow := *openAPIMeta.AdditionalProperties
		s.Additional = &model.Additional{Allow: &allow}
	}
	if openAPIMeta.Nullable != nil {
		s.Nullable = *openAPIMeta.Nullable
	}
}

// applyDefaultValue reads the default tag from metadata and applies it to the schema.
func (g *SchemaGenerator) applyDefaultValue(fs *model.Schema, fieldMeta schema.FieldMetadata) {
	defaultMeta, ok := schema.GetTagMetadata[*metadata.DefaultMetadata](&fieldMeta, g.tagCfg.Default)
	if !ok || defaultMeta.Value == nil {
		return
	}

	fs.Default = defaultMeta.Value
}

// applyValidateMetadata applies validation constraints from ValidateMetadata to a schema.
func (g *SchemaGenerator) applyValidateMetadata(fs *model.Schema, fieldMeta schema.FieldMetadata) {
	validateMeta, ok := schema.GetTagMetadata[*metadata.ValidateMetadata](&fieldMeta, g.tagCfg.Validate)
	if !ok {
		return
	}

	// Handle minimum/maximum based on type
	applyMinMaxConstraints(fs, validateMeta)

	// Exclusive numeric constraints (convert to Bound format)
	if validateMeta.ExclusiveMinimum != nil {
		fs.Minimum = &model.Bound{Value: *validateMeta.ExclusiveMinimum, Exclusive: true}
	}
	if validateMeta.ExclusiveMaximum != nil {
		fs.Maximum = &model.Bound{Value: *validateMeta.ExclusiveMaximum, Exclusive: true}
	}
	fs.MultipleOf = validateMeta.MultipleOf

	// String-specific constraints
	fs.Pattern = validateMeta.Pattern
	if fs.Format == "" {
		fs.Format = validateMeta.Format
	}

	// Handle enum
	applyEnumConstraints(fs, validateMeta)
}

// applyMinMaxConstraints applies minimum and maximum constraints based on schema type.
func applyMinMaxConstraints(fs *model.Schema, validateMeta *metadata.ValidateMetadata) {
	switch fs.Type {
	case TypeString:
		applyStringMinMax(fs, validateMeta)
	case TypeInteger, TypeNumber:
		applyNumericMinMax(fs, validateMeta)
	case TypeArray:
		applyArrayMinMax(fs, validateMeta)
	case TypeObject:
		applyObjectMinMax(fs, validateMeta)
	}
}

// applyStringMinMax applies min/max length constraints for string types.
func applyStringMinMax(fs *model.Schema, validateMeta *metadata.ValidateMetadata) {
	if validateMeta.Minimum != nil {
		minLen := int(*validateMeta.Minimum)
		fs.MinLength = &minLen
	}
	if validateMeta.Maximum != nil {
		maxLen := int(*validateMeta.Maximum)
		fs.MaxLength = &maxLen
	}
}

// applyNumericMinMax applies min/max value constraints for numeric types.
func applyNumericMinMax(fs *model.Schema, validateMeta *metadata.ValidateMetadata) {
	if validateMeta.Minimum != nil {
		fs.Minimum = &model.Bound{Value: *validateMeta.Minimum, Exclusive: false}
	}
	if validateMeta.Maximum != nil {
		fs.Maximum = &model.Bound{Value: *validateMeta.Maximum, Exclusive: false}
	}
}

// applyArrayMinMax applies min/max item count constraints for array types.
func applyArrayMinMax(fs *model.Schema, validateMeta *metadata.ValidateMetadata) {
	if validateMeta.Minimum != nil {
		minItems := int(*validateMeta.Minimum)
		fs.MinItems = &minItems
	}
	if validateMeta.Maximum != nil {
		maxItems := int(*validateMeta.Maximum)
		fs.MaxItems = &maxItems
	}
}

// applyObjectMinMax applies min/max property count constraints for object types.
func applyObjectMinMax(fs *model.Schema, validateMeta *metadata.ValidateMetadata) {
	if validateMeta.Minimum != nil {
		minProps := int(*validateMeta.Minimum)
		fs.MinProperties = &minProps
	}
	if validateMeta.Maximum != nil {
		maxProps := int(*validateMeta.Maximum)
		fs.MaxProperties = &maxProps
	}
}

// applyEnumConstraints applies enum or const constraints to the schema.
func applyEnumConstraints(fs *model.Schema, validateMeta *metadata.ValidateMetadata) {
	target := fs
	if fs.Type == TypeArray && fs.Items != nil {
		target = fs.Items
	}

	if len(validateMeta.Enum) == 1 {
		target.Const = validateMeta.Enum[0]
	} else {
		target.Enum = validateMeta.Enum
	}
}

// applyDependentRequired applies requires metadata to the dependentRequired map.
func (g *SchemaGenerator) applyDependentRequired(dependentRequired map[string][]string, fieldMeta schema.FieldMetadata, fieldName string) {
	reqMeta, ok := schema.GetTagMetadata[*metadata.RequiresMetadata](&fieldMeta, g.tagCfg.Requires)
	if !ok || len(reqMeta.Fields) == 0 {
		return
	}

	dependentRequired[fieldName] = reqMeta.Fields
}

// applyNullableForScalar sets nullable for scalar types if isPointer is true.
func applyNullableForScalar(s *model.Schema, isPointer bool) {
	if s.Type == TypeBoolean || s.Type == TypeInteger || s.Type == TypeNumber || s.Type == TypeString {
		s.Nullable = isPointer
	}
}
