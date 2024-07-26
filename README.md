# xup
a tool that takes masscan output through a pipeline, extracts IP and port information, and passes it to other tools.

# install
```
go install -v github.com/musana/xup@latest
```
# usage

**to extract ip:port pair** 

`masscan IP:PORT -p80|xup`

**to extract only IP** 

`masscan IP:PORT -p80|xup -onlyip`
