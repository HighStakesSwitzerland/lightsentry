package lightsentry

import (
  "fmt"
  "github.com/mitchellh/go-homedir"
  "github.com/tendermint/tendermint/libs/log"
  tmos "github.com/tendermint/tendermint/libs/os"
  tmstrings "github.com/tendermint/tendermint/libs/strings"
  "github.com/tendermint/tendermint/p2p"
  "github.com/tendermint/tendermint/p2p/pex"
  "github.com/tendermint/tendermint/version"
  "net"
  "os"
  "path/filepath"
  "time"
)

var (
	logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "config")
)

func StartSeedNode(config *P2PConfig, nodeKey *p2p.NodeKey) *p2p.Switch {
	logger.Info("Starting Sentry Node for chain " + config.ChainId)

	protocolVersion :=
		p2p.NewProtocolVersion(
			version.P2PProtocol,
			version.BlockProtocol,
			0,
		)

	// NodeInfo gets info on your node
	nodeInfo := p2p.DefaultNodeInfo{
		ProtocolVersion: protocolVersion,
		DefaultNodeID:   nodeKey.ID(),
		ListenAddr:      config.ListenAddress,
		Network:         config.ChainId,
		Version:         "1.0.0",
		Channels:        []byte{byte(0x00)},
		Moniker:         fmt.Sprintf("%s-lightsentry", config.ChainId),
	}

	addr, err := p2p.NewNetAddressString(p2p.IDAddressString(nodeInfo.DefaultNodeID, nodeInfo.ListenAddr))
	if err != nil {
		panic(err)
	}

	// set conn settings
	config.RecvRate = 5120000
	config.SendRate = 5120000
	config.MaxPacketMsgPayloadSize = 1024
	config.FlushThrottleTimeout = 100 * time.Millisecond
	config.AllowDuplicateIP = true
	config.DialTimeout = 30 * time.Second
	config.HandshakeTimeout = 20 * time.Second
	config.SeedMode = false
	config.PexReactor = true
  config.AddrBookStrict = false

	transport := p2p.NewMultiplexTransport(nodeInfo, *nodeKey, p2p.MConnConfig(&config.P2PConfig))
	if err := transport.Listen(*addr); err != nil {
		panic(err)
	}

	userHomeDir, _ := homedir.Dir()
	addrBookFilePath := filepath.Join(userHomeDir, ".lightsentry", "addrbook-"+config.ChainId+".json")
	addrBook := pex.NewAddrBook(addrBookFilePath, config.AddrBookStrict)

	pexReactor := pex.NewReactor(addrBook, &pex.ReactorConfig{

		SeedMode:                     config.SeedMode,
		Seeds:                        tmstrings.SplitAndTrim(config.Seeds, ",", " "),
		SeedDisconnectWaitPeriod:     5 * time.Minute, // default is 28 hours, we just want to harvest as many addresses as possible
		PersistentPeersMaxDialPeriod: 5 * time.Minute, // use exponential back-off
	})

	sw := p2p.NewSwitch(&config.P2PConfig, transport)

	var configuredLogger log.Logger
	switch config.LogLevel {
	case "none":
		configuredLogger = log.NewNopLogger()
	case "info":
		configuredLogger = log.NewFilter(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), log.AllowInfo())
	case "error":
		configuredLogger = log.NewFilter(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), log.AllowError())
	case "debug":
		configuredLogger = log.NewFilter(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), log.AllowDebug())
	default:
		configuredLogger = logger
	}

	sw.SetLogger(configuredLogger.With("module", "switch"))
	addrBook.SetLogger(configuredLogger.With("module", "addrbook", "chain", config.ChainId))
	pexReactor.SetLogger(configuredLogger.With("module", "pex"))

	sw.SetNodeKey(nodeKey)
	sw.SetAddrBook(addrBook)
	sw.AddReactor("pex", pexReactor)
	sw.SetNodeInfo(nodeInfo)
	err = sw.Start()

	if err != nil {
		panic(err)
	}

	netAddress := p2p.NetAddress{
		ID:   "437464640c5f3e0a6b1fb6756df21ef115b26ffc",
		IP:   net.ParseIP("195.201.225.154"),
		Port: 26656,
	}

	err = sw.DialPeerWithAddress(&netAddress)
	if err != nil {
		logger.Error("dial persistent peer", err)
	}

	tmos.TrapSignal(logger, func() {
		logger.Info("shutting down addrbooks...")
		_ = addrBook.Stop()
		_ = sw.Stop()
	})

	return sw
}
