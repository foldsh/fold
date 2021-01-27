module.exports = {
  testEnvironment: "node",
  roots: ["<rootDir>/src"],
  transform: {
    "^.+\\.tsx?$": "ts-jest",
  },
  testRegex: "(\\.|/)(test|spec)\\.tsx?$",
  testPathIgnorePatterns: [
    "(\\.|/)(integration.test|integration.spec)\\.tsx?$",
  ],
  moduleFileExtensions: ["ts", "tsx", "js", "jsx", "json", "node"],
  globals: {
    "ts-jest": {
      diagnostics: {
        ignoreCodes: [151001],
      },
    },
  },
};
