pragma solidity ^0.4.19;

contract Gift_1_ETH {
  bool passHasBeenSet = false;
  bytes32 hashPass;

  function SetPass(bytes32 hash) payable {
    if (!passHasBeenSet && (msg.value >= 1 ether))
      hashPass = hash;
  }
  function GetGift(bytes pass) returns (bytes32){
    if (hashPass == sha3(pass))
      msg.sender.transfer(this.balance);
    return sha3(pass);
  }
  function PassHasBeenSet(bytes32 hash) {
    if (hash == hashPass) passHasBeenSet = true;
  }
}
