## Go bigcache "bug" demo

Demonstrates how item eviction does not work properly on [bigcache](https://github.com/allegro/bigcache) when using more  
than 1 shards.  


The issue has been reported [here](https://github.com/allegro/bigcache/issues/31) and has been documented with a test  
case [here](https://github.com/allegro/bigcache/pull/33/files) but has not been fixed at the time of writing this.