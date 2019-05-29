const path = require('path')
const webpack = require('webpack')
const HtmlWebpackPlugin = require('html-webpack-plugin')
const MiniCssExtractPlugin = require('mini-css-extract-plugin') 
const UglifyjsPlugin = require('uglifyjs-webpack-plugin')
const isProduction = process.env.NODE_ENV === 'production'

const config = {
    mode: process.env.NODE_ENV,
    entry: {
        app: path.join(__dirname, '../app.js')
    },
    output: {
        filename: 'js/[name].[hash:8].js',
        path: path.resolve(__dirname, '../dist'),
        chunkFilename: 'js/[name].[hash:8].js',
        publicPath: '/'
    },
    resolve: {
        extensions: ['.js', '.jsx']
    },
    module: {
        rules: [
            {
                test: /\.(js|jsx)$/,
                exclude: /node_modules/,
                use: {
                    loader: 'babel-loader',
                }
            },
            {
                test: /\.css$/,
                use: [
                    {loader: isProduction ? MiniCssExtractPlugin.loader : 'style-loader' },
                    {
                        loader: 'css-loader', 
                        options: {
                            modules: true,
                            localIdentName: '[path][name]__[local]--[hash:base64:5]',
                        }
                    }
                ]
            },
            {
                test: /\.scss$/,
                use: [
                    {loader: isProduction ? MiniCssExtractPlugin.loader : 'style-loader' },
                    {
                        loader: 'css-loader', 
                        options: {
                            modules: true,
                            localIdentName: '[path][name]__[local]--[hash:base64:5]',
                        }
                    },
                    {
                        loader: 'sass-loader'
                    }
                ]
            },
            {
                test: /\.(png|jpg|jpeg|gif)$/,
                use: [
                    { 
                        loader: 'url-loader',
                        options: {
                           limit: 10000,
                           outputPath: 'static'
                        }
                    }
                ]
            }
        ]
    },
    plugins: [
        new MiniCssExtractPlugin({
           filename: 'css/[name].[hash].css',
            chunkFilename: 'css/[id].[hash].css'
        }),
        new HtmlWebpackPlugin({
            filename: 'index.html',
            template: path.join(__dirname, '../index.html')
        }),
    ]
}

if(!isProduction) {
    config.devServer = {
        contentBase: path.join(__dirname, '../dist'),
        host: 'localhost',
        port: '8080',
        hot: true,
        overlay: {
            error: true
        },
        historyApiFallback: true,
        clientLogLevel: 'none'
    }
    config.devtool = '#cheap-module-eval-source-map'
    config.plugins.push(new webpack.HotModuleReplacementPlugin())
} else {
    config.optimization = {
        splitChunks: {
            chunks: 'async',
            minSize: 1000,
            maxSize: 0,
            cacheGroups: {
                vendor: {
                    test: /[\\/]node_modules[\\]/,
                    name: 'vendor',
                    chunks: 'async'
                }
            }
        },
        runtimeChunk: {
            name: 'manifest'
        },
        namedModules: true,
        namedChunks: true,
    }
    config.plugins.push(
        new webpack.DefinePlugin({
            'process.env.NODE_ENV': JSON.stringify('production')
        })
    )
    config.plugins.push(
        new UglifyjsPlugin({
            test: /\.(js|jsx)$/,
            exclude: /\/node_modules/
        })
    )
}

module.exports = config
