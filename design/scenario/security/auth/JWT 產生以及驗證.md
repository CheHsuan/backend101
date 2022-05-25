# JWT 產生以及驗證

###### tags: `JWT`, `encryption`

## 如何產生token

**JWT Token**

```
token = header.payload.signature
```

**Header**

```
header = Base64UrlEncode({...})
```

**Payload**

```
payload = Base64UrlEncode({...})
```

**Signature**

```
signature = Base64UrlEncode(Sign(header.payload, your-secret))
```

## 簽章演算法
不同的簽名演算法可以設計出不同的驗證方法。

1. 對稱式加密
2. 非對稱式加密

### 對稱式加密

常見演算法：HMAC

使用同一把密鑰簽署以及驗證jwt token，因此auth驗證流程如下。

``` mermaid
sequenceDiagram
User->>Auth: ask for JWT token
Auth->>MySQL: verify user and password
Note over Auth: generate JWT token with secret key
Auth-->>User: JWT token
User->>Service: query with JWT token
Service->>Auth: verify token
Auth-->>Service: verified
Service-->>User: OK
```

可以看到這個流程當中，都是由Auth服務進行token分發和驗證。

**優點**
專職專責，所有Auth事務歸Auth服務管理，架構明確。

**缺點**
每次驗證都須藉由Auth服務來進行，壓力較大，另外多一層API訪問較耗時。


### 非對稱式加密

常見演算法：RSA

使用私鑰進行jwt token簽署，使用公鑰進行驗證。私鑰必須自行保留，公鑰可以公開出去，由其他方進行驗證。因此除了上述的流程外，我們也可以如下設計。

``` mermaid
sequenceDiagram
User->>Auth: ask for JWT token
Auth->>MySQL: verify user and password
Note over Auth: generate JWT token with private key
Auth-->>User: JWT token
User->>Service: query with JWT token
Service->>Auth: query JWKS API
Auth-->>Service: public key
Note over Service: verify JWT token with public key
Note over Service: keep th public key in local cache
Service-->>User: OK
User->>Service: query agagin with JWT token
Note over Service: verify JWT token with public key from cache
Service-->>User: OK
```

在這個流程當中，由Auth充當token生成方，Service為驗證方。

**優點**
分散壓力，減少API呼叫。

**缺點**
非對稱式加密本身比對稱式加密來得緩慢。

> 其實不管是對稱式或非對稱式的方法，都是可以藉由加入緩存的方式進行快速驗證的，不需重新執行解密流程，從而加速整體速度。
