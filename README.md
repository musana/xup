# xup
a tool that takes masscan output through a pipeline, extracts IP and port information, and passes it to chain other tools.

# install
```
go install -v github.com/musana/xup@latest
```
# usage

**to extract ip:port pair** 

```
masscan 1.1.1.1/28 -p 80|xup

output:
1.1.1.10:80
1.1.1.9:80
1.1.1.12:80
1.1.1.4:80
1.1.1.7:80

```


**to extract only IP** 

```
masscan 1.1.1.1/28 -p 80|xup -onlyip

output:
1.1.1.10
1.1.1.9
1.1.1.12
1.1.1.4
1.1.1.7

```
