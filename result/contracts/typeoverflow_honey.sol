pragma solidity ^0.4.19;

contract TypeDeductionOverflowHoneypot {
    uint8 public balance;

    function TypeDeductionOverflowHoneypot() public {
        balance = 100;
    }

  function Test() payable public {
    if (msg.value > 0.1 ether) {
      uint256 multi = 0;
      uint256 amountToTransfer = 0;
      for (var i = 0; i < 2 * msg.value; i++) {
        multi = i*2;
        if (multi < amountToTransfer) break;
        amountToTransfer = multi;
      }
      msg.sender.transfer(amountToTransfer);
    }
  }
}

