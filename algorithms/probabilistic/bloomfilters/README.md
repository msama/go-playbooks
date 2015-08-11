#Identify duplicates using a bloom filter

Bloom filters are great for identifying cache miss/hits in a probabilistic way. 
The benefit of using a bloom filter is that with tiny footprint it is possible to estimate the unicity of element sina dataset.

In this experiment the bloom filter stored 10M UUIDs with an error rate of 0.009730%

```
Added 10000000 elements, 973 false hits.
Miss rate with 10000000 elements was 0.009730. Estimated was 0.003280.
```

Given that each UUID (in its string representation `6ba7b810-9dad-11d1-80b4-00c04fd430c8`) uses 36B, storing all of them would require 360MB. Instead the size of the bloom filter is only 25MB, achieving a 93% memory saving. 


```
(pprof) top5    
24928.14kB of 24928.14kB total (  100%)
Dropped 15 nodes (cum <= 124.64kB)
Showing top 5 nodes out of 8 (cum >= 24416kB)
      flat  flat%   sum%        cum   cum%
   24416kB 97.95% 97.95%    24416kB 97.95%  github.com/willf/bitset.New
  512.14kB  2.05%   100%   512.14kB  2.05%  mcommoninit
         0     0%   100%    24416kB 97.95%  github.com/willf/bloom.New
         0     0%   100%    24416kB 97.95%  main.main
         0     0%   100%    24416kB 97.95%  runtime.goexit
```