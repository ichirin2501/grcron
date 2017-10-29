## grcron導入方法

### grcronとは

- keepalivedのVRRP機能を利用したcron冗長化ツール
- keepalivedのmaster/backupの状態をもとにプログラム実行を判断

### 必要なこと

- keepalivedの準備
- grcronの設置
- cronの設定変更

#### keepalivedの準備
設定例  
master側の設定例  
```
vrrp_instance r0 {
    state MASTER
    priority 150
    interface eth0
    garp_master_delay 5
    virtual_router_id 45  # 必ずuniq、設定時は注意
    advert_int 1
    authentication {
        auth_type PASS
        auth_pass TEST
    }
    notify_backup "/bin/echo passive > /var/lib/grcron/state"
    notify_master "/bin/echo active > /var/lib/grcron/state"
    notify_fault "/bin/echo passive > /var/lib/grcron/state"
}
```
backup側の設定例  
```
vrrp_instance r0 {
    state BACKUP
    priority 100
    interface eth0
    garp_master_delay 5
    virtual_router_id 45
    advert_int 1
    authentication {
        auth_type PASS
        auth_pass TEST
    }
    notify_backup "/bin/echo passive > /var/lib/grcron/state"
    notify_master "/bin/echo active > /var/lib/grcron/state"
    notify_fault "/bin/echo passive > /var/lib/grcron/state"
}
```

状態の出力先のディレクトリは事前作成が必要。  
keepalivedのunicast設定でも動くかもしれない（未検証）

#### grcronの設置
download: https://github.com/ichirin2501/grcron/releases  
```
$ wget https://github.com/ichirin2501/grcron/releases/download/v0.0.5/grcron_linux_amd64.zip
$ unzip grcron_linux_amd64.zip
$ sudo mv grcron_linux_amd64/grcron /usr/bin/
```

#### cronの設定変更
実行プログラムの先頭にgrcronを挟むだけで良い。  
before  
```
*/5 * * * * root echo "nihaha" > /tmp/test.txt 2>&1 | logger -t grcrontest
```
after  
```
*/5 * * * * root /usr/bin/grcron echo "nihaha" > /tmp/test.txt 2>&1 | logger -t grcrontest
```

dry-run機能もあるので事前に動作確認すると良い。  
```
$ /usr/bin/grcron -h
Usage of /usr/bin/grcron:
  -dryrun
        dry-run.
  -f string
        grcron state file. (default "/var/lib/grcron/state")
  -n    dry-run.
  -s string
        grcron default state. (default "passive")
  -v    show version number.
  -version
        show version number.
```

