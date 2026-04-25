package sqlboiler

import (
	"context"
	"fmt"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/hdget/sdk/common/bizctx"
	"github.com/hdget/sdk/common/provider"
	loggerUtils "github.com/hdget/utils/logger"
)

type Transactor interface {
	Finalize(err error)
}

type trans struct {
	tx     boil.Transactor
	ctx    context.Context
	errLog func(msg string, kvs ...any)
}

func NewTransactor(ctx context.Context, logger provider.Logger) (Transactor, error) {
	errLog := loggerUtils.Error
	if logger != nil {
		errLog = logger.Error
	}

	t := bizctx.GetTransactor(ctx)
	var transactor boil.Transactor
	if v, ok := t.GetTx().(boil.Transactor); ok {
		transactor = v
	} else {
		transactor, err := boil.BeginTx(context.Background(), nil)
		if err != nil {
			return nil, err
		}
		t.Ref(transactor)
	}

	return &trans{tx: transactor, ctx: ctx, errLog: errLog}, nil
}

func (t *trans) Finalize(err error) {
	tx := bizctx.GetTransactor(t.ctx)
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
		tx.Unref()
	}()

	if needFinalize := tx.ReachRoot(); !needFinalize {
		return
	}

	if err != nil {
		e := t.tx.Rollback()
		t.errLog("db roll back", "err", err, "rollback", e)
		return
	}

	e := t.tx.Commit()
	if e != nil {
		t.errLog("db commit", "err", e)
	}
}
