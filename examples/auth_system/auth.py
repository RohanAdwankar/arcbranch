from fastapi import FastAPI, Depends, HTTPException
from fastapi.security import OAuth2PasswordRequestForm

app = FastAPI()

@app.post("/login")
def login(form_data: OAuth2PasswordRequestForm = Depends()):
    # 1. Access the username and password from the form_data object.
    
    # 2. Check if the username is "user" and the password is "pass".
    
    # 3. If valid, return a JSON response: {"message": "Login successful"}.
    
    # 4. If invalid, raise an HTTPException with status code 401 and detail "Invalid credentials".
    ...
