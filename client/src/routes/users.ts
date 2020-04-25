import {Response, Request, Router} from "express";

const router = Router();

/* GET users listing. */
router.get("/", function(req: Request, res: Response) {
    res.send("respond with a resource");
});

export = router;
