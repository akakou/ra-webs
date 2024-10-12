# RA-WEBs
[![Go](https://github.com/akakou/ra-webs/actions/workflows/go.yml/badge.svg)](https://github.com/akakou/ra-webs/actions/workflows/go.yml)

RA-WEBs is a protocol that enables browsers to verify proof of Remote Attestation while maintaining compatibility.

### Dependencies

- An Azure instance with Intel SGX (for running the example TA)
- Ubuntu 22.04

### How to Deploy the Test Environment

#### 1. Clone the Repository

```bash
git clone https://github.com/akakou/RA-WEBs
cd RA-WEBs
```


#### 2. Configure the Verifier Environment Files

Copy the templates and fill in each parameter.

```sh
cp test/env/verifier.env.template test/env/verifier.env
cp test/env/common.env.template test/env/common.env
```


#### 3. Run the Verifier


```sh
docker compose -f verifier up
```

#### 4. Configure the TA Environment Files

Copy the templates and fill in each parameter.

```sh
cp test/env/ta.env.template test/env/ta.env
```


#### 5. Run the Example TA

```sh
docker compose -f ta up
```


### NOTE 

The functionality was verified using the following:

Google Chrome 129.0.6668.58
DC1s v2 (1 vCPU, 4 GiB memory)