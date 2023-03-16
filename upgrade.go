package main

import (
	"context"
	"encoding/base64"
	"fmt"
	uint128 "github.com/eteu-technologies/golang-uint128"
	"github.com/mapprotocol/near-api-go/pkg/client"
	"github.com/mapprotocol/near-api-go/pkg/types"
	"github.com/mapprotocol/near-api-go/pkg/types/action"
	"github.com/mapprotocol/near-api-go/pkg/types/key"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	// defaultPathName is the default config dir name
	defaultPathName = ".mcs_upgrade_tool"
	// defaultPathRoot is the path to the default config dir location.
	defaultPathRoot = "~/" + defaultPathName
	// envDir is the environment variable used to change the path root.
	envDir = "MCS_UPGRADE_TOOL_PATH"
	// Config name
	configName = "upgrade.json"
)

func upgradeCMD() *cli.Command {
	return &cli.Command{
		Name:   "upgrade",
		Usage:  "Upgrade MCS contract",
		Action: upgrade,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "config",
				Value: "./",
				Usage: "mcs upgrade tool config file upgrade.json path",
			},
		},
	}
}

type Config struct {
	Sender           string `json:"sender"`
	MultisigAccount  string `json:"multisig_account"`
	MCSAccount       string `json:"mcs_account"`
	MCSWasmFile      string `json:"mcs_wasm_file"`
	SenderPrivateKey string `json:"sender_private_key"`
	NearRPCUrl       string `json:"near_rpc_url"`
}

func UnmarshalConfig(repoRoot string) (*Config, error) {
	viper.SetConfigFile(filepath.Join(repoRoot, configName))
	viper.SetConfigType("json")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("MCS_UPGRADE_TOOL")
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	fmt.Println(config)

	return config, nil
}

func upgrade(cliCtx *cli.Context) error {
	path := cliCtx.String("config")
	fmt.Println("path", path)
	config, err := UnmarshalConfig(path)
	if err != nil {
		return fmt.Errorf("unmarshal config file: %v", err)
	}

	rpc, err := client.NewClient(config.NearRPCUrl)
	if err != nil {
		return fmt.Errorf("new near client: %v", err)
	}

	fmt.Println("private key", config.SenderPrivateKey)
	keyPair, err := key.NewBase58KeyPair(config.SenderPrivateKey)
	if err != nil {
		return fmt.Errorf("new base58 key pair with %s: %v", config.SenderPrivateKey, err)
	}

	code, err := ioutil.ReadFile(config.MCSWasmFile)
	if err != nil {
		return fmt.Errorf("read mcs contract file: %v", err)
	}

	encodedCode := base64.StdEncoding.EncodeToString(code)
	upgradeArgs := fmt.Sprintf("{\"code\": \"%s\"}", encodedCode)
	upgradeArgsEncode := base64.StdEncoding.EncodeToString([]byte(upgradeArgs))
	payload := fmt.Sprintf(
		"{\"request\":{\"receiver_id\":\"%s\",\"actions\":[{\"type\":\"FunctionCall\",\"method_name\":\"upgrade_self\",\"args\":\"%s\",\"deposit\":\"0\",\"gas\":\"180000000000000\"}]}}",
		config.MCSAccount,
		upgradeArgsEncode,
	)

	ctx := client.ContextWithKeyPair(context.Background(), keyPair)
	res, err := rpc.TransactionSendAwait(
		ctx,
		config.Sender,
		config.MultisigAccount,
		[]action.Action{action.NewFunctionCall("add_request_and_confirm", []byte(payload), 30*10000000000000, types.Balance(uint128.From64(1)))},
		client.WithLatestBlock(),
		client.WithKeyPair(keyPair),
	)
	if err != nil {
		return fmt.Errorf("TransactionSendAwait: %v", err)
	}

	fmt.Printf("https://explorer.testnet.near.org/transactions/%s\n", res.Transaction.Hash)

	return nil
}

func PathRoot() (string, error) {
	dir := os.Getenv(envDir)
	var err error
	if len(dir) == 0 {
		dir, err = homedir.Expand(defaultPathRoot)
	}
	return dir, err
}
