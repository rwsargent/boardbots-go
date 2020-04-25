import express from "express";
const router = express.Router();

/* GET home page. */
router.get("/", function(req: express.Request, res: express.Response) {
    res.render("game/game",
        {
            title: "Express Title" ,
            entryPoint: "game.bundle.js"});
});

export = router;