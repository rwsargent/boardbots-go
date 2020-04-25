// @ts-ignore
import {BoardbotsServiceClient} from "../pb/boardbots_grpc_pb";
import grpc, {Client} from "grpc";
let instance: Client;


function newClient(): Client {
    if (process.env["NODE_ENV"] === "development") {
        return new BoardbotsServiceClient(process.env["BOARDBOTS_ADDRESS"], grpc.credentials.createInsecure());
    } else {
        // TODO(rwsargent) make this a secure connection for production.
        return new BoardbotsServiceClient(process.env["BOARDBOTS_ADDRESS"], grpc.credentials.createInsecure());
    }
}

export function getClient(): Client {
    if(typeof instance === "undefined") {
        instance = newClient();
    }
    return instance;
}