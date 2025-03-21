package sourcecode

import (
	"fmt"
	"github.com/hdget/common/protobuf"
	"github.com/hdget/sdk/dapr"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"os"
	"path/filepath"
	"strings"
)

// 检测路由注解，必须在patch动作完成后重新来执行该动作，
// 否则之前动作导入的patch内容不生效
type findRouteAnnotationHandler struct {
	sc *sourceCodeManagerImpl
}

const (
	fileRoutes = ".routes.json" // 路由信息
)

func newFindRouteAnnotationHandler(sc *sourceCodeManagerImpl) Handler {
	return &findRouteAnnotationHandler{
		sc: sc,
	}
}

func (h *findRouteAnnotationHandler) Handle() error {
	fmt.Println("")
	fmt.Println("=== find route annotations ===")
	fmt.Println("")

	meta, err := newMetaDataManager(h.sc.srcDir).Load()
	if err != nil {
		return err
	}

	if len(meta.ModulePaths) == 0 || meta.ModulePaths["INVOCATION_MODULE"] == "" {
		return errors.New("invocation module path not found")
	}

	routeItems := make([]*protobuf.RouteItem, 0)
	absModulePath := filepath.Join(h.sc.srcDir, meta.ModulePaths["INVOCATION_MODULE"])
	for _, m := range dapr.GetInvocationModules() {
		routeAnnotations, err := m.GetRouteAnnotations(absModulePath, h.sc.handlerMatchers...)
		if err != nil {
			return err
		}

		for _, ann := range routeAnnotations {
			for _, httpMethod := range ann.Methods {
				isPublic := int32(0)
				if ann.IsPublic {
					isPublic = 1
				}

				isRawResponse := int32(0)
				if ann.IsRawResponse {
					isRawResponse = 1
				}

				routeItems = append(routeItems, &protobuf.RouteItem{
					App:           m.GetApp(),
					ModuleVersion: int32(m.GetModuleInfo().ModuleVersion),
					ModuleName:    m.GetModuleInfo().ModuleName,
					Handler:       ann.HandlerAlias,
					Endpoint:      ann.Endpoint,
					HttpMethod:    httpMethod,
					Permissions:   ann.Permissions,
					Origin:        ann.Origin,
					IsPublic:      isPublic,
					IsRawResponse: isRawResponse,
					Comment:       strings.Join(ann.Comments, "\r"),
				})
			}
		}
	}

	err = h.sc.assetManager.Store(fileRoutes, routeItems)
	if err != nil {
		return err
	}

	h.printRoutes(routeItems)

	return nil
}

func (h *findRouteAnnotationHandler) printRoutes(routeItems []*protobuf.RouteItem) {
	moduleName2handlerNames := make(map[string][]string)
	for _, routeItem := range routeItems {
		k := fmt.Sprintf("v%d_%s", routeItem.ModuleVersion, routeItem.ModuleName)
		moduleName2handlerNames[k] = append(moduleName2handlerNames[k], routeItem.Handler)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"MODULE_NAME", "TOTAL", "HANDLERS"})
	table.SetRowLine(true)
	for moduleName, handlerNames := range moduleName2handlerNames {
		table.Append([]string{moduleName, cast.ToString(len(handlerNames)), strings.Join(handlerNames, ", ")})
	}
	table.Render() // Send output
}
