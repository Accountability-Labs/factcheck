# FactCheck backend

## Development

1. Set up a PostgreSQL database.
2. Go to the `sql/schema` directory and run:

       goose postgres <DB_URL> UP
3. Customize your `.env` file.
4. Compile and run the service:

       make && export $(cat .env | xargs) && ./factcheck