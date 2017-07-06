# KingDB #
A hierarchical in-memory data store with a RESTful API.

## Why ##
Some applications can represent their data model as a hierarchy of entities. A specialized data store can take advantage of optimized data structures internally, simplify the query model and reduce development time and complexity.

Realtional database can be slow since joins are performed at query time. Moreover, they often require the application to make a trade off between making multiple round trips if there are many tables involved or repeating joinned data.

Graph databases perform fast joins, but provide too much flexibility in how entities are related. For a hierarchy data model, the client application is left with the responsibility of defining the edges that connect the vertices in every query. Graph databases also need to parse is a mostly-neglible, but still existent query.

Key-Value stores provide the fastest reads when all data is stored as a JSON blob, but frequently modifyiing a deeply nested value becomes an inneficient, and sometimes difficult task.

Document stores, like relational databases, need to make multiple round trips to get the child entities. The alternative is to store all the data into a single, highly-nested document, which is inneficient to update, and difficult to work with.

## Features ##
- Queries can return any entity in the hierarchy and queries can specify the depth level to reduce the amount of data returned.
- Each entity can be updated independently of it's parent or child entities.
- In memory storage provides fast atomic reads and writes.
- A RESTful HTTP interface combines the application and the data store. Putting KingDB behind an API Gateway makes writing REST APIs unnecessary, simplifies deployments, and reduces the number of servers.

## Setting up KingDB ##

#### Prerequisites ####
Go > 1.8

#### Compiling into a binary ####
```
go get bitbucket.org/enticusa/kingdb
cd $GOPATH/bitbucket.org/enticusa/kingdb
go install
```

#### Running the database ####
```
$GOPATH/bin/kingdb --labels parent,child,grandchild
```

#### CLI Options ###
```
  -addr string
        The binding address (default ":6789")
  -labels string
        A comma-separated, ordered list of the labels of elements in the hieararchy
```

## REST API ##