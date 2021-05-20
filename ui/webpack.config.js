// Webpack uses this to work with directories
const path = require('path');

// This is the main configuration object.
// Here you write different options and tell Webpack what to do
module.exports = {
  entry: './src/graph.js',

  output: {
    path: path.resolve(__dirname, 'dist'),
    filename: 'bundle.js'
  },

  // Default mode for Webpack is production.
  // Depending on mode Webpack will apply different things
  // on final bundle. For now we don't need production's JavaScript
  // minifying and other thing so let's set mode to development
  mode: 'development'
};