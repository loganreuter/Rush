# Rush
## NoSQL Database Written Purely in Go

## Installation
```bash
go get github.com/lreuter2020/Rush
```
or 

Add the library to import statement and run
```bash
go mod tidy
```

## Usage
### Create a database connection:
```go
func main(){
    /* 
        Creates a connection to database.
        First argument is relative path
        Second argument is the name of the database
    */
    conn, err := rush.Connect("/", "db")
    if err != nil {
        panic(err)
    }

    

    /*
    
    */
    users.Member("Guest").Create(map[string]interface{}{
        "Name": "Guest",
        "Age": 18
    })
}
```

### Create A New Group
_(Equivalent of creating a table in SQL)_
```go
    ...
    
    /*
        Creates a collection titled Users.
    */
    users := conn.Group("Users")
    
    ...
```

### Create A New Member
```go
    ...

    /*
        Method 1:
    */
    users.Member("Guest")

    /*
        Method 2:
    */
    conn.Group("Users").Member("Guest")

    ...
```