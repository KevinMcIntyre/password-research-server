# Temple University Password Research Tool

## Please Note:
- Go is required to build this project. Although any version >= 1.5.x should work fine, 1.7.x is recommended.
- PostgreSQL is required to use this project. Version >= 1.9.5 is recommended.

There are currently no config files. This project requires a database called "tupwresearch" owned by the postgres user whose password is "password". Until config file support is created, changes to these settings must be made in main.go in the "setupDatabase" function.

Build instructions:

1. Use the 2.sql file in the sql directory to seed the database.

2. Build the executable in the main directory:

    $ go build main.go

3. Run the executable:

    $ ./main