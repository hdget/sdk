package svc

import (
	"github.com/pkg/errors"
	"regexp"
)

var (
	moduleRegistry       = make(map[string]Module)
	errInvalidModuleName = errors.New("invalid module name, it should be: v<number>_name, e,g: v1_abc")
)

var (
	regModuleName = regexp.MustCompile(`^[vV]([0-9]+)_([a-zA-Z0-9]+)`)
)

func RegisterAsDaprModule(app string, svcHolder any, args ...map[string]any) error {
	module, err := NewDaprModule(app, svcHolder)
	if err != nil {
		return err
	}

	// 注册handlers
	var handlers map[string]any
	if len(args) > 0 {
		handlers = args[0]
	} else {
		handlers, err = module.DiscoverHandlers()
		if err != nil {
			return errors.Wrap(err, "discover handlers")
		}
	}
	module.RegisterHandlers(handlers)

	return nil
}

func GetRegistry() map[string]Module {
	return moduleRegistry
}

func addRegistry(name string, module Module) {
	moduleRegistry[name] = module
}

//func RegisterDaprModule(app string, version int, svcHolder any, options ...Option) error {
//	err := utils.StructSetComplexField(svcHolder, (*Module)(nil), NewDaprModule(app, moduleName, version, options...))
//	if err != nil {
//		return errors.Wrapf(err, "set base module for: %s ", moduleName)
//	}
//
//	module, ok := svcHolder.(Module)
//	if !ok {
//		return errors.New("invalid module")
//	}
//
//	// 将实际的struct实例保存进去
//	module.Super(svcHolder)
//
//	// 注册handlers
//	moduleRegistry[module.GetName()] = module
//
//	return nil
//}