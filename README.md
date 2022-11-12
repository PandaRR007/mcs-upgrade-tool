# MCS Upgrade Tool

A tool for multisig member to add "upgrade MCS contract" request and confirm it.

## Pre-requisite

**Go (version 1.18 or later)**

Go are available at the [Go installation page](https://go.dev/doc/install) and necessary bundles can be downloaded from the [Go download page](https://go.dev/dl/).


## Build
```shell
git clone https://github.com/PandaRR007/mcs-upgrade-tool.git
cd mcs-upgrade-tool
make
```

You can see a binary **mcs-upgrade-tool** is generated.

## Usage

1. Configure upgrade.json
```json
{
  "nearRpcUrl": "https://archival-rpc.testnet.near.org",   // near rpc url
  "sender": "member1.map002.testnet",                      // multisig member account to add the request
  "senderPrivateKey": "ed25519:...",                       // private key of the multisig member
  "multisigAccount": "multisig.mfac.map004.testnet",       // multisig contract account
  "mcsAccount": "mos2.mfac.map004.testnet",                // MCS contract account to be upgraded
  "mcsWasmFile": "/path/to/mcs.wasm"                       // MCS contract wasm file to upgrade the MCS contract
}
```

2. Execute upgrade command

```shell
./mcs-upgrade-tool upgrade
```

3. Get the request ID

After the command is executed successfully, you can see a NEAR blockchain explorer link at the last line of the output.
![img.png](res/img.png)

Click it, and you can get the request ID from the transaction details.

![img.png](res/img1.png)