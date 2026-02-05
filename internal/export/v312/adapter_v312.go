package v312

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/talav/openapi/debug"
	"github.com/talav/openapi/internal/model"
)

//go:embed schema_v312.json
var schemaV312JSON []byte

type AdapterV312 struct{}

func (a *AdapterV312) Version() string {
	return "3.1.2"
}

func (a *AdapterV312) SchemaJSON() []byte {
	return schemaV312JSON
}

func (a *AdapterV312) View(spec *model.Spec) (any, debug.Warnings, error) {
	if spec == nil {
		return nil, nil, fmt.Errorf("nil spec")
	}

	var warnings debug.Warnings

	result := &ViewV312{
		OpenAPI:      a.Version(),
		Info:         a.transformInfo(spec.Info),
		Servers:      a.transformServers(spec.Servers),
		Paths:        a.transformPaths(spec.Paths, &warnings),
		Components:   a.transformComponents(spec.Components, &warnings),
		Security:     a.transformSecurity(spec.Security),
		Tags:         a.transformTags(spec.Tags),
		ExternalDocs: a.transformExternalDocs(spec.ExternalDocs),
		Webhooks:     a.transformWebhooks(spec.Webhooks, &warnings),
		Extensions:   spec.Extensions,
	}

	if err := validateViewV312(result); err != nil {
		return nil, nil, err
	}

	return result, warnings, nil
}

// validateViewV312 validates a ViewV312 instance according to OpenAPI 3.1.2 requirements.
func validateViewV312(result *ViewV312) error {
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

	for key := range result.Extensions {
		if err := validateExtensionKey(key, "root"); err != nil {
			return err
		}
	}

	for key := range result.Info.Extensions {
		if err := validateExtensionKey(key, "info"); err != nil {
			return err
		}
	}

	return nil
}

func validateExtensionKey(key, placement string) error {
	if !strings.HasPrefix(key, "x-") {
		return fmt.Errorf("openapi: %s extension key must start with 'x-': %s", placement, key)
	}
	if strings.HasPrefix(key, "x-oai-") || strings.HasPrefix(key, "x-oas-") {
		return fmt.Errorf("openapi: %s extension key uses reserved prefix (x-oai- or x-oas-): %s", placement, key)
	}

	return nil
}

func (a *AdapterV312) transformInfo(in model.Info) *InfoV31 {
	info := &InfoV31{
		Title:          in.Title,
		Summary:        in.Summary,
		Description:    in.Description,
		TermsOfService: in.TermsOfService,
		Version:        in.Version,
		Extensions:     in.Extensions,
	}

	if in.Contact != nil {
		info.Contact = &ContactV31{
			Name:       in.Contact.Name,
			URL:        in.Contact.URL,
			Email:      in.Contact.Email,
			Extensions: in.Contact.Extensions,
		}
	}

	if in.License != nil {
		info.License = &LicenseV31{
			Name:       in.License.Name,
			Identifier: in.License.Identifier,
			URL:        in.License.URL,
			Extensions: in.License.Extensions,
		}
	}

	return info
}

func (a *AdapterV312) transformServers(in []model.Server) []*ServerV31 {
	if len(in) == 0 {
		return nil
	}

	servers := make([]*ServerV31, 0, len(in))
	for _, s := range in {
		server := &ServerV31{
			URL:         s.URL,
			Description: s.Description,
			Extensions:  s.Extensions,
		}

		if len(s.Variables) > 0 {
			server.Variables = make(map[string]*ServerVariableV31, len(s.Variables))
			for name, v := range s.Variables {
				server.Variables[name] = &ServerVariableV31{
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

func (a *AdapterV312) transformTags(in []model.Tag) []*TagV31 {
	if len(in) == 0 {
		return nil
	}

	tags := make([]*TagV31, 0, len(in))
	for _, t := range in {
		tag := &TagV31{
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

func (a *AdapterV312) transformSecurity(in []model.SecurityRequirement) []SecurityRequirementV31 {
	if len(in) == 0 {
		return nil
	}

	security := make([]SecurityRequirementV31, 0, len(in))
	for _, s := range in {
		security = append(security, SecurityRequirementV31(s))
	}

	return security
}

func (a *AdapterV312) transformExternalDocs(in *model.ExternalDocs) *ExternalDocsV31 {
	if in == nil {
		return nil
	}

	return &ExternalDocsV31{
		Description: in.Description,
		URL:         in.URL,
		Extensions:  in.Extensions,
	}
}

func (a *AdapterV312) transformPaths(in map[string]*model.PathItem, warnings *debug.Warnings) PathsV31 {
	if len(in) == 0 {
		return make(PathsV31)
	}

	paths := make(PathsV31, len(in))
	for path, item := range in {
		paths[path] = a.transformPathItem(item, warnings)
	}

	return paths
}

func (a *AdapterV312) transformWebhooks(in map[string]*model.PathItem, warnings *debug.Warnings) PathsV31 {
	if len(in) == 0 {
		return nil
	}

	webhooks := make(map[string]*PathItemV31, len(in))
	for name, item := range in {
		webhooks[name] = a.transformPathItem(item, warnings)
	}

	return webhooks
}

func (a *AdapterV312) transformPathItem(in *model.PathItem, warnings *debug.Warnings) *PathItemV31 {
	if in == nil {
		return nil
	}

	// Handle $ref case
	if in.Ref != "" {
		return &PathItemV31{Ref: in.Ref}
	}

	item := &PathItemV31{
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

	// Transform Servers
	if len(in.Servers) > 0 {
		item.Servers = a.transformServers(in.Servers)
	}

	return item
}

func (a *AdapterV312) transformParameters(in []model.Parameter, warnings *debug.Warnings) []*ParameterV31 {
	out := make([]*ParameterV31, 0, len(in))
	for _, param := range in {
		p := a.transformParameter(param, warnings)
		out = append(out, &p)
	}

	return out
}

func (a *AdapterV312) transformParameter(in model.Parameter, warnings *debug.Warnings) ParameterV31 {
	// Handle $ref case
	if in.Ref != "" {
		return ParameterV31{Ref: in.Ref}
	}

	param := ParameterV31{
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
		param.Examples = make(map[string]*ExampleV31, len(in.Examples))
		for k, v := range in.Examples {
			param.Examples[k] = a.transformExample(v, warnings)
		}
	}

	if len(in.Content) > 0 {
		param.Content = make(map[string]*MediaTypeV31, len(in.Content))
		for ct, mt := range in.Content {
			param.Content[ct] = a.transformMediaType(mt, warnings)
		}
	}

	return param
}

func (a *AdapterV312) transformExample(in *model.Example, warnings *debug.Warnings) *ExampleV31 {
	if in == nil {
		return nil
	}

	// Handle $ref case
	if in.Ref != "" {
		return &ExampleV31{Ref: in.Ref}
	}

	// Per OpenAPI spec, value and externalValue are mutually exclusive
	out := &ExampleV31{
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

func (a *AdapterV312) transformOperation(in *model.Operation, warnings *debug.Warnings) *OperationV31 {
	if in == nil {
		return nil
	}

	op := &OperationV31{
		Tags:        append([]string(nil), in.Tags...),
		Summary:     in.Summary,
		Description: in.Description,
		OperationID: in.OperationID,
		Deprecated:  in.Deprecated,
		Extensions:  in.Extensions,
	}

	if in.ExternalDocs != nil {
		op.ExternalDocs = a.transformExternalDocs(in.ExternalDocs)
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

	if len(in.Callbacks) > 0 {
		op.Callbacks = make(map[string]*CallbackV31, len(in.Callbacks))
		for name, cb := range in.Callbacks {
			op.Callbacks[name] = a.transformCallback(cb, warnings)
		}
	}

	return op
}

func (a *AdapterV312) transformRequestBody(in *model.RequestBody, warnings *debug.Warnings) *RequestBodyV31 {
	if in == nil {
		return nil
	}

	// Handle $ref case
	if in.Ref != "" {
		return &RequestBodyV31{Ref: in.Ref}
	}

	rb := &RequestBodyV31{
		Description: in.Description,
		Required:    in.Required,
		Extensions:  in.Extensions,
	}

	if len(in.Content) > 0 {
		rb.Content = make(map[string]*MediaTypeV31, len(in.Content))
		for ct, mt := range in.Content {
			rb.Content[ct] = a.transformMediaType(mt, warnings)
		}
	}

	return rb
}

func (a *AdapterV312) transformMediaType(in *model.MediaType, warnings *debug.Warnings) *MediaTypeV31 {
	if in == nil {
		return nil
	}

	mt := &MediaTypeV31{
		Example:    in.Example,
		Extensions: in.Extensions,
	}

	mt.Schema = a.transformSchema(in.Schema, warnings)

	if len(in.Examples) > 0 {
		mt.Examples = make(map[string]*ExampleV31, len(in.Examples))
		for k, ex := range in.Examples {
			mt.Examples[k] = a.transformExample(ex, warnings)
		}
	}

	if len(in.Encoding) > 0 {
		mt.Encoding = make(map[string]*EncodingV31, len(in.Encoding))
		for name, enc := range in.Encoding {
			mt.Encoding[name] = a.transformEncoding(enc, warnings)
		}
	}

	return mt
}

func (a *AdapterV312) transformEncoding(in *model.Encoding, warnings *debug.Warnings) *EncodingV31 {
	if in == nil {
		return nil
	}

	enc := &EncodingV31{
		ContentType:   in.ContentType,
		Style:         in.Style,
		Explode:       in.Explode,
		AllowReserved: in.AllowReserved,
		Extensions:    in.Extensions,
	}

	if len(in.Headers) > 0 {
		enc.Headers = make(map[string]*HeaderV31, len(in.Headers))
		for name, h := range in.Headers {
			enc.Headers[name] = a.transformHeader(h, warnings)
		}
	}

	return enc
}

//nolint:cyclop,gocognit
func (a *AdapterV312) transformComponents(in *model.Components, warnings *debug.Warnings) *ComponentsV31 {
	if in == nil {
		return nil
	}

	comp := &ComponentsV31{
		Extensions: in.Extensions,
	}

	if len(in.Schemas) > 0 {
		comp.Schemas = make(map[string]*SchemaV31, len(in.Schemas))
		for name, schema := range in.Schemas {
			comp.Schemas[name] = a.transformSchema(schema, warnings)
		}
	}

	if len(in.Responses) > 0 {
		comp.Responses = make(map[string]*ResponseV31, len(in.Responses))
		for name, r := range in.Responses {
			comp.Responses[name] = a.transformResponse(r, warnings)
		}
	}

	if len(in.Parameters) > 0 {
		comp.Parameters = make(map[string]*ParameterV31, len(in.Parameters))
		for name, param := range in.Parameters {
			pv := a.transformParameter(*param, warnings)
			comp.Parameters[name] = &pv
		}
	}

	if len(in.Examples) > 0 {
		comp.Examples = make(map[string]*ExampleV31, len(in.Examples))
		for name, ex := range in.Examples {
			comp.Examples[name] = a.transformExample(ex, warnings)
		}
	}

	if len(in.RequestBodies) > 0 {
		comp.RequestBodies = make(map[string]*RequestBodyV31, len(in.RequestBodies))
		for name, rb := range in.RequestBodies {
			comp.RequestBodies[name] = a.transformRequestBody(rb, warnings)
		}
	}

	if len(in.Headers) > 0 {
		comp.Headers = make(map[string]*HeaderV31, len(in.Headers))
		for name, h := range in.Headers {
			comp.Headers[name] = a.transformHeader(h, warnings)
		}
	}

	if len(in.SecuritySchemes) > 0 {
		comp.SecuritySchemes = make(map[string]*SecuritySchemeV31, len(in.SecuritySchemes))
		for name, ss := range in.SecuritySchemes {
			comp.SecuritySchemes[name] = a.transformSecurityScheme(ss)
		}
	}

	if len(in.Links) > 0 {
		comp.Links = make(map[string]*LinkV31, len(in.Links))
		for name, link := range in.Links {
			comp.Links[name] = a.transformLink(link)
		}
	}

	if len(in.Callbacks) > 0 {
		comp.Callbacks = make(map[string]*CallbackV31, len(in.Callbacks))
		for name, cb := range in.Callbacks {
			comp.Callbacks[name] = a.transformCallback(cb, warnings)
		}
	}

	// PathItems are supported in 3.1.2
	if len(in.PathItems) > 0 {
		comp.PathItems = make(map[string]*PathItemV31, len(in.PathItems))
		for name, pi := range in.PathItems {
			comp.PathItems[name] = a.transformPathItem(pi, warnings)
		}
	}

	return comp
}

func (a *AdapterV312) transformResponse(in *model.Response, warnings *debug.Warnings) *ResponseV31 {
	if in == nil {
		return nil
	}

	// Handle $ref case
	if in.Ref != "" {
		return &ResponseV31{Ref: in.Ref}
	}

	r := &ResponseV31{
		Description: in.Description,
		Extensions:  in.Extensions,
	}

	if len(in.Content) > 0 {
		r.Content = make(map[string]*MediaTypeV31, len(in.Content))
		for ct, mt := range in.Content {
			r.Content[ct] = a.transformMediaType(mt, warnings)
		}
	}

	if len(in.Headers) > 0 {
		r.Headers = make(map[string]*HeaderV31, len(in.Headers))
		for name, h := range in.Headers {
			r.Headers[name] = a.transformHeader(h, warnings)
		}
	}

	if len(in.Links) > 0 {
		r.Links = make(map[string]*LinkV31, len(in.Links))
		for name, link := range in.Links {
			r.Links[name] = a.transformLink(link)
		}
	}

	return r
}

func (a *AdapterV312) transformHeader(in *model.Header, warnings *debug.Warnings) *HeaderV31 {
	if in == nil {
		return nil
	}

	// Handle $ref case
	if in.Ref != "" {
		return &HeaderV31{Ref: in.Ref}
	}

	h := &HeaderV31{
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
		h.Examples = make(map[string]*ExampleV31, len(in.Examples))
		for k, ex := range in.Examples {
			h.Examples[k] = a.transformExample(ex, warnings)
		}
	}

	if len(in.Content) > 0 {
		h.Content = make(map[string]*MediaTypeV31, len(in.Content))
		for ct, mt := range in.Content {
			h.Content[ct] = a.transformMediaType(mt, warnings)
		}
	}

	return h
}

func (a *AdapterV312) transformSecurityScheme(in *model.SecurityScheme) *SecuritySchemeV31 {
	if in == nil {
		return nil
	}

	// Handle $ref case
	if in.Ref != "" {
		return &SecuritySchemeV31{Ref: in.Ref}
	}

	out := &SecuritySchemeV31{
		Type:             in.Type,
		Description:      in.Description,
		Name:             in.Name,
		In:               in.In,
		Scheme:           in.Scheme,
		BearerFormat:     in.BearerFormat,
		OpenIDConnectURL: in.OpenIDConnectURL,
		Extensions:       in.Extensions,
	}

	if in.Flows != nil {
		out.Flows = a.transformOAuthFlows(in.Flows)
	}

	return out
}

func (a *AdapterV312) transformOAuthFlows(in *model.OAuthFlows) *OAuthFlowsV31 {
	if in == nil {
		return nil
	}

	flows := &OAuthFlowsV31{
		Extensions: in.Extensions,
	}

	if in.Implicit != nil {
		flows.Implicit = a.transformOAuthFlow(in.Implicit)
	}
	if in.Password != nil {
		flows.Password = a.transformOAuthFlow(in.Password)
	}
	if in.ClientCredentials != nil {
		flows.ClientCredentials = a.transformOAuthFlow(in.ClientCredentials)
	}
	if in.AuthorizationCode != nil {
		flows.AuthorizationCode = a.transformOAuthFlow(in.AuthorizationCode)
	}

	return flows
}

func (a *AdapterV312) transformOAuthFlow(in *model.OAuthFlow) *OAuthFlowV31 {
	if in == nil {
		return nil
	}

	return &OAuthFlowV31{
		AuthorizationURL: in.AuthorizationURL,
		TokenURL:         in.TokenURL,
		RefreshURL:       in.RefreshURL,
		Scopes:           in.Scopes,
		Extensions:       in.Extensions,
	}
}

func (a *AdapterV312) transformLink(in *model.Link) *LinkV31 {
	if in == nil {
		return nil
	}

	// Handle $ref case
	if in.Ref != "" {
		return &LinkV31{Ref: in.Ref}
	}

	link := &LinkV31{
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

func (a *AdapterV312) transformCallback(in *model.Callback, warnings *debug.Warnings) *CallbackV31 {
	if in == nil {
		return nil
	}

	// Handle $ref case
	if in.Ref != "" {
		return &CallbackV31{Ref: in.Ref}
	}

	cb := &CallbackV31{
		PathItems:  make(map[string]*PathItemV31, len(in.PathItems)),
		Extensions: in.Extensions,
	}

	for path, item := range in.PathItems {
		cb.PathItems[path] = a.transformPathItem(item, warnings)
	}

	return cb
}

func (a *AdapterV312) transformResponses(in map[string]*model.Response, warnings *debug.Warnings) map[string]*ResponseV31 {
	if len(in) == 0 {
		return nil
	}

	responses := make(map[string]*ResponseV31, len(in))
	for code, response := range in {
		responses[code] = a.transformResponse(response, warnings)
	}

	return responses
}

//nolint:cyclop,gocognit,gocyclo,unparam
func (a *AdapterV312) transformSchema(in *model.Schema, warnings *debug.Warnings) *SchemaV31 {
	if in == nil {
		return nil
	}

	// Handle $ref case
	if in.Ref != "" {
		return &SchemaV31{Ref: in.Ref}
	}

	out := &SchemaV31{
		Title:            in.Title,
		Description:      in.Description,
		Format:           in.Format,
		Deprecated:       in.Deprecated,
		ReadOnly:         in.ReadOnly,
		WriteOnly:        in.WriteOnly,
		ContentEncoding:  in.ContentEncoding,
		ContentMediaType: in.ContentMediaType,
		Extensions:       in.Extensions,
	}

	// Handle type - in 3.1.2, nullable is represented as type: ["T", "null"]
	//nolint:gocritic
	if in.Nullable && in.Type != "" {
		// Convert to array type with null
		out.Type = []any{in.Type, "null"}
	} else if in.Nullable && in.Type == "" {
		// If no type specified but nullable, use ["null"]
		out.Type = []any{"null"}
	} else if in.Type != "" {
		out.Type = in.Type
	}

	// Handle examples - 3.1.2 supports both single example and examples array
	if in.Example != nil {
		out.Example = in.Example
	}
	if len(in.Examples) > 0 {
		out.Examples = append([]any(nil), in.Examples...)
	}

	// Handle enum
	if len(in.Enum) > 0 {
		out.Enum = append([]any(nil), in.Enum...)
	}

	// Handle const (3.1.2 feature)
	if in.Const != nil {
		out.Const = in.Const
	}

	// Handle numeric constraints
	out.MultipleOf = in.MultipleOf

	// Handle bounds - in 3.1.2, exclusive bounds are numbers, not booleans
	if in.Minimum != nil {
		if in.Minimum.Exclusive {
			out.ExclusiveMinimum = &in.Minimum.Value
		} else {
			out.Minimum = &in.Minimum.Value
		}
	}
	if in.Maximum != nil {
		if in.Maximum.Exclusive {
			out.ExclusiveMaximum = &in.Maximum.Value
		} else {
			out.Maximum = &in.Maximum.Value
		}
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
		out.Properties = make(map[string]*SchemaV31, len(in.Properties))
		for name, prop := range in.Properties {
			out.Properties[name] = a.transformSchema(prop, warnings)
		}
	}
	if len(in.Required) > 0 {
		out.Required = append([]string(nil), in.Required...)
	}
	out.MinProperties = in.MinProperties
	out.MaxProperties = in.MaxProperties

	// Handle pattern properties (3.1.2 feature)
	if len(in.PatternProps) > 0 {
		out.PatternProperties = make(map[string]*SchemaV31, len(in.PatternProps))
		for pattern, schema := range in.PatternProps {
			out.PatternProperties[pattern] = a.transformSchema(schema, warnings)
		}
	}

	// Handle additional properties
	if in.Additional != nil {
		if in.Additional.Allow != nil {
			out.AdditionalProperties = *in.Additional.Allow
		} else {
			out.AdditionalProperties = a.transformSchema(in.Additional.Schema, warnings)
		}
	}

	// Handle unevaluated properties (3.1.2 feature)
	if in.Unevaluated != nil {
		out.UnevaluatedProperties = a.transformSchema(in.Unevaluated, warnings)
	}

	// Handle composition
	if len(in.AllOf) > 0 {
		out.AllOf = make([]*SchemaV31, 0, len(in.AllOf))
		for _, schema := range in.AllOf {
			out.AllOf = append(out.AllOf, a.transformSchema(schema, warnings))
		}
	}
	if len(in.AnyOf) > 0 {
		out.AnyOf = make([]*SchemaV31, 0, len(in.AnyOf))
		for _, schema := range in.AnyOf {
			out.AnyOf = append(out.AnyOf, a.transformSchema(schema, warnings))
		}
	}
	if len(in.OneOf) > 0 {
		out.OneOf = make([]*SchemaV31, 0, len(in.OneOf))
		for _, schema := range in.OneOf {
			out.OneOf = append(out.OneOf, a.transformSchema(schema, warnings))
		}
	}
	out.Not = a.transformSchema(in.Not, warnings)

	// Handle default value
	out.Default = in.Default

	// Handle discriminator
	if in.Discriminator != nil {
		out.Discriminator = &DiscriminatorV31{
			PropertyName: in.Discriminator.PropertyName,
			Mapping:      in.Discriminator.Mapping,
		}
	}

	// Handle XML
	if in.XML != nil {
		out.XML = &XMLV31{
			Name:      in.XML.Name,
			Namespace: in.XML.Namespace,
			Prefix:    in.XML.Prefix,
			Attribute: in.XML.Attribute,
			Wrapped:   in.XML.Wrapped,
		}
	}

	// Handle external docs
	if in.ExternalDocs != nil {
		out.ExternalDocs = a.transformExternalDocs(in.ExternalDocs)
	}

	return out
}
