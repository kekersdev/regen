// Copyright 2026 Alexey Zakharov. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found in the LICENSE.txt file.

package regen

import (
	"bytes"
	"fmt"
	"io"
	"math/rand/v2"
	"regexp/syntax"
)

type xegerGenerator struct {
	expressions []*syntax.Regexp
	unboundMax int
}

const DefaultUnboundLimit int = 32

func NewGenerator(pattern string) (g *xegerGenerator) {
	g, _ = NewGeneratorAdvanced([]string{pattern}, syntax.Perl, false)
	return
}

func NewGeneratorAdvanced(patterns []string, mode syntax.Flags, simplify bool) (*xegerGenerator, error) {
	g := xegerGenerator{ unboundMax: DefaultUnboundLimit }
	for _, pattern := range patterns {
		if expr, err := syntax.Parse(pattern, mode); err != nil {
			return nil, fmt.Errorf("error while parsing regular expression `%q`: %w", pattern, err)
		} else {
			if simplify {
				g.expressions = append(g.expressions, expr.Simplify())
			} else {
				g.expressions = append(g.expressions, expr)
			}
		}
	}
	if len(g.expressions) == 0 {
		return nil, fmt.Errorf("error while creating generator: at least one pattern must be provided")
	}
	return &g, nil
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
	expr := g.expressions[rand.IntN(len(g.expressions))]
	e := genString(&b, expr, g.unboundMax)
	if e != nil && e != io.EOF {
		err = fmt.Errorf("error while generating string: %w", e)
		str = ""
	} else {
		str = b.String()
	}

	return
}
