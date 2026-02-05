package v304

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/talav/openapi/debug"
	"github.com/talav/openapi/internal/model"
)

//go:embed schema_v304.json
var schemaV304JSON []byte

type AdapterV304 struct{}

func (a *AdapterV304) Version() string {
	return "3.0.4"
}

func (a *AdapterV304) SchemaJSON() []byte {
	return schemaV304JSON
}

func (a *AdapterV304) View(spec *model.Spec) (any, debug.Warnings, error) {
	if spec == nil {
		return nil, nil, fmt.Errorf("nil spec")
	}

	var warnings debug.Warnings

	// Warn about Webhooks (3.1-only)
	if len(spec.Webhooks) > 0 {
		warnings = append(warnings, debug.NewWarning(debug.WarnDegradationWebhooks, "#/webhooks", "webhooks are 3.1-only; dropped"))
	}

	result := &ViewV304{
		OpenAPI:      a.Version(),
		Info:         a.transformInfo(spec.Info, &warnings),
		Servers:      a.transformServers(spec.Servers),
		Paths:        a.transformPaths(spec.Paths, &warnings),
		Components:   a.transformComponents(spec.Components, &warnings),
		Security:     a.transformSecurity(spec.Security),
		Tags:         a.transformTags(spec.Tags),
		ExternalDocs: a.transformExternalDocs(spec.ExternalDocs),
		Extensions:   spec.Extensions,
	}

	if err := validateViewV304(result); err != nil {
		return nil, nil, err
	}

	return result, warnings, nil
}

// validateViewV304 validates a ViewV304 instance according to OpenAPI 3.0.4 requirements.
func validateViewV304(result *ViewV304) error {
	if result.Info.Title == "" {
		return fmt.Errorf("openapi: title is required")
	}
	if result.Info.Version == "" {
		return fmt.Errorf("openapi: version is required")
	}

	// Validate servers: variables require a server URL
	for i, server := range result.Servers {
		if len(server.Variables) > 0 && server.URL == "" {
			return fmt.Errorf("openapi: server[%d]: server variables require a server URL", i)
		}
	}

	// Validate extension keys must start with "x-"
	for key := range result.Extensions {
		if !strings.HasPrefix(key, "x-") {
			return fmt.Errorf("openapi: extension key must start with 'x-': %s", key)
		}
	}

	// Validate info extension keys must start with "x-"
	for key := range result.Info.Extensions {
		if !strings.HasPrefix(key, "x-") {
			return fmt.Errorf("openapi: info extension key must start with 'x-': %s", key)
		}
	}

	return nil
}

func (a *AdapterV304) transformInfo(in model.Info, warnings *debug.Warnings) *InfoV30 {
	info := &InfoV30{
		Title:          in.Title,
		Description:    in.Description,
		TermsOfService: in.TermsOfService,
		Version:        in.Version,
		Extensions:     in.Extensions,
	}

	// Drop Summary (3.1-only)
	if in.Summary != "" {
		*warnings = append(*warnings, debug.NewWarning(debug.WarnDegradationInfoSummary, "#/info/summary", "info.summary is 3.1-only; dropped"))
	}

	if in.Contact != nil {
		info.Contact = &ContactV30{
			Name:       in.Contact.Name,
			URL:        in.Contact.URL,
			Email:      in.Contact.Email,
			Extensions: in.Contact.Extensions,
		}
	}

	if in.License != nil {
		info.License = &LicenseV30{
			Name:       in.License.Name,
			URL:        in.License.URL,
			Extensions: in.License.Extensions,
		}
		// Drop Identifier (3.1-only)
		if in.License.Identifier != "" {
			*warnings = append(*warnings, debug.NewWarning(debug.WarnDegradationLicenseIdentifier, "#/info/license", "license identifier is 3.1-only; dropped (use url instead)"))
		}
	}

	return info
}

func (a *AdapterV304) transformServers(in []model.Server) []*ServerV30 {
	if len(in) == 0 {
		return nil
	}

	servers := make([]*ServerV30, 0, len(in))
	for _, s := range in {
		server := &ServerV30{
			URL:         s.URL,
			Description: s.Description,
			Extensions:  s.Extensions,
		}

		if len(s.Variables) > 0 {
			server.Variables = make(map[string]*ServerVariableV30, len(s.Variables))
			for name, v := range s.Variables {
				server.Variables[name] = &ServerVariableV30{
					Enum:        v.Enum,
					Default:     v.Default,
					Description: v.Description,
					Extensions:  v.Extensions,
				}
			}
		}

		servers = append(servers, server)
	}

	return servers
}

func (a *AdapterV304) transformTags(in []model.Tag) []*TagV30 {
	if len(in) == 0 {
		return nil
	}

	tags := make([]*TagV30, 0, len(in))
	for _, t := range in {
		tag := &TagV30{
			Name:        t.Name,
			Description: t.Description,
			Extensions:  t.Extensions,
		}

		if t.ExternalDocs != nil {
			tag.ExternalDocs = a.transformExternalDocs(t.ExternalDocs)
		}

		tags = append(tags, tag)
	}

	return tags
}

func (a *AdapterV304) transformSecurity(in []model.SecurityRequirement) []SecurityRequirementV30 {
	if len(in) == 0 {
		return nil
	}

	security := make([]SecurityRequirementV30, 0, len(in))
	for _, s := range in {
		security = append(security, SecurityRequirementV30(s))
	}

	return security
}

func (a *AdapterV304) transformExternalDocs(in *model.ExternalDocs) *ExternalDocsV30 {
	if in == nil {
		return nil
	}

	return &ExternalDocsV30{
		Description: in.Description,
		URL:         in.URL,
		Extensions:  in.Extensions,
	}
}

func (a *AdapterV304) transformPaths(in map[string]*model.PathItem, warnings *debug.Warnings) PathsV30 {
	if len(in) == 0 {
		return make(PathsV30)
	}

	paths := make(PathsV30, len(in))
	for path, item := range in {
		paths[path] = a.transformPathItem(item, warnings)
	}

	return paths
}

func (a *AdapterV304) transformPathItem(in *model.PathItem, warnings *debug.Warnings) *PathItemV30 {
	if in == nil {
		return nil
	}

	// Handle $ref case
	if in.Ref != "" {
		return &PathItemV30{Ref: in.Ref}
	}

	item := &PathItemV30{
		Summary:     in.Summary,
		Description: in.Description,
		Extensions:  in.Extensions,
	}

	// Transform Parameters
	if len(in.Parameters) > 0 {
		item.Parameters = a.transformParameters(in.Parameters, warnings)
	}

	// Transform Operations
	item.Get = a.transformOperation(in.Get, warnings)
	item.Put = a.transformOperation(in.Put, warnings)
	item.Post = a.transformOperation(in.Post, warnings)
	item.Delete = a.transformOperation(in.Delete, warnings)
	item.Options = a.transformOperation(in.Options, warnings)
	item.Head = a.transformOperation(in.Head, warnings)
	item.Patch = a.transformOperation(in.Patch, warnings)
	item.Trace = a.transformOperation(in.Trace, warnings)

	return item
}

func (a *AdapterV304) transformParameters(in []model.Parameter, warnings *debug.Warnings) []*ParameterV30 {
	out := make([]*ParameterV30, 0, len(in))
	for _, param := range in {
		p := a.transformParameter(param, warnings)
		out = append(out, &p)
	}

	return out
}

func (a *AdapterV304) transformParameter(in model.Parameter, warnings *debug.Warnings) ParameterV30 {
	// Handle $ref case
	if in.Ref != "" {
		return ParameterV30{Ref: in.Ref}
	}

	param := ParameterV30{
		Name:            in.Name,
		In:              in.In,
		Description:     in.Description,
		Required:        in.Required,
		Deprecated:      in.Deprecated,
		AllowEmptyValue: in.AllowEmptyValue,
		Style:           in.Style,
		Explode:         in.Explode,
		AllowReserved:   in.AllowReserved,
		Example:         in.Example,
		Extensions:      in.Extensions,
	}

	param.Schema = a.transformSchema(in.Schema, warnings)

	if len(in.Examples) > 0 {
		param.Examples = make(map[string]*ExampleV30, len(in.Examples))
		for k, v := range in.Examples {
			param.Examples[k] = a.transformExample(v, warnings)
		}
	}

	return param
}

func (a *AdapterV304) transformExample(in *model.Example, warnings *debug.Warnings) *ExampleV30 {
	if in == nil {
		return nil
	}

	// Handle $ref case
	if in.Ref != "" {
		return &ExampleV30{Ref: in.Ref}
	}

	// Per OpenAPI spec, value and externalValue are mutually exclusive
	out := &ExampleV30{
		Summary:     in.Summary,
		Description: in.Description,
		Extensions:  in.Extensions,
	}
	if in.ExternalValue != "" {
		out.ExternalValue = in.ExternalValue
		// Warn if both are set (spec violation)
		if in.Value != nil && warnings != nil {
			*warnings = append(*warnings, debug.NewWarning(
				debug.WarnInvalidExampleMutualExclusivity,
				"#/components/examples",
				"example has both value and externalValue set; using externalValue only (spec requires mutual exclusivity)",
			))
		}
	} else {
		out.Value = in.Value
	}

	return out
}

func (a *AdapterV304) transformOperation(in *model.Operation, warnings *debug.Warnings) *OperationV30 {
	if in == nil {
		return nil
	}

	op := &OperationV30{
		Tags:        append([]string(nil), in.Tags...),
		Summary:     in.Summary,
		Description: in.Description,
		OperationID: in.OperationID,
		Deprecated:  in.Deprecated,
		Extensions:  in.Extensions,
	}

	if len(in.Parameters) > 0 {
		op.Parameters = a.transformParameters(in.Parameters, warnings)
	}

	op.RequestBody = a.transformRequestBody(in.RequestBody, warnings)
	op.Security = a.transformSecurity(in.Security)
	op.Servers = a.transformServers(in.Servers)

	if len(in.Responses) > 0 {
		op.Responses = a.transformResponses(in.Responses, warnings)
	}

	return op
}

func (a *AdapterV304) transformRequestBody(in *model.RequestBody, warnings *debug.Warnings) *RequestBodyV30 {
	if in == nil {
		return nil
	}

	// Handle $ref case
	if in.Ref != "" {
		return &RequestBodyV30{Ref: in.Ref}
	}

	rb := &RequestBodyV30{
		Description: in.Description,
		Required:    in.Required,
		Extensions:  in.Extensions,
	}

	if len(in.Content) > 0 {
		rb.Content = make(map[string]*MediaTypeV30, len(in.Content))
		for ct, mt := range in.Content {
			rb.Content[ct] = a.transformMediaType(mt, warnings)
		}
	}

	return rb
}

func (a *AdapterV304) transformMediaType(in *model.MediaType, warnings *debug.Warnings) *MediaTypeV30 {
	if in == nil {
		return nil
	}

	mt := &MediaTypeV30{
		Example:    in.Example,
		Extensions: in.Extensions,
	}

	mt.Schema = a.transformSchema(in.Schema, warnings)

	if len(in.Examples) > 0 {
		mt.Examples = make(map[string]*ExampleV30, len(in.Examples))
		for k, ex := range in.Examples {
			mt.Examples[k] = a.transformExample(ex, warnings)
		}
	}

	return mt
}

//nolint:cyclop
func (a *AdapterV304) transformComponents(in *model.Components, warnings *debug.Warnings) *ComponentsV30 {
	if in == nil {
		return nil
	}

	comp := &ComponentsV30{
		Extensions: in.Extensions,
	}

	if len(in.Schemas) > 0 {
		comp.Schemas = make(map[string]*SchemaV30, len(in.Schemas))
		for name, schema := range in.Schemas {
			comp.Schemas[name] = a.transformSchema(schema, warnings)
		}
	}

	if len(in.Responses) > 0 {
		comp.Responses = make(map[string]*ResponseV30, len(in.Responses))
		for name, r := range in.Responses {
			comp.Responses[name] = a.transformResponse(r, warnings)
		}
	}

	if len(in.Parameters) > 0 {
		comp.Parameters = make(map[string]*ParameterV30, len(in.Parameters))
		for name, param := range in.Parameters {
			pv := a.transformParameter(*param, warnings)
			comp.Parameters[name] = &pv
		}
	}

	if len(in.Examples) > 0 {
		comp.Examples = make(map[string]*ExampleV30, len(in.Examples))
		for name, ex := range in.Examples {
			comp.Examples[name] = a.transformExample(ex, warnings)
		}
	}

	if len(in.RequestBodies) > 0 {
		comp.RequestBodies = make(map[string]*RequestBodyV30, len(in.RequestBodies))
		for name, rb := range in.RequestBodies {
			comp.RequestBodies[name] = a.transformRequestBody(rb, warnings)
		}
	}

	if len(in.Headers) > 0 {
		comp.Headers = make(map[string]*HeaderV30, len(in.Headers))
		for name, h := range in.Headers {
			comp.Headers[name] = a.transformHeader(h, warnings)
		}
	}

	if len(in.SecuritySchemes) > 0 {
		comp.SecuritySchemes = make(map[string]*SecuritySchemeV30, len(in.SecuritySchemes))
		for name, ss := range in.SecuritySchemes {
			comp.SecuritySchemes[name] = a.transformSecurityScheme(ss)
		}
	}

	if len(in.Links) > 0 {
		comp.Links = make(map[string]*LinkV30, len(in.Links))
		for name, link := range in.Links {
			comp.Links[name] = a.transformLink(link)
		}
	}

	if len(in.Callbacks) > 0 {
		comp.Callbacks = make(map[string]*CallbackV30, len(in.Callbacks))
		for name, cb := range in.Callbacks {
			comp.Callbacks[name] = a.transformCallback(cb, warnings)
		}
	}

	// Warn about PathItems (3.1-only)
	if len(in.PathItems) > 0 {
		*warnings = append(*warnings, debug.NewWarning(debug.WarnDegradationPathItems, "#/components/pathItems", "pathItems in components are 3.1-only; dropped"))
	}

	return comp
}

func (a *AdapterV304) transformResponse(in *model.Response, warnings *debug.Warnings) *ResponseV30 {
	if in == nil {
		return nil
	}

	// Handle $ref case
	if in.Ref != "" {
		return &ResponseV30{Ref: in.Ref}
	}

	r := &ResponseV30{
		Description: in.Description,
		Extensions:  in.Extensions,
	}

	if len(in.Content) > 0 {
		r.Content = make(map[string]*MediaTypeV30, len(in.Content))
		for ct, mt := range in.Content {
			r.Content[ct] = a.transformMediaType(mt, warnings)
		}
	}

	if len(in.Headers) > 0 {
		r.Headers = make(map[string]*HeaderV30, len(in.Headers))
		for name, h := range in.Headers {
			r.Headers[name] = a.transformHeader(h, warnings)
		}
	}

	return r
}

func (a *AdapterV304) transformHeader(in *model.Header, warnings *debug.Warnings) *HeaderV30 {
	if in == nil {
		return nil
	}

	// Handle $ref case
	if in.Ref != "" {
		return &HeaderV30{Ref: in.Ref}
	}

	h := &HeaderV30{
		Description:     in.Description,
		Required:        in.Required,
		Deprecated:      in.Deprecated,
		AllowEmptyValue: in.AllowEmptyValue,
		Style:           in.Style,
		Explode:         in.Explode,
		Example:         in.Example,
		Extensions:      in.Extensions,
	}

	h.Schema = a.transformSchema(in.Schema, warnings)

	if len(in.Examples) > 0 {
		h.Examples = make(map[string]*ExampleV30, len(in.Examples))
		for k, ex := range in.Examples {
			h.Examples[k] = a.transformExample(ex, warnings)
		}
	}

	return h
}

func (a *AdapterV304) transformSecurityScheme(in *model.SecurityScheme) *SecuritySchemeV30 {
	if in == nil {
		return nil
	}

	// Handle $ref case
	if in.Ref != "" {
		return &SecuritySchemeV30{Ref: in.Ref}
	}

	out := &SecuritySchemeV30{
		Type:             in.Type,
		Description:      in.Description,
		Name:             in.Name,
		In:               in.In,
		Scheme:           in.Scheme,
		BearerFormat:     in.BearerFormat,
		OpenIDConnectURL: in.OpenIDConnectURL,
		Extensions:       in.Extensions,
	}

	return out
}

func (a *AdapterV304) transformLink(in *model.Link) *LinkV30 {
	if in == nil {
		return nil
	}

	// Handle $ref case
	if in.Ref != "" {
		return &LinkV30{Ref: in.Ref}
	}

	link := &LinkV30{
		OperationRef: in.OperationRef,
		OperationID:  in.OperationID,
		Parameters:   in.Parameters,
		RequestBody:  in.RequestBody,
		Description:  in.Description,
		Extensions:   in.Extensions,
	}

	if in.Server != nil {
		servers := a.transformServers([]model.Server{*in.Server})
		if len(servers) > 0 {
			link.Server = servers[0]
		}
	}

	return link
}

func (a *AdapterV304) transformCallback(in *model.Callback, warnings *debug.Warnings) *CallbackV30 {
	if in == nil {
		return nil
	}

	// Handle $ref case
	if in.Ref != "" {
		return &CallbackV30{Ref: in.Ref}
	}

	cb := &CallbackV30{
		PathItems:  make(map[string]*PathItemV30, len(in.PathItems)),
		Extensions: in.Extensions,
	}

	for path, item := range in.PathItems {
		cb.PathItems[path] = a.transformPathItem(item, warnings)
	}

	return cb
}

func (a *AdapterV304) transformResponses(in map[string]*model.Response, warnings *debug.Warnings) ResponsesV30 {
	if len(in) == 0 {
		return nil
	}

	responses := make(ResponsesV30, len(in))
	for code, response := range in {
		responses[code] = a.transformResponse(response, warnings)
	}

	return responses
}

//nolint:cyclop
func (a *AdapterV304) transformSchema(in *model.Schema, warnings *debug.Warnings) *SchemaV30 {
	if in == nil {
		return nil
	}

	// Handle $ref case
	if in.Ref != "" {
		return &SchemaV30{Ref: in.Ref}
	}

	out := &SchemaV30{
		Title:       in.Title,
		Description: in.Description,
		Format:      in.Format,
		Deprecated:  in.Deprecated,
		ReadOnly:    in.ReadOnly,
		WriteOnly:   in.WriteOnly,
		Nullable:    in.Nullable,
		Type:        in.Type,
		Extensions:  in.Extensions,
	}

	// Handle examples - prefer single example for 3.0
	if in.Example != nil {
		out.Example = in.Example
	} else if len(in.Examples) > 0 {
		out.Example = in.Examples[0] // Use first example for 3.0
		if len(in.Examples) > 1 {
			*warnings = append(*warnings, debug.NewWarning(debug.WarnDegradationMultipleExamples, "#/components/schemas/...", "multiple examples collapsed to first example only"))
		}
	}

	// Handle enum
	if len(in.Enum) > 0 {
		out.Enum = append([]any(nil), in.Enum...)
	}

	// Handle const (3.1 feature) - convert to enum
	if in.Const != nil {
		out.Enum = []any{in.Const}
		// Clear type to avoid conflicts (const value may not match schema type)
		out.Type = ""
		*warnings = append(*warnings, debug.NewWarning(debug.WarnDegradationConstToEnum, "#/components/schemas/...", "const converted to enum"))
	}

	// Handle numeric constraints
	out.MultipleOf = in.MultipleOf

	// Handle bounds - convert Bound structs to float64 and boolean flags
	if in.Minimum != nil {
		out.Minimum = &in.Minimum.Value
		out.ExclusiveMinimum = in.Minimum.Exclusive
	}
	if in.Maximum != nil {
		out.Maximum = &in.Maximum.Value
		out.ExclusiveMaximum = in.Maximum.Exclusive
	}

	// Handle string constraints
	out.MinLength = in.MinLength
	out.MaxLength = in.MaxLength
	out.Pattern = in.Pattern

	// Handle array constraints
	out.MinItems = in.MinItems
	out.MaxItems = in.MaxItems
	out.UniqueItems = in.UniqueItems
	out.Items = a.transformSchema(in.Items, warnings)

	// Handle object constraints
	if len(in.Properties) > 0 {
		out.Properties = make(map[string]*SchemaV30, len(in.Properties))
		for name, prop := range in.Properties {
			out.Properties[name] = a.transformSchema(prop, warnings)
		}
	}
	if len(in.Required) > 0 {
		out.Required = append([]string(nil), in.Required...)
	}
	out.MinProperties = in.MinProperties
	out.MaxProperties = in.MaxProperties

	// Handle additional properties
	if in.Additional != nil {
		if in.Additional.Allow != nil {
			out.AdditionalProperties = *in.Additional.Allow
		} else {
			out.AdditionalProperties = a.transformSchema(in.Additional.Schema, warnings)
		}
	}

	// Handle composition
	if len(in.AllOf) > 0 {
		out.AllOf = make([]*SchemaV30, 0, len(in.AllOf))
		for _, schema := range in.AllOf {
			out.AllOf = append(out.AllOf, a.transformSchema(schema, warnings))
		}
	}
	if len(in.AnyOf) > 0 {
		out.AnyOf = make([]*SchemaV30, 0, len(in.AnyOf))
		for _, schema := range in.AnyOf {
			out.AnyOf = append(out.AnyOf, a.transformSchema(schema, warnings))
		}
	}
	if len(in.OneOf) > 0 {
		out.OneOf = make([]*SchemaV30, 0, len(in.OneOf))
		for _, schema := range in.OneOf {
			out.OneOf = append(out.OneOf, a.transformSchema(schema, warnings))
		}
	}
	out.Not = a.transformSchema(in.Not, warnings)

	// Handle default value
	out.Default = in.Default

	// Warn about 3.1-only features that are dropped in 3.0
	if in.ContentEncoding != "" {
		*warnings = append(*warnings, debug.NewWarning(debug.WarnDegradationContentEncoding, "#/components/schemas/...", "contentEncoding dropped (3.1-only)"))
	}
	if in.ContentMediaType != "" {
		*warnings = append(*warnings, debug.NewWarning(debug.WarnDegradationContentMediaType, "#/components/schemas/...", "contentMediaType dropped (3.1-only)"))
	}
	if in.Unevaluated != nil {
		*warnings = append(*warnings, debug.NewWarning(debug.WarnDegradationUnevaluatedProperties, "#/components/schemas/...", "unevaluatedProperties dropped (3.1-only)"))
	}

	return out
}
