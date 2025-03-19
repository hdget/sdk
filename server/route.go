package server

import (
	"encoding/json"
	"github.com/hdget/common/protobuf"
)

func (impl *appServerImpl) GetRoutes() ([]*protobuf.RouteItem, error) {
	data, err := impl.assetManager.Load(fileRoutes)
	if err != nil {
		return nil, err
	}

	var routes []*protobuf.RouteItem
	err = json.Unmarshal(data, &routes)
	if err != nil {
		return nil, err
	}
	return routes, nil
}
