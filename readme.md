# Order statistic tree
An implement of order statistic tree in Go.  
This implement is augmented from [petar/GoLLRB](https://github.com/petar/GoLLRB).   
The augmenting algorithm is from an article on [ustc.edu.cn](http://staff.ustc.edu.cn/~csli/graduate/algorithms/book6/chap15.htm). 


### Changes:
* Add a func to retrieve an element with a given rank (LLRB_GetByRank) 
* Add a func to determine the rank of an element (LLRB_GetRankOf)
