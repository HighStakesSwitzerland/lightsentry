package main

import (
	"github.com/highstakesswitzerland/lightsentry/internal/lightsentry"
	"github.com/tendermint/tendermint/libs/log"
	"os"
)

var (
	logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "main")
)

func main() {
	seedConfig, nodeKey := lightsentry.InitConfigs()
	sentryNode := lightsentry.NewSentryNode(seedConfig, nodeKey)
	sentryNode.DialPersistentPeers()
	sentryNode.AddPrivatePeerIDs()
	sentryNode.Switch.Wait() // block
}
