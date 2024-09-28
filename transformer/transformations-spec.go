package transformer

import (
	"github.com/diegosz/go-archetype/inputs"
	"github.com/diegosz/go-archetype/operations"
)

const (
	TransformationTypeInclude = "include"
	TransformationTypeReplace = "replace"
	TransformationTypeRename  = "rename"
)

type transformationsSpec struct {
	Ignore          []string             `yaml:"ignore"`
	Inputs          []inputs.InputSpec   `yaml:"inputs"`
	Transformations []transformationSpec `yaml:"transformations"`
	Before          operations.Spec      `yaml:"before"`
	After           operations.Spec      `yaml:"after"`
}

type transformationSpec struct {
	Name         string   `yaml:"name"`
	Type         string   `yaml:"type"`
	Pattern      string   `yaml:"pattern"`
	Replacement  string   `yaml:"replacement"`
	Files        []string `yaml:"files"`
	Condition    string   `yaml:"condition"`
	RegionMarker string   `yaml:"region_marker"`
}
