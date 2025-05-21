from fastapi import FastAPI
from pydantic import BaseModel

app = FastAPI()

# Define a request model with a single field: interest (string)
class Profile(BaseModel):
    interest: str

@app.post("/recommend")
def recommend(profile: Profile):
    # 1. Access the 'interest' value from the request body.
    
    # 2. If the interest is "music", return a recommendation to "Listen to some chill beats".
    
    # 3. Otherwise, return a recommendation to "Read a good book".
    ...
