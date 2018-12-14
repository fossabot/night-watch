## Night Watch



### Depend
need golang version >= 1.11.
we used go module
so maybe you need to set ```export GO111MODULE=on```


### 
build
```bash
export GO111MODULE=on make build_linux
```


### TODO
- [ ] support "b*/**/z*.txt"
```
1. use filebeat code [doubleStarPatternDepth]
2. https://github.com/mattn/go-zglob
```
- [ ] use https://github.com/karrick/godirwalk to improve performance