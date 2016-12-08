// Copyright 2016 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/googleapis/openapi-compiler/printer"
)

func (classes *ClassCollection) generateCompiler(packageName string, license string) string {
	code := printer.Code{}
	code.Print(license)
	code.Print("// THIS FILE IS AUTOMATICALLY GENERATED.")
	code.Print()
	code.Print("package %s", packageName)
	code.Print()
	code.Print("import (")
	imports := []string{
		"errors",
		"fmt",
		"encoding/json",
		"github.com/googleapis/openapi-compiler/helpers",
	}
	for _, filename := range imports {
		code.Print("\"" + filename + "\"")
	}
	code.Print(")")
	code.Print()

	code.Print("func Version() string {")
	code.Print("  return \"%s\"", packageName)
	code.Print("}")
	code.Print()

	classNames := classes.sortedClassNames()

	// constructors
	for _, className := range classNames {
		code.Print("func New%s(in interface{}) (*%s, error) {", className, className)

		classModel := classes.ClassModels[className]
		parentClassName := className

		if classModel.IsStringArray {
			code.Print("value, ok := in.(string)")
			code.Print("x := &TypeItem{}")
			code.Print("if ok {")
			code.Print("x.Value = make([]string, 0)")
			code.Print("x.Value = append(x.Value, value)")
			code.Print("} else {")
			code.Print("return nil, errors.New(fmt.Sprintf(\"unexpected value for string array: %%+v\", in))")
			code.Print("}")
			code.Print("return x, nil")
			code.Print("}")
			code.Print()
			continue
		}

		if classModel.IsItemArray {
			code.Print("m, ok := helpers.UnpackMap(in)")
			code.Print("if (!ok) {")
			code.Print("return nil, errors.New(fmt.Sprintf(\"unexpected value for item array: %%+v\", in))")
			code.Print("}")

			code.Print("x := &ItemsItem{}")
			code.Print("if ok {")
			code.Print("x.Schema = make([]*Schema, 0)")
			code.Print("y, err := NewSchema(m)")
			code.Print("if err != nil {return nil, err}")
			code.Print("x.Schema = append(x.Schema, y)")
			code.Print("} else {")
			code.Print("return nil, errors.New(fmt.Sprintf(\"unexpected value for item array: %%+v\", in))")
			code.Print("}")
			code.Print("return x, nil")
			code.Print("}")
			code.Print()
			continue
		}

		if classModel.IsBlob {
			code.Print("x := &Any{}")
			code.Print("bytes, _ := json.Marshal(in)")
			code.Print("x.Value = string(bytes)")
			code.Print("return x, nil")
			code.Print("}")
			code.Print()
			continue
		}

		if classModel.Name == "StringArray" {
			code.Print("a, ok := in.([]interface{})")
			code.Print("if ok {")
			code.Print("x := &StringArray{}")
			code.Print("x.Value = make([]string, 0)")
			code.Print("for _, s := range a {")
			code.Print("x.Value = append(x.Value, s.(string))")
			code.Print("}")
			code.Print("return x, nil")
			code.Print("} else {")
			code.Print("return nil, nil")
			code.Print("}")
			code.Print("}")
			code.Print()
			continue
		}
		code.Print("m, ok := helpers.UnpackMap(in)")
		code.Print("if (!ok) {")
		code.Print("return nil, errors.New(fmt.Sprintf(\"unexpected value for %s section: %%+v\", in))", className)
		code.Print("}")
		oneOfWrapper := classModel.OneOfWrapper

		if len(classModel.Required) > 0 {
			// verify that map includes all required keys
			keyString := ""
			sort.Strings(classModel.Required)
			for _, k := range classModel.Required {
				if keyString != "" {
					keyString += ","
				}
				keyString += "\""
				keyString += k
				keyString += "\""
			}
			code.Print("requiredKeys := []string{%s}", keyString)
			code.Print("if !helpers.MapContainsAllKeys(m, requiredKeys) {")
			code.Print("return nil, errors.New(\"%s does not contain all required properties (%s)\")",
				classModel.Name,
				strings.Replace(keyString, "\"", "'", -1))
			code.Print("}")
		}

		if !classModel.Open {
			// verify that map has no unspecified keys
			allowedKeys := make([]string, 0)
			for _, property := range classModel.Properties {
				if !property.Implicit {
					allowedKeys = append(allowedKeys, property.Name)
				}
			}
			sort.Strings(allowedKeys)
			allowedKeyString := ""
			for _, allowedKey := range allowedKeys {
				if allowedKeyString != "" {
					allowedKeyString += ","
				}
				allowedKeyString += "\""
				allowedKeyString += allowedKey
				allowedKeyString += "\""
			}
			allowedPatternString := ""
			if classModel.OpenPatterns != nil {
				for _, pattern := range classModel.OpenPatterns {
					if allowedPatternString != "" {
						allowedPatternString += ","
					}
					allowedPatternString += "\""
					allowedPatternString += pattern
					allowedPatternString += "\""
				}
			}
			// verify that map includes only allowed keys and patterns
			code.Print("allowedKeys := []string{%s}", allowedKeyString)
			code.Print("allowedPatterns := []string{%s}", allowedPatternString)
			code.Print("if !helpers.MapContainsOnlyKeysAndPatterns(m, allowedKeys, allowedPatterns) {")
			code.Print("return nil, errors.New(\nfmt.Sprintf(\"%s includes properties not in (%s) or (%s): %%+v\",\nhelpers.SortedKeysForMap(m)))",
				classModel.Name,
				strings.Replace(allowedKeyString, "\"", "'", -1),
				strings.Replace(allowedPatternString, "\"", "'", -1))
			code.Print("}")
		}

		code.Print("  x := &%s{}", className)

		var fieldNumber = 0
		for _, propertyModel := range classModel.Properties {
			propertyName := propertyModel.Name
			fieldNumber += 1
			propertyType := propertyModel.Type
			if propertyType == "int" {
				propertyType = "int64"
			}
			var displayName = propertyName
			if displayName == "$ref" {
				displayName = "_ref"
			}
			if displayName == "$schema" {
				displayName = "_schema"
			}
			displayName = camelCaseToSnakeCase(displayName)

			var line = fmt.Sprintf("%s %s = %d;", propertyType, displayName, fieldNumber)
			if propertyModel.Repeated {
				line = "repeated " + line
			}
			code.Print("// " + line)

			fieldName := strings.Title(propertyName)
			if propertyName == "$ref" {
				fieldName = "XRef"
			}

			classModel, classFound := classes.ClassModels[propertyType]
			if classFound && !classModel.IsPair {
				if propertyModel.Repeated {
					code.Print("v%d := helpers.MapValueForKey(m, \"%s\")", fieldNumber, propertyName)
					code.Print("if (v%d != nil) {", fieldNumber)
					code.Print("// repeated class %s", classModel.Name)
					code.Print("x.%s = make([]*%s, 0)", fieldName, classModel.Name)
					code.Print("a, ok := v%d.([]interface{})", fieldNumber)
					code.Print("if ok {")
					code.Print("for _, item := range a {")
					code.Print("y, err := New%s(item)", classModel.Name)
					code.Print("if err != nil {return nil, err}")
					code.Print("x.%s = append(x.%s, y)", fieldName, fieldName)
					code.Print("}")
					code.Print("}")
					code.Print("}")
				} else {
					if oneOfWrapper {
						code.Print("{")
						code.Print("// errors are ok here, they mean we just don't have the right subtype")
						code.Print("t, _ := New%s(m)", classModel.Name)
						code.Print("if t != nil {")
						code.Print("x.Oneof = &%s_%s{%s: t}", parentClassName, classModel.Name, classModel.Name)
						code.Print("}")
						code.Print("}")
					} else {
						code.Print("v%d := helpers.MapValueForKey(m, \"%s\")", fieldNumber, propertyName)
						code.Print("if (v%d != nil) {", fieldNumber)
						code.Print("var err error")
						code.Print("x.%s, err = New%s(v%d)", fieldName, classModel.Name, fieldNumber)
						code.Print("if err != nil {return nil, helpers.ExtendError(\"%s\", err)}", classModel.Name)
						code.Print("}")
					}
				}
			} else if propertyType == "string" {
				if propertyModel.Repeated {
					code.Print("v%d := helpers.MapValueForKey(m, \"%s\")", fieldNumber, propertyName)
					code.Print("if (v%d != nil) {", fieldNumber)
					code.Print("v, ok := v%d.([]interface{})", fieldNumber)
					code.Print("if ok {")
					code.Print("x.%s = helpers.ConvertInterfaceArrayToStringArray(v)", fieldName)
					code.Print("} else {")
					code.Print(" return nil, errors.New(fmt.Sprintf(\"unexpected value for %s property: %%+v\", in))", propertyName)
					code.Print("}")
					code.Print("}")
				} else {
					code.Print("v%d := helpers.MapValueForKey(m, \"%s\")", fieldNumber, propertyName)
					code.Print("if (v%d != nil) {", fieldNumber)
					code.Print("x.%s = v%d.(string)", fieldName, fieldNumber)
					code.Print("}")
				}
			} else if propertyType == "float" {
				code.Print("v%d := helpers.MapValueForKey(m, \"%s\")", fieldNumber, propertyName)
				code.Print("if (v%d != nil) {", fieldNumber)
				code.Print("x.%s = v%d.(float64)", fieldName, fieldNumber)
				code.Print("}")
			} else if propertyType == "int64" {
				code.Print("v%d := helpers.MapValueForKey(m, \"%s\")", fieldNumber, propertyName)
				code.Print("if (v%d != nil) {", fieldNumber)
				code.Print("x.%s = v%d.(int64)", fieldName, fieldNumber)
				code.Print("}")
			} else if propertyType == "bool" {
				code.Print("v%d := helpers.MapValueForKey(m, \"%s\")", fieldNumber, propertyName)
				code.Print("if (v%d != nil) {", fieldNumber)
				code.Print("x.%s = v%d.(bool)", fieldName, fieldNumber)
				code.Print("}")
			} else {
				mapTypeName := propertyModel.MapType
				isMap := mapTypeName != ""
				if isMap {
					code.Print("// MAP: %s %s", mapTypeName, propertyModel.Pattern)
					if mapTypeName == "string" {
						code.Print("x.%s = make([]*NamedString, 0)", fieldName)
					} else {
						code.Print("x.%s = make([]*Named%s, 0)", fieldName, mapTypeName)
					}
					code.Print("for _, item := range m {")
					code.Print("k := item.Key.(string)")
					code.Print("v := item.Value")
					if propertyModel.Pattern != "" {
						code.Print("if helpers.PatternMatches(\"%s\", k) {", propertyModel.Pattern)
					}
					code.Print("pair := &Named" + strings.Title(mapTypeName) + "{}")
					code.Print("pair.Name = k")
					if mapTypeName == "string" {
						code.Print("pair.Value = v.(string)")
					} else {
						code.Print("var err error")
						code.Print("pair.Value, err = New%v(v)", mapTypeName)
						code.Print("if err != nil {return nil, err}")
					}
					code.Print("x.%s = append(x.%s, pair)", fieldName, fieldName)
					if propertyModel.Pattern != "" {
						code.Print("}")
					}
					code.Print("}")
				} else {
					code.Print("// TODO: %s", propertyType)
				}
			}
		}
		code.Print("  return x, nil")
		code.Print("}\n")
	}

	// ResolveReferences() methods
	for _, className := range classNames {
		code.Print("func (m *%s) ResolveReferences(root string) (interface{}, error) {", className)
		//code.Print("  log.Printf(\"%s.ResolveReferences(%%+v)\", m)", className)

		classModel := classes.ClassModels[className]
		if classModel.OneOfWrapper {
			// call ResolveReferences on whatever is in the Oneof.
			for _, propertyModel := range classModel.Properties {
				propertyType := propertyModel.Type
				code.Print("{")
				code.Print("p, ok := m.Oneof.(*%s_%s)", className, propertyType)
				code.Print("if ok {")
				if propertyType == "JsonReference" { // Special case for OpenAPI
					code.Print("info, err := p.%s.ResolveReferences(root)", propertyType)
					code.Print("if err != nil {")
					code.Print("  return nil, err")
					code.Print("} else if info != nil {")
					code.Print("  n, err := New%s(info)", className)
					code.Print("  if err != nil {")
					code.Print("    return nil, err")
					code.Print("  } else if n != nil {")
					code.Print("    *m = *n")
					code.Print("    return nil, nil")
					code.Print("  }")
					code.Print("}")
				} else {
					code.Print("p.%s.ResolveReferences(root)", propertyType)
				}
				code.Print("}")
				code.Print("}")
			}
		} else {
			for _, propertyModel := range classModel.Properties {
				propertyName := propertyModel.Name
				var displayName = propertyName
				if displayName == "$ref" {
					displayName = "_ref"
				}
				if displayName == "$schema" {
					displayName = "_schema"
				}
				displayName = camelCaseToSnakeCase(displayName)

				fieldName := strings.Title(propertyName)
				if propertyName == "$ref" {
					fieldName = "XRef"
					code.Print("if m.XRef != \"\" {")
					//code.Print("log.Printf(\"%s reference to resolve %%+v\", m.XRef)", className)
					code.Print("info := helpers.ReadInfoForRef(root, m.XRef)")
					//code.Print("log.Printf(\"%%+v\", info)")

					if len(classModel.Properties) > 1 {
						code.Print("if info != nil {")
						code.Print("replacement, _ := New%s(info)", className)
						code.Print("*m = *replacement")
						code.Print("return m.ResolveReferences(root)")
						code.Print("}")
					} else {
						code.Print("return info, nil")
					}

					code.Print("return info, nil")
					code.Print("}")
				}

				if !propertyModel.Repeated {
					propertyType := propertyModel.Type
					classModel, classFound := classes.ClassModels[propertyType]
					if classFound && !classModel.IsPair {
						code.Print("if m.%s != nil {", fieldName)
						code.Print("m.%s.ResolveReferences(root)", fieldName)
						code.Print("}")
					}
				} else {
					propertyType := propertyModel.Type
					_, classFound := classes.ClassModels[propertyType]
					if classFound {
						code.Print("for _, item := range m.%s {", fieldName)
						code.Print("if item != nil {")
						code.Print("item.ResolveReferences(root)")
						code.Print("}")
						code.Print("}")
					}

				}
			}
		}
		code.Print("  return nil, nil")
		code.Print("}\n")
	}

	return code.String()
}
