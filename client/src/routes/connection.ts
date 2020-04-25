import express from "express";
import {getClient} from "../bb-client/client";
import {GameRequest, GameResponse, UUID} from "../pb/boardbots_pb";
import {InterceptingCall, Metadata, RequesterBuilder} from "grpc";

const router = express.Router();

router.get("/connect", function(req: express.Request, res: express.Response) {
    console.log(req.query.id);
    const gameRequest = new GameRequest();
    const uuid = new UUID();
    uuid.setValue(req.query.id as string);
    gameRequest.setGameId(uuid);

    const interceptors: any[] = [
        (options: any, nextCall: (n: any) => InterceptingCall) => {
            return new InterceptingCall(nextCall(options), new RequesterBuilder().withSendMessage(
                (message, next) => {
                    console.log("This is my new interceptor!");
                    next(message);
                }
            ).withStart((md: Metadata, listener: object, next: (md: Metadata, listener: object) => any) => {
                md.add("token", "druid" + req.query.token);
                next(md, listener);
            }).build());
        }
    ];
    const opts = {
        interceptors: interceptors,
    };

    // @ts-ignore
    getClient().getGames(gameRequest, opts, (err: any, response: GameResponse) => {
        if (err) {
            console.log(err);
            res.send({error: err});
        }
        res.send({ uuid : response.getGameId().getValue()});
    });
});

export = router;