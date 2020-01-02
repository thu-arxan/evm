pragma solidity >=0.4.0 <0.7.0;

contract SimpleBlock {
    function getAddress() public view returns (address) {
        return address(this);
    }
    function getBalance() public view returns (uint) {
        return address(this).balance;
    }
    function getOrigin() public view returns (address) {
        return tx.origin;
    }
    function getGasprice() public view returns (uint) {
        return tx.gasprice;
    }
    function getCoinbase() public view returns (address) {
        return block.coinbase;
    }
    function getTimestamp() public view returns (uint) {
        return block.timestamp;

    }
    function getNumber() public view returns (uint) {
        return block.number;
    }
    function getDifficulty() public view returns (uint) {
        return block.difficulty;
    }
    function getGaslimit() public view returns (uint) {
        return block.gaslimit;
    }
    function getChainID() external view returns (uint256) {
        uint256 id;
        assembly {
            id := chainid()
        }
        return id;
    }
}