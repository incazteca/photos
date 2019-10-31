# photos

Keep your photos locally and view them on your browser. This is currently a WIP
and not ready for use.

## Development Setup

In postgres create a database user with the following:

```bash
  createuser -s -P photos_user
  createdb photos
```

migrations are run with the [migrate](https://github.com/golang-migrate/migrate) module
Run migrations with:
```bash
migrate -source file:./db/migrations -database postgres://photos_user:{password}@127.0.0.1:5432/photos up
```

Set the following environment variables to configure the application:

DB_USER
DB_PASSWORD
DB_HOST
DB_PORT
