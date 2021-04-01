package template

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"

	radix "github.com/armon/go-radix"
	"github.com/x-mod/dir"
)

type Template struct {
	name       string
	ext        string
	delimLeft  string
	delimRight string
	fns        map[string]interface{}
	dir        string
	root       *dir.Dir
	nameByPath bool
	rtree      *radix.Tree
	*template.Template
}

type Option func(*Template)

func Name(name string) Option {
	return func(t *Template) {
		t.name = name
	}
}

func Dir(dir string) Option {
	return func(t *Template) {
		t.dir = dir
	}
}

func Extension(ext string) Option {
	return func(t *Template) {
		t.ext = ext
	}
}

func Delims(left string, right string) Option {
	return func(t *Template) {
		t.delimLeft = left
		t.delimRight = right
	}
}

func Function(name string, fn interface{}) Option {
	return func(t *Template) {
		t.fns[name] = fn
	}
}

func NameByPath(flag bool) Option {
	return func(t *Template) {
		t.nameByPath = flag
	}
}

func New(opts ...Option) *Template {
	t := &Template{
		dir:   ".",
		ext:   ".tpl",
		fns:   make(map[string]interface{}),
		rtree: radix.New(),
	}
	for _, opt := range opts {
		opt(t)
	}
	return t
}

func (t *Template) Find(uri string) (string, error) {
	dst := strings.TrimPrefix(uri, "/")
	if dst == "" {
		return "index", nil
	}
	if name, _, ok := t.rtree.LongestPrefix(dst); ok {
		return name, nil
	}
	return "", fmt.Errorf("<%s> not matched", uri)
}

func (t *Template) AddFunc(name string, fn interface{}) {
	t.Template.Funcs(map[string]interface{}{
		name: fn,
	})
}

func (t *Template) Open() error {
	root := dir.New(dir.Root(t.dir))
	if err := root.Open(); err != nil {
		return err
	}
	t.root = root

	t.Template = template.New(t.name)
	t.Template.Delims(t.delimLeft, t.delimRight)
	t.Template.Funcs(t.fns)

	files := []string{}
	walk := func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, t.ext) {
			files = append(files, path)
		}
		return nil
	}
	if err := filepath.Walk(t.root.Path(), walk); err != nil {
		return err
	}
	//!nameByPath
	if !t.nameByPath {
		_, err := t.Template.ParseFiles(files...)
		return err
	}

	//nameByPath
	for _, file := range files {
		rel, err := filepath.Rel(t.root.Path(), file)
		if err != nil {
			return fmt.Errorf("filepath relative: %w", err)
		}

		name := filepath.ToSlash(rel)
		name = strings.TrimSuffix(name, t.ext)

		fd, err := os.Open(file)
		if err != nil {
			return err
		}
		bytes, err := io.ReadAll(fd)
		if err != nil {
			return err
		}
		if _, err := t.Template.New(name).Parse(string(bytes)); err != nil {
			return fmt.Errorf("template parse <%s> : %w", file, err)
		}
		t.rtree.Insert(name, struct{}{})
	}

	return nil
}
