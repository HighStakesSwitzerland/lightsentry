package lightsentry

const defaultConfigTemplate = `# This is a TOML config file.
# For more information, see https://github.com/toml-lang/toml

# NOTE: Any path below can be absolute (e.g. "/var/myawesomeapp/data") or
# relative to the home directory (e.g. "data"). The home directory is
# "$HOME/.tendermint" by default, but could be changed via $TMHOME env variable
# or --home cmd flag.

#######################################################
###     Lightsentry Server Configuration Options      ###
#######################################################

# Output level for logging: "none", info", "error", "debug". debug will enable pex and addrbook (very) verbose logs
log_level = "{{ .LogLevel }}"
# the chain id of the network
chain_id = "{{ .ChainId }}"
# the listen address
laddr = "{{ .ListenAddress }}"
# the seed list
seeds = "{{ .Seeds }}"
# the network of persistent peers
persistent_peers = "{{ .PersistentPeers }}"
# the private node ids not to gossip about
private_peer_ids = "{{ .PrivatePeerIDs }}"
`
