import express from "express";
import connection from "../connection/connection";
import {constants} from "http2";

// Middleware to validate every request is authenticated. If the request doesn't have credentials associated with it,
// redirect to the login screen.
// This middleware is not designed to be attached at root.
export = (req: express.Request, res: express.Response, next: express.NextFunction) => {
    if (req.path.startsWith("/public")) {
        return next();
    }

    if (!req.cookies.SESSION) {
        if (!req.path.startsWith("/auth")) {
            return res.redirect("/auth/login");
        }
        // We're on the process for validation
        return next();
    }
    const server = connection(req);
    server.post("/auth/validate")
        .then(response => {
            if (response.status === constants.HTTP_STATUS_OK) {
                next();
            } else {
                return res.sendStatus(constants.HTTP_STATUS_UNAUTHORIZED);
            }
        }).catch(err => {
            console.log(err.response.data);
            return res.sendStatus(constants.HTTP_STATUS_UNAUTHORIZED);
        });
}