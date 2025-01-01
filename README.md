To run the application locally: DATABASE_URL="postgres://postgres:9363T@localhost:5433/evolve?sslmode=disable" PORT="3000" go run main.go

This implementation assumes:
that external apis are being called for certain things
that a one day cache is created
that authentication is taken care of
that account creation is implemented already and picks from the db
How I can Improve on it if I have more time:
<style> .button { background-color: orangered; border-radius: 4px; border: none; padding: 10px 5px 10px 5px; } .button:hover { background-color: black; } .a { color: white; } .p { padding-top: 40px; } </style>
API DOCUMENTATION

[Documentation](https://documenter.getpostman.com/view/16161718/2s93eWys9i)
