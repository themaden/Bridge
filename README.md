# 🌉 BTC-ETH Liquidity Bridge (Prototype)

![Go](https://img.shields.io/badge/Go-1.21-blue) ![Solidity](https://img.shields.io/badge/Solidity-0.8-gray) ![Ethereum](https://img.shields.io/badge/Network-Anvil%2FLocal-green)

A decentralized cross-chain bridge prototype that monitors the **Bitcoin Network** for new blocks and automatically mints **synthetic xBTC tokens** on the **Ethereum Network** (Local Anvil Chain).

Built with **Golang (Geth)** for the relayer service and **Solidity** for the destination smart contracts.

---

## 🏗 Architecture

The system consists of three main components:

1.  **The Watcher (Golang):** Monitors the Bitcoin mainnet via mempool.space API for new blocks.
2.  **The Relayer (Golang):** Validates the data and signs a transaction using ECDSA.
3.  **The Vault (Solidity):** An ERC-20 compliant smart contract that mints tokens upon receiving valid signals.

### Workflow
`[Bitcoin Network] -> (New Block) -> [Go Relayer] -> (Tx Sign) -> [Ethereum Smart Contract] -> (Mint xBTC)`

---

## 🚀 Tech Stack

* **Backend / Relayer:** Golang, go-ethereum (Geth) library.
* **Smart Contracts:** Solidity (ERC-20 Standard).
* **Local Blockchain:** Foundry (Anvil).
* **Data Source:** Mempool.space API (Bitcoin).

---

## 🛠 Installation & Run

### Prerequisites
* Go (1.19+)
* Foundry (Anvil)
* VS Code

### 1. Start Local Ethereum Chain
```bash
anvil
