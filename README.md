# RA-WEBs
[![Go](https://github.com/akakou/ra-webs/actions/workflows/go.yml/badge.svg)](https://github.com/akakou/ra-webs/actions/workflows/go.yml)

RA-WEBs is a protocol that enables browsers to verify proof of Remote Attestation while maintaining compatibility.

### Dependencies

- An Azure instance with Intel SGX (for running the example TA)
- Ubuntu 20.04

### How to Deploy the Test Environment

#### 1. Clone the Repository

```bash
git clone https://github.com/akakou/ra-webs
cd ra-webs/test
```

#### 2. Run the Cloudflare Tunnel (Optional)


```sh
docker compose -f compose.test.yaml --profile tunnel
```

#### 3. Configure the Environment Files

Copy the templates and fill in each parameter.

```sh
cp env/common.env.template env/common.env
cp env/ta.env.template env/ta.env
cp env/monitor.env.template env/monitor.env
```


#### 4. Run the servers

```sh
docker compose -f compose.test.yaml --profile ta up
```


### NOTE 

The deployability was verified using DC1ds v3 (1 vCPU, 8 GiB memory).
