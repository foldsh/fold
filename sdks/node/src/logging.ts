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
  if (foldStage == "DEV") {
    return winston.createLogger(baseConfig);
  } else if (foldStage == "PROD") {
    baseConfig.level = "info";
    return winston.createLogger(baseConfig);
  } else {
    baseConfig.format = combine(
      errors({ stack: true }),
      timestamp(),
      prettyPrint()
    );
    return winston.createLogger(baseConfig);
  }
}
