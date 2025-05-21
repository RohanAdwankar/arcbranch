import pytest
from fastapi.testclient import TestClient
from jokes_api import app

client = TestClient(app)

def test_joke_endpoint_returns_success():
    """Test that the joke endpoint returns a successful status code."""
    response = client.get("/joke")
    assert response.status_code == 200

def test_joke_endpoint_returns_json():
    """Test that the joke endpoint returns valid JSON."""
    response = client.get("/joke")
    assert response.headers["Content-Type"] == "application/json"
    # This will succeed only if valid JSON is returned
    response.json()

def test_joke_contains_required_fields():
    """Test that the joke contains all required fields."""
    response = client.get("/joke")
    joke = response.json()
    
    assert "id" in joke
    assert "type" in joke
    assert "setup" in joke
    assert "punchline" in joke

def test_joke_api_has_rate_limiting():
    """Test that the API implements rate limiting (this will fail)."""
    # Make 10 rapid requests to test rate limiting
    for _ in range(10):
        response = client.get("/joke")
    
    # Make one more request that should be rate limited
    response = client.get("/joke")
    
    # This will fail if API doesn't implement rate limiting
    assert response.status_code == 429
    assert "rate limit exceeded" in response.json().get("detail", "").lower()