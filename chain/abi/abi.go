package abi

const VaultABI = `[
		{
			"anonymous": false,
			"inputs": [
				{
					"indexed": false,
					"internalType": "address",
					"name": "previousAdmin",
					"type": "address"
				},
				{
					"indexed": false,
					"internalType": "address",
					"name": "newAdmin",
					"type": "address"
				}
			],
			"name": "AdminChanged",
			"type": "event"
		},
		{
			"anonymous": false,
			"inputs": [
				{
					"indexed": true,
					"internalType": "address",
					"name": "beacon",
					"type": "address"
				}
			],
			"name": "BeaconUpgraded",
			"type": "event"
		},
		{
			"anonymous": false,
			"inputs": [],
			"name": "ChequeBounced",
			"type": "event"
		},
		{
			"anonymous": false,
			"inputs": [
				{
					"indexed": true,
					"internalType": "address",
					"name": "beneficiary",
					"type": "address"
				},
				{
					"indexed": true,
					"internalType": "address",
					"name": "recipient",
					"type": "address"
				},
				{
					"indexed": true,
					"internalType": "address",
					"name": "caller",
					"type": "address"
				},
				{
					"indexed": false,
					"internalType": "uint256",
					"name": "totalPayout",
					"type": "uint256"
				},
				{
					"indexed": false,
					"internalType": "uint256",
					"name": "cumulativePayout",
					"type": "uint256"
				},
				{
					"indexed": false,
					"internalType": "uint256",
					"name": "callerPayout",
					"type": "uint256"
				}
			],
			"name": "ChequeCashed",
			"type": "event"
		},
		{
			"anonymous": false,
			"inputs": [
				{
					"indexed": true,
					"internalType": "address",
					"name": "implementation",
					"type": "address"
				}
			],
			"name": "Upgraded",
			"type": "event"
		},
		{
			"anonymous": false,
			"inputs": [
				{
					"indexed": true,
					"internalType": "address",
					"name": "from",
					"type": "address"
				},
				{
					"indexed": false,
					"internalType": "uint256",
					"name": "amount",
					"type": "uint256"
				}
			],
			"name": "VaultDeposit",
			"type": "event"
		},
		{
			"anonymous": false,
			"inputs": [
				{
					"indexed": true,
					"internalType": "address",
					"name": "from",
					"type": "address"
				},
				{
					"indexed": false,
					"internalType": "uint256",
					"name": "amount",
					"type": "uint256"
				}
			],
			"name": "VaultWithdraw",
			"type": "event"
		},
		{
			"inputs": [],
			"name": "CHEQUE_TYPEHASH",
			"outputs": [
				{
					"internalType": "bytes32",
					"name": "",
					"type": "bytes32"
				}
			],
			"stateMutability": "view",
			"type": "function"
		},
		{
			"inputs": [],
			"name": "EIP712DOMAIN_TYPEHASH",
			"outputs": [
				{
					"internalType": "bytes32",
					"name": "",
					"type": "bytes32"
				}
			],
			"stateMutability": "view",
			"type": "function"
		},
		{
			"inputs": [],
			"name": "bounced",
			"outputs": [
				{
					"internalType": "bool",
					"name": "",
					"type": "bool"
				}
			],
			"stateMutability": "view",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "address",
					"name": "recipient",
					"type": "address"
				},
				{
					"internalType": "uint256",
					"name": "cumulativePayout",
					"type": "uint256"
				},
				{
					"internalType": "bytes",
					"name": "issuerSig",
					"type": "bytes"
				}
			],
			"name": "cashChequeBeneficiary",
			"outputs": [],
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "uint256",
					"name": "amount",
					"type": "uint256"
				}
			],
			"name": "deposit",
			"outputs": [],
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"inputs": [],
			"name": "implementation",
			"outputs": [
				{
					"internalType": "address",
					"name": "impl",
					"type": "address"
				}
			],
			"stateMutability": "view",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "address",
					"name": "_issuer",
					"type": "address"
				},
				{
					"internalType": "address",
					"name": "_token",
					"type": "address"
				}
			],
			"name": "init",
			"outputs": [],
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"inputs": [],
			"name": "issuer",
			"outputs": [
				{
					"internalType": "address",
					"name": "",
					"type": "address"
				}
			],
			"stateMutability": "view",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "address",
					"name": "",
					"type": "address"
				}
			],
			"name": "paidOut",
			"outputs": [
				{
					"internalType": "uint256",
					"name": "",
					"type": "uint256"
				}
			],
			"stateMutability": "view",
			"type": "function"
		},
		{
			"inputs": [],
			"name": "proxiableUUID",
			"outputs": [
				{
					"internalType": "bytes32",
					"name": "",
					"type": "bytes32"
				}
			],
			"stateMutability": "view",
			"type": "function"
		},
		{
			"inputs": [],
			"name": "token",
			"outputs": [
				{
					"internalType": "contract ERC20",
					"name": "",
					"type": "address"
				}
			],
			"stateMutability": "view",
			"type": "function"
		},
		{
			"inputs": [],
			"name": "totalPaidOut",
			"outputs": [
				{
					"internalType": "uint256",
					"name": "",
					"type": "uint256"
				}
			],
			"stateMutability": "view",
			"type": "function"
		},
		{
			"inputs": [],
			"name": "totalbalance",
			"outputs": [
				{
					"internalType": "uint256",
					"name": "",
					"type": "uint256"
				}
			],
			"stateMutability": "view",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "address",
					"name": "newImplementation",
					"type": "address"
				}
			],
			"name": "upgradeTo",
			"outputs": [],
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "address",
					"name": "newImplementation",
					"type": "address"
				},
				{
					"internalType": "bytes",
					"name": "data",
					"type": "bytes"
				}
			],
			"name": "upgradeToAndCall",
			"outputs": [],
			"stateMutability": "payable",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "uint256",
					"name": "amount",
					"type": "uint256"
				}
			],
			"name": "withdraw",
			"outputs": [],
			"stateMutability": "nonpayable",
			"type": "function"
		}
]`

const VaultFactoryABI = `[
		{
			"inputs": [
				{
					"internalType": "address",
					"name": "_TokenAddress",
					"type": "address"
				}
			],
			"stateMutability": "nonpayable",
			"type": "constructor"
		},
		{
			"anonymous": false,
			"inputs": [
				{
					"indexed": false,
					"internalType": "address",
					"name": "issuer",
					"type": "address"
				},
				{
					"indexed": false,
					"internalType": "address",
					"name": "contractAddress",
					"type": "address"
				},
				{
					"indexed": false,
					"internalType": "string",
					"name": "id",
					"type": "string"
				}
			],
			"name": "VaultDeployed",
			"type": "event"
		},
		{
			"inputs": [],
			"name": "TokenAddress",
			"outputs": [
				{
					"internalType": "address",
					"name": "",
					"type": "address"
				}
			],
			"stateMutability": "view",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "address",
					"name": "issuer",
					"type": "address"
				},
				{
					"internalType": "address",
					"name": "_logic",
					"type": "address"
				},
				{
					"internalType": "bytes32",
					"name": "salt",
					"type": "bytes32"
				},
				{
					"internalType": "string",
					"name": "id",
					"type": "string"
				},
				{
					"internalType": "bytes",
					"name": "_data",
					"type": "bytes"
				}
			],
			"name": "deployVault",
			"outputs": [
				{
					"internalType": "address",
					"name": "",
					"type": "address"
				}
			],
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "address",
					"name": "",
					"type": "address"
				}
			],
			"name": "deployedContracts",
			"outputs": [
				{
					"internalType": "bool",
					"name": "",
					"type": "bool"
				}
			],
			"stateMutability": "view",
			"type": "function"
		},
		{
			"inputs": [],
			"name": "master",
			"outputs": [
				{
					"internalType": "address",
					"name": "",
					"type": "address"
				}
			],
			"stateMutability": "view",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "string",
					"name": "",
					"type": "string"
				}
			],
			"name": "peerVaultAddress",
			"outputs": [
				{
					"internalType": "address",
					"name": "",
					"type": "address"
				}
			],
			"stateMutability": "view",
			"type": "function"
		}
]`

const Erc20ABI = `[
	{
		"constant": true,
		"inputs": [],
		"name": "name",
		"outputs": [
			{
				"name": "",
				"type": "string"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "guy",
				"type": "address"
			},
			{
				"name": "wad",
				"type": "uint256"
			}
		],
		"name": "approve",
		"outputs": [
			{
				"name": "",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "totalSupply",
		"outputs": [
			{
				"name": "",
				"type": "uint256"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "src",
				"type": "address"
			},
			{
				"name": "dst",
				"type": "address"
			},
			{
				"name": "wad",
				"type": "uint256"
			}
		],
		"name": "transferFrom",
		"outputs": [
			{
				"name": "",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "wad",
				"type": "uint256"
			}
		],
		"name": "withdraw",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "decimals",
		"outputs": [
			{
				"name": "",
				"type": "uint8"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "",
				"type": "address"
			}
		],
		"name": "balanceOf",
		"outputs": [
			{
				"name": "",
				"type": "uint256"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "symbol",
		"outputs": [
			{
				"name": "",
				"type": "string"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "dst",
				"type": "address"
			},
			{
				"name": "wad",
				"type": "uint256"
			}
		],
		"name": "transfer",
		"outputs": [
			{
				"name": "",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [],
		"name": "deposit",
		"outputs": [],
		"payable": true,
		"stateMutability": "payable",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "",
				"type": "address"
			},
			{
				"name": "",
				"type": "address"
			}
		],
		"name": "allowance",
		"outputs": [
			{
				"name": "",
				"type": "uint256"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"payable": true,
		"stateMutability": "payable",
		"type": "fallback"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"name": "src",
				"type": "address"
			},
			{
				"indexed": true,
				"name": "guy",
				"type": "address"
			},
			{
				"indexed": false,
				"name": "wad",
				"type": "uint256"
			}
		],
		"name": "Approval",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"name": "src",
				"type": "address"
			},
			{
				"indexed": true,
				"name": "dst",
				"type": "address"
			},
			{
				"indexed": false,
				"name": "wad",
				"type": "uint256"
			}
		],
		"name": "Transfer",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"name": "dst",
				"type": "address"
			},
			{
				"indexed": false,
				"name": "wad",
				"type": "uint256"
			}
		],
		"name": "Deposit",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"name": "src",
				"type": "address"
			},
			{
				"indexed": false,
				"name": "wad",
				"type": "uint256"
			}
		],
		"name": "Withdrawal",
		"type": "event"
	}
]`

// Current VaultFactory bytecode
const FactoryDeployedBin = "0x608060405234801561001057600080fd5b50600436106100575760003560e01c80633695ddb21461005c578063c2cba3061461008c578063c70242ad1461009f578063e4e36723146100d2578063ee97f7f314610106575b600080fd5b61006f61006a366004610461565b610119565b6040516001600160a01b0390911681526020015b60405180910390f35b60035461006f906001600160a01b031681565b6100c26100ad366004610447565b60006020819052908152604090205460ff1681565b6040519015158152602001610083565b61006f6100e0366004610500565b80516020818301810180516001825292820191909301209152546001600160a01b031681565b60025461006f906001600160a01b031681565b6000806001600160a01b03166001846040516101359190610567565b908152604051908190036020019020546001600160a01b0316146101985760405162461bcd60e51b81526020600482015260156024820152741d985d5b1d08185b195c98591e4818dc99585d1959605a1b60448201526064015b60405180910390fd5b600254604080513360208201529081018690526000916101dc916001600160a01b0390911690606001604051602081830303815290604052805190602001206102e9565b60405163c0d91eaf60e01b81529091506001600160a01b0382169063c0d91eaf9061020d90899087906004016105b8565b600060405180830381600087803b15801561022757600080fd5b505af115801561023b573d6000803e3d6000fd5b5050506001600160a01b03821660009081526020819052604090819020805460ff191660019081179091559051839250610276908790610567565b90815260405190819003602001812080546001600160a01b03939093166001600160a01b0319909316929092179091557fb2e029ab52a406bb52fae42d1d1e5dca2793900c250a242af032038b02eb8659906102d790899084908890610583565b60405180910390a19695505050505050565b6000604051733d602d80600a3d3981f3363d3d373d3d3d363d7360601b81528360601b60148201526e5af43d82803e903d91602b57fd5bf360881b6028820152826037826000f59150506001600160a01b0381166103895760405162461bcd60e51b815260206004820152601760248201527f455243313136373a2063726561746532206661696c6564000000000000000000604482015260640161018f565b92915050565b600067ffffffffffffffff808411156103aa576103aa61060c565b604051601f8501601f19908116603f011681019082821181831017156103d2576103d261060c565b816040528093508581528686860111156103eb57600080fd5b858560208301376000602087830101525050509392505050565b80356001600160a01b038116811461041c57600080fd5b919050565b600082601f830112610431578081fd5b6104408383356020850161038f565b9392505050565b600060208284031215610458578081fd5b61044082610405565b600080600080600060a08688031215610478578081fd5b61048186610405565b945061048f60208701610405565b935060408601359250606086013567ffffffffffffffff808211156104b2578283fd5b6104be89838a01610421565b935060808801359150808211156104d3578283fd5b508601601f810188136104e4578182fd5b6104f38882356020840161038f565b9150509295509295909350565b600060208284031215610511578081fd5b813567ffffffffffffffff811115610527578182fd5b61053384828501610421565b949350505050565b600081518084526105538160208601602086016105dc565b601f01601f19169290920160200192915050565b600082516105798184602087016105dc565b9190910192915050565b6001600160a01b038481168252831660208201526060604082018190526000906105af9083018461053b565b95945050505050565b6001600160a01b03831681526040602082018190526000906105339083018461053b565b60005b838110156105f75781810151838201526020016105df565b83811115610606576000848401525b50505050565b634e487b7160e01b600052604160045260246000fdfea264697066735822122059d06d999e706dafb3c116ebf57456505059026b103b38cf7520feb64ead230b64736f6c63430008030033"

const FactoryDeployedBinV1 = "0x608060405234801561001057600080fd5b50600436106100575760003560e01c80633695ddb21461005c578063c2cba3061461008c578063c70242ad1461009f578063e4e36723146100d2578063ee97f7f314610106575b600080fd5b61006f61006a3660046103e6565b610119565b6040516001600160a01b0390911681526020015b60405180910390f35b60035461006f906001600160a01b031681565b6100c26100ad3660046103cc565b60006020819052908152604090205460ff1681565b6040519015158152602001610083565b61006f6100e0366004610485565b80516020818301810180516001825292820191909301209152546001600160a01b031681565b60025461006f906001600160a01b031681565b60025460408051336020820152908101859052600091829161015d916001600160a01b0316906060016040516020818303038152906040528051906020012061026a565b60405163c0d91eaf60e01b81529091506001600160a01b0382169063c0d91eaf9061018e908990879060040161053d565b600060405180830381600087803b1580156101a857600080fd5b505af11580156101bc573d6000803e3d6000fd5b5050506001600160a01b03821660009081526020819052604090819020805460ff1916600190811790915590518392506101f79087906104ec565b90815260405190819003602001812080546001600160a01b03939093166001600160a01b0319909316929092179091557fb2e029ab52a406bb52fae42d1d1e5dca2793900c250a242af032038b02eb86599061025890899084908890610508565b60405180910390a19695505050505050565b6000604051733d602d80600a3d3981f3363d3d373d3d3d363d7360601b81528360601b60148201526e5af43d82803e903d91602b57fd5bf360881b6028820152826037826000f59150506001600160a01b03811661030e5760405162461bcd60e51b815260206004820152601760248201527f455243313136373a2063726561746532206661696c6564000000000000000000604482015260640160405180910390fd5b92915050565b600067ffffffffffffffff8084111561032f5761032f610591565b604051601f8501601f19908116603f0116810190828211818310171561035757610357610591565b8160405280935085815286868601111561037057600080fd5b858560208301376000602087830101525050509392505050565b80356001600160a01b03811681146103a157600080fd5b919050565b600082601f8301126103b6578081fd5b6103c583833560208501610314565b9392505050565b6000602082840312156103dd578081fd5b6103c58261038a565b600080600080600060a086880312156103fd578081fd5b6104068661038a565b94506104146020870161038a565b935060408601359250606086013567ffffffffffffffff80821115610437578283fd5b61044389838a016103a6565b93506080880135915080821115610458578283fd5b508601601f81018813610469578182fd5b61047888823560208401610314565b9150509295509295909350565b600060208284031215610496578081fd5b813567ffffffffffffffff8111156104ac578182fd5b6104b8848285016103a6565b949350505050565b600081518084526104d8816020860160208601610561565b601f01601f19169290920160200192915050565b600082516104fe818460208701610561565b9190910192915050565b6001600160a01b03848116825283166020820152606060408201819052600090610534908301846104c0565b95945050505050565b6001600160a01b03831681526040602082018190526000906104b8908301846104c0565b60005b8381101561057c578181015183820152602001610564565b8381111561058b576000848401525b50505050565b634e487b7160e01b600052604160045260246000fdfea26469706673582212208c36175c3335854e47a8af80ec2a9b56ffcbd20d6864139d65f9b167c8c7d64164736f6c63430008030033"
const FactoryDeployedBinV2 = "0x608060405234801561001057600080fd5b50600436106100575760003560e01c80633695ddb21461005c578063c2cba3061461008c578063c70242ad1461009f578063e4e36723146100d2578063ee97f7f314610106575b600080fd5b61006f61006a366004610461565b610119565b6040516001600160a01b0390911681526020015b60405180910390f35b60035461006f906001600160a01b031681565b6100c26100ad366004610447565b60006020819052908152604090205460ff1681565b6040519015158152602001610083565b61006f6100e0366004610500565b80516020818301810180516001825292820191909301209152546001600160a01b031681565b60025461006f906001600160a01b031681565b6000806001600160a01b03166001846040516101359190610567565b908152604051908190036020019020546001600160a01b0316146101985760405162461bcd60e51b81526020600482015260156024820152741d985d5b1d08185b195c98591e4818dc99585d1959605a1b60448201526064015b60405180910390fd5b600254604080513360208201529081018690526000916101dc916001600160a01b0390911690606001604051602081830303815290604052805190602001206102e9565b60405163c0d91eaf60e01b81529091506001600160a01b0382169063c0d91eaf9061020d90899087906004016105b8565b600060405180830381600087803b15801561022757600080fd5b505af115801561023b573d6000803e3d6000fd5b5050506001600160a01b03821660009081526020819052604090819020805460ff191660019081179091559051839250610276908790610567565b90815260405190819003602001812080546001600160a01b03939093166001600160a01b0319909316929092179091557fb2e029ab52a406bb52fae42d1d1e5dca2793900c250a242af032038b02eb8659906102d790899084908890610583565b60405180910390a19695505050505050565b6000604051733d602d80600a3d3981f3363d3d373d3d3d363d7360601b81528360601b60148201526e5af43d82803e903d91602b57fd5bf360881b6028820152826037826000f59150506001600160a01b0381166103895760405162461bcd60e51b815260206004820152601760248201527f455243313136373a2063726561746532206661696c6564000000000000000000604482015260640161018f565b92915050565b600067ffffffffffffffff808411156103aa576103aa61060c565b604051601f8501601f19908116603f011681019082821181831017156103d2576103d261060c565b816040528093508581528686860111156103eb57600080fd5b858560208301376000602087830101525050509392505050565b80356001600160a01b038116811461041c57600080fd5b919050565b600082601f830112610431578081fd5b6104408383356020850161038f565b9392505050565b600060208284031215610458578081fd5b61044082610405565b600080600080600060a08688031215610478578081fd5b61048186610405565b945061048f60208701610405565b935060408601359250606086013567ffffffffffffffff808211156104b2578283fd5b6104be89838a01610421565b935060808801359150808211156104d3578283fd5b508601601f810188136104e4578182fd5b6104f38882356020840161038f565b9150509295509295909350565b600060208284031215610511578081fd5b813567ffffffffffffffff811115610527578182fd5b61053384828501610421565b949350505050565b600081518084526105538160208601602086016105dc565b601f01601f19169290920160200192915050565b600082516105798184602087016105dc565b9190910192915050565b6001600160a01b038481168252831660208201526060604082018190526000906105af9083018461053b565b95945050505050565b6001600160a01b03831681526040602082018190526000906105339083018461053b565b60005b838110156105f75781810151838201526020016105df565b83811115610606576000848401525b50505050565b634e487b7160e01b600052604160045260246000fdfea264697066735822122059d06d999e706dafb3c116ebf57456505059026b103b38cf7520feb64ead230b64736f6c63430008030033"
