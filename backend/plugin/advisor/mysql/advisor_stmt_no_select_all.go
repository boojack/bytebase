package mysql

import (
	"context"
	"fmt"

	"github.com/antlr4-go/antlr/v4"
	mysql "github.com/bytebase/mysql-parser"
	"github.com/pkg/errors"

	"github.com/bytebase/bytebase/backend/common"
	storepb "github.com/bytebase/bytebase/backend/generated-go/store"
	"github.com/bytebase/bytebase/backend/plugin/advisor"
	mysqlparser "github.com/bytebase/bytebase/backend/plugin/parser/mysql"
)

var (
	_ advisor.Advisor = (*NoSelectAllAdvisor)(nil)
)

func init() {
	advisor.Register(storepb.Engine_MYSQL, advisor.MySQLNoSelectAll, &NoSelectAllAdvisor{})
	advisor.Register(storepb.Engine_MARIADB, advisor.MySQLNoSelectAll, &NoSelectAllAdvisor{})
	advisor.Register(storepb.Engine_OCEANBASE, advisor.MySQLNoSelectAll, &NoSelectAllAdvisor{})
}

// NoSelectAllAdvisor is the advisor checking for no "select *".
type NoSelectAllAdvisor struct {
}

// Check checks for no "select *".
func (*NoSelectAllAdvisor) Check(_ context.Context, checkCtx advisor.Context) ([]*storepb.Advice, error) {
	root, ok := checkCtx.AST.([]*mysqlparser.ParseResult)
	if !ok {
		return nil, errors.Errorf("failed to convert to mysql parse result")
	}

	level, err := advisor.NewStatusBySQLReviewRuleLevel(checkCtx.Rule.Level)
	if err != nil {
		return nil, err
	}
	checker := &noSelectAllChecker{
		level: level,
		title: string(checkCtx.Rule.Type),
	}
	for _, stmtNode := range root {
		checker.baseLine = stmtNode.BaseLine
		antlr.ParseTreeWalkerDefault.Walk(checker, stmtNode.Tree)
	}

	return checker.adviceList, nil
}

type noSelectAllChecker struct {
	*mysql.BaseMySQLParserListener

	baseLine   int
	adviceList []*storepb.Advice
	level      storepb.Advice_Status
	title      string
	text       string
}

func (checker *noSelectAllChecker) EnterQuery(ctx *mysql.QueryContext) {
	checker.text = ctx.GetParser().GetTokenStream().GetTextFromRuleContext(ctx)
}

func (checker *noSelectAllChecker) EnterSelectItemList(ctx *mysql.SelectItemListContext) {
	if ctx.MULT_OPERATOR() != nil {
		checker.adviceList = append(checker.adviceList, &storepb.Advice{
			Status:        checker.level,
			Code:          advisor.StatementSelectAll.Int32(),
			Title:         checker.title,
			Content:       fmt.Sprintf("\"%s\" uses SELECT all", checker.text),
			StartPosition: common.ConvertANTLRLineToPosition(checker.baseLine + ctx.GetStart().GetLine()),
		})
	}
}
