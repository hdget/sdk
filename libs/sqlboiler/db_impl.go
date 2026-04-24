package sqlboiler

import (
	"context"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/hdget/sdk/common/biz"
)

type Db interface {
	Copier() DbCopier
	Executor() boil.Executor
}

type dbImpl struct {
	ctx    context.Context
	copier DbCopier
}

func (impl *dbImpl) Executor() boil.Executor {
	if tx, ok := biz.GetTransactor(impl.ctx).GetTx().(boil.Executor); ok {
		return tx
	}
	return boil.GetDB()
}

func (impl *dbImpl) Copier() DbCopier {
	return impl.copier
}
