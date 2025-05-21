import pytest
from fastapi.testclient import TestClient
from event_tracker import app
import uuid
from datetime import datetime

# Create test client
client = TestClient(app)

# Clean up before each test by resetting the API state
@pytest.fixture(autouse=True)
def clear_events():
    # Use the API to clear events (you'll need to add this endpoint)
    client.delete("/track")
    yield

def test_track_event_basic():
    """Test basic event tracking functionality"""
    response = client.post("/track", json={"action": "page_view"})
    
    # Check response
    assert response.status_code == 200
    assert "event_id" in response.json()
    assert response.json()["message"] == "Event logged"
    
    # Check that event was stored via the GET endpoint
    events_response = client.get("/track")
    events = events_response.json()
    assert len(events) == 1
    assert events[0]["action"] == "page_view"
    assert "event_id" in events[0]
    assert "timestamp" in events[0]
    assert isinstance(events[0]["properties"], dict)

def test_track_event_with_properties():
    """Test event tracking with additional properties"""
    response = client.post(
        "/track", 
        json={
            "action": "button_click", 
            "properties": {"button_id": "submit", "page": "checkout"}
        }
    )
    
    # Check response
    assert response.status_code == 200
    
    # Check that event was stored with properties
    events_response = client.get("/track")
    events = events_response.json()
    assert len(events) == 1
    assert events[0]["action"] == "button_click"
    assert events[0]["properties"]["button_id"] == "submit"
    assert events[0]["properties"]["page"] == "checkout"

def test_track_event_missing_action():
    """Test error handling when action is missing"""
    response = client.post("/track", json={"properties": {"test": "value"}})
    
    # Should get validation error
    assert response.status_code == 422
    assert "action" in response.text.lower()

def test_get_events_empty():
    """Test retrieving events when none exist"""
    response = client.get("/track")
    
    assert response.status_code == 200
    assert response.json() == []

def test_get_events_multiple():
    """Test retrieving multiple events"""
    # Add a few events
    client.post("/track", json={"action": "event1"})
    client.post("/track", json={"action": "event2"})
    client.post("/track", json={"action": "event3"})
    
    # Get events
    response = client.get("/track")
    
    assert response.status_code == 200
    events_data = response.json()
    assert len(events_data) == 3
    assert events_data[0]["action"] == "event1"
    assert events_data[1]["action"] == "event2"
    assert events_data[2]["action"] == "event3"

def test_track_event_id_generation():
    """Test UUID generation for events"""
    response = client.post("/track", json={"action": "test_action"})
    event_id = response.json()["event_id"]
    
    # Verify ID is a valid UUID
    try:
        uuid_obj = uuid.UUID(event_id)
        assert str(uuid_obj) == event_id
    except ValueError:
        pytest.fail("Event ID is not a valid UUID")

def test_track_event_timestamp():
    """Test timestamp generation for events"""
    client.post("/track", json={"action": "test_action"})
    
    # Get the event and verify timestamp format
    response = client.get("/track")
    events = response.json()
    
    timestamp_str = events[0]["timestamp"]
    try:
        # Try to parse the timestamp
        timestamp = datetime.fromisoformat(timestamp_str)
        assert timestamp.tzinfo is not None  # Ensure it's timezone-aware
    except ValueError:
        pytest.fail("Event timestamp is not in valid ISO format")

def test_event_persistence():
    """Test that events persist between requests"""
    # Add first event
    client.post("/track", json={"action": "first_event"})
    
    # Add second event
    client.post("/track", json={"action": "second_event"})
    
    # Check that both events exist
    response = client.get("/track")
    assert len(response.json()) == 2
    
    # Verify correct order
    events_data = response.json()
    assert events_data[0]["action"] == "first_event"
    assert events_data[1]["action"] == "second_event"