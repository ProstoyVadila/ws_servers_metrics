# Results

## Initial setup

rust_app: **2Gb 2cpu** \
python_app: **2Gb 2cpu**

## Load Scenarios

### 1. 3 Direct Messages (constVUWith3DirectMessage)

Virtual users: 80

Each user connects, sends, recieves and sends again.

### 2. 100 broadcast Messages (rampUp50With100Messages)

Virtual users: 3-200

Each user produces 100 messages for broadcast

### 3. Single Message (rampUp50WithSingleMessage)

Virtual users: 3-750

Each user produces 1 message for broadcast

## K6 results

```
     checks................: 50.00% ✓ 120979     ✗ 120979
     data_received.........: 32 MB  38 kB/s
     data_sent.............: 32 MB  38 kB/s
     iteration_duration....: avg=1.59s    min=1.4s     med=1.4s     max=1m11s  p(90)=1.57s    p(95)=1.71s
     iterations............: 120979 143.690523/s
     vus...................: 351    min=0        max=749
     vus_max...............: 750    min=750      max=750
     ws_connecting.........: avg=51.04ms  min=302.95µs med=4.8ms    max=19.33s p(90)=140.07ms p(95)=265.43ms
     ws_msgs_sent..........: 151103 179.469735/s
     ws_session_duration...: avg=884.46ms min=700.42ms med=708.06ms max=1m10s  p(90)=864.48ms p(95)=1s
     ws_sessions...........: 121103 143.837802/s


running (14m01.9s), 000/750 VUs, 120979 complete and 124 interrupted iterations
constVUWith3DirectMessage ✓ [======================================] 80 VUs       45s
rampUp50With100Messages   ✓ [======================================] 124/200 VUs  3m30s
rampUp50WithSingleMessage ✓ [======================================] 000/750 VUs  9m0s
```

![k6 results](/results/screenshot_results_k6.png)

## Services results

![results](/results/screenshot_results.png)

1. **Memory**: ~10x avg difference
2. **CPU**: ~2-3x avg difference
3. **Latency**: ~same
4. **Broadcast Latency**: ~100x avg difference (python won)

## Errors peak in Python App

The first peak started with **~240** users per service on **100 broadcast messages per user**.

The second peak started with **~300** users per service on **single broadcast message**

![python peak](/results/screenshot_peak.png)

![python peak2](/results/screenshot_peak2.png)
