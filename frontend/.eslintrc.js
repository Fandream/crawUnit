module.exports = {
  root: true,
  env: {
    node: true
  },
  'extends': [
    'plugin:vue/essential',
    '@vue/standard'
  ],
  rules: {
    'no-console': process.env.NODE_ENV === 'production' ? 'error' : 'off',
    'no-debugger': process.env.NODE_ENV === 'production' ? 'error' : 'off',
    'semi': 0, // 去掉分号检查和未使用变量检查
    'no-unused-vars': 0
  },
  parserOptions: {
    parser: 'babel-eslint'
  },
  globals: {
    '_hmt': 1
  }
}
