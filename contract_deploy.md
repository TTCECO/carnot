# TTC合约部署及初始设定

1. 编译合约，生成ABI及Bytecode

编译生成ABI和Bytecode ， 将bytecode中的内容记录在 carnot/x/tcchan/contract/doc.go当中

2. 将ABI中的内容拷贝到 carnot/x/tcchan/contract/tcchan.abi中

3. 将tcchan.abi编译成go文件
cd carnot/x/tcchan/contract/
abigen --abi tcchan.abi --pkg contract --type Contract --out contract.go

4. 将Bytecode文件中object字段中的内容部署到主链

eth.sendTransaction({from:your_address,gas:9999999,data:"0x606060405260006001556003600255620f4240600355336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506118b8806100646000396000f3006060604052600436106100f1576000357c......"})

5. 根据hash值，检查是否部署成功

eth.getTransactionReceipt("t03f4fb3e135498a40784e3e6d27ce32184ad5961c2d6137a4796f92ef20aec5ab").status 是否为 "0x1"

6. 记录合约地址

eth.getTransactionReceipt("t03f4fb3e135498a40784e3e6d27ce32184ad5961c2d6137a4796f92ef20aec5ab").contractAddress
t0e04b30b03f4371958ecd29bcdbf0d8b4493ec344
并将得到的结果，保存到 carnot/x/tcchan/params.go 中给 contractAddress 赋值。

7. 创建多个ttc地址，给cosmos的见证人使用
t007573C3F5c21373B3430998F809BCFDAca38Fe28  (keyfile_1.json)
t0B7c4565B1210054CAc3a0F08eD4BD631ec1C8cC9  (keyfile_1.json)
t0cC2a7F0a041e0975c0B7854364e154cdA059a9F0  (keyfile_1.json)

8. 初始化设定
myabi = eth.contract(abi文件内容)
co = myabi.at("t0e04b30b03f4371958ecd29bcdbf0d8b4493ec344")

co.setMinConfirmNum(2,{from:eth.accounts[0]})
co.minConfirmNum() // 检查结果
注意，同时应该 carnot/x/tcchan/params.go 中给minValidatorCount赋值


co.addValidator("t007573C3F5c21373B3430998F809BCFDAca38Fe28",{from:eth.accounts[0]})
co.validators("t007573C3F5c21373B3430998F809BCFDAca38Fe28") // 检查结果

co.ownerChargeFund({from:eth.accounts[0],value:web3.toWei(200)}) // fortest
给三个地址分别充值，当gas

eth.sendTransaction({from:eth.accounts[0],to:"t007573C3F5c21373B3430998F809BCFDAca38Fe28",value:web3.toWei(100)})
eth.getBalance("t0cC2a7F0a041e0975c0B7854364e154cdA059a9F0") // 检查结果



