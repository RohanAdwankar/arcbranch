# Below is the skeleton for the event tracker API:
from fastapi import FastAPI
from pydantic import BaseModel

app = FastAPI()

# store the events somewhere

class Event(BaseModel):
    # define action and optional properties
    ...

@app.post("/track")
def track_event(event: Event):
    # turn event into a record with uuid and timestamp
    # store it
    # return message and event id
    ...

@app.get("/track")
def get_events():
    # return all stored events
    ...
