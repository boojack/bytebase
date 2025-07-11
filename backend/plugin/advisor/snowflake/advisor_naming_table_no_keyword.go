// Package snowflake is the advisor for snowflake database.
package snowflake

import (
	"context"
	"fmt"

	"github.com/antlr4-go/antlr/v4"
	parser "github.com/bytebase/snowsql-parser"
	"github.com/pkg/errors"

	storepb "github.com/bytebase/bytebase/backend/generated-go/store"
	"github.com/bytebase/bytebase/backend/plugin/advisor"
	snowsqlparser "github.com/bytebase/bytebase/backend/plugin/parser/snowflake"
)

var (
	_ advisor.Advisor = (*NamingTableNoKeywordAdvisor)(nil)
)

func init() {
	advisor.Register(storepb.Engine_SNOWFLAKE, advisor.SnowflakeTableNamingNoKeyword, &NamingTableNoKeywordAdvisor{})
}

// NamingTableNoKeywordAdvisor is the advisor checking for table naming convention without keyword.
type NamingTableNoKeywordAdvisor struct {
}

// Check checks for table naming convention without keyword.
func (*NamingTableNoKeywordAdvisor) Check(_ context.Context, checkCtx advisor.Context) ([]*storepb.Advice, error) {
	tree, ok := checkCtx.AST.(antlr.Tree)
	if !ok {
		return nil, errors.Errorf("failed to convert to Tree")
	}

	level, err := advisor.NewStatusBySQLReviewRuleLevel(checkCtx.Rule.Level)
	if err != nil {
		return nil, err
	}

	listener := &namingTableNoKeywordChecker{
		level: level,
		title: string(checkCtx.Rule.Type),
	}

	antlr.ParseTreeWalkerDefault.Walk(listener, tree)

	return listener.generateAdvice()
}

// namingTableNoKeywordChecker is the listener for table naming convention without keyword.
type namingTableNoKeywordChecker struct {
	*parser.BaseSnowflakeParserListener

	level storepb.Advice_Status
	title string

	adviceList []*storepb.Advice
}

// generateAdvice returns the advices generated by the listener, the advices must not be empty.
func (l *namingTableNoKeywordChecker) generateAdvice() ([]*storepb.Advice, error) {
	return l.adviceList, nil
}

// EnterCreate_table is called when production create_table is entered.
func (l *namingTableNoKeywordChecker) EnterCreate_table(ctx *parser.Create_tableContext) {
	tableName := snowsqlparser.NormalizeSnowSQLObjectNamePart(ctx.Object_name().GetO())
	if snowsqlparser.IsSnowflakeKeyword(tableName, false) {
		l.adviceList = append(l.adviceList, &storepb.Advice{
			Status:  l.level,
			Code:    advisor.NameIsKeywordIdentifier.Int32(),
			Title:   l.title,
			Content: fmt.Sprintf("Table name %q is a keyword identifier and should be avoided.", tableName),
		})
	}
}

// EnterAlter_table is called when production alter_table is entered.
func (l *namingTableNoKeywordChecker) EnterAlter_table(ctx *parser.Alter_tableContext) {
	if ctx.RENAME() == nil {
		return
	}

	tableName := snowsqlparser.NormalizeSnowSQLObjectNamePart(ctx.Object_name(1).GetO())
	if snowsqlparser.IsSnowflakeKeyword(tableName, false) {
		l.adviceList = append(l.adviceList, &storepb.Advice{
			Status:  l.level,
			Code:    advisor.NameIsKeywordIdentifier.Int32(),
			Title:   l.title,
			Content: fmt.Sprintf("Table name %q is a keyword identifier and should be avoided.", tableName),
		})
	}
}
