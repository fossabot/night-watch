## Night Watch



### Depend
need golang version >= 1.11.
we used go module
so maybe you need to set ```export GO111MODULE=on```


### 
build
```bash
make build_linux
```


### TEST
```
On Macbook pro 15 2017 
100 file * 10byte/s
get 995 - 998 byte/s ,  0.2 -- 0.5%

1000 file * 10byte/s
get 9931 - 9963 byte/s,   0.4 -- 0.7 %
```





### TODO
- [ ] support "b*/**/z*.txt"
```
1. use filebeat code [doubleStarPatternDepth]
2. https://github.com/mattn/go-zglob
```
- [ ] use https://github.com/karrick/godirwalk to improve performance