const path = require("path");
const webpack = require("webpack");
const nodeExternals = require("webpack-node-externals");

const shouldWatch = (process.argv[2] || false) === "-w";

module.exports = {
    entry: {
        game: "./src/public/js/game.ts",
        auth: "./src/public/js/auth.ts"
    },
    devtool: 'inline-source-map',
    output: {
        path: path.join(__dirname, "dist/public/js"),
        publicPath: "/",
        filename: "[name].bundle.js"
    },
    mode: "development",
    watch: shouldWatch,
    target: "web",
    node: {
        // Need this when working with express, otherwise the build fails
        __dirname: false,   // if you don't put this is, __dirname
        __filename: false,  // and __filename return blank or /
    },
    externals: [nodeExternals()], // Need this to avoid error when working with Express
    module: {
        rules: [
            {
                test: /\.ts$/,
                use: [
                    {
                        loader: "ts-loader",
                        options: {
                            configFile: "tsconfig.browser.json"
                        }
                    }
                ],
                exclude: /node_modules/
            }
        ]
    },
    resolve: {
        extensions: [".ts", ".js"]
    }
};