# Lightsentry

This project is **NOT** functional, still thinking about it. And actually thinking it may no be achievable :/

## Lightsentry is a tool that uses Tendermint's implementation to run a sentry node without the need of the data folder.

As a sentry node is *in fine* just a node isolating the validator behind it, there is actually no need to run all the other process a full node runs (ie. all the reactors like mempool, blockchain, consensus, tx indexer, state sync etc etc)

This implementation only runs the pex reactor, address book and the switch to correctly discover peers and have a nice connectivity for the private peers linked to it.

All this with 0 data!

### Configuration

```bash
git clone https://github.com/Terran-Stakers/lightsentry
go mod tidy
npm install
npm run build
go install .
./lightsentry
```

A file `$HOME/.lightsentry/config.toml` will be generated if it doesn't exist yet, with some default parameters,
and the program will exit.

You need to fill the `seeds` and `chain_id` and other params and run it again.

## License

[Blue Oak Model License 1.0.0](https://blueoakcouncil.org/license/1.0.0)
