import http from "k6/http";
import { check, sleep } from "k6";

export default function() {
    const BASE_URL = "http://localhost:8000";
    const SIGNUP_ENDPOINT = "/auth/sign-up";
    const SIGNIN_ENDPOINT = "/auth/sign-in";
    const API_ENDPOINT = "/api/list";

    // Step 1: Sign-up with random credentials
    const login = "user" + Math.floor(Math.random() * 999999999);
    const password = "password" + Math.floor(Math.random() * 999999999);

    const signupPayload = JSON.stringify({
        login: login,
        password: password
    });

    const signupParams = {
        headers: {
            "Content-Type": "application/json"
        }
    };

    const signupResponse = http.post(
        BASE_URL + SIGNUP_ENDPOINT,
        signupPayload,
        signupParams
    );

    check(signupResponse, {
        "signup status is 200": (r) => r.status === 200
    });

    // Step 2: Log-in with JSON credentials and parse JWT token
    const signinPayload = JSON.stringify({
        login: login,
        password: password
    });

    const signinParams = {
        headers: {
            "Content-Type": "application/json"
        }
    };

    const signinResponse = http.post(
        BASE_URL + SIGNIN_ENDPOINT,
        signinPayload,
        signinParams
    );

    check(signinResponse, {
        "signin status is 200": (r) => r.status === 200
    });

    const authToken = signinResponse.json("token");

    // Step 3: Send POST request to API endpoint with Bearer token and JSON body
    const apiPayload = JSON.stringify({
        title: "randomtitle"
    });

    const apiParams = {
        headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${authToken}`
        }
    };

    const apiResponse = http.post(
        BASE_URL + API_ENDPOINT,
        apiPayload,
        apiParams
    );

    check(apiResponse, {
        "api status is 200": (r) => r.status === 200
    });

    for (let i = 0; i < 10000; i++) {
        let listPOST = http.post(BASE_URL + API_ENDPOINT, apiPayload, apiParams);
        check(listPOST, { "API request successful": (r) => r.status === 200 });

        let listGET = http.get(BASE_URL + API_ENDPOINT, apiParams);
        check(listGET, { "API request successful": (r) => r.status === 200 });


    }
}