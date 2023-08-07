# gogoat

A terminal user interface (TUI) way to view stats from the [goatcounter](https://www.goatcounter.com/). It provides a fast and easy way to check-up on your running website stats.

## What's goatcounter

GoatCounter is a lightweight web analytics platform that provides insights into website traffic while respecting user privacy. It distinguishes itself by its simplicity and focus on privacy, avoiding the use of cookies and not collecting personal data. GoatCounter offers essential tracking features, such as page views, referrers, and popular pages, through an intuitive and clean interface. It's particularly suitable for individuals and small businesses looking for a privacy-conscious way to understand their website's performance without compromising on user data security.

## How to use

First, you have to set up the goatcounter dashboard and have two things ready:

1. Site `prefix` - for example if your goatcounter site is `stats.goatcounter.com` the prefix is `stats`
2. The API token

Both of those should be put in an `.env` file in the proect directory:

```env
TOKEN="tokentokentokentokentokentokentokentokentokentoken"
SITE_PREFIX="stats"
```

Then you can run `go run main.go` and get basic page hits. Consult the `client/data.go` file for details.
