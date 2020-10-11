proxy forward for wsl
----

## method use powershell

- get wsl ip by bash
- portproxy by netsh (admin required)

```pwsh
$prt = 1234
$wip = (bash -c "ip addr show eth0| grep 'inet '| awk '{print `$2}'| awk -F'/' '{print `$1}'")
netsh interface `
  portproxy add v4tov4 `
  listenport=$prt listenaddress=0.0.0.0 `
  connectport=$prt connectaddress=$wip
```
