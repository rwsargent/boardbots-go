import express from "express";
import connection from "../connection/connection";

const router = express.Router();

// Simple path to login - Only render login form if the app is in DEV mode.
router.get("/login", (req: express.Request, res: express.Response) => {
    if(process.env["NODE_ENV"] !== "dev") {
        return res.sendStatus(403);
    }
    res.render("login",
        {
            title: "BoardBots" ,
            entryPoint: "auth.bundle.js",
            failed : Boolean(req.query.auth)});
});

// Authenticates dev credentials against backend server. Redirects to homepage on success.
router.post("/auth", (req: express.Request, res: express.Response) => {
    const username = req.body.name;
    const password = req.body.password;
    const conn = connection(req);
    conn.post("/auth/login ", {username : username, password: password})
        .then(response => {
            if (response.status === 200) {
                res.cookie("SESSION", response.data.token, {maxAge: 900000, httpOnly: true});
                res.redirect("/");
            }
        }).catch(err => {
            const code = err?.response?.status || 500;
            const string = err?.response?.data || "error";
            res.redirect("/auth/login?auth=fail")
        });
});

export = router;