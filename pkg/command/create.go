package command

import (
	"github.com/urfave/cli/v2"
	"github.com/weaveworks-liquidmetal/flintlock/api/types"

	"github.com/warehouse-13/hammertime/pkg/client"
	"github.com/warehouse-13/hammertime/pkg/config"
	"github.com/warehouse-13/hammertime/pkg/defaults"
	"github.com/warehouse-13/hammertime/pkg/flags"
	"github.com/warehouse-13/hammertime/pkg/microvm"
	"github.com/warehouse-13/hammertime/pkg/utils"
)

func createCommand() *cli.Command {
	cfg := &config.Config{
		ClientConfig: config.ClientConfig{
			ClientBuilderFunc: client.New,
		},
	}

	return &cli.Command{
		Name:    "create",
		Usage:   "create a new microvm",
		Aliases: []string{"c"},
		Before:  flags.ParseFlags(cfg),
		Flags: flags.CLIFlags(
			flags.WithGRPCAddressFlag(),
			flags.WithNameAndNamespaceFlags(),
			flags.WithJSONSpecFlag(),
			flags.WithSSHKeyFlag(),
			flags.WithBasicAuthToken(),
		),
		Action: func(c *cli.Context) error {
			return CreateFn(cfg)
		},
	}
}

func CreateFn(cfg *config.Config) error {
	client, err := cfg.ClientBuilderFunc(cfg.GRPCAddress)
	if err != nil {
		return err
	}

	defer client.Close()

	var mvm *types.MicroVMSpec

	if utils.IsSet(cfg.JSONFile) {
		mvm, err = utils.LoadSpecFromFile(cfg.JSONFile)
		if err != nil {
			return err
		}
	} else {
		mvm, err = newMicroVM(cfg.MvmName, cfg.MvmNamespace, cfg.SSHKeyPath)
		if err != nil {
			return err
		}
	}

	res, err := client.Create(mvm)
	if err != nil {
		return err
	}

	return utils.PrettyPrint(res)
}

func newMicroVM(name, namespace, sshPath string) (*types.MicroVMSpec, error) {
	mvm := defaults.BaseMicroVM()

	metaData, err := microvm.CreateMetadata(name, namespace)
	if err != nil {
		return nil, err
	}

	userData, err := microvm.CreateUserData(name, sshPath)
	if err != nil {
		return nil, err
	}

	mvm.Id = name
	mvm.Namespace = namespace
	mvm.Metadata = map[string]string{
		"meta-data": metaData,
		"user-data": userData,
	}

	return mvm, nil
}
