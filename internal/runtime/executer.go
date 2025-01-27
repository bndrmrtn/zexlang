package runtime

import (
	"fmt"

	"github.com/bndrmrtn/zexlang/internal/builtin"
	"github.com/bndrmrtn/zexlang/internal/errs"
	"github.com/bndrmrtn/zexlang/internal/models"
	"github.com/bndrmrtn/zexlang/internal/tokens"
)

type Executer struct {
	runtime *Runtime
	parent  *Executer

	fns    map[string]*models.Node
	vars   map[string]*models.Node
	blocks map[string]*models.Node
}

func NewExecuter(r *Runtime, parent *Executer) *Executer {
	return &Executer{
		runtime: r,
		parent:  parent,
		fns:     make(map[string]*models.Node),
		vars:    make(map[string]*models.Node),
		blocks:  make(map[string]*models.Node),
	}
}

func (e *Executer) Bind(variable *models.Node) {
	e.vars[variable.Content] = variable
}

func (e *Executer) Execute(ts []*models.Node) ([]*builtin.FuncReturn, error) {
	for _, token := range ts {
		switch token.Type {
		case tokens.Function:
			e.fns[token.Content] = token
		case tokens.Let, tokens.Const:
			if _, ok := e.vars[token.Content]; ok {
				return nil, errs.WithDebug(fmt.Errorf("%w: %v", errs.VariableAlreadyExists, token.Content), token.Debug)
			}
			e.vars[token.Content] = token
		case tokens.Define:
			e.blocks[token.Content] = token
		case tokens.FuncCall:
			e.executeFn(token)
		}
	}
	return nil, nil
}

func (e *Executer) ExecuteFn(name string, args []*builtin.Variable) ([]*builtin.FuncReturn, error) {
	fn, ok := e.fns[name]
	if !ok {
		return nil, errs.WithDebug(fmt.Errorf("function %v not found", name), nil)
	}

	return NewExecuter(e.runtime, e).Execute(fn.Children)
}
