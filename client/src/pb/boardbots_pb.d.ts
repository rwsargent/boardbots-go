// package: 
// file: boardbots.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_timestamp_pb from "google-protobuf/google/protobuf/timestamp_pb";

export class UUID extends jspb.Message {
  getValue(): string;
  setValue(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UUID.AsObject;
  static toObject(includeInstance: boolean, msg: UUID): UUID.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UUID, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UUID;
  static deserializeBinaryFromReader(message: UUID, reader: jspb.BinaryReader): UUID;
}

export namespace UUID {
  export type AsObject = {
    value: string,
  }
}

export class GameRequest extends jspb.Message {
  hasGameId(): boolean;
  clearGameId(): void;
  getGameId(): UUID | undefined;
  setGameId(value?: UUID): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GameRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GameRequest): GameRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GameRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GameRequest;
  static deserializeBinaryFromReader(message: GameRequest, reader: jspb.BinaryReader): GameRequest;
}

export namespace GameRequest {
  export type AsObject = {
    gameId?: UUID.AsObject,
  }
}

export class PlayerState extends jspb.Message {
  getPlayerName(): string;
  setPlayerName(value: string): void;

  hasPawnPosition(): boolean;
  clearPawnPosition(): void;
  getPawnPosition(): Position | undefined;
  setPawnPosition(value?: Position): void;

  getBarriers(): number;
  setBarriers(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PlayerState.AsObject;
  static toObject(includeInstance: boolean, msg: PlayerState): PlayerState.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: PlayerState, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PlayerState;
  static deserializeBinaryFromReader(message: PlayerState, reader: jspb.BinaryReader): PlayerState;
}

export namespace PlayerState {
  export type AsObject = {
    playerName: string,
    pawnPosition?: Position.AsObject,
    barriers: number,
  }
}

export class Position extends jspb.Message {
  getRow(): number;
  setRow(value: number): void;

  getCol(): number;
  setCol(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Position.AsObject;
  static toObject(includeInstance: boolean, msg: Position): Position.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Position, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Position;
  static deserializeBinaryFromReader(message: Position, reader: jspb.BinaryReader): Position;
}

export namespace Position {
  export type AsObject = {
    row: number,
    col: number,
  }
}

export class Piece extends jspb.Message {
  getType(): Piece.TypeMap[keyof Piece.TypeMap];
  setType(value: Piece.TypeMap[keyof Piece.TypeMap]): void;

  hasPosition(): boolean;
  clearPosition(): void;
  getPosition(): Position | undefined;
  setPosition(value?: Position): void;

  getOwner(): number;
  setOwner(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Piece.AsObject;
  static toObject(includeInstance: boolean, msg: Piece): Piece.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Piece, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Piece;
  static deserializeBinaryFromReader(message: Piece, reader: jspb.BinaryReader): Piece;
}

export namespace Piece {
  export type AsObject = {
    type: Piece.TypeMap[keyof Piece.TypeMap],
    position?: Position.AsObject,
    owner: number,
  }

  export interface TypeMap {
    BARRIER: 0;
    PAWN: 1;
  }

  export const Type: TypeMap;
}

export class GameResponse extends jspb.Message {
  hasGameId(): boolean;
  clearGameId(): void;
  getGameId(): UUID | undefined;
  setGameId(value?: UUID): void;

  clearPlayersList(): void;
  getPlayersList(): Array<PlayerState>;
  setPlayersList(value: Array<PlayerState>): void;
  addPlayers(value?: PlayerState, index?: number): PlayerState;

  getCurrentTurn(): number;
  setCurrentTurn(value: number): void;

  hasStartDate(): boolean;
  clearStartDate(): void;
  getStartDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setStartDate(value?: google_protobuf_timestamp_pb.Timestamp): void;

  hasEndDate(): boolean;
  clearEndDate(): void;
  getEndDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setEndDate(value?: google_protobuf_timestamp_pb.Timestamp): void;

  getWinner(): number;
  setWinner(value: number): void;

  clearBoardList(): void;
  getBoardList(): Array<Piece>;
  setBoardList(value: Array<Piece>): void;
  addBoard(value?: Piece, index?: number): Piece;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GameResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GameResponse): GameResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GameResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GameResponse;
  static deserializeBinaryFromReader(message: GameResponse, reader: jspb.BinaryReader): GameResponse;
}

export namespace GameResponse {
  export type AsObject = {
    gameId?: UUID.AsObject,
    playersList: Array<PlayerState.AsObject>,
    currentTurn: number,
    startDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    endDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    winner: number,
    boardList: Array<Piece.AsObject>,
  }
}

export class AuthRequest extends jspb.Message {
  getUsername(): string;
  setUsername(value: string): void;

  getPassword(): string;
  setPassword(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AuthRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AuthRequest): AuthRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AuthRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AuthRequest;
  static deserializeBinaryFromReader(message: AuthRequest, reader: jspb.BinaryReader): AuthRequest;
}

export namespace AuthRequest {
  export type AsObject = {
    username: string,
    password: string,
  }
}

export class AuthResponse extends jspb.Message {
  getToken(): string;
  setToken(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AuthResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AuthResponse): AuthResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AuthResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AuthResponse;
  static deserializeBinaryFromReader(message: AuthResponse, reader: jspb.BinaryReader): AuthResponse;
}

export namespace AuthResponse {
  export type AsObject = {
    token: string,
  }
}

