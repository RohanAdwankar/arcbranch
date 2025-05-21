from fastapi import FastAPI
from pydantic import BaseModel
from typing import Optional, Dict, Any
from datetime import datetime, timezone
import uuid

app = FastAPI()

# In-memory store for events; each event should be a dictionary with the following keys:
# - event_id (UUID string)
# - timestamp (ISO 8601 string in UTC)
# - action (string from request)
# - properties (dict from request, default to empty dict if None)
events = []

# Request model for incoming events
class Event(BaseModel):
    action: str
    properties: Optional[Dict[str, Any]] = None

@app.post("/track")
def track_event(event: Event):
    # 1. Generate a UUID for the event_id.
    # 2. Get current UTC timestamp and format it as ISO 8601.
    # 3. Create a dictionary with keys: event_id, timestamp, action, properties.
    # 4. If properties is None, store an empty dict.
    # 5. Append this event dictionary to the global 'events' list.
    # 6. Return a dict: {"message": "Event logged", "event_id": <event_id>}
    ...

@app.get("/track")
def get_events():
    # Return the raw list of event dictionaries (not wrapped in another dict).
    ...

@app.delete("/track")
def clear_events():
    # Clear the global 'events' list.
    # Return: {"message": "All events cleared"}
    ...
