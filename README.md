# Blog-Aggregator
RSS feed aggregator in GO
A multi-player command line tool for aggregating RSS feeds and viewing the posts.

## Installation

Make sure you have the laterest [Go toolchain](https://golang.org/dl/) installed as well as a local Postgres database. You can then install the Aggregator with:

```bash 
go install github.com/tepidmilk/gator 
```

## Config 

Create a `.gatorconfig.json` file in your home directory with the following structure:

```json
{
    "db_url": "postgres://username:@localhost:5432/database?sslmode=disable"
}
```

Replace the value with your database connection string.

## Usage

Create a new user:

```bash
gator register <name>
```

Add a feed:

```bash
gator addfeed <url>
```

Start the aggregator:

```bash
gator agg 30s
```

View the posts:

```bash
gator browse [limit]
```

There are a few other commands you'll need as well:

- `gator login <name>` - Log in as a user that already exists
- `gator users` - List all users
- `gator feeds` - List all feeds
- `gator follow <url>` - Follow a feed that already exists in the database
- `gator unfollow <url>` - Unfollow a feed that already exists in the database
