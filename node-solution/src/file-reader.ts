import { createReadStream } from 'fs';
import { createInterface } from 'readline';
import { Fn } from './types';
import Logger from './logger';

export function createStreamFileReader(path: string) {
  const stream = createReadStream(path, {
    encoding: 'utf-8',
    flags: 'r'
  });

  const reader = createInterface({
    input: stream,
    crlfDelay: Infinity,
    historySize: 0,
  });

  stream.on('close', () => Logger.debug('Stream closing 🎬'));
  stream.on('error', (error) => {
    Logger.error(error);
    process.exit(1);
  });

  return {
    pause: () => reader.pause(),
    resume: () => reader.resume(),
    getData: (callback: (data: string) => void) => reader.on('line', callback),
    onEnd: (callback: Fn) => reader.on('close', callback)
  };
}