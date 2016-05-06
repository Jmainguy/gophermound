# gophermound
A golang daemon to display zookeeper information.

##How to use
Currently gomound listens on 8080 by default.
It provides 4 endpoints
/mntr which provides output of mntr command to zookeeper
/stat which provides output of stat command to zookeeper
/ruok which provides output of ruok command to zookeeper
/connections which provides current connections (included in /stat)
