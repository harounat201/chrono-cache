# Semantic Proxy Layer for LLMs

Required:
- Find a way to vectorize queries and store them in a database
- Find a way to perform similarity operations on new queries
- Create a procedure to vectorize outgoing queries, run similarity, and if exists in a KV store, return the associated query


Specifications:
- What data structures will work at scale?
- How will "similarity" be determined?
- What is the actual data structure associated with a KV store that will work at scale?
- What conditions do queries need to meet in order to exist in the KV store?
- How are queries added?


Workflow
--> Embed queries 
--> run cosine similarity 
--> build KV based on conditions: (latency too high? durable/binary question answer? [intent] query frequency?) Build read replicas on hot keys to mitigate traffic congestion 
--> 

File system
