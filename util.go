// Copyright 2026 Alexey Zakharov. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found in the LICENSE.txt file.

// Package regen is a tool for generating random strings from Go/RE2 regular expressions.
// It is based on homonymous CLI tool originally developed by by Noel Cower.
package regen

import (
	"bytes"
	"fmt"
	"io"
	"math/rand/v2"
	"regexp/syntax"
)

type XegerGenerator struct {
	expressions []*syntax.Regexp
	unboundMax int
}

func FromPattern(pattern string) *XegerGenerator {
	g := &XegerGenerator{ unboundMax: DefaultUnboundLimit }
	err := g.AddPattern(pattern, syntax.Perl, false)
	if err != nil { return nil } else { return g }
}

func (g *XegerGenerator) AddPattern(pattern string, mode syntax.Flags, simplify bool) error {
	expr, e := syntax.Parse(pattern, mode)
	if e != nil { return fmt.Errorf("failed to parse pattern `%q`: %w", pattern, e) }
	if simplify { expr = expr.Simplify() }
	g.expressions = append(g.expressions, expr)
	return nil
}

func (g *XegerGenerator) SetUnboundLimit(limit int) error {
	if limit <= 0 { return fmt.Errorf("failed to set new unbound limit: limit must be positive") }
	g.unboundMax = limit
	return nil
}

func (g *XegerGenerator) MustGenerate() string {
	str, err := g.Generate()
	if err != nil { panic(err) } else { return str }
}

func (g *XegerGenerator) Generate() (str string, err error) {
	str, e := genStringWrapper(g.expressions, g.unboundMax)
	if e != nil {
		err = fmt.Errorf("failed to generate string: %w", e)
	}
	return
}

func genStringWrapper(expressionList []*syntax.Regexp, unboundMax int) (str string, err error) {
	defer func() { if cause := recover(); cause != nil { err = fmt.Errorf("%v", cause) } }()

	if len(expressionList) == 0 { return "", fmt.Errorf("no patterns provided") }
	if unboundMax == 0 { unboundMax = DefaultUnboundLimit }

	expr := expressionList[rand.IntN(len(expressionList))]
	var b bytes.Buffer
	err = genString(&b, expr, unboundMax)
	if err != nil && err != io.EOF { str = b.String() }
	return
}
