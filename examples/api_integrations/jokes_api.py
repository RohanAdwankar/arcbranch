from fastapi import FastAPI
import requests

app = FastAPI()

@app.get("/joke")
def get_joke():
    # 1. Send a GET request to the Official Joke API's /jokes/random endpoint.
    
    # 2. Parse the JSON response.
    
    # 3. Return the parsed JSON directly as the response.
    ...