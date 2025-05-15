import pytest
from fastapi.testclient import TestClient
from auth import app

client = TestClient(app)

def test_valid_login():
    """Test successful login with correct credentials."""
    response = client.post(
        "/login",
        data={"username": "user", "password": "pass"},
        headers={"Content-Type": "application/x-www-form-urlencoded"}
    )
    assert response.status_code == 200
    assert response.json() == {"message": "Login successful"}

def test_invalid_credentials():
    """Test login failure with incorrect credentials."""
    response = client.post(
        "/login",
        data={"username": "wrong", "password": "wrong"},
        headers={"Content-Type": "application/x-www-form-urlencoded"}
    )
    assert response.status_code == 401
    assert "Invalid credentials" in response.json()["detail"]

def test_admin_login():
    """Test admin user login with appropriate permissions."""
    response = client.post(
        "/login",
        data={"username": "admin", "password": "admin123"},
        headers={"Content-Type": "application/x-www-form-urlencoded"}
    )
    assert response.status_code == 200
    assert response.json() == {"message": "Login successful", "role": "admin"}

def test_account_lockout():
    """Test account lockout after multiple failed attempts."""
    # First make 3 failed login attempts
    for _ in range(3):
        client.post(
            "/login",
            data={"username": "user", "password": "wrong"},
            headers={"Content-Type": "application/x-www-form-urlencoded"}
        )
    
    # Then try with correct credentials
    response = client.post(
        "/login",
        data={"username": "user", "password": "pass"},
        headers={"Content-Type": "application/x-www-form-urlencoded"}
    )
    assert response.status_code == 429
    assert "Account locked" in response.json()["detail"]