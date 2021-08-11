const path = require("path");
const HtmlWebpackPlugin = require("html-webpack-plugin");

module.exports = {
  entry: "./src/js/index.js",

  output: {
    path: path.resolve(__dirname, "dist"),
    filename: "app.bundle.js",
  },

  module: {
    rules: [
      {
        test: /\.elm$/,
        loader: "elm-webpack-loader",
      },
    ],
  },

  plugins: [
    new HtmlWebpackPlugin({
      hash: true,
      filename: "index.html",
      template: "./src/index.html",
      title: "Orbital Data Platform",
    }),
  ],

  devServer: {
    inline: true,
    stats: { colors: true },
  },
};
