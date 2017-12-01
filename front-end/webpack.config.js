module.exports = {
  entry: './index',
  output: {
    path: __dirname + 'dist',
    filename: 'bundle.js'
  },
  devtool: 'source-map',
  resolve: {
    extensions: ['*', '.js', '.jsx']
  },
  node: {
    fs: 'empty'
  },
  module: {
    rules: [
      {
        test: /\.jsx?$/,
        exclude: /(node_modules|bower_components)/,
        use: {
          loader: 'babel-loader',
          query: {
            presets: ['es2015', 'es2017', 'stage-0', 'react'],
            plugins: ['transform-runtime']
          }
        }
      }
    ]
  }
};
