import express from "express";
const router = express.Router();

/* GET home page. */
router.get("/", function(req: express.Request, res: express.Response) {
    res.render("index",
        {
            title: "Boardbots" ,
            entryPoint: "game.bundle.js",
            user: req.cookies.SESSION});
});

export = router;