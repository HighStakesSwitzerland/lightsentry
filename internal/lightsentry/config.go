package lightsentry

import (
	"bytes"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/config"
	tmos "github.com/tendermint/tendermint/libs/os"
	"github.com/tendermint/tendermint/p2p"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type P2PConfig struct {
	config.P2PConfig `mapstructure:",squash"`
	LogLevel         string `mapstructure:"log_level"`
	ChainId          string `mapstructure:"chain_id"`
}

var configTemplate *template.Template

func init() {
	var err error
	tmpl := template.New("configFileTemplate").Funcs(template.FuncMap{
		"StringsJoin": strings.Join,
	})
	if configTemplate, err = tmpl.Parse(defaultConfigTemplate); err != nil {
		panic(err)
	}
}

func InitConfigs() (*P2PConfig, *p2p.NodeKey) {
	var tsConfig *P2PConfig

	userHomeDir, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	// init config directory & files if they don't exists yet
	homeDir := filepath.Join(userHomeDir, ".lightsentry")
	if err = os.MkdirAll(homeDir, os.ModePerm); err != nil {
		panic(err)
	}

	configFilePath := filepath.Join(homeDir, "config.toml")
	viper.SetConfigName("config")
	viper.AddConfigPath(homeDir)

	if err := viper.ReadInConfig(); err == nil {
		logger.Info(fmt.Sprintf("Loading config file: %s", viper.ConfigFileUsed()))
		err := viper.Unmarshal(&tsConfig)
		if err != nil {
			panic(err)
		}
	} else if _, ok := err.(viper.ConfigFileNotFoundError); ok { // ignore not found error, return other errors
		logger.Info("No existing configuration found, generating one")
		tsConfig = defaultP2PConfig()
		writeConfigFile(configFilePath, tsConfig)
		os.Exit(0)
	} else {
		panic(err)
	}

	nodeKeyFilePath := filepath.Join(homeDir, "node_key.json")
	nodeKey, err := p2p.LoadOrGenNodeKey(nodeKeyFilePath)
	if err != nil {
		panic(err)
	}

	logger.Info("Node key: ", "nodeId", nodeKey.ID())
  logger.Info("Listen Addr: ", "tcp", tsConfig.ListenAddress)

	return tsConfig, nodeKey
}

func defaultP2PConfig() *P2PConfig {
	p := &P2PConfig{
		P2PConfig: *config.DefaultP2PConfig(),
		ChainId:   "",
    LogLevel: "info",
	}
	p.ListenAddress = fmt.Sprintf("tcp://0.0.0.0:%d", 26656)
	return p
}

// WriteConfigFile renders config using the template and writes it to configFilePath.
func writeConfigFile(configFilePath string, config *P2PConfig) {
	var buffer bytes.Buffer

	if err := configTemplate.Execute(&buffer, config); err != nil {
		panic(err)
	}

	tmos.MustWriteFile(configFilePath, buffer.Bytes(), 0644)
}
