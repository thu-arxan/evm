pragma solidity ^0.6.0;

contract OpCode {
    function signExtend() public pure returns (uint r) {
        assembly { r := signextend(10, 1) }
	return r;
    }
    function signExtend2() public pure returns (uint r) {
        assembly { r := signextend(35, 1) }
    }
    // not test yet
    function callDataCopy() public pure returns (string memory r) {
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

    // not test yet
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

    function mStore8() public pure returns (uint s) {
        assembly {
            mstore8(0, 128888)
            s := mload(0)
        }
    }

    function testPC() public pure returns (uint p) {
        assembly {
            p := pc()
        }
    }

    function testMSize() public pure returns (uint s) {
            assembly {
                s := msize()
            }
    }

    function testGas() public view returns (uint g) {
        assembly {
            g := gas()
        }
    }

    function testRevert() public pure {
        assembly {
            revert(0, 0)
        }
    }

    function testInvalid() public pure {
        assembly {
            invalid()
        }
    }

    function testStop() public pure {
        assembly {
            stop()
        }
    }
}
