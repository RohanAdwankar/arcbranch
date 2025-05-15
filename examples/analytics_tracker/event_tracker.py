from fastapi import FastAPI, HTTPException
from pydantic import BaseModel, Field
from typing import Dict, Any, List, Optional
from datetime import datetime, timezone  # Import timezone separately
import uuid

app = FastAPI(title="Simple Analytics Tracker")

# In-memory event storage
events = []

class Event(BaseModel):
    """Data model for analytics events"""
    action: str = Field(..., description="The action performed by the user")
    properties: Optional[Dict[str, Any]] = Field(default=None, description="Additional event properties")
    
    def to_record(self) -> Dict[str, Any]:
        """Convert to a storage record with metadata"""
        return {
            "id": str(uuid.uuid4()),
            "timestamp": datetime.now(timezone.utc).isoformat(),  # Use timezone.utc instead
            "action": self.action,
            "properties": self.properties or {}
        }

@app.post("/track")
def track_event(event: Event):
    """Track a new analytics event"""
    # Create a complete event record with metadata
    event_record = event.to_record()
    
    # Store the event
    events.append(event_record)
    
    return {
        "message": "Event logged",
        "event_id": event_record["id"]
    }

@app.get("/track")
def get_events():
    """Retrieve all tracked events"""
    return events