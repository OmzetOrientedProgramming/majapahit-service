[![pipeline status](https://gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/badges/master/pipeline.svg)](https://gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/-/commits/master) [![coverage report](https://gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/badges/master/coverage.svg)](https://gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/-/commits/master)

# Majapahit Service

## Description

This is backend main repository for Software Engineering Project

## Installation

1. Clone this repository
    ```
    git clone https://gitlab.com/omzet-oriented-programming/majapahit-service.git 
    ```
2. Installing all dependency
    ```
    go install
    ```
3. Create .env file using .env.example and fill the environment variable you need (database data, port, etc)
4. Activate postgreSQL
    1. If you have docker-compose you can run:
         ```
         docker-compose up -d
         ```
        * note: The default user is **postgres** and the password is **root**. You must make sure the user and password
          is the same in the .env file with the docker-compose.yml file.
    2. If you don't have it, you can install it [here](https://docs.docker.com/engine/install/) or you can manually
       install PostgreSQL and set up the user and password and make sure they are the same as the .env file
       configuration
5. Go to your postgreSQL console and create the database using this command
    ```
    CREATE DATABASE <your_database_name>
    ```
    * note: make sure the name of database is the same as the .env file configuration. Default name is majapahitdb
6. If you are not installed `go-migrate-cli` yet. You need to install it first to do migrations. You can check
   it [here](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate). Then just run:
    ```
   make migrate-up [N]
   ```
   or
    ```
   make migrate-down [N]
   ```
    * note : N is an optional argument to do how many migrations you want to. Leave it blank if you want to do all.
7. Run the server
    ```
   go run main.go
   ```

### Using pgAdmin

If you are using docker I have set up the pgAdmin as database client GUI, you can use it if you want to.

1. Make sure you already run the docker-compose service, you can run:
    ```
   docker-compose up -d
   ```
2. Go to the pgAdmin URL, default URL will be `http://localhost:5050`
3. Create new server and fill some information
    1. Name : `<your-server-name>`
    2. Host : `host.docker.internal`
    3. Port : `5432`
    4. Maintenance Database : `majapahitdb`
    5. Username : `postgres`
    6. Password : `root`
4. Save

## Development

1. Before development create new branch from master or staging.
   ```
   git checkout master && git pull origin master
   ```
   or
   ```
   git checkout staging && git pull origin staging
   ```
2. Create new branch using naming convention `PBI-<pbi_number>-OOP-<jira_ticket_id>-<name_of_your_feature>`. Example
    ```
   git checkout -b PBI-1-OOP-1-registration
   ```
3. Start developing you feature with TDD
    1. Red phase use `[RED]` tag as prefix of your commit
       ```
       git commit -m "[RED] create success test case for registration"
       ```
    2. Green phase use `[GREEN]` tag as prefix of your commit
       ```
       git commit -m "[GREEN] create implementation for success registration"
       ```
    3. Refactor phase use `[REFACTOR]` tag as prefix of your commit
       ```
       git commit -m "[REFACTOR] improve clean code for registration"
       ```
    4. Non related TDD phase will use `[CHORES]` tag as prefix of your commit
       ```
       git commit -m "[CHORES] update readme"
       ```
4. Run all the pre-deployment step
    1. Format
       ```
        make fmt
       ```
    2. Lint
       ```
        make lint
       ```
    3. Test
       ```
        make fmt
       ```
    4. Coverage
       ```
        make fmt
       ```
5. If all the step is passed, you can push your commit to remote
   ```
   git push
   ```
   