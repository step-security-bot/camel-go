package wasm

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	"gopkg.in/yaml.v3"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	camelerrors "github.com/lburgazzoli/camel-go/pkg/core/errors"
	"github.com/lburgazzoli/camel-go/pkg/util/registry"
	"github.com/lburgazzoli/camel-go/pkg/wasm"
)

func New() *Wasm {
	w := Wasm{}
	return &w
}

func NewWithValue(value string) *Wasm {
	w := Wasm{}
	_ = w.UnmarshalText([]byte(value))

	return &w
}

type Definition struct {
	Path  string `yaml:"path"`
	Image string `yaml:"image,omitempty"`
}

type Wasm struct {
	Definition `yaml:",inline"`
}

func (l *Wasm) UnmarshalYAML(value *yaml.Node) error {
	switch value.Kind {
	case yaml.ScalarNode:
		return l.UnmarshalText([]byte(value.Value))
	case yaml.MappingNode:
		return value.Decode(&l.Definition)
	default:
		return fmt.Errorf("unsupported node kind: %v (line: %d, column: %d)", value.Kind, value.Line, value.Column)
	}
}

func (l *Wasm) UnmarshalText(text []byte) error {
	in := string(text)
	parts := strings.Split(in, "?")

	switch len(parts) {
	case 1:
		l.Path = parts[0]
	case 2:
		l.Image = parts[0]
		l.Path = parts[1]
	default:
		return camelerrors.InvalidParameterf("wasm", "unsupported wasm reference '%s'", in)
	}

	return nil
}

func (l *Wasm) Processor(ctx context.Context, _ camel.Context) (camel.Processor, error) {
	if l.Path == "" {
		return nil, camelerrors.MissingParameterf("wasm.path", "failure configuring wasm processor")
	}

	rootPath := ""

	if l.Image != "" {
		fp, err := registry.Pull(ctx, l.Image)
		if err != nil {
			return nil, err
		}

		rootPath = fp
	}

	defer func() {
		if rootPath != "" {
			_ = os.RemoveAll(rootPath)
		}
	}()

	r, err := wasm.NewRuntime(ctx, wasm.Options{})
	if err != nil {
		return nil, err
	}

	f, err := r.Load(ctx, path.Join(rootPath, l.Path))
	if err != nil {
		return nil, err
	}

	p := func(ctx context.Context, m camel.Message) error {
		result, err := f.Invoke(ctx, m)
		if err != nil {
			return err
		}

		return result.CopyTo(m)
	}

	return p, nil
}
