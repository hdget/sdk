package server

import (
	"encoding/json"
	"fmt"
	"github.com/hdget/sdk/dapr"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cast"
	"os"
	"strings"
)

const (
	fileSourceCodeInfo = ".source.json" // 保存源代码信息
	fileRoutes         = ".routes.json" // 保存找到的路由信息
)

func (impl *appServerImpl) PatchSourceCode(srcDir string, skipDirs ...string) error {
	// patch source codes
	sourceCodeInfo, err := dapr.NewSourceCodeHandler(srcDir).Discover(impl.serverImportPath, impl.serverRunFuncName, skipDirs...)
	if err != nil {
		return err
	}

	err = impl.assetManager.Store(fileSourceCodeInfo, sourceCodeInfo)
	if err != nil {
		return err
	}

	fmt.Println("")
	fmt.Printf("=== discover source code ===")
	fmt.Println("")

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Path"})
	table.SetRowLine(true)
	table.Append([]string{
		"ServerEntry", sourceCodeInfo.ServerEntry,
	})
	for k, v := range sourceCodeInfo.ModulePaths {
		table.Append([]string{k, v})
	}
	table.Render() // Send output

	fmt.Println("")
	fmt.Println("=== patch source code ===")
	fmt.Println("")

	return dapr.NewSourceCodeHandler(srcDir).Patch(sourceCodeInfo)
}

func (impl *appServerImpl) InspectSourceCode(srcDir string, skipDirs ...string) error {
	content, err := impl.assetManager.Load(fileSourceCodeInfo)
	if err != nil {
		return err
	}

	var sourceCodeInfo *dapr.SourceCodeInfo
	err = json.Unmarshal(content, &sourceCodeInfo)
	if err != nil {
		return err
	}

	fmt.Println("")
	fmt.Println("=== inspect source code ===")
	fmt.Println("")

	// inspect routes
	routes, err := dapr.NewSourceCodeHandler(srcDir).Inspect(sourceCodeInfo)
	if err != nil {
		return err
	}

	err = impl.assetManager.Store(fileRoutes, routes)
	if err != nil {
		return err
	}

	moduleName2handlerNames := make(map[string][]string)
	for _, routeItem := range routes {
		k := fmt.Sprintf("v%d_%s", routeItem.ModuleVersion, routeItem.ModuleName)
		moduleName2handlerNames[k] = append(moduleName2handlerNames[k], routeItem.Handler)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ModuleName", "Total", "Handlers"})
	table.SetRowLine(true)
	for moduleName, handlerNames := range moduleName2handlerNames {
		table.Append([]string{moduleName, cast.ToString(len(handlerNames)), strings.Join(handlerNames, ", ")})
		//fmt.Printf(" * invocation module: %s\ttotal: %-5d\tfunctions: [%s]\n", moduleName, len(handlerNames), strings.Join(handlerNames, ", "))
	}
	table.Render() // Send output

	return nil
}
