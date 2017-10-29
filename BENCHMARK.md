# simple-server
```
Running 30s test @ http://localhost:8081/index.html
  12 threads and 400 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    20.35ms   31.59ms 384.83ms   97.54%
    Req/Sec   212.61    108.91   633.00     67.08%
  10447 requests in 30.05s, 2.01MB read
  Socket errors: connect 157, read 10655, write 0, timeout 0
Requests/sec:    347.68
Transfer/sec:     68.59KB
```
# in-memory-cache
```
Running 30s test @ http://localhost:8081/index.html
  12 threads and 400 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     1.69ms    6.39ms 668.67ms   97.04%
    Req/Sec   766.78      1.19k    4.01k    84.04%
  15102 requests in 30.11s, 2.91MB read
  Socket errors: connect 157, read 15333, write 0, timeout 0
Requests/sec:    501.64
Transfer/sec:     98.96KB
```
