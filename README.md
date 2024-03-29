![MonarchStore Logo](/logo.png)

[![Build Status](https://travis-ci.org/MonarchStore/monarchs.svg?branch=master)](https://travis-ci.org/MonarchStore/monarchs)
![GitHub tag](https://img.shields.io/github/tag/MonarchStore/monarchs.svg)
![Docker Pulls](https://img.shields.io/docker/pulls/monarchstore/monarchs.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/MonarchStore/monarchs)](https://goreportcard.com/report/github.com/MonarchStore/monarchs)


A hierarchical, NoSQL, in-memory data store with a RESTful API.

## What is a hierarchical data store? ##
A hierarchical data store is a service which facilitates CRUD actions on data organized in a tree structure with a pre-determined depth.

An example domain may define these relationships:

* Continents have countries
* Countries may have states (or provinces)
* States have cities

The hierarchy is then defined as:
```mermaid
graph TD;
continents --> countries --> states --> cities
```

A hierarchical data store enforces rules to maintain the integrity of the tree structure. All elements need have a parent of the pre-defined type. In our example domain described above, countries cannot exist without a parent continent, and cities cannot be created as direct children of a country.

A hierarchical data store makes it easy to query the stored entities at any depth and retrieve the parent nodes or the children nodes in a single request. For example, to get a country's properties along with the states of that country and the cities in those states, send a request with the `countries` label, the country's ID and the `depth` parameter set to 2.

## Why use a hierarchical data store? ##
Some applications can represent their data model as a hierarchy. A specialized data store can take advantage of optimized data structures for fast reads.
A hierarchichal data store also provides a simple interface to write and query nodes in the hierarchy. This reduces code complexity and development time.

#### How does MonarchStore compare against a relational database? ####
Relational databases, such as MySQL or PostgreSQL, can be slow, since joins are computed at query time. Moreover, they often require the application to make a trade-off between making multiple round trips or repeating the joinned data when a query involves many tables. A hierarchical database can provide the data of an entire hierarchy with a single request, and keeps references to related entities in memory to bypass the join computation at query time.

#### How does MonarchStore compare against a graph database? ####
Graph databases, such as Neo4J or OrientDB, store hard-links between entities and don't have to compute joins at query time. Yet, graph databases provide too much flexibility. In a graph database, any vertex can be linked to any other vertex of any type and vertices can be created even without a parent vertex to attach to. A hierarchical data store enforces every newly created entity to have a parent unless they are meant to be directly under the root.

Applications querying a graph database must construct a query string. These graph queries, in essence, redundantly redeclare the relationships among the entities in the hierarchy. There is a neglible but existent step for the graph databases to parse this query on every request. Entities in a hierarchical data store can only have one parent and zero or more children, and this information only needs to be provided at creation time.

#### How does MonarchStore compare against a key-value store? ####
Key-Value stores, such as Redis or Riak KV, provide the fastest reads when all data is stored as a JSON blob, but frequently modifying a deeply nested value is inneficient. The application must pull the entire hierarchy from the root to the leaves, deserialize it, make changes, serialize it, and update the data store. A hierarchical data store allows applications to efficiently create or update a single nested entity without reading all other entities in the hierarchy.

#### How does MonarchStore compare against a document store? ####
Document stores, such as MongoDB or Elasticsearch, need to make multiple round trips to get the child entities. The alternative is to store all the data into a single, large, highly-nested document, which is inneficient to fetch and update, and so large that it becomes too difficult to handle.

## MonarchStore Features ##
- Queries can return any entity in any of the hierarchy levels.
- Queries can specify the depth level of child entities to reduce the amount of data returned.
- Entities may have custom properties.
- Each entity can be updated independently of it's parent or child entities.
- In memory storage provides speedy reads and writes atomically.
- The RESTful HTTP interface combines the application relational logic and the data store. When there is no additional domain logic to add, putting MonarchStore behind a protected proxy makes writing REST APIs unnecessary. Application architectures can be reduced from a database server + application server to a MonarchStore single server.

## Setting up MonarchStore ##

#### Prerequisites ####
Go > 1.8

#### Compiling into a binary ####
```
go get github.com/MonarchStore/monarchs
cd $GOPATH/github.com/MonarchStore/monarchs
go install
```

#### Running the database ####
```
$GOPATH/bin/monarchs
```

#### Environment ###
```
LISTEN_PORT=":6789"
```

#### CLI Options ###
```
  -addr string
        The binding address (default ":6789")
```

## REST API ##

#### Sample Usage ####
Below are sample REST CRUD actions for a registry of `continents` -> `countries` -> `states` -> `cities`. These actions are available as a [Postman collection](demo_postman_collection.json).

##### Define the "locations" hierarchy
```http
POST http://localhost:6789/locations
```
```json
["continents", "countries", "states", "cities"]
```


##### View the "locations" hiearchy we just created
```http
GET http://localhost:6789/locations
```


##### Create a "continent" document under the "root" document
```http
POST http://localhost:6789/locations/continents/north_america?parent=root
```
```json
{"name": "North America"}
```


##### Create a "country" document
```http
POST http://localhost:6789/locations/countries/usa?parent=north_america
```
```json
{"name": "United States of America", "capital": "Washington, DC", "code": "usa"}
```


##### Create a "state" document
```http
POST http://localhost:6789/locations/states/ny?parent=usa
```
```json
{"name": "New York", "abbr": "NY"}
```


##### Create another "state" document
```http
POST http://localhost:6789/locations/states/ca?parent=usa
```
```json
{"name": "California", "abbr": "CA"}
```


##### Create a "city" document
```http
POST http://localhost:6789/locations/cities/nyc?parent=ny
```
```json
{"name": "New York City"}
```


##### Update a "city" document
```http
PUT http://localhost:6789/locations/cities/nyc
```
```json
{"name": "New York City", "stats": {"population_in_millions": 8.491}}
```


##### Read the "root" document and all the elements in the hierarchy
```http
GET http://localhost:6789/locations/root/root?depth=4
```
The steps above would result in this hierarchy:
```mermaid
graph LR;
root:::current
    --> north_america[North America]:::child
    --> usa[United States of America]:::child
    --> ny[New York State]:::child
    --> nyc[New York City]:::child

usa
    --> ca[California]:::child

classDef current fill:#E3371E
classDef child fill:#0593A2
classDef parent fill:#103778
```



##### Read a "country" document, with all of its "state" documents and "city" documents
```http
GET http://localhost:6789/locations/countries/usa?depth=2
```
Which would result in this hierarchy:
```mermaid
graph LR;
usa[United States of America]:::current
    --> ny[New York State]:::child
    --> nyc[New York City]:::child

usa
    --> ca[California]:::child

classDef current fill:#E3371E
classDef child fill:#0593A2
classDef parent fill:#103778
```


##### Read a "city" document, and the parent "state" document, and the grandparent "country" document
```http
GET http://localhost:6789/locations/cities/nyc?depth=0&parents=2
```
Which would result in this hierarchy:
```mermaid
graph LR;
usa[United States of America]:::parent
    --> ny[New York State]:::parent
    --> nyc[New York City]:::current

classDef current fill:#E3371E
classDef child fill:#0593A2
classDef parent fill:#103778
```

##### Read a "country" document, the parent "continent" document, and the children "state" documents
```http
GET http://localhost:6789/locations/countries/usa?depth=1&parents=1
```
Which would result in this hierarchy:
```mermaid
graph LR;
north_america[North America]:::parent
    --> usa[United States of America]:::current
    --> ny[New York State]:::child
usa
    --> ca[California]:::child

classDef current fill:#E3371E
classDef child fill:#0593A2
classDef parent fill:#103778
```




##### Delete a "country" document. All children "state" documents and grandchilden "city" documents are deleted as well
```http
DELETE http://localhost:6789/locations/countries/usa
```


##### Delete the "locations" hierarchy
```http
DELETE http://localhost:6789/locations
```
