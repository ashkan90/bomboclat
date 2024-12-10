package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"sort"
	"strings"
	"text/template"
)

type Generator struct {
	config  GeneratorConfig
	parser  *Parser
	builder *CodeBuilder
}

func NewGenerator(config GeneratorConfig) (*Generator, error) {
	parser, err := NewParser(config.InputPath)
	if err != nil {
		return nil, err
	}

	return &Generator{
		config:  config,
		parser:  parser,
		builder: &CodeBuilder{buf: &bytes.Buffer{}},
	}, nil
}

func (g *Generator) Generate() error {
	interfaces, err := g.parser.ParseInterface(g.config.Interfaces...)
	if err != nil {
		return fmt.Errorf("parsing interfaces: %w", err)
	}

	features, err := g.validateAndGetAllFeatures()
	if err != nil {
		return fmt.Errorf("validating features: %w", err)
	}

	customTypes, _, cErr := g.parser.CollectTypes(interfaces)
	if cErr != nil {
		return fmt.Errorf("collection types: %w", cErr)
	}

	var allTemplates []string
	imports := make(map[string]bool)

	for _, feature := range features {
		if config, ok := featureConfigs[feature]; ok {
			allTemplates = append(allTemplates, config.Template)
			for _, imp := range config.RequiredImports {
				imports[imp] = true
			}
		}
	}

	// Template data oluştur
	data := TemplateData{
		Package:    g.config.Package,
		Interfaces: interfaces,
		Imports:    Keys(imports),
		Types:      customTypes,
		Features:   features,
	}

	// Tüm template'leri birleştir ve işle
	combined := strings.Join(allTemplates, "\n\n")
	tmpl, err := template.New("generator").
		Funcs(g.templateFuncs()).
		Parse(combined)
	if err != nil {
		return fmt.Errorf("parsing template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("executing template: %w", err)
	}

	// Format ve yaz
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("formatting generated code: %w", err)
	}

	return os.WriteFile(g.config.OutputPath, formatted, 0644)
}

func (g *Generator) generateInterface(iface Interface) error {
	// Generate spawn types
	if err := g.generateSpawnTypes(iface); err != nil {
		return err
	}

	// Generate spawn functions
	if err := g.generateSpawnFuncs(iface); err != nil {
		return err
	}

	// Generate execution handlers
	if err := g.generateHandlers(iface); err != nil {
		return err
	}

	return nil
}

func (g *Generator) generateSpawnTypes(iface Interface) error {
	for _, method := range iface.Methods {
		g.builder.Printf("// %sSpawnOption represents spawn configuration for %s.%s\n",
			method.Name, iface.Name, method.Name)

		g.builder.Printf("type %sSpawnOption struct {\n", method.Name)
		g.builder.Printf("\tfn %s\n", method.FuncType())
		g.builder.Printf("\targs []interface{}\n")
		g.builder.Printf("\tconfig SpawnConfig\n")
		g.builder.Printf("}\n\n")

		g.builder.Printf("func (o %sSpawnOption) opt() {}\n\n", method.Name)
	}
	return nil
}

func (g *Generator) generateSpawnFuncs(iface Interface) error {
	for _, method := range iface.Methods {
		g.builder.Printf("// %sSpawn creates a spawn option for %s.%s\n",
			method.Name, iface.Name, method.Name)

		g.builder.Printf("func %sSpawn(fn %s, args ...interface{}) Option {\n",
			method.Name, method.FuncType())

		g.builder.Printf("\treturn &%sSpawnOption{\n", method.Name)
		g.builder.Printf("\t\tfn: fn,\n")
		g.builder.Printf("\t\targs: args,\n")
		g.builder.Printf("\t\tconfig: DefaultSpawnConfig,\n")
		g.builder.Printf("\t}\n")
		g.builder.Printf("}\n\n")
	}
	return nil
}

func (g *Generator) generateHandlers(iface Interface) error {
	for _, method := range iface.Methods {
		g.builder.Printf("type %sHandler struct {}\n\n", method.Name)

		g.builder.Printf("func (h *%sHandler) Execute(ctx context.Context, opt Option) (interface{}, error) {\n", method.Name)
		g.builder.Printf("\to, ok := opt.(*%sSpawnOption)\n", method.Name)
		g.builder.Printf("\tif !ok {\n")
		g.builder.Printf("\t\treturn nil, fmt.Errorf(\"invalid option type: %%T\", opt)\n")
		g.builder.Printf("\t}\n\n")

		// Parameter validation
		g.builder.Printf("\tif len(o.args) != %d {\n", len(method.Params))
		g.builder.Printf("\t\treturn nil, fmt.Errorf(\"expected %%d arguments, got %%d\", %d, len(o.args))\n", len(method.Params))
		g.builder.Printf("\t}\n\n")

		// Type assertions for arguments
		for i, param := range method.Params {
			g.builder.Printf("\targ%d, ok := o.args[%d].(%s)\n", i, i, param.Type)
			g.builder.Printf("\tif !ok {\n")
			g.builder.Printf("\t\treturn nil, fmt.Errorf(\"argument %%d must be %s\", %d)\n", i, param.Type, i)
			g.builder.Printf("\t}\n\n")
		}

		// Function call
		g.builder.Printf("\treturn o.fn(")
		for i := range method.Params {
			if i > 0 {
				g.builder.Printf(", ")
			}
			g.builder.Printf("arg%d", i)
		}
		g.builder.Printf(")\n")
		g.builder.Printf("}\n\n")
	}
	return nil
}

func (g *Generator) validateAndGetAllFeatures() ([]Feature, error) {
	requested := make(map[Feature]bool)

	// Base feature'ı her zaman ekle
	requested[FeatureBase] = true

	// Kullanıcının istediği feature'ları ekle
	for _, f := range g.config.Features {
		feat := Feature(f)
		if _, ok := featureConfigs[feat]; !ok {
			return nil, fmt.Errorf("unknown feature: %s", f)
		}
		requested[feat] = true
	}

	// Bağımlılıkları ekle
	allFeatures := make(map[Feature]bool)
	for f := range requested {
		allFeatures[f] = true
		for _, dep := range featureDependencies[f] {
			allFeatures[dep] = true
		}
	}

	// Feature'ları sırala (deterministic output için)
	var result []Feature
	for f := range allFeatures {
		result = append(result, f)
	}
	sort.Slice(result, func(i, j int) bool {
		return string(result[i]) < string(result[j])
	})

	return result, nil
}

func (g *Generator) templateFuncs() template.FuncMap {
	return templateFuncs
}

var templateFuncs = template.FuncMap{
	"last": func(i int, slice []Parameter) bool {
		return i == len(slice)-1
	},
	"lastIndex": func(slice []Parameter) int {
		return len(slice) - 1
	},
	"hasError": func(returns []Parameter) bool {
		for _, r := range returns {
			if r.Type == "error" {
				return true
			}
		}
		return false
	},
	"isError": func(param Parameter) bool {
		return param.Type == "error"
	},
	"nonErrorReturns": func(returns []Parameter) []int {
		var indices []int
		for i, r := range returns {
			if r.Type != "error" {
				indices = append(indices, i)
			}
		}
		return indices
	},
	"hasAnyFeature": func(features []Feature, required ...string) bool {
		featureMap := make(map[string]bool)
		for _, f := range features {
			featureMap[string(f)] = true
		}

		for _, r := range required {
			if featureMap[r] {
				return true
			}
		}
		return false
	},
	"hasFeature": func(features []Feature, feature string) bool {
		for _, f := range features {
			if string(f) == feature {
				return true
			}
		}
		return false
	},
}
