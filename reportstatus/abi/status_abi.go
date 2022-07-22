package abi

const StatusHeartABI = `[
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
		"inputs": [
			{
				"indexed": false,
				"internalType": "uint8",
				"name": "version",
				"type": "uint8"
			}
		],
		"name": "Initialized",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"internalType": "address",
				"name": "previousOwner",
				"type": "address"
			},
			{
				"indexed": true,
				"internalType": "address",
				"name": "newOwner",
				"type": "address"
			}
		],
		"name": "OwnershipTransferred",
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
				"indexed": false,
				"internalType": "address",
				"name": "lastSignAddress",
				"type": "address"
			},
			{
				"indexed": false,
				"internalType": "address",
				"name": "currentSignAddress",
				"type": "address"
			}
		],
		"name": "signAddressChanged",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"internalType": "string",
				"name": "peer",
				"type": "string"
			},
			{
				"indexed": false,
				"internalType": "uint32",
				"name": "createTime",
				"type": "uint32"
			},
			{
				"indexed": false,
				"internalType": "string",
				"name": "version",
				"type": "string"
			},
			{
				"indexed": false,
				"internalType": "uint32",
				"name": "Nonce",
				"type": "uint32"
			},
			{
				"indexed": false,
				"internalType": "address",
				"name": "bttcAddress",
				"type": "address"
			},
			{
				"indexed": false,
				"internalType": "uint32",
				"name": "signedTime",
				"type": "uint32"
			},
			{
				"indexed": false,
				"internalType": "uint32",
				"name": "lastNonce",
				"type": "uint32"
			},
			{
				"indexed": false,
				"internalType": "uint32",
				"name": "nowTime",
				"type": "uint32"
			},
			{
				"indexed": false,
				"internalType": "uint16[30]",
				"name": "hearts",
				"type": "uint16[30]"
			}
		],
		"name": "statusReported",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"internalType": "string",
				"name": "currentVersion",
				"type": "string"
			},
			{
				"indexed": false,
				"internalType": "string",
				"name": "version",
				"type": "string"
			}
		],
		"name": "versionChanged",
		"type": "event"
	},
	{
		"inputs": [],
		"name": "currentVersion",
		"outputs": [
			{
				"internalType": "string",
				"name": "",
				"type": "string"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "string",
				"name": "peer",
				"type": "string"
			},
			{
				"internalType": "uint32",
				"name": "createTime",
				"type": "uint32"
			},
			{
				"internalType": "string",
				"name": "version",
				"type": "string"
			},
			{
				"internalType": "uint32",
				"name": "Nonce",
				"type": "uint32"
			},
			{
				"internalType": "address",
				"name": "bttcAddress",
				"type": "address"
			},
			{
				"internalType": "uint32",
				"name": "signedTime",
				"type": "uint32"
			}
		],
		"name": "genHashExt",
		"outputs": [
			{
				"internalType": "bytes32",
				"name": "",
				"type": "bytes32"
			}
		],
		"stateMutability": "pure",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "getHighScoreHost",
		"outputs": [
			{
				"components": [
					{
						"internalType": "uint32",
						"name": "createTime",
						"type": "uint32"
					},
					{
						"internalType": "string",
						"name": "version",
						"type": "string"
					},
					{
						"internalType": "uint32",
						"name": "lastNonce",
						"type": "uint32"
					},
					{
						"internalType": "uint32",
						"name": "lastSignedTime",
						"type": "uint32"
					},
					{
						"internalType": "bytes",
						"name": "lastSigned",
						"type": "bytes"
					},
					{
						"internalType": "uint16[30]",
						"name": "hearts",
						"type": "uint16[30]"
					}
				],
				"internalType": "struct BtfsStatus.info[]",
				"name": "",
				"type": "tuple[]"
			}
		],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "string",
				"name": "peer",
				"type": "string"
			}
		],
		"name": "getStatus",
		"outputs": [
			{
				"internalType": "string",
				"name": "",
				"type": "string"
			},
			{
				"internalType": "uint32",
				"name": "",
				"type": "uint32"
			},
			{
				"internalType": "string",
				"name": "",
				"type": "string"
			},
			{
				"internalType": "uint32",
				"name": "",
				"type": "uint32"
			},
			{
				"internalType": "uint32",
				"name": "",
				"type": "uint32"
			},
			{
				"internalType": "bytes",
				"name": "",
				"type": "bytes"
			},
			{
				"internalType": "uint16[30]",
				"name": "",
				"type": "uint16[30]"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "signAddress",
				"type": "address"
			}
		],
		"name": "initialize",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "owner",
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
		"inputs": [
			{
				"internalType": "bytes32",
				"name": "hash",
				"type": "bytes32"
			},
			{
				"internalType": "bytes",
				"name": "sig",
				"type": "bytes"
			}
		],
		"name": "recoverSignerExt",
		"outputs": [
			{
				"internalType": "address",
				"name": "",
				"type": "address"
			}
		],
		"stateMutability": "pure",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "renounceOwnership",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "string",
				"name": "peer",
				"type": "string"
			},
			{
				"internalType": "uint32",
				"name": "createTime",
				"type": "uint32"
			},
			{
				"internalType": "string",
				"name": "version",
				"type": "string"
			},
			{
				"internalType": "uint32",
				"name": "Nonce",
				"type": "uint32"
			},
			{
				"internalType": "address",
				"name": "bttcAddress",
				"type": "address"
			},
			{
				"internalType": "uint32",
				"name": "signedTime",
				"type": "uint32"
			},
			{
				"internalType": "bytes",
				"name": "signed",
				"type": "bytes"
			}
		],
		"name": "reportStatus",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "string",
				"name": "ver",
				"type": "string"
			}
		],
		"name": "setCurrentVersion",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "addr",
				"type": "address"
			}
		],
		"name": "setSignAddress",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "totalStat",
		"outputs": [
			{
				"internalType": "uint64",
				"name": "total",
				"type": "uint64"
			},
			{
				"internalType": "uint64",
				"name": "totalUsers",
				"type": "uint64"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "newOwner",
				"type": "address"
			}
		],
		"name": "transferOwnership",
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
	}
]`
