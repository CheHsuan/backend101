# Logging架構

## **流程**

``` mermaid
flowchart LR
   資料收集A & 資料收集B & 資料收集C--> 資料轉發 --> 資料儲存[(資料儲存)] --> 資料展示
```



一般使用的元件如下

**資料收集**

Filebeat, Logstash, Fluentd

**資料轉發**

Logstash, Fluentd

**資料儲存**

Elasticsearch

**資料展示**

Kibana

常見的架構為ELK, EFK等，但這樣的架構隨著接入的業務等級擴增而難以應變，而尖峰時刻更是會造成資料轉發和資料儲存層的壓力大增，造成高延遲，影響整體服務運作。

## **優化**

``` mermaid
flowchart LR
   資料收集A & 資料收集B & 資料收集C--> 資料緩存 --> 資料轉發 --> 資料儲存[(資料儲存)] --> 資料展示
```

在架構中加入消息佇列，隔離了資料收集和資料轉發，這種方式降低了尖峰時期對後端處理系統的壓力，讓資料轉發端可以依照穩定的速度處理logging訊息。

那實際的架構及其元件對應如下。

``` mermaid
flowchart LR
   FilebeatA & FilebeatB & FilebeatC --> Kafka --> Logstash/Fluentd--> Elasticsearch[(Elasticsearch)] --> Kibana
```

