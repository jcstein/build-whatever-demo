# Start your light node

## Dependencies

- celestia-node v0.13.5
- golang go1.22.3

## Clone this repo

```bash
git clone git@github.com:jcstein/build-whatever-demo.git
cd build-whatever-demo
```

## Get trusted height \& hash

> Note: this is only required if you have not fully synced your node or if you don't want to wait for it to sync.

[On Celenium](https://celenium.io)

## Set the trusted height \& hash

In a separate terminal, open your `config.toml`\:

```bash
code ~/.celestia-light/config.toml
```

Set `DASer.SampleFrom` to the trusted height\.

## Initialize the node store

```bash
celestia light init
```

## Run the node with core access \& skip auth

```bash
celestia light start --core.ip consensus.lunaroasis.net --rpc.skip-auth --headers.trusted-hash <hash>
```

## Check the sampling stats

In a new terminal, check the sampling stats\:

```bash
celestia das sampling-stats --node.store ~/.celestia-light | jq '.result | {head_of_sampled_chain, head_of_catchup, network_head_height}'
```

When they look like this\, you\'re synced\:

```bash
{
  "head_of_sampled_chain": 1512232,
  "head_of_catchup": 1512232,
  "network_head_height": 1512232
}
```

## Run main\.go to watch the 0xdeadbeef namespace

```bash
cd ~/build-whatever-demo/rollkit-monitor && go run main.go
```

See the results \(you shouldn\'t see any blobs yet\)\:

```bash
📡 Subscribed to headers. Waiting for new headers...
📝 New header received: Height 1513753, Hash 39443444324644313031343545314545384444433134323637384433443043393643413730304334423144333244324430313245423846354536354138374444
⚠️ Error fetching blobs: getting blobs for namespace(00000000000000000000000000000000000000000000000000deadbeef): blob: not found
blob: not found
⊞ EDS fetched at height 1513753: &{0x14000695300 0x140000b0020 16}
...
```

## Run a Rollkit rollup

### Set up gm rollup

```bash
cd $HOME && bash -c "$(curl -sSL https://rollkit.dev/install-gm-rollup.sh)"
```

Build the rollup\:

```bash
cd ~/gm && bash init.sh
```

### Set the block height

```bash
DA_BLOCK_HEIGHT=$(curl https://rpc.celestia.pops.one/block | jq -r '.result.block.header.height')
echo -e "\n Your DA_BLOCK_HEIGHT is $DA_BLOCK_HEIGHT \n"
```

### Set the auth token

```bash
AUTH_TOKEN=$(celestia light auth admin)
echo -e "\n Your DA AUTH_TOKEN is $AUTH_TOKEN \n"
```

### Set the namespace

```bash
DA_NAMESPACE=00000000000000000000000000000000000000000000000000deadbeef
```

### Run the rollup

```bash
gmd start \
    --rollkit.aggregator \
    --rollkit.da_auth_token $AUTH_TOKEN \
    --rollkit.da_namespace $DA_NAMESPACE \
    --rollkit.da_start_height $DA_BLOCK_HEIGHT \
    --minimum-gas-prices="0.025stake"
```

## Watch the namespace of the rollup in main\.go

```bash
...
📝 New header received: Height 1513754, Hash 46444536373932393944354136344232354536323639333842304144343431423632343930393539433143393134383834324338393530353635373341354345
🟣 Found 15 blobs at height 1513754 in 0xdeadbeef namespace
⊞ EDS fetched at height 1513754: &{0x14000800380 0x140000b0020 32}
📝 New header received: Height 1513755, Hash 38434343323046393042463543434246424641423530363544323532384432333046443138453133323734393038394639323833454446323239444339363243
🟣 Found 15 blobs at height 1513755 in 0xdeadbeef namespace
⊞ EDS fetched at height 1513755: &{0x14000fa2000 0x140000b0020 16}
```

## Send a transaction on the rollup

First\, list the keys\:

```bash
gmd keys list --keyring-backend test
```

Then set 2 of them to variables\:

```bash
export KEY1=
export KEY2=
```

Send between the 2 keys\:

```bash
gmd tx bank send $KEY1 $KEY2 42069stake --keyring-backend test --chain-id gm --fees 5000stake
```

Check the balance\:

```bash
gmd query bank balances $KEY2
```

And of KEY1\:

```bash
gmd query bank balances $KEY1
```

## Run the go code with transaction submission

```bash
cd ~/build-whatever-demo/celestia-monitor && go run main.go
```

```bash
📡 Subscribed to headers. Waiting for new headers...
🟢 Blob was included at height 1517550
🧐 Blobs are equal? true
✅ New blob submitted and verified successfully
🧊 New header received: Height 1517550, Hash 35413946443831393934433144424242463736343142443941313441373038314346313236323041423130373243353746343130453537373338433634423937
🟣 Found 1 blobs at height 1517550 in 0xdeadbeef namespace
🟩 EDS fetched at height 1517550: &{0x140001dae00 0x1400019c000 16}
```

[View the namespace on Celenium](https://celenium.io/namespace/000000000000000000000000000000000000000000000000deadbeef?tab=Blobs)\.
