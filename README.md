# ssm2dotenv

## Usage

Suppose that you have the following parameters in aws SSM:

| key                  | value                       |
|----------------------|-----------------------------|
| /app-api/test/DB_DSN | postgres://xxx:yyyy/test_db |
| /app-api/prod/DB_DSN | postgres://xxx:yyyy/prod_db |

The sample env file
```dotenv
ENV=test
DB_META_DSNS=ssm:///app-api/${env}/DB_DSN
PORT=8080
```

Then, you can run the following command to get the env file

```bash
$ ssm2dotenv --env test -i ./sample.env -o ./.env
```

You will get the following env file

```dotenv
ENV=test
DB_META_DSNS=postgres://xxx:yyyy/test_db
PORT=8080
```
