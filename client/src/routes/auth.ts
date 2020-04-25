import express from "express";
import {AuthRequest, AuthResponse} from "../pb/boardbots_pb";
import promises from "fs";
import {getClient} from "../bb-client/client";

const router = express.Router();


function newAuthRequest(username: string, password: string): AuthRequest {
    const authReq = new AuthRequest();
    authReq.setUsername(username);
    authReq.setPassword(password);
    return authReq;
}

function clientSideAuthentication(username: string, password: string): string {
    let path = process.env["USER_FILE_PATH"];
    console.log("promises: " + promises);
    const fileData = promises.readFileSync(path, "utf8");
    const users = JSON.parse(fileData);
    if (users[username] && users[username].password === password) {
        return "FAKETOKEN" + users[username].token;
    }
    return "";
}

function setCookie(res: express.Response, token: string) {
    res.cookie("bb-token", token, {
        httpOnly: false,
        expires : new Date()
    });
}

/* GET home page. */
router.get("/login", (req: express.Request, res: express.Response) => {
    res.render("login",
        {
            title: "BoardBots" ,
            entryPoint: "auth.bundle.js",
            failed : Boolean(req.query.auth)});
});

router.post("/auth", (req: express.Request, res: express.Response) => {
    const username = req.body.name;
    const password = req.body.password;
    //@ts-ignore
    getClient().authenticate(newAuthRequest(username, password), {}, (err, authRes: AuthResponse) => {
        if(err || (authRes && !authRes.getToken())) {
            let token = clientSideAuthentication(username, password);
            if(token) {
                setCookie(res, token);
                res.redirect("/");
            } else {
                res.redirect("/auth/login?auth=failed&err=notoken");
            }
        } else {
            const token = authRes.getToken();
            setCookie(res, token);
            res.redirect("/");
        }
    });
});

export = router;