[BROKEN] by design due to Zabbix agent plugins loading limitations.

# zabbix-mongo
Zabbix agent (3.0) native plugin for Mongodb monitoring

This plugin is an agent Loadable Module (https://www.zabbix.com/documentation/3.0/manual/config/items/loadablemodules) .

It can connect to a Mongo server, run basic queries and return a simple value usable by Zabbix server.

An XML template with a few useful example queries that you can import in your Zabbix server (Configuration -> Templates -> Import) is provided.


## Install

_Note: Until there's an easier way to distribute this module as a package, an **archive** containing a **already compiled** module (x86-64) and the XML template is available [here](https://share.zabbix.com/component/mtree/dir-libraries/zabbix-loadable-modules/mongodb-monitoring-loadable-module?Itemid=_)_

Build from source: has only been tested with Go 1.6. 
The quickest and easiest way to build the project is to use GVM.

### Prepare the Go environment
Go >= 1.5 is required. But to install it we first need a working Go 1.4 compiler.


 - Install GVM:

   `bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)`

 - Load GVM commands and environment:

   `source ${HOME}/.gvm/scripts/gvm`

 - Install Go 1.4

   `gvm install go1.4`

 - Enter the Go 1.4 env and install Go 1.6

   `gvm use go1.4`

    `export GOROOT_BOOTSTRAP=$GOROOT`

 - Install a fresh Go 1.6 env:

   `gvm install go1.6`
   

### Build the project
 
 - Clone this repo, fetch Go dependencies, then build it:

   ` git clone https://github.com/m-barthelemy/zabbix-mongo.git`

   ` cd zabbix-mongo`

   ` gvm use go1.6`

   `go get gopkg.in/mgo.v2 gopkg.in/cavaliercoder/g2z.v3 github.com/mattn/go-scan`

   `make`


If the build succeeds, it creates the `zbx_mongo.so` library, that can be loaded by the Zabbix agent.


## Zabbix configuration

Copy the built `zbx_mongo.so` to a server having a Zabbix Agent installed, for example into `/etc/zabbix/zbx_mongo.so`.

Edit the Agent configuration file (`/etc/zabbix/zabbix_agentd.conf`) to tell it where are the loadable modules and which ones should be loaded:

    LoadModulePath=/etc/zabbix
    LoadModule=zbx_mongo.so

Now restart the Zabbix Agent. In its log file, you should see the confirmation that the zbx_mongo module has been successfully loaded:

    4255:20160430:201554.976 using configuration file: /etc/zabbix/zabbix_agentd.conf
    4255:20160430:201555.028 loaded modules: zbx_mongo.so


## Usage

The module is called by defining a regular Zabbix agent item :

`mongo.run[<mongo_url>,<command>,<query>,<result_path>]`

with:


### `mongo_url `:

The Mongo URL to connect, authenticate, select the database:

 Format: `[mongodb://][user:pass@]host1[:port1][,host2[:port2],...][/database][?options]`

 Example: `mongodb://127.0.0.1/mydb`
 



### `command` and `query`: 

 - `command` is the command name such as `serverStatus`, `dbStats`, `find`, `count`,...

 -  `query` is the complete JSON query (including command). It **must** be double-quoted, and internal quotes must be escaped using `\`

     Its format closely follows the format of `db.RunCommand()`. See [MongoDB documentation](https://docs.mongodb.org/manual/reference/command/)


 


 Examples:

    `find()` on the "colltest" collection: `"{\"find\": \"colltest\"}"`

    `find()` on the "colltest" collection, only retrieve 1 document: `"{\"find\": \"colltest\", \"limit\": 1}"`

    get `dbStats()` : `"{\"dbStats\": 1}"`

    _Note: using `find` requires Mongo >= 3.2_

 


### `result_path`:

The path to the wanted value. If empty, the complete result will be returned as a JSON string. if the path points to something not being a simple value, the content is returned as a JSON string.

 Format: `/property/subproperty`
 
 Examples:
 
 If we query `"{\"dbStats\": 1}"`, we get the following JSON object:
 
     {"avgObjSize":0,"collections":0,"dataSize":0,"db":"test","fileSize":0,"indexSize":0,"indexes":0,"numExtents":0,"objects":0,"ok":1,"storageSize":0}

 If we want to get the value of the `dataSize ` property, then `wanted_value` will be `/dataSize `.

 If we want a single value from an array, we can fetch it by its index : `/path/to/array[0]`

 More examples in the `zbx_mongo_template.xml` file.

## 


####Complete examples 

`mongo.run[mongodb://127.0.0.1/myDb, dbStats,  "{\"dbStats\": 1}", /dataSize]`

`mongo.run[mongodb://127.0.0.1/myDb, serverStatus, "{\"serverStatus\":1, \"repl\":0, \"metrics\":0}",/connections/totalCreated]`


## Roadmap / TODO
 
 - More documentation examples

 - Discovery item `mongo.discover[....]` to have a super easy way to populate databases, collections in Zabbix and monitor them individually.

 - XML Zabbix template improvements
