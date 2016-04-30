# zabbix-mongo
Zabbix agent (3.0) native plugin for Mongodb monitoring

This plugin is an agent Loadable Module (https://www.zabbix.com/documentation/3.0/manual/config/items/loadablemodules) .

It can connect to a Mongo server, run basic queries and return a simple value usable by Zabbix server.

We provide an XML template with a few useful example queries that you can import in your Zabbix server (Configuration -> Templates -> Import) .


## Install



## Configure


## Usage

The module is called by defining a regular Zabbix agent item :

`mongo.query[<mongo_connection_url>,<bson_query>,<wanted_value>]`

with:

 - `mongo_connection_url` : the Mongo URL to connect, authenticate, select the database:

 Format: [mongodb://][user:pass@]host1[:port1][,host2[:port2],...][/database][?options]

 Example: mongodb://127.0.0.1/mydb
 

 - `bson_query` : the query in the same format as `db.RunCommand()`. The query must be double-quoted, and internal quotes must be escaped using `\`
 
 Format: https://docs.mongodb.org/manual/reference/command/

 Examples:

`find()` on the "colltest" collection: `"{\"find\": \"colltest\"}"`

`find()` on the "colltest" collection, only retrieve 1 document: `"{\"find\": \"colltest\", \"limit\": 1}"`

get `dbStats()` : `"{\"dbStats\": 1}"`

_Note: using `find` requires Mongo >= 3.2_

 
 - `wanted_value` : The value to pick from the result data. If empty, the complete result will be returned as a JSON string.

 Format: `property.subproperty`
 
 Examples:
 
 If we query `"{\"dbStats\": 1}"`, we get the following JSON object:
 
     {"avgObjSize":0,"collections":0,"dataSize":0,"db":"test","fileSize":0,"indexSize":0,"indexes":0,"numExtents":0,"objects":0,"ok":1,"storageSize":0}

 If we want to get the value of the `objects` property, then `wanted_value` will be `objects`.
