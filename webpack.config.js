const path = require('path');
const UglifyJsPlugin = require('uglifyjs-webpack-plugin');

module.exports = (env, argv) => ({
  entry: {
    'welcome-bundle': './web/wwwroot/js/welcome.entry.js',
  },
  //devtool: 'source-map',
  output: {
    filename: '[name].js',
    path: path.resolve(__dirname, './web/wwwroot/js/')
  },
  optimization: {
    minimize: argv.mode === 'production',
    minimizer: argv.mode === 'production' ? [
        new UglifyJsPlugin({
            cache: true,
            parallel: true,
            uglifyOptions: {
                compress: false,
                ecma: 6,
                mangle: true
            },
            sourceMap: true
        })
    ] : []
  },
  module: {
    rules: [
      {
        test: /\.html$/,
        use: [{
            loader: 'html-loader',
            options: { minimize: true }
        }]
      }
    ]
  },
  resolve: {
    alias: {
      'vue$': 'vue/dist/vue.esm.js'
    }
  }
});
