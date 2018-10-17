# MonarchStore #

[![Build Status](https://travis-ci.org/arturom/monarchs.svg?branch=master)](https://travis-ci.org/arturom/monarchs)


A hierarchical, NoSQL, in-memory data store with a RESTful API.

## What is a hierarchical data store? ##
A hierarchical data store is a service which facilitates CRUD actions on data organized in predefined structures with parent-child relationships.

An example domain may define these relationships:

* Continents have countries
* Countries may have states (or subdivisions)
* States have cities

The hierarchy is then defined as: `continents` -> `countries` -> `states` -> `cities`

A hierarchical data store enforces all elements to have a parent of the predefined type. In our example domain described above, countries cannot exist without a parent continent, and cities cannot be created as direct children of a country.

A hierarchical data store makes it simple to query the stored entities and their relationships. To get a country's properties along with the states of that country and the cities in those states, an application issues a query with the `countries` label, the country's ID and the `depth` parameter set to 2.

## Why use a hierarchical data store? ##
Some applications can represent their data model as a hierarchy of entities. A specialized data store can take advantage of optimized data structures, simplify the write and query model, reduce development time, and reduce code complexity.

#### How does it compare against a relational database? ####
Relational databases, such as MySQL or PostgreSQL, can be slow, since joins are computed at query time. Moreover, they often require the application to make a trade-off between making multiple round trips or repeating the joinned data when a query involves many tables. A hierarchical database can provide the data of an entire hierarchy with a single request, and keeps references to related entities in memory to bypass the join computation at query time.

#### How does it compare against a graph database? ####
Graph databases, such as Neo4J or OrientDB, store hard-links between entities and don't have to compute joins at query time. Yet, graph databases provide too much flexibility. In a graph database, any vertex can be linked to any other vertex of any type and vertices can be created even without a parent vertex to attach to. A hierarchical data store enforces every newly created entity to have a parent unless they are meant to be directly under the root.

Applications querying a graph database must construct a query string. These graph queries, in essence, redundantly redeclare the relationships among the entities in the hierarchy. There is a neglible but existent step for the graph databases to parse this query on every request. Entities in a hierarchical data store can only have one parent and zero or more children, and this information only needs to be provided at creation time.

#### How does it compare against a key-value store? ####
Key-Value stores, such as Redis or Riak KV, provide the fastest reads when all data is stored as a JSON blob, but frequently modifying a deeply nested value is inneficient. The application must pull the entire hierarchy from the root to the leaves, deserialize it, make changes, serialize it, and update the data store. A hierarchical data store allows applications to efficiently create or update a single nested entity without reading all other entities in the hierarchy.

#### How does it compare against a document store? ####
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
go get github.com/arturom/monarchs
cd $GOPATH/github.com/arturom/monarchs
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
```
# Create the "locations" schema
POST localhost:6789/locations
["continents", "countries", "states", "cities"]

# Read the "locations" schema we just created
GET localhost:6789/locations

# Create a "continent" under the "root"
POST localhost:6789/locations/continents/north_america?parent=root
{"name": "North America"}

# Create a "country"
POST localhost:6789/locations/countries/usa?parent=north_america
{"name": "United States of America", "capital": "Washington, DC", "code": "usa"}

# Create a "state"
POST localhost:6789/locations/states/ny?parent=usa
{"name": "New York", "abbr": "NY"}

# Create another "state"
POST localhost:6789/locations/states/ca?parent=usa
{"name": "California", "abbr": "CA"}

# Create a "city"
POST localhost:6789/locations/cities/nyc?parent=ny
{"name": "New York City"}

# Update a "city"
PUT localhost:6789/locations/cities/nyc
{"name": "New York City", "stats": {"population_in_millions": 8.491}}

# Read a "country", with all its "states" and "cities"
GET localhost:6789/locations/countries/usa?depth=2

# Read a "city", and its parent "state" and "country"
GET localhost:6789/locations/cities/nyc?depth=0&parents=2

# Read the entire hierarchy root
GET localhost:6789/locations/root/root?depth=4

# Read the entire hierarchy root
GET localhost:6789/locations/root/root?depth=4

# Delete a "country". All child "states" and "cities" are deleted as well
DELETE localhost:6789/locations/countries/usa

# Delete the "locations" schema
DELETE localhost:6789/locations
```