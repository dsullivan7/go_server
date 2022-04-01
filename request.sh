#!/bin/bash

curl -X POST 'https://www.connectebt.com/nyebtclient/siteLogonClient.recip' \
  -H 'authority: www.connectebt.com' \
  -H 'pragma: no-cache' \
  -H 'cache-control: no-cache' \
  -H 'sec-ch-ua: " Not A;Brand";v="99", "Chromium";v="99", "Google Chrome";v="99"' \
  -H 'sec-ch-ua-mobile: ?0' \
  -H 'sec-ch-ua-platform: "macOS"' \
  -H 'origin: https://www.connectebt.com' \
  -H 'upgrade-insecure-requests: 1' \
  -H 'dnt: 1' \
  -H 'content-type: application/x-www-form-urlencoded' \
  -H 'user-agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.84 Safari/537.36' \
  -H 'accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9' \
  -H 'sec-fetch-site: same-origin' \
  -H 'sec-fetch-mode: navigate' \
  -H 'sec-fetch-user: ?1' \
  -H 'sec-fetch-dest: document' \
  -H 'referer: https://www.connectebt.com/nyebtclient/siteLogonClient.recip' \
  -H 'accept-language: en-US,en;q=0.9,es;q=0.8' \
  -H 'cookie: JSESSIONID=28816A350481768ABEE3CB53A1CAC5D9; nlbi_1287184=U4gvJWjNAGIF++phUN/a+gAAAACpa92pNkscXnG7DpGGu3Uo; incap_ses_7225_1287184=AgLnMsKG639PBiaEo1tEZBQFR2IAAAAAslBe2GO9bx3AJXxgbnA2Gg==; reese84=3:TNlB0Ik/xxl0alGFERIMnw==:RCplFTLgEX1ZSNogoMILLHp4VJV0p4crqgsXaUnoJS+EoIxKgCh0DkZJuB5tIC/Ze7DmnbgRQcf5TLFQ4nkv92Z4CEmIy3uDDvU7lNkm1aWTDKmpOYXsWOtdvufMFlGxj5GqKMoSVMVcGLeR2EaejdrhQXkIDZfRDmZDYU42DB2p5fEMbVJ4KG1XjS0wpBMchToyyCdw3Sc6JA0xqonE9rRd9WS7qZ94sVpwsYh4MoohOjXhMOlE8rvOHZjQyhIam4itpmRGyUemUWAl2vlobZ4NFakSbHEkPcR3AXgsUFUmAiiz8QOxNxQJjEmBPHui/xYtSBXX38n2y6UctWg2FVlJvysfy6h34B8Ucwr0gp9co4QanCCtnlfh0Ar+nSp77WOR65KVhpdRl0Q0LV+U+YcHGiU4oSyCJqqGFNKqjcc=:C4dFnwKeEOcQVA+6o1i4dnCfnIbwgXcMZ9Wn0KjsUMs=; incap_sh_1287184=RwhHYgAAAAD5jN5YDAAIx5CckgYQxJCckgbytc2p3Zv5dezuVf3ZpjxI; visid_incap_1287184=BNtM7I6VSim8J3j68i11ZwBlP2IAAAAARUIPAAAAAACANFGjAUen+FNfXZmigt4f2rXwwuq3O0pA; nlbi_1287184_2147483392=pNo3GdoUBz0ChHtoUN/a+gAAAACxRWtD4DzR/ilBX2YdqPIh' \
  --data-raw 'login=SOMEUSER&password=somelogin' \
  --compressed
