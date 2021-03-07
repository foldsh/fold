import * as winston from "winston";

export type LogMethod = (message: any) => void;

export interface Logger {
  debug: LogMethod;
  info: LogMethod;
  warn: LogMethod;
  error: LogMethod;
  crit: LogMethod;
}

const { combine, timestamp, prettyPrint, errors, json } = winston.format;

export function getLogger(service: string): Logger {
  const foldStage = process.env.FOLD_STAGE!;
  let baseConfig = {
    level: "debug",
    format: combine(
      errors({ stack: true }),
      timestamp(),
      prettyPrint(),
      json()
    ),
    defaultMeta: { service: service },
    transports: [new winston.transports.Console()],
  };
  switch (foldStage) {
    case "DEV":
      return winston.createLogger(baseConfig);
    case "PROD":
      baseConfig.level = "info";
      return winston.createLogger(baseConfig);
    case "TEST_LOCAL":
      return mockLogger();
    default:
      baseConfig.format = combine(
        errors({ stack: true }),
        timestamp(),
        prettyPrint()
      );
      return winston.createLogger(baseConfig);
  }
}

export function mockLogger(): Logger {
  return {
    debug: (_: any) => {},
    info: (_: any) => {},
    warn: (_: any) => {},
    error: (_: any) => {},
    crit: (_: any) => {},
  };
}
