from fastapi import FastAPI
from pydantic import BaseModel

app = FastAPI()

class Profile(BaseModel):
    interest: str

@app.post("/recommend")
def recommend(profile: Profile):
    if profile.interest == "music":
        return {"recommendation": "Listen to some chill beats"}
    return {"recommendation": "Read a good book"}