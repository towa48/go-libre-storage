const path = require('path');

//const addWebpackAlias = alias => config => {
//    if (!config.resolve) {
//        config.resolve = {};
//    }
//    if (!config.resolve.alias) {
//        config.resolve.alias = {};
//    }
//    Object.assign(config.resolve.alias, alias);
//    return config;
//};

module.exports = {
    webpack: function (config, env) {
        const mainEntry = config.entry;
        config.entry = {
            'main': mainEntry, // the main key is required by react-scripts
            'welcome': path.resolve(__dirname, 'web/src/welcome.jsx')
        }
        return config;
    },
    jest: function (config) {
      return config;
    },
    devServer: function (configFunction) {
      return function (proxy, allowedHost) {
        return config;
      };
    },
    paths: function (paths, env) {
        paths['appPublic'] = path.resolve(__dirname, 'web/public');
        paths['appHtml'] = path.resolve(__dirname, 'web/public/index.html');
        paths['appIndexJs'] = path.resolve(__dirname, 'web/src/index.jsx');
        paths['appSrc'] = path.resolve(__dirname, 'web/src/');
        paths['appBuild'] = path.resolve(__dirname, 'web/build/');
        return paths;
    },
};