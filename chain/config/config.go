package config

import (
	"github.com/ethereum/go-ethereum/common"
)

var (
	// chain ID
	ethChainID      = int64(5)
	tronChainID     = int64(100)
	bttcChainID     = int64(199)
	bttcTestChainID = int64(1029)
	testChainID     = int64(1337)
	// start block
	ethStartBlock = uint64(10000)

	tronStartBlock = uint64(4933174)
	bttcStartBlock = uint64(100)
	// factory address
	ethFactoryAddress = common.HexToAddress("0x5E6802d9e7C8CD43BB7C96524fDD50FE8460B92c")
	ethOracleAddress  = common.HexToAddress("0xFB6a65aF1bb250EAf3f58C420912B0b6eA05Ea7a")
	ethBatchAddress   = common.HexToAddress("0xFB6a65aF1bb250EAf3f58C420912B0b6eA05Ea7a")

	tronFactoryAddress = common.HexToAddress("0x0c9de531dcb38b758fe8a2c163444a5e54ee0db2")
	tronOracleAddress  = common.HexToAddress("0x0c9de531dcb38b758fe8a2c163444a5e54ee0db2")
	tronBatchAddress   = common.HexToAddress("0x0c9de531dcb38b758fe8a2c163444a5e54ee0db2")

	bttcTestFactoryAddress = common.HexToAddress("0x392E3CE1B56c34ae0102C70cE3D640FFB9048A79")
	bttcTestOracleAddress  = common.HexToAddress("0xb2C746a9C81564bEF8382e885AF11e73De4a9E15")
	bttcTestBatchAddress   = common.HexToAddress("0x0c9de531dcb38b758fe8a2c163444a5e54ee0db2")

	bttcFactoryAddress = common.HexToAddress("0x107742EB846b86CEaAF7528D5C85cddcad3e409A")
	bttcOracleAddress  = common.HexToAddress("0x70fD2b6b6fEd65c8BC0D9Fd0656502Ffd05B6B0E")
	bttcBatchAddress   = common.HexToAddress("0x0c9de531dcb38b758fe8a2c163444a5e54ee0db2")

	// deploy gas
	ethDeploymentGas      = "10"
	tronDeploymentGas     = "10"
	bttcDeploymentGas     = "300000000000000"
	bttcTestDeploymentGas = "300000000000000"
	testDeploymentGas     = "10"

	//endpoint
	ethEndpoint      = ""
	tronEndpoint     = ""
	bttcEndpoint     = "https://rpc.bittorrentchain.io/"
	bttcTestEndpoint = "https://pre-rpc.bt.io/"
	testEndpoint     = "http://18.144.29.246:8110"

	DefaultChain = bttcTestChainID
)

type ChainConfig struct {
	StartBlock         uint64
	CurrentFactory     common.Address
	PriceOracleAddress common.Address
	BatchAddress       common.Address
	DeploymentGas      string
	Endpoint           string
}

func GetChainConfig(chainID int64) (*ChainConfig, bool) {
	var cfg ChainConfig
	switch chainID {
	case ethChainID:
		cfg.StartBlock = ethStartBlock
		cfg.CurrentFactory = ethFactoryAddress
		cfg.PriceOracleAddress = ethOracleAddress
		cfg.DeploymentGas = ethDeploymentGas
		cfg.Endpoint = ethEndpoint
		cfg.BatchAddress = ethBatchAddress
		return &cfg, true
	case tronChainID:
		cfg.StartBlock = tronStartBlock
		cfg.CurrentFactory = tronFactoryAddress
		cfg.PriceOracleAddress = tronOracleAddress
		cfg.DeploymentGas = tronDeploymentGas
		cfg.Endpoint = tronEndpoint
		cfg.BatchAddress = tronBatchAddress
		return &cfg, true
	case bttcChainID:
		cfg.StartBlock = bttcStartBlock
		cfg.CurrentFactory = bttcFactoryAddress
		cfg.PriceOracleAddress = bttcOracleAddress
		cfg.DeploymentGas = bttcDeploymentGas
		cfg.Endpoint = bttcEndpoint
		cfg.BatchAddress = bttcBatchAddress
		return &cfg, true
	case bttcTestChainID:
		cfg.StartBlock = bttcStartBlock
		cfg.CurrentFactory = bttcTestFactoryAddress
		cfg.PriceOracleAddress = bttcTestOracleAddress
		cfg.DeploymentGas = bttcTestDeploymentGas
		cfg.Endpoint = bttcTestEndpoint
		cfg.BatchAddress = bttcTestBatchAddress
		return &cfg, true
	case testChainID:
		cfg.StartBlock = ethStartBlock
		cfg.CurrentFactory = ethFactoryAddress
		cfg.PriceOracleAddress = ethOracleAddress
		cfg.DeploymentGas = testDeploymentGas
		cfg.Endpoint = testEndpoint
		cfg.BatchAddress = ethBatchAddress
		return &cfg, true

	default:
		cfg.StartBlock = bttcStartBlock
		cfg.CurrentFactory = bttcTestFactoryAddress
		cfg.PriceOracleAddress = bttcTestOracleAddress
		cfg.DeploymentGas = bttcTestDeploymentGas
		cfg.Endpoint = bttcTestEndpoint
		cfg.BatchAddress = bttcTestBatchAddress
		return &cfg, true
	}
}
