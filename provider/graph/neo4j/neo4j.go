package neo4j

import (
	"github.com/fatih/structs"
	"github.com/hdget/hdsdk/v2/intf"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/spf13/cast"
)

type neo4jProvider struct {
	logger intf.LoggerProvider
	config *neo4jProviderConfig
	driver neo4j.Driver
}

func New(configProvider intf.ConfigProvider, logger intf.LoggerProvider) (Provider, error) {
	c, err := newConfig(configProvider)
	if err != nil {
		return nil, err
	}

	provider := &neo4jProvider{
		logger: logger,
		config: c,
	}

	err = provider.Init()
	if err != nil {
		logger.Fatal("init neo4j provider", "err", err)
	}

	return provider, nil
}

// Init	initialize neo4j driver
func (p *neo4jProvider) Init(args ...any) error {
	var err error
	p.driver, err = p.newNeo4jDriver()
	if err != nil {
		return err
	}
	p.logger.Debug("init neo4j", "uri", p.config.VirtualUri)

	return nil
}

func (p *neo4jProvider) Exec(workFuncs []neo4j.TransactionWork, bookmarks ...string) (string, error) {
	session := p.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func(session neo4j.Session) {
		_ = session.Close()
	}(session)

	for _, fn := range workFuncs {
		if _, err := session.WriteTransaction(fn); err != nil {
			return "", err
		}
	}

	return session.LastBookmark(), nil
}

func (impl *neo4jProvider) Get(cypher string, args ...interface{}) (interface{}, error) {
	session := impl.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer func(session neo4j.Session) {
		_ = session.Close()
	}(session)

	var param map[string]interface{}
	if len(args) > 0 {
		param = structs.Map(args[0])
	}

	result, err := session.Run(cypher, param)
	if err != nil {
		return nil, err
	}

	var ret interface{}
	if result.Next() {
		ret = result.Record().Values[0]
	}

	if err = result.Err(); err != nil {
		return nil, err
	}

	return ret, nil
}

func (impl *neo4jProvider) Select(cypher string, args ...interface{}) ([]interface{}, error) {
	session := impl.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer func(session neo4j.Session) {
		_ = session.Close()
	}(session)

	var param map[string]interface{}
	if len(args) > 0 {
		param = structs.Map(args[0])
	}

	result, err := session.Run(cypher, param)
	if err != nil {
		return nil, err
	}

	rets := make([]interface{}, 0)
	for result.Next() {
		rets = append(rets, result.Record().Values[0])
	}

	if err = result.Err(); err != nil {
		return nil, err
	}

	return rets, nil
}

func (impl *neo4jProvider) Reader(bookmarks ...string) neo4j.Session {
	return impl.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead, Bookmarks: bookmarks})
}

func (impl *neo4jProvider) Writer(bookmarks ...string) neo4j.Session {
	return impl.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite, Bookmarks: bookmarks})
}

func (p *neo4jProvider) newNeo4jDriver() (neo4j.Driver, error) {
	// Address resolver is only valid for neo4j uri
	driver, err := neo4j.NewDriver(
		p.config.VirtualUri,
		neo4j.BasicAuth(p.config.Username, p.config.Password, ""),
		func(conf *neo4j.Config) {
			conf.AddressResolver = func(address neo4j.ServerAddress) []neo4j.ServerAddress {
				serverAddresses := make([]neo4j.ServerAddress, 0)
				for _, server := range p.config.Servers {
					serverAddresses = append(serverAddresses, neo4j.NewServerAddress(server.Host, cast.ToString(server.Port)))
				}
				return serverAddresses
			}
			conf.MaxConnectionPoolSize = p.config.MaxPoolSize
		})
	if err != nil {
		return nil, err
	}

	// check if neo4j can be connected or not
	_, err = driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead}).Run(
		`CALL dbms.components() YIELD name, versions, edition RETURN name, versions, edition`,
		nil)
	if err != nil {
		return nil, err
	}

	return driver, nil
}
