# Golang GMP模型

###### tags: `Golang`

## 預備知識
``` markmap
# Buzzword
## 單/多核心 與 單/多工
## kernel space 和 user space
## kernel mode 和 user mode
## process 和 thread
## kernel thread 和 user thread
## threading model
```

### 單/多核心與單/多工
#### 單核心單工
Task A先執行，而Task B必須等待Task A執行完畢，Task C必須等待Task B執行完畢。
``` mermaid
gantt
    title Task A/B/C為順序執行
    dateFormat hh:mm:ss.SSS
    axisFormat %M:%S
    section CPU
    Task A:t1, 12:00:00.000, 2s
    Task B:t2, after t1, 1s
    Task C:after t2, 3s
```
#### 單核心多工
Task A先執行一部分，之後換B，再來換C，就下來又換A執行剩下的部分。
``` mermaid
gantt
    title Task A/B/C為併發執行
    dateFormat hh:mm:ss.SSS
    axisFormat %M:%S
    section CPU
    Task A-1:t1, 12:00:00.000, 1s
    Task B-1:t2, after t1, 1s
    Task C-1:t3, after t2, 1s
    Task A-2:t4, after t3, 1s
    Task B-2:t5, after t4, 1s
    Task C-2:after t5, 1s
```
注意每個Task可以執行的時間（區塊）都是一個一個時間切面。
#### 多核心多工
``` mermaid
gantt
    title Task A/B/C/D為併發執行，任務分散到CPU 1/2上
    dateFormat hh:mm:ss.SSS
    axisFormat %M:%S
    section CPU 1
    Task A-1:t1, 12:00:00.000, 1s
    Task B-1:t2, after t1, 1s
    Task A-2:after t2, 1s
    section CPU 2
    Task C-1:t3, 12:00:00.000, 1s
    Task D-1:t4, after t3, 1s
```
#### 任務切換
在上面多工的圖當中其實沒有畫出細節的部分，也就是併發執行當中，從一個task切換到另一個task的過程，也就是context switch。
``` mermaid
gantt
    title Context Switch
    dateFormat hh:mm:ss.SSS
    axisFormat %M:%S
    section CPU
    Task A-1:t1, 12:00:00.000, 1s
    context swtich: done, c1, 12:00:01.000, 12:00:01.200
    Task B-1:after c1, 1s
    context swtich: done, c2, 12:00:02.200, 12:00:02.400
    Task C-1:after c2, 1s
    context swtich: done, c3, 12:00:03.400, 12:00:03.600
    ...: after c3, 1s
```

### kernel mode 和 user mode
* kernel mode：是最高權限的模式，可以存取硬碟，可以控制I/O，可以分配資源。
* user mode： 是最沒有權限的模式，只能做做基本的運算。process運行時，是藉由syscall來從user mode切換到kernel mode。

#### kernel的類型
* monolithic：大部分的系統功能(記憶體管理，排程管理，檔案系統，網路痛訊...等等)都是寫在kernel裡面，kernel較肥，因為大部分功能都在kernel裡，功能模組和功能模組間的呼叫只是func call而已，較有效率。
    * 如Linux
* microkernel：相對於monolithic，有些功能從kernel mode移出，成為了user mode的process，因此在呼叫對應功能時需要透過IPC交換訊息，IPC會帶來context switch和mode change，呼叫成本較高，但相對安全，因process掛掉不至於影響其他process的運作。
    * 如L4

* 但不管kernel類型是什麼，都還是會區分成kernel mode或是user mode。

### kernel space 和 user space
* memory分成kernel space和user space
* kernel space：kernel code儲存和kernel運作的。
* user space：一般常見user process運作的地方，而這個地方是由kernel來管理，避免user process亂存取到他人process。

**對memory存取的權限**
|可以存取 | kernel space | user space |
| -------- | -------- | -------- |
|   kernel mode   | 可以     | 可以     |
|   user mode   | 不行     | 可以     |


### process 和 thread

#### Process
* 為OS分配資源的最小單位
* 一個process至少有一個thread
* 一個process可以有很多thread
* process間的address space（記憶體中所能夠使用與控制的位址區段）是獨立的
* 持有自己的address space，包含code,data,stack和heap空間
* 持有自己的OS資源 

#### Thread
* 又名light weight process
* 為OS排程的最小單位，無法獨立存在
* 同一個process中的thread，共享code,data區塊，但有自己的stack, register和PC
* thread可以建立其他thread

#### Task
* 在Linux系統當中，process或是thread只有概念上差別，但具體表現process和thread都是同一個struct，叫做task_struct，所以可稱呼process和thread為task。
* task_struct裡面涵蓋描述process/thread的訊息，如pid, 狀態，記憶體配置, files,signal...等資訊
* 具體上區分process的方法是看task_struct裡面的指向mm_struct類型的指標，這個指標指到的記憶體位置是否是共用的。

>  mm_struct 裡面就是存放著process的記憶體資訊，如code section, data section...的記憶體位置

### kernel-level thread 和 user-level thread

#### Kernel-level thread
* 上述提到light weight process，就是kernel-level thread，會被scheduler來排程放置在CPU上運行的thread。也就是說，如果要建立一個kernel level thread，是需要呼叫system call並且由kernel來建立。
* 可被搶佔

#### User-level thread
* 由語言自行實做，運行在user mode上的thread，依賴於kernel-level thread，OS無法感知這種thread，其排程也是由user mode上的排程器來進行管理，因此在做這層級的context-switch時，成本是比較小。
    * 如pthread, java thread
* 不可被搶佔



### threading model
指的是在user-level thread如何和kernel-level thread去做對應。分為以下四種。

#### One-to-one
每建立一個user-level thread也會跟著生成kernel-level thread並且將之綁定。
``` graphviz
digraph one2one {
        subgraph cluster_level1{
            label ="Kernel space";
            kernel_thread1 [label="Thread1" shape=circle];
            kernel_thread2 [label="Thread2" shape=circle];
            kernel_thread3 [label="Thread3" shape=circle];
        }
        subgraph cluster_level2{
            label ="User space";
            thread1 [label="Thread 1" shape=circle];
            thread2 [label="Thread 2" shape=circle];
            thread3 [label="Thread 3" shape=circle];
        }
        thread1 -> kernel_thread1
        thread2 -> kernel_thread2
        thread3 -> kernel_thread3
}
```
#### Many-to-one
多對一的方式進行綁定。
``` graphviz
digraph many2one {
        subgraph cluster_level1{
            label ="Kernel space";
            kernel_thread [label="Thread" shape=circle];
        }
        subgraph cluster_level2{
            label ="User space";
            thread1 [label="Thread 1" shape=circle];
            thread2 [label="Thread 2" shape=circle];
            thread3 [label="Thread 3" shape=circle];
        }
        thread1 -> kernel_thread
        thread2 -> kernel_thread
        thread3 -> kernel_thread
}
```
#### Many-to-many
user-level thread可以任意分配到kernel-level thread上執行。
``` graphviz
digraph many2many {
        subgraph cluster_level1{
            label ="Kernel space";
            kernel_thread1 [label="Thread 1" shape=circle];
            kernel_thread2 [label="Thread 2" shape=circle];
        }
        subgraph cluster_level2{
            label ="User space";
            thread1 [label="Thread 1" shape=circle];
            thread2 [label="Thread 2" shape=circle];
            thread3 [label="Thread 3" shape=circle];
        }
        thread1 -> kernel_thread1
        thread2 -> kernel_thread1
        thread3 -> kernel_thread1
        thread1 -> kernel_thread2
        thread2 -> kernel_thread2
        thread3 -> kernel_thread2   
}
```
#### Two-level
與many-to-many相似，但可以允許指定做one-to-one
``` graphviz
digraph many2many {
        subgraph cluster_level1{
            label ="Kernel space";
            kernel_thread1 [label="Thread 1" shape=circle];
            kernel_thread2 [label="Thread 2" shape=circle];
            kernel_thread3 [label="Thread 3" shape=circle];
        }
        subgraph cluster_level2{
            label ="User space";
            thread1 [label="Thread 1" shape=circle];
            thread2 [label="Thread 2" shape=circle];
            thread3 [label="Thread 3" shape=circle];
            thread4 [label="Thread 4" shape=circle];

        }
        thread1 -> kernel_thread1
        thread2 -> kernel_thread1
        thread3 -> kernel_thread1
        thread1 -> kernel_thread2
        thread2 -> kernel_thread2
        thread3 -> kernel_thread2
        thread4 -> kernel_thread3
}
```

## 正題
### Goroutine
Golang很好的封裝了user-level thread，並且提供了`go`關鍵字來建立一個thread，稱呼為Goroutine，同時在`runtime` package裡面定義了Go scheduler，用來排程和管理Goroutine們。

### 什麼是GMP

go scheduler所採用的調度模型稱呼為GMP（以前是GM），由G, M和P這三種元件構成。

#### G
就是Goroutine，可以視為user-level thread。如同 kernel-level thread是由OS做context switch在CPU core上運作，Goroutine是由Go scheduler做context switch在M上面運作。

Goroutine有三個狀態
* Waiting：代表Goroutine正在等待systemcall執行完畢，或是因為正在等待鎖。
* Runnable：代表Goroutine想要在M上執行指令。
* Executing：代表Goroutine正在M上執行指令當中。

#### M
Machine的縮寫，代表的是kernel-level thread，由OS管理，是執行Goroutine的實體，預設最多可有10,000個M。
#### P
Processor，是一個邏輯上意義的處理器，不代表真實的CPU core數量，這個數量在程序啟動時就被決定，這個也被代表著在process運行期間，最多同時只有#P個goroutine在運作，我們可以藉由P來限制process的併發程度。
#### LRQ
Local Run Queue, 用來放置G，每個P都會擁有自己的LRQ。
#### GRQ
Global Run Queue，也是用來放置G，當有些LRQ滿了之後，無法塞進更多G時，就會把G到GRQ裡面。
``` graphviz
digraph GMP {
	node[shape=record]
    rankdir=TB;
    
    GRQ [label="{ GRQ |{<Head>G|G|G|G|<Tail>}}"];
    
    new [label="go\ func()" shape=none];
    newG [label="G'" shape=circle];
    new -> newG [label="new"];
    newG -> P1:Tail [label="push"]; 
    
    P1 [label="<P> P |  { LRQ |{||||<Tail>}}"];
    M1 [label="M" shape=triangle];
    P1:P -> M1 [label="attach"];
    RunG [label="G" shape=circle];
    RunG -> M1 [label="run on"];
    RunG -> new;
    
    GRQ ->P1:Tail [label="pull & push"];
    
    P2 [label="<P>P | { LRQ |{<Head>G|G|G|G|G}}"];
    M2 [label="M" shape=triangle];
    P2:P -> M2;
    P2:Head->M2 [label="pull"];
    
    Scheduler [label="OS Scheduler"];
    M1 -> Scheduler;
    M2 -> Scheduler [label="managed by"];
    
    Core1 [label="Core" shape=hexagon];
    Core2 [label="Core" shape=hexagon];
    Core3 [label="Core" shape=hexagon];
    Core4 [label="Core" shape=hexagon];
    Scheduler -> Core1;
    Scheduler -> Core2;
    Scheduler -> Core3;
    Scheduler -> Core4;
}
```


### 運作機制
* M必須先取得一個P後，才能從P的LRQ中取得G來執行
* 若是LRQ為空，則會從GRQ或是其他P的LRQ拿G來放到本地P的LRQ裡面
* M拿到G，執行G，執行到某個時間點後，會進行context switch，把G放回LRQ，並從P拿下一個G執行，一直重複上述步驟。


#### Handover
當G1正在M1上執行時，遇到了需要呼叫blocking system call時，為提高系統效能，Go scheduler會執行一種**Handover**的機制，如以下。 
``` graphviz
digraph Handover1 {
	node[shape=record]
    rankdir=TB;
 
    P [label="<P> P |  { LRQ |{G2|G3|G4|G5}}"];
    M1 [label="M1" shape=triangle];
    P:P -> M1 ;
    G1 [label="G1" shape=circle];
    G1 -> M1;
}
```
這時G1想要呼叫blocking system call，這時會將M1和P解除關聯，並取得一個新的M(M2)來和原P進行關聯，並繼續執行LRQ內的G2。而M1則是持續負責G1的執行。
``` graphviz
digraph Handover2 {
	node[shape=record]
    rankdir=TB;
 
    P [label="<P> P |  { LRQ |{G3|G4|G5|}}"];
    M1 [label="M1" shape=triangle];
    M2 [label="M2" shape=triangle];
    P:P -> M2 ;
    G2 -> M2;
    G1 [label="G1" shape=circle];
    G1 -> M1;
}
```
當G1執行system call完畢後，會被塞回LRQ，而M1則變成閒置狀態，並等待之後使用。
``` graphviz
digraph Handover2 {
	node[shape=record]
    rankdir=TB;
 
    P [label="<P> P |  { LRQ |{G3|G4|G5|G1}}"];
    M1 [label="M1" shape=triangle];
    M2 [label="M2" shape=triangle];
    P:P -> M2 ;
    G2 -> M2;
}
```

#### Work stealing
如果有個P的LRQ已經空了，且M也沒有正在執行的G，那這個M因為進入到閒置狀態很有可能被OS scheduler context switch掉，即使這個process中還有其他待執行的G。為了避免上述的情況發生，這時Go scheduler會執行work stealing，從其他的P LRQ中或是GRQ中偷G。

一開始執行狀態如下圖。

``` graphviz
digraph workstealing1 {
    node[shape=record]
    rankdir=TB;
 
    P1 [label="<P> P1 |  { LRQ |{G3|G5|G7}}"];
    M1 [label="M1" shape=triangle];
    P1:P -> M1 ;
    G1 [label="G1" shape=circle];
    G1 -> M1;
    
    P2 [label="<P> P2 |  { LRQ |{G4|G6|G8}}"];
    M2 [label="M2" shape=triangle];
    P2:P -> M2 ;
    G2 [label="G2" shape=circle];
    G2 -> M2;
    
    GRQ [label="{ GRQ |{G9||}}"];
}
```
而當M1及其對應的P1 LRQ為空時。
``` graphviz
digraph workstealing1 {
    node[shape=record]
    rankdir=TB;
 
    P1 [label="<P> P1 |  { LRQ |{||}}"];
    M1 [label="M1" shape=triangle];
    P1:P -> M1 ;
    
    P2 [label="<P> P2 |  { LRQ |{G4|G6|G8}}"];
    M2 [label="M2" shape=triangle];
    P2:P -> M2 ;
    G2 [label="G2" shape=circle];
    G2 -> M2;
    
    GRQ [label="{ GRQ |{G9||}}"];
}
```
P1嘗試從P2的LRQ偷一半的G，結果如下。
``` graphviz
digraph workstealing1 {
    node[shape=record]
    rankdir=TB;
 
    P1 [label="<P> P1 |  { LRQ |{G6||}}"];
    M1 [label="M1" shape=triangle];
    P1:P -> M1 ;
    G4 [label="G4" shape=circle];
    G4 -> M1;
    
    P2 [label="<P> P2 |  { LRQ |{G8||}}"];
    M2 [label="M2" shape=triangle];
    P2:P -> M2 ;
    G2 [label="G2" shape=circle];
    G2 -> M2;
    
    GRQ [label="{ GRQ |{G9||}}"];
}
```
而當M2把P2的G都執行完畢時，且這時候P1的LRQ也空了。
``` graphviz
digraph workstealing1 {
    node[shape=record]
    rankdir=TB;
 
    P1 [label="<P> P1 |  { LRQ |{||}}"];
    M1 [label="M1" shape=triangle];
    P1:P -> M1 ;
    G6 [label="G6" shape=circle];
    G6 -> M1;
    
    P2 [label="<P> P2 |  { LRQ |{||}}"];
    M2 [label="M2" shape=triangle];
    P2:P -> M2 ;
    
    GRQ [label="{ GRQ |{G9||}}"];
}
```
P2這時候會嘗試從GRQ偷，如下。
``` graphviz
digraph workstealing1 {
    node[shape=record]
    rankdir=TB;
 
    P1 [label="<P> P1 |  { LRQ |{||}}"];
    M1 [label="M1" shape=triangle];
    P1:P -> M1 ;
    G6 [label="G6" shape=circle];
    G6 -> M1;
    
    P2 [label="<P> P2 |  { LRQ |{||}}"];
    M2 [label="M2" shape=triangle];
    P2:P -> M2 ;
    G9 [label="G9" shape=circle];
    G9 -> M1;
    
    GRQ [label="{ GRQ |{||}}"];
}
```
<details> 
<summary></summary>
custom_mark10
  digraph G {
    size ="4,4"
    main [shape=box]
    main -> parse [weight=8]
    parse -> execute
    main -> init [style=dotted]
    main -> cleanup
    execute -> { make_string; printf}
    init -> make_string
    edge [color=red]
    main -> printf [style=bold,label="100 times"]
    make_string [label="make a string"]
    node [shape=box,style=filled,color=".7 .3 1.0"]
    execute -> compare
  }
custom_mark10
</details>


## Reference
1. http://www.it.uu.se/education/course/homepage/os/vt18/module-4/implementing-threads/
2. https://www.ardanlabs.com/blog/2018/08/scheduling-in-go-part1.html
