# btctxbuilder
A Bitcoin toolkit with your own keys, independent of Wallets.

## Quick Start
```bash
make run
```

## Features
- Generate addresses
- Build and sign transactions
- Broadcast transactions

## Supported Transaction Types
| Account Type | Generate Account   | Send Transaction |
|--------------|--------------------|------------------|
| P2PK         | ✅                 | ✅              |
| P2PKH        | ✅                 | ✅              |
| P2WPKH       | ✅                 | ✅              |
| NP2WPKH      | ✅                 | ❌              |
| P2TR(Spend)  | ✅                 | ✅              |

## Network Support
- Bitcoin Mainnet
- Bitcoin Testnet3
- Bitcoin Testnet4
- Bitcoin Signet

## Contributing
Contributions are always welcome!  
If you find a bug, have a feature idea, or just want to improve the project, feel free to open an issue or submit a pull request.