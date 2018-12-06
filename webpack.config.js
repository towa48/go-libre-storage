const path = require('path');
const webpack = require('webpack');
const UglifyJsPlugin = require('uglifyjs-webpack-plugin');
const VueLoaderPlugin = require('vue-loader/lib/plugin');

const sassResourcesLoader = {
  loader: 'sass-resources-loader',
  options: {
    resources: [
      path.resolve(__dirname, 'node_modules/bootstrap/scss/_functions.scss'),
      path.resolve(__dirname, 'node_modules/bootstrap/scss/_mixins.scss'),
      path.resolve(__dirname, 'node_modules/bootstrap/scss/_variables.scss'),
      path.resolve(__dirname, 'web/wwwroot/styles/_variables.scss'),
    ]
  }
}

const scssLoaders = [
  'vue-style-loader',
  'css-loader',
  'sass-loader',
  sassResourcesLoader
]

module.exports = (env, argv) => ({
  entry: {
    'welcome-bundle': './web/wwwroot/js/welcome.entry.js',
    'main-bundle': './web/wwwroot/js/main.entry.js',
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
        test: /\.css$/,
        use: [
          'vue-style-loader',
          'css-loader'
        ],
      },
      {
        test: /\.scss$/,
        use: scssLoaders,
      },
      {
        test: /\.vue$/,
        loader: 'vue-loader',
        options: {
          loaders: {
            'scss': scssLoaders
          }
        }
      }
    ]
  },
  plugins: argv.mode === 'production' ? [
    new webpack.DefinePlugin({
      'process.env': {
        NODE_ENV: '"production"'
      }
    }),
    new VueLoaderPlugin()
  ] : [
    new VueLoaderPlugin()
  ],
  resolve: {
    alias: {
      'vue$': 'vue/dist/vue.esm.js'
    }
  }
});
