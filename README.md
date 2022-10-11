# ssm2dotenv

Inspired by [ssm-env](https://github.com/remind101/ssm-env). 

Usually, if you are using libraries like [godotenv](https://github.com/joho/godotenv) to load environment variables from a `.env` file, but you won't like to place the secrets in the `.env` file, you can use `ssm2dotenv` command to load the secrets from AWS SSM Parameter Store.

## Usage

Suppose that you have the following parameters in aws SSM:

| key                  | value                       |
|----------------------|-----------------------------|
| /app-api/test/DB_DSN | postgres://xxx:yyyy/test_db |
| /app-api/prod/DB_DSN | postgres://xxx:yyyy/prod_db |

And you have the following `sample.env` template file:
```dotenv
ENV=test
DB_META_DSNS=ssm:///app-api/${env}/DB_DSN
PORT=8080
```

Then, you can run the following command to get the env file

```bash
$ ssm2dotenv --env test -i ./sample.env -o ./.env
```

So you can get the env file including secrets just only existing during the local development or your CI/CD pipeline. 

```dotenv
ENV=test
DB_META_DSNS=postgres://xxx:yyyy/test_db
PORT=8080
```
