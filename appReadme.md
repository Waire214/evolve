<b>To run the application locally: </b> <i>DATABASE_URL="postgres://postgres:9363T@localhost:5433/evolve?sslmode=disable" PORT="3000" go run main.go</i>

<ul>
<b>This implementation assumes: </b>
<li>that external apis are being called for certain things</li>
<li>that a one day cache is created</li>
<li>that authentication is taken care of</li>
<li>that account creation is implemented already and picks from the db</li>
</ul>

<ul>
<b>How I can Improve on it if I have more time: </b>
</ul>

<style>
.button {
    background-color: orangered;
    border-radius: 4px;
    border: none;
    padding: 10px 5px 10px 5px;
}
.button:hover {
    background-color: black;
}
.a {
    color: white;
}
.p {
padding-top: 40px;
}
</style>
<p class="p">API DOCUMENTATION</p>
<button class="button"><a href="https://documenter.getpostman.com/view/16161718/2s93eWys9i" class="a">Documentation</a></button>