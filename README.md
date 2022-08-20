# go언어로 만든 리더보드
[![codecov](https://codecov.io/gh/JeongMinSik/go-leaderboard/branch/master/graph/badge.svg?token=NT5G079D2B)](https://codecov.io/gh/JeongMinSik/go-leaderboard)
[![golangci-lint](https://github.com/JeongMinSik/go-leaderboard/actions/workflows/0_golangci-lint.yaml/badge.svg)](https://github.com/JeongMinSik/go-leaderboard/actions/workflows/0_golangci-lint.yaml)
[![test](https://github.com/JeongMinSik/go-leaderboard/actions/workflows/1_test.yaml/badge.svg)](https://github.com/JeongMinSik/go-leaderboard/actions/workflows/1_test.yaml)
[![build](https://github.com/JeongMinSik/go-leaderboard/actions/workflows/2_build.yaml/badge.svg)](https://github.com/JeongMinSik/go-leaderboard/actions/workflows/2_build.yaml)
[![load-test](https://github.com/JeongMinSik/go-leaderboard/actions/workflows/3_load-test.yaml/badge.svg)](https://github.com/JeongMinSik/go-leaderboard/actions/workflows/3_load-test.yaml)

# 개발 목적
- CRUD 기능을 갖춘 Go 언어 API서버 개발

# 실행
```
docker compose up
```

# OpenApi Document
- https://app.swaggerhub.com/apis-docs/JeongMinSik/leaderboard-api/1.0
- 실행 후 http://localhost:6025/swagger/index.html 에서도 확인

# 기술 스택
### 언어
- __Go__
    - [echo framework](https://github.com/labstack/echo) 사용

### DB
- __Redis__
    - [ZSet](https://redis.io/docs/data-types/sorted-sets/) 사용하여 순위 관리
    - 테스트 코드에서는 [go-redismock](https://github.com/go-redis/redismock) 패키지 사용

### Log
- __Elasticsearch__
    - 날짜별 인덱스 생성: "api-log-YYYY-mm-dd"
    - [elogrus](https://github.com/sohlich/elogrus) 패키지를 사용하여 [logrus](https://github.com/sirupsen/logrus)의 hook에 elasticsearch 연결
- __Kibana__
    - log 시각화 툴로 사용

### Tools
- __Docker Compose__
    - go api server, redis, elasticsearch, kibana 간편하게 구동 가능
- __Swagger__
    - [swaggo](https://github.com/swaggo/swag) 패키지를 사용하여 주석으로 api 정의
- __golangci-lint__
    - [.golangci.yaml](https://github.com/JeongMinSik/go-leaderboard/blob/master/app/.golangci.yaml)로 rules 정의

### CI(Github Actions)
- __lint check__
- __test__
    - [Codecov](https://app.codecov.io/gh/JeongMinSik/go-leaderboard)로 Test Coverage 확인
- __build__
- __load test__
    - docker compose로 서버 실행 후
    - [iter8](https://github.com/iter8-tools/iter8) Github Action을 사용하여 간단한 부하테스트(SLO 설정)
    - example:
    ```
  Metric                     |value
  -------                    |-----
  http/error-count           |0.00
  http/error-rate            |0.00
  http/latency-max (msec)    |0.32
  http/latency-mean (msec)   |0.15
  http/latency-min (msec)    |0.05
  http/latency-p50 (msec)    |0.18
  http/latency-p75 (msec)    |0.25
  http/latency-p90 (msec)    |0.29
  http/latency-p95 (msec)    |0.31
  http/latency-p99 (msec)    |0.32
  http/latency-p99.9 (msec)  |0.32
  http/latency-stddev (msec) |0.08
  http/request-count         |100.00
    ```

