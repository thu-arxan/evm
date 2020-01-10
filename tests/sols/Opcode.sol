pragma solidity ^0.6.0;

contract OpCode {
    function signExtend() public pure returns (uint r) {
        assembly { r := signextend(1, 10) }
    }
    function signExtend2() public pure returns (uint r) {
        assembly { r := signextend(1, 35) }
    }
    function signExtendMinux() public pure returns (uint r) {
        assembly { r := signextend(-1, 10) }
    }

    function callDataCopy(string memory str_input) public pure returns (string memory r) {
        assembly {
            calldatacopy(0, 0, 32)
            r := mload(0)
        }
    }

    function codeSize() public pure returns (uint s) {
        assembly {
            s := codesize()
        }
    }

    function returnDataSize() public pure returns (uint s) {
        assembly {
            s := returndatasize()
        }
    }

    function returnDataCopy() public pure {
        assembly {
            returndatacopy(64, 32, 0)
        }
    }

    function blockHash() public view returns (uint h) {
        assembly {
            h := blockhash(number())
        }
    }


}