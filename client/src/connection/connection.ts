import axios, {AxiosInstance} from "axios";
import express from "express";

// Returns a AxiosInstance with default settings to the backend server.
const connection = (req: express.Request ): AxiosInstance => {
    return axios.create({
        baseURL: process.env["BACKEND"],
        headers :{
            "Content-Type" : "application/json",
            "Authorization" : "Bearer " + req.cookies.SESSION,
            "X-Client-Id" : process.env["CLIENT_ID"]
        }
    });
};

export = connection;