# bench-tools

This code is used to access the RA-WEBs Monitor.
It allows you to check the latency of the monitor API through the monitor logs.


## how to works
### 1. Set URL
Update the URL in the code:

```
URL = "https://kevin-fares-prospective-fame.trycloudflare.com/"
```

### 2. Run
Use the following command:

```
nix-shell --command "python3 main.py"
```