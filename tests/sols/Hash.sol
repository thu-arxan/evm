pragma solidity ^0.6.0;

contract Hash {

    function SHA256(string memory input) public pure returns (bytes32) {
        bytes32 result = sha256(abi.encode(input));
        return result;
    }

    function KECCAK256(string memory input) public pure returns (bytes32) {
        bytes32 result = keccak256(abi.encode(input));
        return result;
    }
}