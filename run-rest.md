# Run REST routes

Now that you tested your CLI queries and transactions, time to test same things in the REST server. Leave the `tcd` that you had running earlier and start by gathering your addresses:

```bash
$ tccli keys show jack --address
```

Now its time to start the `rest-server` in another terminal window:

```bash
$ tccli rest-server --chain-id tctestchain --trust-node
```

Then you can construct and run the following queries:

> NOTE: Be sure to substitute your password and buyer/owner addresses for the ones listed below!

```bash
# Get the sequence and account numbers for jack to construct the below requests
$ curl -s http://localhost:1317/auth/accounts/$(tccli keys show jack -a)
# > {
  "type": "auth/Account",
  "value": {
    "address": "cosmos1eg37masx8qk20lfz0sfcu2tsfea7pcxkcst40g",
    "coins": [
      {
        "denom": "cttc",
        "amount": "800"
      }
    ],
    "public_key": {
      "type": "tendermint/PubKeySecp256k1",
      "value": "A/7aLwoApKQM19OqenU8yVg16cgYe5rfzUIWFpD0fP1+"
    },
    "account_number": "0",
    "sequence": "4"
  }


# Deposit coin from Cosmos to TTC
# NOTE: Be sure to specialize this request for your specific environment
$ curl -XPOST -s http://localhost:1317/tcchan/deposit --data-binary '{"base_req":{"from":"cosmos1eg37masx8qk20lfz0sfcu2tsfea7pcxkcst40g","password":"12345678","chain_id":"tctestchain","sequence":"7","account_number":"0","amount":"66cttc"},"target":"t0c233eC8cB98133Bf202DcBAF07112C6Abb058B89","amount":"50cttc","name":"jack","pass":"12345678"}'
# > {
  "height": "0",
  "txhash": "62CABFDE45F4BE3B4B078282D41E0B6D9BAE221E935825512352DE3482FEE190"
}

# Query the current status
$ curl -s http://localhost:1317/tcchan/current
# {
  "maxOrderNumber": "4",
  "currentDeposit": [
    {
      "orderID": "1",
      "step": "0"
    },
    {
      "orderID": "2",
      "step": "0"
    },
    {
      "orderID": "3",
      "step": "0"
    },
    {
      "orderID": "4",
      "step": "0"
    }
  ],
  "currentWithdraw": null
}

# Query order details
$ curl -s http://localhost:1317/tcchan/order/1
# {
  "orderId": "1",
  "blockNumber": "0",
  "accAddress": "cosmos1eg37masx8qk20lfz0sfcu2tsfea7pcxkcst40g",
  "ttcAddress": "t0c233eC8cB98133Bf202DcBAF07112C6Abb058B89",
  "value": {
    "denom": "",
    "amount": "0"
  },
  "deposit": false,
  "status": "0"
}
```

### Request Schemas:

#### `POST /tcchan/names` BuyName Request Body:
```json
{
  "base_req": {
    "from": "string",
    "chain_id": "string",
    "sequence": "number",
    "account_number": "number",
    "gas": "string,not_req",
    "gas_adjustment": "string,not_req",
  },
  "target": "string",
  "amount": "string"
  "name": "string",
  "password": "string"
}
```

#### `PUT /tcchan/names` SetName Request Body:
```json
{
  "base_req": {
    "from": "string",
    "chain_id": "string",
    "sequence": "number",
    "account_number": "number",
    "gas": "string,not_req",
    "gas_adjustment": "strin,not_reqg"
  },
  "target": "string",
  "amount": "string",
  "name": "string",
  "password": "string"
}
```

### [Back to start of tutorial](./README.md)
