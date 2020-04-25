import createError from "http-errors";
import express from "express";
import path from "path";
import cookieParser from "cookie-parser";
import logger from "morgan";
import sassMiddleware from "node-sass-middleware";

import indexRouter from "./routes/index";
import usersRouter from "./routes/users";
import connectRouter from "./routes/connection";
import authRouter from "./routes/auth";

const app = express();

// view engine setup
console.log(__dirname);
console.log("view:" + path.join(__dirname, "../views"));
app.set("views", path.join(__dirname, "../views"));
app.set("view engine", "pug");

app.use(logger("dev"));
app.use(express.json());
app.use(express.urlencoded({ extended: false }));
app.use(cookieParser());
app.use(sassMiddleware({
    src: path.join(__dirname, "public"),
    dest: path.join(__dirname, "public"),
    indentedSyntax: true, // true = .sass and false = .scss
    sourceMap: true
}));
app.use(express.static(path.join(__dirname, "public")));

app.use("/", indexRouter);
app.use("/auth", authRouter);
app.use("/users", usersRouter);
app.use("/connect", connectRouter);

// catch 404 and forward to error handler
app.use(function(req, res, next) {
    console.log(`404: ${req.originalUrl}`);
    next(createError(404));
});

// error handler
app.use(function(req, res) {
    // render the error page
    res.status(500);
    res.render("error");
});

export = app;
