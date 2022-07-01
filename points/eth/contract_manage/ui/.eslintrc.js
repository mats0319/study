module.exports = {
  root: true,
  env: {
    node: true
  },
  'extends': [
    'plugin:vue/essential',
    'eslint:recommended',
    '@vue/typescript/recommended'
  ],
  parserOptions: {
    ecmaVersion: 2020
  },
  rules: {
    "prettier/prettier": "off",
    "no-console": "off",
    "no-debugger": process.env.NODE_ENV === "production" ? "error" : "off",
    "no-inner-declarations": "off",
    "no-case-declarations": "off",
    "@typescript-eslint/no-var-requires": "off", // mainly for vue.config.js
    "@typescript-eslint/no-explicit-any": "off", // for ignore 'Promise<any>'
    "@typescript-eslint/no-unused-vars": "off", // for dev
    "@typescript-eslint/ban-ts-comment": "off", // for use '@ts-ignore'
    "@typescript-eslint/no-inferrable-types": "off", // for define class
    "prefer-const": "off", // for use 'let' at first without error
  }
}
