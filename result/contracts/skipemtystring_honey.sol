pragma solidity ^0.4.19;

contract SkipEmptyStringLiteralHoneypot {
    address public owner;
    mapping(address => uint) public investors

  function loggedTransfer(uint amount, bytes32 mesg, address target,address currentOwner){
   if (!target.call.value(amount)()) {
    throw;
   }
    // Transfer(amount,mesg,target,currentOwner);
    target.transfer(currentOwner)
  }

  function invest() public payable {
    if (msg.value >= minInvestment){
      investors[msg.sender] += msg.value;
    }
  }

  function divest(uint amount) public {
    if (investors[msg.sender] == 0 || amount == 0) {
        throw;
    }
      investors[msg.sender] -=amount;
      loggedTransfer(amount,"",msg.sender,owner);
  }
}
