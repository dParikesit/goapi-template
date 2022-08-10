# Golang Template

This is a template to get started building on Coherence (withcoherence.com)

## To use

- Sign up for an account at app.withcoherence.com
- Import this template into a new repo
- Follow the onboarding steps on the Coherence site to set up your Cloud IDE, automatic preview environments, CI/CD pipelines, and managed cloud infrastructure

## To connect to the database

- Run the toolbox from the Coherence UI
- Run a terminal in the toolbox
- Use cocli to run commands in the running instance. Example:

To dump a database:

```console
cocli exec backend 'DB_SOCKET_NAME=$(compgen -A variable | grep _SOCKET) && PGPASSWORD="$DB_PASSWORD" pg_dump -h ${!DB_SOCKET_NAME} -U $DB_USER $DB_NAME'
```

To run the psql cli:

```console
cocli exec backend 'DB_SOCKET_NAME=$(compgen -A variable | grep _SOCKET) && PGPASSWORD="$DB_PASSWORD" psql -h ${!DB_SOCKET_NAME} -U $DB_USER $DB_NAME'
Going to run command (backend): [DB_SOCKET_NAME=$(compgen -A variable | grep _SOCKET) && PGPASSWORD="$DB_PASSWORD" psql -h ${!DB_SOCKET_NAME} -U $DB_USER $DB_NAME]
psql (13.7 (Debian 13.7-1.pgdg100+1))
Type "help" for help.

main=> \l
                                                List of databases
     Name      |       Owner       | Encoding |  Collate   |   Ctype    |            Access privileges            
---------------+-------------------+----------+------------+------------+-----------------------------------------
 cloudsqladmin | cloudsqladmin     | UTF8     | en_US.UTF8 | en_US.UTF8 | 
 main          | cloudsqlsuperuser | UTF8     | en_US.UTF8 | en_US.UTF8 | 
 postgres      | cloudsqlsuperuser | UTF8     | en_US.UTF8 | en_US.UTF8 | 
 template0     | cloudsqladmin     | UTF8     | en_US.UTF8 | en_US.UTF8 | =c/cloudsqladmin                       +
               |                   |          |            |            | cloudsqladmin=CTc/cloudsqladmin
 template1     | cloudsqlsuperuser | UTF8     | en_US.UTF8 | en_US.UTF8 | =c/cloudsqlsuperuser                   +
               |                   |          |            |            | cloudsqlsuperuser=CTc/cloudsqlsuperuser
(5 rows)

main=>
main=> SELECT * FROM message;
     value     
---------------
 Message one
 Message two
 Message three
(3 rows)

main=> 
```
