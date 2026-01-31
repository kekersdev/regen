// Copyright 2026 Alexey Zakharov. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found in the LICENSE.txt file.

package regen

import (
	"bytes"
	"fmt"
	"io"
	"regexp/syntax"
)

type xegerGenerator struct {
	expression *syntax.Regexp	// TODO: multiple expressions
	unboundMax int
}

const DefaultUnboundLimit int = 32

func NewGenerator(pattern string) (g *xegerGenerator) {
	g, _ = NewGeneratorAdvanced(pattern, syntax.Perl, false)
	return
}

func NewGeneratorAdvanced(pattern string, mode syntax.Flags, simplify bool) (g *xegerGenerator, err error) {
	g = &xegerGenerator{ unboundMax: DefaultUnboundLimit }
	g.expression, err = syntax.Parse(pattern, mode)
	if err != nil {
		return nil, fmt.Errorf("error while parsing regular expression `%q`: %w", pattern, err)
	}
	if simplify { g.expression = g.expression.Simplify() }
	return
}

func (g *xegerGenerator) SetUnboundLimit(limit int) (ok bool) {
	if limit <= 0 {
		ok = false
	} else {
		g.unboundMax = limit
		ok = true
	}
	return
}

func (g *xegerGenerator) MustGenerate() string {
	str, err := g.Generate()
	if err != nil { panic(err) } else { return str }
}

func (g *xegerGenerator) Generate() (str string, err error) {
	defer func() {
		if cause := recover(); cause != nil {
			err = fmt.Errorf("error while generating string: %v", cause)
			str = ""
		}
	}()

	var b bytes.Buffer
	e := genString(&b, g.expression, g.unboundMax)
	if e != nil && e != io.EOF {
		err = fmt.Errorf("error while generating string: %w", e)
		str = ""
	} else {
		str = b.String()
	}

	return
}
