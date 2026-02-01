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
	"sync"
)

// A Xeger is a regex-based string generator.
//
// The zero value for Xeger holds no patterns and has the default limit for unbound quantifiers expansion.
// 
// All Xeger methods are safe for concurrent use by multiple goroutines.
type Xeger struct {
	mutex sync.RWMutex
	expressions []*syntax.Regexp
	unboundLimit int
}

// FromPattern attempts to parse the pattern assuming it has Perl-like syntax.
// 
// Returns a ready-to-use instance of Xeger or nil if expression can not be parsed.
func FromPattern(pattern string) *Xeger {
	g := &Xeger{ unboundLimit: DefaultUnboundLimit }
	err := g.AddPattern(pattern, syntax.Perl, false)
	if err != nil { return nil } else { return g }
}

// AddPattern attempts to parse the pattern and add it to the Xeger instance.
// 
//   - mode is used to specify parser flags
//   - simplify controls whether parsed expression should be simplified using [syntax.Regexp.Simplify]
func (g *Xeger) AddPattern(pattern string, mode syntax.Flags, simplify bool) error {
	expr, e := syntax.Parse(pattern, mode)
	if e != nil { return fmt.Errorf("failed to parse pattern `%q`: %w", pattern, e) }
	if simplify { expr = expr.Simplify() }
	g.mutex.Lock()
	g.expressions = append(g.expressions, expr)
	g.mutex.Unlock()
	return nil
}

// SetUnboundLimit is used to update the maximum number of repetitions allowed for patterns with
// unbound quantifiers.
func (g *Xeger) SetUnboundLimit(limit int) error {
	if limit <= 0 { return fmt.Errorf("failed to set new unbound limit: limit must be positive") }
	g.mutex.Lock()
	g.unboundLimit = limit
	g.mutex.Unlock()
	return nil
}

// MustGenerate attempts to generate a string using a pattern from Xeger instance. If multiple patterns are
// held by Xeger instance, a random one is selected.
// 
// Returns an empty string if any errors occur (including the case when Xeger instance does not hold any pattern)
func (g *Xeger) MustGenerate() string {
	g.mutex.RLock()
	str, err := g.Generate()
	g.mutex.RUnlock()
	if err != nil { panic(err) } else { return str }
}

// Generate attempts to generate a string using a pattern from Xeger instance. If multiple patterns are held by
// Xeger instance, a random one is selected.
func (g *Xeger) Generate() (str string, err error) {
	g.mutex.RLock()
	lim := g.unboundLimit
	expr, e := g.selectExpression()
	g.mutex.RUnlock()
	if e == nil { str, e = genStringWrapper(expr, lim) }
	if e != nil { err = fmt.Errorf("failed to generate string: %w", e) }
	return
}

func (g *Xeger) selectExpression() (*syntax.Regexp, error) {
	if len(g.expressions) == 0 { return nil, fmt.Errorf("generator does not hold any patterns") }
	return g.expressions[rand.IntN(len(g.expressions))], nil
}

// genStringWrapper is used for handling errors that mught arise when calling the original genString method
func genStringWrapper(expression *syntax.Regexp, unboundLimit int) (str string, err error) {
	defer func() { if cause := recover(); cause != nil { err = fmt.Errorf("%v", cause) } }()
	if unboundLimit == 0 { unboundLimit = DefaultUnboundLimit }
	var b bytes.Buffer
	err = genString(&b, expression, unboundLimit)
	if err == nil || err == io.EOF { str = b.String() }
	return
}
