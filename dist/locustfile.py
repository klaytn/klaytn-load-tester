from locust import HttpUser


class WebsiteUser(HttpUser):
    host = "http://localhost:8551"
    maxsize = 10000

    min_wait = 0
    max_wait = 0
