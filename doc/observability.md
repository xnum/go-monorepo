可觀測性
===

# logging

在logging包中實作了送往fluentd的集中化log管理

# healthy probe

配合k8s的probe在health包中實作了相關程式，提供`/alive` `/ready` `/vars`功能

針對每一個長時間運行的task或goroutine，如果我們想觀測其運行狀態或健康度，

就會讓他持有一個health.Info物件，其中 `NewInfo(name, d, pt)`

name代表這個task的唯一可識別名稱，用來為這個task命名及作為unique key使用。

d則是這個task理應在每多少時間內進行一次回報，0代表該功能關閉。

pt則是這個Info聯繫的ProbeType。

none = 不聯繫到任何Probe上，也不處理task退出的異常狀況，僅做vars回報用。

alive = 這是個關鍵task，如果出了問題，k8s會將整個pod重啟

ready = 這是個服務task，如果出了問題，k8s不會將service流量導入，但也不會重啟。

一般流程如下：

```
info := health.NewInfo("slipping", N+5 seconds, Ready)
defer info.Down()

SetupTask()

for {
  sleep N seconds
  if err := DoTask(); err != nil {
    info.Pause()
  } else {
    info.Up()
  }

  vars := CollectTaskStatus()
  info.UpdateVars(vars)
}
```

當task啟動後進入init狀態，並在init成功後進入running狀態，

如果task超過時間d沒有回報資訊，則被視為該task已經失聯，

這時對應的ProbeType會被設為false。在所有Info中只要有任何一個被認定false，

整個alive或ready Probe會回傳false。

如果Info被初始化後遲遲沒有呼叫任何一次Up()或呼叫了Down()，

代表這個task從未啟動成功或因為不可預期的錯誤已經永久退出，

只要有任何一個Info處於這個狀態，這時候alive probe會回傳false。

UpdateVars()則是該task主動蒐集資訊並回報到server，

避免server詢問狀態時產生大量的lock動作。

每個process至少應該註冊一個Info。

# metric

使用promauto的constructor在建立時順便註冊到DefaultRegistry中，

health server啟動時就會提供`/metrics`給prometheus使用。
