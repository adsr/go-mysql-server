// Copyright 2021 Dolthub, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package analyzer

import (
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/expression"
	"github.com/dolthub/go-mysql-server/sql/plan"
)

func applyHashIn(ctx *sql.Context, a *Analyzer, n sql.Node, scope *Scope) (sql.Node, error) {
	return plan.TransformUpCtx(n, nil, func(c plan.TransformContext) (sql.Node, error) {
		filter, ok := c.Node.(*plan.Filter)
		if !ok {
			return c.Node, nil
		}

		e, err := expression.TransformUp(filter.Expression, func(expr sql.Expression) (sql.Expression, error) {
			switch e := expr.(type) {
			case *expression.InTuple:
				switch e.Left().(type) {
				// cannot HASH IN *plan.Subquery
				case expression.Tuple, *expression.Literal, *expression.GetField:
					return expression.NewHashInTuple(e.Left(), e.Right())
				default:
				}
			default:
			}
			return expr, nil
		})
		if err != nil {
			return nil, err
		}
		return filter.WithExpressions(e)
	})
}
