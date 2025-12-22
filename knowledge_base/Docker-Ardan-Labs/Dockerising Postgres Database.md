#docker #database #postgres

### db

This is a PostgreSQL database.

The database must be initialized with the schema (database and tables)
and the data (used by the application).

The file `words.sql` contains all the SQL commands necessary to create
the schema and load the data.

```
# cat words.sql
CREATE TABLE nouns (word TEXT NOT NULL);
CREATE TABLE verbs (word TEXT NOT NULL);
CREATE TABLE adjectives (word TEXT NOT NULL);

INSERT INTO nouns(word) VALUES
('cloud'),
('elephant'),
('gø language'),
('laptøp'),
('cøntainer'),
('micrø-service'),
('turtle'),
('whale'),
('gøpher'),
('møby døck'),
('server'),
('bicycle'),
('viking'),
('mermaid'),
('fjørd'),
('legø'),
('flødebolle'),
('smørrebrød');

INSERT INTO verbs(word) VALUES
('will drink'),
('smashes'),
('smøkes'),
('eats'),
('walks tøwards'),
('løves'),
('helps'),
('pushes'),
('debugs'),
('invites'),
('hides'),
('will ship');

INSERT INTO adjectives(word) VALUES
('the exquisite'),
('a pink'),
('the røtten'),
('a red'),
('the serverless'),
('a brøken'),
('a shiny'),
('the pretty'),
('the impressive'),
('an awesøme'),
('the famøus'),
('a gigantic'),
('the gløriøus'),
('the nørdic'),
('the welcøming'),
('the deliciøus');
```

Additional information:

- we strongly suggest using the official PostgreSQL image that can be found on the Docker Hub (it's called `postgres`)
- if we check the [page of that official image](https://hub.docker.com/_/postgres) on the Docker Hub, we will find a lot of documentation; the section "Initialization scripts" is particularly useful to understand how to load `words.sql`
- it is advised to set up password authentication for the database; but in this case, to make our lives easier, we will simply authorize all connections (by setting environment variable `POSTGRES_HOST_AUTH_METHOD=trust`)

### Dockerfile
```
FROM postgres
COPY words.sql /docker-entrypoint-initdb.d
```
need to include an environment variable for `POSTGRES_HOST_AUTH_METHOD=trust`:

```
FROM postgres
COPY words.sql /docker-entrypoint-initdb.d
ENV POSTGRES_HOST_AUTH_METHOD= trust
```
-  other option:
```sh
$ docker build . -t db
$ docker run -e POSTGRES_HOST_AUTH_METHOD=trust db
```
