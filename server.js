const express = require('express');
const { MConnect } = require('@mstock-mirae-asset/nodetradingapi-typeb');
const { authenticator } = require('otplib');
const { LRUCache } = require('lru-cache');
const axios = require('axios');
const app = express();
app.use(express.json());


const options = {
    max: 500, // Maximum 500 active trading sessions
    ttl: 1000 * 60 * 60 * 12, // Auto-delete sessions after 12 hours (sessions expire anyway)
};

// Use this instead of new Map()
const userSessions = new LRUCache(options);

// --- 1. LOGIN ROUTE (Multi-User) ---
app.post('/api/auth/login', async (req, res) => {
    const { apiKey, clientCode, password, totpSecret } = req.body;

    try {
        // Create a fresh instance for this specific user
        const userClient = new MConnect('https://api.mstock.trade', apiKey);

        // Generate TOTP using the user's specific secret
        const totp = authenticator.generate(totpSecret);

        const loginResponse = await userClient.login({
            clientcode: clientCode,
            password: password,
            totp: totp,
            state: 'multi-user-env'
        });

        if (loginResponse.status) {
            // Set the JWT inside this instance
            userClient.setAccessToken(loginResponse.data.jwtToken);
            userSessions.set(clientCode, userClient);
            res.json({ message: `Session started for ${clientCode}`, status: "Success", jwtToken: loginResponse.data.jwtToken });
        } else {
            res.status(401).json({ error: loginResponse.message });
        }
    } catch (err) {
        res.status(500).json({ error: err.message });
    }
});

// --- 2. TRADING ROUTE (Routed by clientCode) ---
app.post('/api/trade/order', async (req, res) => {
    const { clientCode, symbol, qty, side, token } = req.body;

    const userClient = userSessions.get(clientCode);

    if (!userClient) {
        return res.status(404).json({ error: "No active session found for this user." });
    }

    try {
        const order = await userClient.placeOrder({
            variety: 'NORMAL',
            tradingsymbol: symbol,
            symboltoken: token,
            exchange: 'NSE',
            transactiontype: side,
            ordertype: 'MARKET',
            quantity: qty.toString(),
            producttype: 'DELIVERY',
            price: '0'
        });
        res.json(order);
    } catch (err) {
        console.log(err);
        res.status(500).json({ error: err.message });
    }
});


app.listen(3000, () => console.log('Multi-user Trading API running on port 3000'));