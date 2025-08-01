package tidb

// Framework code is generated by the generator.

import (
	"context"
	"fmt"

	"github.com/pingcap/tidb/pkg/parser/ast"
	"github.com/pkg/errors"

	"github.com/bytebase/bytebase/backend/common"
	storepb "github.com/bytebase/bytebase/backend/generated-go/store"
	"github.com/bytebase/bytebase/backend/plugin/advisor"
)

var (
	_ advisor.Advisor = (*InsertDisallowOrderByRandAdvisor)(nil)
	_ ast.Visitor     = (*insertDisallowOrderByRandChecker)(nil)
)

func init() {
	advisor.Register(storepb.Engine_TIDB, advisor.MySQLInsertDisallowOrderByRand, &InsertDisallowOrderByRandAdvisor{})
}

// InsertDisallowOrderByRandAdvisor is the advisor checking for to disallow order by rand in INSERT statements.
type InsertDisallowOrderByRandAdvisor struct {
}

// Check checks for to disallow order by rand in INSERT statements.
func (*InsertDisallowOrderByRandAdvisor) Check(_ context.Context, checkCtx advisor.Context) ([]*storepb.Advice, error) {
	stmtList, ok := checkCtx.AST.([]ast.StmtNode)
	if !ok {
		return nil, errors.Errorf("failed to convert to StmtNode")
	}

	level, err := advisor.NewStatusBySQLReviewRuleLevel(checkCtx.Rule.Level)
	if err != nil {
		return nil, err
	}
	checker := &insertDisallowOrderByRandChecker{
		level: level,
		title: string(checkCtx.Rule.Type),
	}

	for _, stmt := range stmtList {
		checker.text = stmt.Text()
		checker.line = stmt.OriginTextPosition()
		(stmt).Accept(checker)
	}

	return checker.adviceList, nil
}

type insertDisallowOrderByRandChecker struct {
	adviceList []*storepb.Advice
	level      storepb.Advice_Status
	title      string
	text       string
	line       int
}

// Enter implements the ast.Visitor interface.
func (checker *insertDisallowOrderByRandChecker) Enter(in ast.Node) (ast.Node, bool) {
	code := advisor.Ok
	if insert, ok := in.(*ast.InsertStmt); ok {
		if insert.Select != nil {
			if selectNode, ok := insert.Select.(*ast.SelectStmt); ok {
				if selectNode.OrderBy != nil {
					for _, item := range selectNode.OrderBy.Items {
						if f, ok := item.Expr.(*ast.FuncCallExpr); ok {
							if f.FnName.L == ast.Rand {
								code = advisor.InsertUseOrderByRand
								break
							}
						}
					}
				}
			}
		}
	}

	if code != advisor.Ok {
		checker.adviceList = append(checker.adviceList, &storepb.Advice{
			Status:        checker.level,
			Code:          code.Int32(),
			Title:         checker.title,
			Content:       fmt.Sprintf("\"%s\" uses ORDER BY RAND in the INSERT statement", checker.text),
			StartPosition: common.ConvertANTLRLineToPosition(checker.line),
		})
	}

	return in, false
}

// Leave implements the ast.Visitor interface.
func (*insertDisallowOrderByRandChecker) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}
