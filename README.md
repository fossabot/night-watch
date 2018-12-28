## Night Watch
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fdashbase%2Fnight-watch.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fdashbase%2Fnight-watch?ref=badge_shield)




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
- [ ] move script/create.go to nightwatch fakeCreate

### How to release

1. checkout your branch and execute `git tag v1.0.0`
2. push your tag `git push origin v1.0.0`


## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fdashbase%2Fnight-watch.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fdashbase%2Fnight-watch?ref=badge_large)