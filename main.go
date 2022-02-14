package main

import (
  "github.com/tendermint/tendermint/libs/log"
  "github.com/terran-stakers/lightsentry/internal/lightsentry"
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
