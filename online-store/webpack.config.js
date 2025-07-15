const webpack = require('webpack');
const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const Proxy = require('./proxy');
const fs = require('fs');
const { VueLoaderPlugin } = require('vue-loader');
// 定义文件夹的路径
const ROOT_PATH = path.resolve(__dirname);

module.exports = {
  devtool: 'source-map',
  entry: {
    index: './src/main.js',
  },
  output: {
    path: path.resolve(__dirname, 'build'),
    filename: '[name].[contenthash].entry.js',
    chunkFilename: '[name].[contenthash].min.js',
  },
  resolve: {
    extensions: ['.js', '.vue', '.json'],
    alias: {
      'vue$': 'vue/dist/vue.esm.js',
    },
  },
  plugins: [
    new VueLoaderPlugin(),
    new webpack.HotModuleReplacementPlugin(),
    new HtmlWebpackPlugin({
      filename: 'index.html',
      template: 'template.html',
      inject: true,
    }),
    new webpack.ProvidePlugin({
      $: 'jquery',
      jQuery: 'jquery',
      'window.jQuery': 'jquery',
      'window.$': 'jquery',
    }),
  ],
  devServer: {
    historyApiFallback: true,
    hot: true,
    port: 8089,
    proxy: Proxy,
    static: {
      directory: path.join(__dirname, 'public'),
      watch: true,
    },
  },
  module: {
    rules: [
      {
        test: /\.vue$/,
        loader: 'vue-loader',
      },
      {
        test: /\.css$/,
        use: ['style-loader', 'css-loader'],
      },
      {
        test: /\.scss$/,
        use: ['style-loader', 'css-loader', 'sass-loader'],
      },
      {
        test: /\.(png|jpe?g|gif|svg|jgp)(\?.*)?$/,
        type: 'asset',
        parser: {
          dataUrlCondition: {
            maxSize: 10 * 1024, // 10kb
          },
        },
        generator: {
          filename: 'static/images/[name].[hash:7][ext]'
        }
      },
      {
        test: /\.(woff2?|eot|ttf|otf)(\?.*)?$/,
        type: 'asset',
        parser: {
          dataUrlCondition: {
            maxSize: 10 * 1024,
          },
        },
        generator: {
          filename: 'static/fonts/[name].[hash:7][ext]'
        }
      },
      {
        test: /\.js$/,
        use: {
          loader: 'babel-loader',
          options: {
            compact: false,
          },
        },
        exclude: /node_modules/,
      },
      {
        test: /\.exec\.js$/,
        use: ['script-loader'],
      },
    ],
  },
};
