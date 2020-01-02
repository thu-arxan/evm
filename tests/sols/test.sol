pragma solidity >=0.4.0 <0.7.0;

contract test {
function getChainID() external view returns (uint256) {
    uint256 id;
    assembly {
        id := chainid()
    }
    return id;
}
}