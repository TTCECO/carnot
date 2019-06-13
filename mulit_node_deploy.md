carnot 多节点部署

1. 添加验证人
carnotcli keys add validator1 --home node1/carnotcli
carnotcli keys add validator2 --home node2/carnotcli
carnotcli keys add validator3 --home node3/carnotcli

2. 初始化三个节点
carnot init  testing --chain-id=testing  --home node1/carnot --p2p.laddr tcp://0.0.0.0:26656 --rpc.laddr tcp://0.0.0.0:26657
carnot init  testing --chain-id=testing  --home node2/carnot --p2p.laddr tcp://0.0.0.0:26659 --rpc.laddr tcp://0.0.0.0:26660
carnot init  testing --chain-id=testing  --home node3/carnot --p2p.laddr tcp://0.0.0.0:26661 --rpc.laddr tcp://0.0.0.0:26662

3. 每个节点创世块增加三个验证人
carnot add-genesis-account --home node1/carnot $(carnotcli keys show validator1 -a --home node1/carnotcli) 1000000000stake;
carnot add-genesis-account --home node1/carnot $(carnotcli keys show validator2 -a --home node2/carnotcli) 1000000000stake;
carnot add-genesis-account --home node1/carnot $(carnotcli keys show validator3 -a --home node3/carnotcli) 1000000000stake;

carnot add-genesis-account --home node2/carnot $(carnotcli keys show validator1 -a --home node1/carnotcli) 1000000000stake;
carnot add-genesis-account --home node2/carnot $(carnotcli keys show validator2 -a --home node2/carnotcli) 1000000000stake;
carnot add-genesis-account --home node2/carnot $(carnotcli keys show validator3 -a --home node3/carnotcli) 1000000000stake;

carnot add-genesis-account --home node3/carnot $(carnotcli keys show validator1 -a --home node1/carnotcli) 1000000000stake;
carnot add-genesis-account --home node3/carnot $(carnotcli keys show validator2 -a --home node2/carnotcli) 1000000000stake;
carnot add-genesis-account --home node3/carnot $(carnotcli keys show validator3 -a --home node3/carnotcli) 1000000000stake;

4. 生成三笔交易
mkdir gentxs
carnot gentx --name validator1 --memo.port 26656 --home node1/carnot --home-client node1/carnotcli --output-document gentxs/node1.json
carnot gentx --name validator2 --memo.port 26659 --home node2/carnot --home-client node2/carnotcli --output-document gentxs/node2.json
carnot gentx --name validator3 --memo.port 26661 --home node3/carnot --home-client node3/carnotcli --output-document gentxs/node3.json

5. 把三笔交易写到创世块
carnot collect-gentxs --home node1/carnot --gentx-dir gentxs/

6. 拷贝node1创世文件到其他节点
cp node1/carnot/config/genesis.json node2/carnot/config/
cp node1/carnot/config/genesis.json node3/carnot/config/

7. 修改 carnot/config/config.toml (如果需要在同一台服务器上部署多个节点)
 allow_duplicate_ip 为 true  （188行）
 prof_laddr 为不同的端口 "localhost:6060"  （52行）
 persistent_peers 为真实值，用逗号分隔 （145）

13cb778bec971708a07964ce95d0fc7b6fb548cb@192.168.80.197:26661,d3cf51ce6640654eebde030e1b71accca6939ed6@192.168.80.197:26659,5288486b9494f753937269b29b4ecb1d8d40fffb@192.168.80.197:26656

carnot cc-start keyfile_1.json 1 validator1 11111111 --home node1/carnot --home-client node1/carnotcli
carnot cc-start keyfile_2.json 1 validator2 22222222 --home node2/carnot --home-client node2/carnotcli
carnot cc-start keyfile_3.json 1 validator3 33333333 --home node3/carnot --home-client node3/carnotcli


8. 配置cli
carnotcli --home node1/carnotcli config chain-id testing
carnotcli --home node1/carnotcli config output json
carnotcli --home node1/carnotcli config indent true
carnotcli --home node1/carnotcli config trust-node true


carnotcli --home node2/carnotcli config chain-id testing
carnotcli --home node2/carnotcli config output json
carnotcli --home node2/carnotcli config indent true
carnotcli --home node2/carnotcli config trust-node true


carnotcli --home node3/carnotcli config chain-id testing
carnotcli --home node3/carnotcli config output json
carnotcli --home node3/carnotcli config indent true
carnotcli --home node3/carnotcli config trust-node true

9. 显示账户信息
carnotcli --home node1/carnotcli keys show validator1
carnotcli --home node2/carnotcli keys show validator2
carnotcli --home node3/carnotcli keys show validator3

carnotcli --home node1/carnotcli q account $(carnotcli keys show validator1 -a --home node1/carnotcli)
carnotcli --home node2/carnotcli q account $(carnotcli keys show validator2 -a --home node2/carnotcli)
carnotcli --home node3/carnotcli q account $(carnotcli keys show validator3 -a --home node3/carnotcli)

10. 提币 withdraw ttc 从ttc主链到cosmos



11. 充值 deposit

carnotcli --home node1/carnotcli tx tcchan deposit t0c233eC8cB98133Bf202DcBAF07112C6Abb058B89 2ttc --from validator1
carnotcli --home node2/carnotcli tx tcchan deposit t0c233eC8cB98133Bf202DcBAF07112C6Abb058B89 3ttc --from validator2
carnotcli --home node3/carnotcli tx tcchan deposit t0c233eC8cB98133Bf202DcBAF07112C6Abb058B89 4ttc --from validator3


