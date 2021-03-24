# RateLimiter

Dcard Backend Intern 作業。

### 題目: 

- 設計一個 middleware 限制每小時來自同一 IP 的請求數量不超過 1000
- 要在 response header 中加入剩餘的請求數量 (X-RateLimit-Remaining) 和 rate limit 歸零的時間(X-RateLimit-Reset)
- 超過限制要回傳 429 (Too Many Requests)

### 想法

使用 Go 的 Gin 框架配合 Redis 資料庫，實做下面的演算法:

- 檢查 IP 是否已被紀錄
  - 如果 IP 還沒紀錄過，將 IP 紀錄下來，把值設定成 1，並設定 TTL。 (SETNX)
  - 如果 IP 已經紀錄過，把它的值加 1。(INCR)
- 取得該 IP 的值、還剩多久時間到期。(TTL)

另外，對 Redis 的操作是利用 Lua script，這是為了避免 Redis 在進行 Incr 的前一刻，key 突然過期，而導致 key 沒有被設置 Expiration time。

考量到在不同的 Route 底下可能會需要不同的 Ratelimiter，所以提供了一個 `Keyprefix` 的選項，以區別相同 IP 在不同 Route 的 key。

### 範例

參考 [example](example)。

