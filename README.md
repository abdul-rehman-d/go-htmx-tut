# Go + HTMX

A tutorial (once again)  
But I deployed it manually on AWS and Docker

## Deployment
- Well I bought my domain and learned how to setup subdomains just on the dashboard to point to EC2 instance.
- I used Caddy as reverse proxy to my Docker container which is running the go app.

## Future plans
Well, not much but:
- CI/CD
- Add DB and users
- Write another go app to vigoursly test the tiny free tier EC2 instances. If I throw money in AWS then I wanna test how much concurrent users it can tackle with random logins, creates, deletes, fetches, logouts. Basically spawn a bunch of users and test the bad boy with simulation of real user behaviour. Possibly reinventing the wheel here but I am just interested to do it.

## References
- [YT Video](https://www.youtube.com/watch?v=x7v6SNIgJpE)
