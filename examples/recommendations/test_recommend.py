import pytest
from fastapi.testclient import TestClient
from recommend import app

client = TestClient(app)

def test_music_recommendations():
    """Test music recommendations functionality."""
    response = client.post("/recommend", json={"interest": "music"})
    assert response.status_code == 200
    assert response.json() == {"recommendation": "Listen to some chill beats"}

def test_movie_recommendations():
    """Test movie recommendations functionality."""
    response = client.post("/recommend", json={"interest": "movies"})
    assert response.status_code == 200
    assert response.json() == {"recommendation": "Watch The Shawshank Redemption"}

def test_sports_recommendations():
    """Test sports recommendations functionality."""
    response = client.post("/recommend", json={"interest": "sports"})
    assert response.status_code == 200
    assert response.json() == {"recommendation": "Try playing basketball"}

def test_input_validation_for_invalid_categories():
    """Test that invalid interest categories are properly rejected."""
    response = client.post("/recommend", json={"interest": "astrology"})
    assert response.status_code == 400
    assert "Interest must be one of" in response.json()["detail"]
