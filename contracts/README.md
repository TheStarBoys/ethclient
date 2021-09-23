# contracts

## How to generate Test contract code?
```bash
solc -o . --abi --bin --overwrite ./test.sol
abigen --abi Test.abi --bin Test.bin --pkg=contracts --out=test_contract.go
```