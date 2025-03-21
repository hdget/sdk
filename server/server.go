package server

//
//type appServerImpl struct {
//	assetManager      asset.Manager
//	actions           []dapr.Action
//	serverImportPath  string // 服务器代码所在的包路径，即server.New().Run所在的包路径
//	serverRunFuncName string // 服务器代码运行的函数名
//}
//
//const ()
//
//func New(logger intf.LoggerProvider, assetFs embed.FS, options ...dapr.Option) AppServer {
//	_, callPath, _, _ := runtime.Caller(1)
//
//	embedAbsPath, embedRelPath := findEmbedPath(callPath)
//
//	srv := &appServerImpl{
//		assetManager:      asset.newAssetManager(embedAbsPath, embedRelPath, assetFs),
//		actions:           make([]dapr.Action, 0),
//		serverImportPath:  defaultAppServerImportPath,
//		serverRunFuncName: defaultAppServerRunFunction,
//	}
//
//	for _, apply := range options {
//		apply(srv)
//	}
//	return srv
//}
//
//func (impl *appServerImpl) Run(appAddress string) error {
//	appServer, err := dapr.NewGrpcServer(logger, appAddress)
//	if err != nil {
//		return errors.Wrap(err, "new app server")
//	}
//
//	if err = appServer.Start(); err != nil {
//		sdk.Logger().Fatal("start app server", "err", err)
//	}
//	return nil
//}
