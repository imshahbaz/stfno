from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
import httpx
import pyotp
import logging

app = FastAPI()

# --- DTOs ---
class FullLoginRequest(BaseModel):
    client_id: str
    password: str
    totp_secret: str

# --- Global Logic to get Public IP ---
async def get_public_ip():
    try:
        async with httpx.AsyncClient(timeout=5.0) as client:
            # We use ipify to get the external outbound IP
            response = await client.get("https://api.ipify.org?format=json")
            response.raise_for_status()
            ip_data = response.json()
            return ip_data.get("ip")
    except Exception as e:
        logging.error(f"Failed to fetch public IP: {e}")
        return "Unknown"

# --- Endpoints ---

@app.get("/system/whitelisted-ip")
async def show_ip():
    """
    Exposes the current outbound public IP of this server.
    Use this to update your Stoxkart Whitelist.
    """
    ip = await get_public_ip()
    print(f"\n-----------------------------------------")
    print(f"WHITELIST THIS IP: {ip}")
    print(f"-----------------------------------------\n")
    return {"public_ip": ip, "note": "Add this to Stoxkart portal"}

@app.post("/auth/generate-token")
async def generate_token(req: FullLoginRequest):
    async with httpx.AsyncClient(timeout=20.0) as client:
        # STEP 1: Initiate Login
        auth_url = "https://superrapi.stoxkart.com/v1/user/auth"
        init_payload = {
            "platform": "api",
            "data": {
                "client_id": req.client_id,
                "password": req.password
            }
        }
        
        resp1 = await client.post(auth_url, json=init_payload)
        data1 = resp1.json()
        
        if data1.get("status") != "success":
            raise HTTPException(status_code=401, detail=f"Step 1 Failed: {data1.get('message')}")
        
        req_token = data1.get("data", {}).get("request_token")
        
        # STEP 2: Generate TOTP and Validate Key
        validate_url = "https://superrapi.stoxkart.com/v1/user/auth/2fa"
        totp_code = pyotp.TOTP(req.totp_secret.replace(" ", "")).now()
        
        final_payload = {
            "platform": "api",
            "data": {
                "client_id": req.client_id,
                "req_token": req_token,
                "action": "api-key-validation",
                "otp": totp_code
            }
        }
        
        resp2 = await client.post(validate_url, json=final_payload)
        data2 = resp2.json()
        
        if data2.get("status") != "success":
            # If this fails, usually it's because the IP isn't whitelisted
            raise HTTPException(status_code=400, detail=f"Step 2 Failed: {data2.get('message')}")
        
        return {
            "status": "success",
            "access_token": data2.get("data", {}).get("access_token")
        }

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, port=8090)