from fastapi import FastAPI
import requests

app = FastAPI()

@app.get("/joke")
def get_joke():
    response = requests.get("https://official-joke-api.appspot.com/jokes/random")
    return response.json()
